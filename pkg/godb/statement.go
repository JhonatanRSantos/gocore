package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Stmt defines a new statement
type Stmt interface {
	// sql

	// Close closes the statement.
	Close() error
	// Exec executes a prepared statement with the given arguments and
	// returns a Result summarizing the effect of the statement.
	// Exec uses context.Background internally; to specify the context, use
	// ExecContext.
	Exec(args ...any) (sql.Result, error)
	// ExecContext executes a prepared statement with the given arguments and
	// returns a Result summarizing the effect of the statement.
	ExecContext(ctx context.Context, args ...any) (sql.Result, error)

	// sqlx

	// Get using the prepared statement.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, args ...interface{}) error
	// GetContext using the prepared statement.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	// MustExec (panic) using this statement. Note that the query portion of the error
	// output will be blank, as Stmt does not expose its query.
	// Any placeholder parameters are replaced with supplied args.
	MustExec(args ...interface{}) sql.Result
	// MustExecContext (panic) using this statement. Note that the query portion of
	// the error output will be blank, as Stmt does not expose its query.
	// Any placeholder parameters are replaced with supplied args.
	MustExecContext(ctx context.Context, args ...interface{}) sql.Result
	// QueryRow using this statement.
	// Any placeholder parameters are replaced with supplied args.
	QueryRow(args ...interface{}) Row
	// QueryRowContext using this statement.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowContext(ctx context.Context, args ...interface{}) Row
	// Query using this statement.
	// Any placeholder parameters are replaced with supplied args.
	Query(args ...interface{}) (Rows, error)
	// QueryContext using this statement.
	// Any placeholder parameters are replaced with supplied args.
	QueryContext(ctx context.Context, args ...interface{}) (Rows, error)
	// Select using the prepared statement.
	// Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, args ...interface{}) error
	// SelectContext using the prepared statement.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
	// Unsafe returns a version of Stmt which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	Unsafe() *sqlx.Stmt
	// Safe returns the underlying Stmt
	Safe() *sqlx.Stmt
}

// customStmt implements the Stmt interface
type customStmt struct {
	stmt    *sqlx.Stmt
	testErr error
}

// pushTestError
func (c *customStmt) pushTestError(err error) {
	c.testErr = err
}

// popTestError
func (c *customStmt) popTestError() error {
	err := c.testErr
	c.testErr = nil
	return err
}

// Close
func (c *customStmt) Close() error {
	return c.stmt.Close()
}

// Exec
func (c *customStmt) Exec(args ...any) (sql.Result, error) {
	return c.stmt.Exec(args...)
}

// ExecContext
func (c *customStmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return c.stmt.ExecContext(ctx, args...)
}

// Get
func (c *customStmt) Get(dest interface{}, args ...interface{}) error {
	return c.stmt.Get(dest, args...)
}

// GetContext
func (c *customStmt) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	return c.stmt.GetContext(ctx, dest, args...)
}

// MustExec
func (c *customStmt) MustExec(args ...interface{}) sql.Result {
	return c.stmt.MustExec(args...)
}

// MustExecContext
func (c *customStmt) MustExecContext(ctx context.Context, args ...interface{}) sql.Result {
	return c.stmt.MustExecContext(ctx, args...)
}

// QueryRow
func (c *customStmt) QueryRow(args ...interface{}) Row {
	return &customRow{row: c.stmt.QueryRowx(args...)}
}

// QueryRowContext
func (c *customStmt) QueryRowContext(ctx context.Context, args ...interface{}) Row {
	return &customRow{row: c.stmt.QueryRowxContext(ctx, args...)}
}

// Query
func (c *customStmt) Query(args ...interface{}) (Rows, error) {
	if rows, err := c.query(args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// query
func (c *customStmt) query(args ...interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.stmt.Queryx(args...)
}

// QueryContext
func (c *customStmt) QueryContext(ctx context.Context, args ...interface{}) (Rows, error) {
	if rows, err := c.queryContext(ctx, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// queryContext
func (c *customStmt) queryContext(ctx context.Context, args ...interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.stmt.QueryxContext(ctx, args...)
}

// Select
func (c *customStmt) Select(dest interface{}, args ...interface{}) error {
	return c.stmt.Select(dest, args...)
}

// SelectContext
func (c *customStmt) SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	return c.stmt.SelectContext(ctx, dest, args...)
}

// Unsafe
func (c *customStmt) Unsafe() *sqlx.Stmt {
	return c.stmt.Unsafe()
}

// Safe
func (c *customStmt) Safe() *sqlx.Stmt {
	return c.stmt
}
