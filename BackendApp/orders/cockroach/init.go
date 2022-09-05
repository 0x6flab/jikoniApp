package cockroach

import (
	"database/sql"
	"fmt"

	"contrib.go.opencensus.io/integrations/ocsql"
	_ "github.com/jackc/pgx/v4/stdlib" // required for SQL access
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance.
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unappeased database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	// Register default views.
	ocsql.RegisterAllViews()

	// Register our ocsql wrapper for the provided Postgres driver.
	driverName, err := ocsql.Register("pgx", ocsql.WithAllTraceOptions())
	if err != nil {
		return nil, err
	}
	// Connect to a Postgres database using the ocsql driver wrapper.
	db, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	// Wrap our *sql.DB with sqlx. use the original db driver name!!!
	dbx := sqlx.NewDb(db, "pgx")
	if err := migrateDB(dbx); err != nil {
		return nil, err
	}

	return dbx, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "jikoni_1",
				Up: []string{
					`CREATE DATABASE IF NOT EXISTS jikoni;`,
					`CREATE TABLE IF NOT EXISTS orders (
						id 			UUID DEFAULT gen_random_uuid() PRIMARY KEY,
						name        STRING NOT NULL,
						price		INT8 NOT NULL,
						metadata    JSONB,
						status      STRING,
						created_at  TIMESTAMP DEFAULT now(),
						updated_at  TIMESTAMP DEFAULT now()
					);`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS orders;`,
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}
