package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Conn defines a new connection
type Conn interface {
	// sql

	// Close returns the connection to the connection pool.
	// All operations after a Close will return with ErrConnDone.
	// Close is safe to call concurrently with other operations and will
	// block until all other operations finish. It may be useful to first
	// cancel any used context and then call close directly after.
	Close() error
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// PingContext verifies the connection to the database is still alive.
	PingContext(ctx context.Context) error
	// Raw executes f exposing the underlying driver connection for the
	// duration of f. The driverConn must not be used outside of f.
	//
	// Once f returns and err is not driver.ErrBadConn, the Conn will continue to be usable
	// until Conn.Close is called.
	Raw(f func(driverConn any) error) (err error)

	// sqlx

	// BeginTx begins a transaction and returns an godb.Tx instead of an *sql.Tx.
	//
	// The provided context is used until the transaction is committed or rolled
	// back. If the context is canceled, the sql package will roll back the
	// transaction. Tx.Commit will return an error if the context provided to
	// BeginContext is canceled.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	// GetContext using this godb.Conn.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// PrepareContext returns an godb.Stmt instead of a sql.Stmt.
	//
	// The provided context is used for the preparation of the statement, not for
	// the execution of the statement.
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	// QueryRowContext queries the database and returns an godb.Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	// QueryContext queries the database and returns an godb.Rows.
	// Any placeholder parameters are replaced with supplied args.
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	// Rebind a query within a Conn's bindvar type.
	Rebind(query string) string
	// SelectContext using this godb.Conn.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// customConn defines a new custom connection
type customConn struct {
	testErr error
	conn    *sqlx.Conn
}

// pushTestError
func (c *customConn) pushTestError(err error) {
	c.testErr = err
}

// popTestError
func (c *customConn) popTestError() error {
	err := c.testErr
	c.testErr = nil
	return err
}

// Close
func (c *customConn) Close() error {
	return c.conn.Close()
}

// ExecContext
func (c *customConn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.conn.ExecContext(ctx, query, args...)
}

// PingContext
func (c *customConn) PingContext(ctx context.Context) error {
	return c.conn.PingContext(ctx)
}

// Raw
func (c *customConn) Raw(f func(driverConn any) error) (err error) {
	return c.conn.Raw(f)
}

// BeginTx
func (c *customConn) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	if tx, err := c.beginTx(ctx, opts); err != nil {
		return &customTx{}, err
	} else {
		return &customTx{tx: tx}, nil
	}
}

// beginTx
func (c *customConn) beginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	if c.testErr != nil {
		return &sqlx.Tx{}, c.popTestError()
	}
	return c.conn.BeginTxx(ctx, opts)
}

// GetContext
func (c *customConn) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.conn.GetContext(ctx, dest, query, args...)
}

// PrepareContext
func (c *customConn) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if stmt, err := c.prepareContext(ctx, query); err != nil {
		return &customStmt{}, err
	} else {
		return &customStmt{stmt: stmt}, nil
	}
}

// prepareContext
func (c *customConn) prepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	if c.testErr != nil {
		return &sqlx.Stmt{}, c.popTestError()
	}
	return c.conn.PreparexContext(ctx, query)
}

// QueryRowContext
func (c *customConn) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return &customRow{row: c.conn.QueryRowxContext(ctx, query, args...)}
}

// QueryContext
func (c *customConn) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if rows, err := c.queryContext(ctx, query, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// queryContext
func (c *customConn) queryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.conn.QueryxContext(ctx, query, args...)
}

// Rebind
func (c *customConn) Rebind(query string) string {
	return c.conn.Rebind(query)
}

// SelectContext
func (c *customConn) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.conn.SelectContext(ctx, dest, query, args...)
}
