package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type DatabaseHandler struct {
	dbHandle *pgxpool.Pool
	lock     sync.Mutex
}

// DbActions interface that the "pgx" library implements (to make it easier to test)
type DbActions interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Release()
}

func (dh *DatabaseHandler) Close() {
	dh.dbHandle.Close()
}

func (dh *DatabaseHandler) Ping() error {
	err := dh.dbHandle.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func InitializePool() (DatabaseHandler, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), configure())
	if err != nil {
		slog.Error("Error while creating connection to the database", "error", err)
		return DatabaseHandler{}, err
	}
	return DatabaseHandler{
		dbHandle: pool,
		lock:     sync.Mutex{},
	}, nil
}

// Acquire acquires a connection from the pool and returns an error if it fails
func (dh *DatabaseHandler) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	dh.lock.Lock()
	defer dh.lock.Unlock()
	return dh.dbHandle.Acquire(ctx)
}

// AcquireWithPanic acquires a connection from the pool and panics if it fails
func (dh *DatabaseHandler) AcquireWithPanic(ctx context.Context, w http.ResponseWriter) *pgxpool.Conn {
	dh.lock.Lock()
	defer dh.lock.Unlock()
	connection, err := dh.dbHandle.Acquire(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		panic(err)
	}
	return connection
}

func configure() *pgxpool.Config {
	const defaultMaxConns = int32(7)
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

	dbConfig.AfterConnect = func(ctx context.Context, connection *pgx.Conn) error {
		slog.Info("Connected to the DB")
		return nil
	}
	dbConfig.BeforeAcquire = func(ctx context.Context, connection *pgx.Conn) bool {
		slog.Debug("Acquiring connection from pool")
		return true
	}
	dbConfig.AfterRelease = func(connection *pgx.Conn) bool {
		slog.Debug("Released connection to the pool")
		return true
	}
	dbConfig.BeforeClose = func(connection *pgx.Conn) {
		slog.Info("Closed DB connection")
	}

	return dbConfig
}
