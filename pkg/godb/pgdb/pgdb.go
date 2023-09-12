package pgdb

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connTimeout           = time.Second * 30
	connMaxLifeTime       = time.Second * 30
	connMaxLifeTimeJitter = time.Second * 2
	connMaxIdleTime       = time.Second * 10
	maxConn               = 10
	minConn               = 2
	healthCheckPeriod     = time.Second * 5
)

type ClientConnection interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Conn() *pgx.Conn
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Release()
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type Client interface {
	Close()
	GetConnection(ctx context.Context) (ClientConnection, error)
}

type ClientConfig struct {
	Host   string
	Port   uint16
	User   string
	Pass   string
	DBName string
}

type DBClient struct {
	connPool *pgxpool.Pool
}

// NewClient Returns a new databse client
func NewClient(ctx context.Context, config ClientConfig) (Client, error) {
	pgPoolConfig, _ := pgxpool.ParseConfig("")

	pgPoolConfig.ConnConfig.Host = config.Host
	pgPoolConfig.ConnConfig.Port = config.Port
	pgPoolConfig.ConnConfig.User = config.User
	pgPoolConfig.ConnConfig.Password = config.Pass
	pgPoolConfig.ConnConfig.Database = config.DBName
	pgPoolConfig.ConnConfig.ConnectTimeout = connTimeout
	pgPoolConfig.ConnConfig.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	pgPoolConfig.MaxConnLifetime = connMaxLifeTime
	pgPoolConfig.MaxConnLifetimeJitter = connMaxLifeTimeJitter
	pgPoolConfig.MaxConnIdleTime = connMaxIdleTime
	pgPoolConfig.MaxConns = maxConn
	pgPoolConfig.MinConns = minConn
	pgPoolConfig.HealthCheckPeriod = healthCheckPeriod
	connPool, err := pgxpool.NewWithConfig(ctx, pgPoolConfig)

	if err != nil {
		return nil, fmt.Errorf("%s. %w", ErrFailedToCreateConnectionPool, err)
	}

	if err = connPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s. %w", ErrFailedToPingDatabase, err)
	}

	if err != nil {
		return nil, err
	}

	return &DBClient{
		connPool: connPool,
	}, nil
}

// Close Closes the connection pool
func (c *DBClient) Close() {
	c.connPool.Close()
}

// GetConnection Get a connection from pool
func (c *DBClient) GetConnection(ctx context.Context) (ClientConnection, error) {
	return c.connPool.Acquire(ctx)
}
