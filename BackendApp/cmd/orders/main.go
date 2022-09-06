package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	fama "github.com/0x6flab/jikoniApp/BackendApp"
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/orders"
	ordersapi "github.com/0x6flab/jikoniApp/BackendApp/orders/api"
	"github.com/0x6flab/jikoniApp/BackendApp/orders/cockroach"
	"github.com/0x6flab/jikoniApp/BackendApp/orders/ocmux"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/sync/errgroup"
)

const (
	stopWaitTime  = 5 * time.Second
	svcName       = "users_service"
	defLogLevel   = "error"
	defDBURL      = "postgresql://root@127.0.0.1:26000/defaultdb?sslmode=disable"
	defHTTPPort   = "8180"
	defServerCert = ""
	defServerKey  = ""
	defZipkinURL  = "http://jikoni-zipkin:9411/api/v2/spans"
	envLogLevel   = "JIKONI_LOG_LEVEL"
	envDBURL      = "JIKONI_DB_URL"
	envHTTPPort   = "JIKONI_HTTP_PORT"
	envServerCert = "JIKONI_SERVER_CERT"
	envServerKey  = "JIKONI_SERVER_KEY"
	envZipkinURL  = "JIKONI_ZIPKIN_URL"
)

type config struct {
	logLevel   string
	dbURL      string
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

	db := connectToDB(cfg, logger)
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
	return config{
		logLevel:   fama.Env(envLogLevel, defLogLevel),
		dbURL:      fama.Env(envDBURL, defDBURL),
		httpPort:   fama.Env(envHTTPPort, defHTTPPort),
		serverCert: fama.Env(envServerCert, defServerCert),
		serverKey:  fama.Env(envServerKey, defServerKey),
		zipkinURL:  fama.Env(envZipkinURL, defZipkinURL),
	}
}

func connectToDB(dbConfig config, logger kitlog.Logger) *sqlx.DB {
	db, err := cockroach.Connect(dbConfig.dbURL)
	if err != nil {
		if err := logger.Log("service", svcName, "message", "Failed to connect to postgres", "error", err); err != nil {
			return nil
		}
		os.Exit(1)
	}
	return db
}

func newService(db *sqlx.DB, logger kitlog.Logger) orders.OrderService {
	ordersRepo := cockroach.NewOrderRepo(db)
	svc := orders.NewOrderService(ordersRepo)
	svc = ordersapi.LoggingMiddleware(svc, kitlog.With(logger, "component", svcName))
	svc = ordersapi.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: svcName,
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: svcName,
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(ctx context.Context, svc orders.OrderService, config config, logger kitlog.Logger) error {
	p := fmt.Sprintf(":%s", config.httpPort)
	errCh := make(chan error)
	router := mux.NewRouter()
	ordersapi.MakeOrdersHandler(svc, router, logger)
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
