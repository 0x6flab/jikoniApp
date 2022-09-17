package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	fama "github.com/0x6flab/jikoniApp/BackendApp"
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	shopsapi "github.com/0x6flab/jikoniApp/BackendApp/shops/api"
	"github.com/0x6flab/jikoniApp/BackendApp/shops/ocmux"
	"github.com/0x6flab/jikoniApp/BackendApp/shops/postgres"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/sync/errgroup"
)

const (
	stopWaitTime     = 5 * time.Second
	svcName          = "jikoni-shops"
	defLogLevel      = "error"
	defDBHost        = "jikoni-db"
	defDBPort        = "5441"
	defDBUser        = "jikoniuser"
	defDBPass        = "jikonipass"
	defDB            = "jikoni"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defHTTPPort      = "8180"
	defServerCert    = ""
	defServerKey     = ""
	defZipkinURL     = "http://jikoni-zipkin:9411/api/v2/spans"
	envLogLevel      = "JIKONI_LOG_LEVEL"
	envDBHost        = "JIKONI_DB_HOST"
	envDBPort        = "JIKONI_ORDERS_DB_PORT"
	envDBUser        = "JIKONI_DB_USER"
	envDBPass        = "JIKONI_DB_PASS"
	envDB            = "JIKONI_DB"
	envDBSSLMode     = "JIKONI_DB_SSL_MODE"
	envDBSSLCert     = "JIKONI_DB_SSL_CERT"
	envDBSSLKey      = "JIKONI_DB_SSL_KEY"
	envDBSSLRootCert = "JIKONI_DB_SSL_ROOT_CERT"
	envHTTPPort      = "JIKONI_SHOPS_HTTP_PORT"
	envServerCert    = "JIKONI_SERVER_CERT"
	envServerKey     = "JIKONI_SERVER_KEY"
	envZipkinURL     = "JIKONI_ZIPKIN_URL"
)

type config struct {
	logLevel   string
	dbConfig   postgres.Config
	httpPort   string
	serverCert string
	serverKey  string
	zipkinURL  string
}

func main() {
	cfg := loadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// Set up our contextual logger.
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
		logger = kitlog.NewSyncLogger(logger)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
		logger = kitlog.With(logger, "svc", svcName)
	}
	defer ocmux.InitOpenCensusWithZipkin(cfg.zipkinURL, svcName, fmt.Sprintf("%s:%s", svcName, cfg.httpPort)).Close()

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	svc := newService(db, logger)

	g.Go(func() error {
		return startHTTPServer(ctx, svc, cfg, logger)
	})

	g.Go(func() error {
		if sig := errors.SignalHandler(ctx); sig != nil {
			cancel()
			if err := logger.Log("service", svcName, "message", fmt.Sprintf("%s service shutdown by signal", svcName), "signal", sig); err != nil {
				return err
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		if err := logger.Log("service", svcName, "message", fmt.Sprintf("%s service terminated", svcName), "error", err); err != nil {
			return
		}
	}
}

func loadConfig() config {
	dbConfig := postgres.Config{
		Host:        fama.Env(envDBHost, defDBHost),
		Port:        fama.Env(envDBPort, defDBPort),
		User:        fama.Env(envDBUser, defDBUser),
		Pass:        fama.Env(envDBPass, defDBPass),
		Name:        fama.Env(envDB, defDB),
		SSLMode:     fama.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     fama.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      fama.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: fama.Env(envDBSSLRootCert, defDBSSLRootCert),
	}
	return config{
		logLevel:   fama.Env(envLogLevel, defLogLevel),
		dbConfig:   dbConfig,
		httpPort:   fama.Env(envHTTPPort, defHTTPPort),
		serverCert: fama.Env(envServerCert, defServerCert),
		serverKey:  fama.Env(envServerKey, defServerKey),
		zipkinURL:  fama.Env(envZipkinURL, defZipkinURL),
	}
}

func connectToDB(dbConfig postgres.Config, logger kitlog.Logger) *sqlx.DB {
	db, err := postgres.Connect(dbConfig)
	if err != nil {
		if err := logger.Log("service", svcName, "message", "Failed to connect to postgres", "error", err); err != nil {
			return nil
		}
		os.Exit(1)
	}
	return db
}

func newService(db *sqlx.DB, logger kitlog.Logger) shops.ShopService {
	shopsRepo := postgres.NewShopRepo(db)
	svc := shops.NewShopService(shopsRepo)
	svc = shopsapi.LoggingMiddleware(svc, kitlog.With(logger, "component", svcName))
	svc = shopsapi.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: strings.Replace(svcName, "-", "_", 1),
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: strings.Replace(svcName, "-", "_", 1),
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(ctx context.Context, svc shops.ShopService, config config, logger kitlog.Logger) error {
	p := fmt.Sprintf(":%s", config.httpPort)
	errCh := make(chan error)
	router := mux.NewRouter()
	shopsapi.MakeShopsHandler(svc, router, logger)
	handler := &ochttp.Handler{Handler: router}
	server := &http.Server{Addr: p, Handler: handler}

	switch {
	case config.serverCert != "" || config.serverKey != "":
		if err := logger.Log("transport", svcName, "message", fmt.Sprintf("%s service started using https", svcName), "exposed_port", config.httpPort, "cert", config.serverCert, "key", config.serverKey); err != nil {
			return err
		}
		go func() {
			errCh <- server.ListenAndServeTLS(config.serverCert, config.serverKey)
		}()
	default:
		if err := logger.Log("transport", svcName, "message", fmt.Sprintf("%s service started using http", svcName), "exposed_port", config.httpPort); err != nil {
			return err
		}
		go func() {
			errCh <- server.ListenAndServe()
		}()
	}

	select {
	case <-ctx.Done():
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), stopWaitTime)
		defer cancelShutdown()
		if err := server.Shutdown(ctxShutdown); err != nil {
			if err := logger.Log("transport", svcName, "message", fmt.Sprintf("%s service error occurred during shutdown", svcName), "exposed_port", config.httpPort, "error", err); err != nil {
				return err
			}
			return fmt.Errorf("%s service occurred during shutdown at %s: %w", svcName, p, err)
		}
		if err := logger.Log("transport", svcName, "message", fmt.Sprintf("%s service shutdown of http", svcName), "exposed_port", config.httpPort); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}
