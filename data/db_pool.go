package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

var fetcher *pgxpool.Pool

// DbActions interface that the "pgx" library implements (to make it easier to test)
type DbActions interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func GetFetcher() *pgxpool.Pool {
	return fetcher
}

func InitializePool() {
	pool, err := pgxpool.NewWithConfig(context.Background(), configure())
	if err != nil {
		slog.Error("Error while creating connection to the database", "error", err)
	}
	fetcher = pool
}

func configure() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// TODO export to env
	const DATABASE_URL string = "postgres://auth_mantle_manager:dudde@localhost:5432/authmantledb?"

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		slog.Error("Failed to create a config, error: ", "error", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		slog.DebugContext(ctx, "Acquiring connection from pool")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		slog.Debug("Released connection to the pool")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		slog.Debug("Closed DB connection")
	}

	return dbConfig
}
