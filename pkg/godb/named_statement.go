package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Defines the named statement
type NamedStmt interface {
	// sqlx

	// Close closes the named statement.
	Close() error
	// Exec executes a named statement using the struct passed.
	// Any named placeholder parameters are replaced with fields from arg.
	Exec(arg interface{}) (sql.Result, error)
	// ExecContext executes a named statement using the struct passed.
	// Any named placeholder parameters are replaced with fields from arg.
	ExecContext(ctx context.Context, arg interface{}) (sql.Result, error)
	// Get using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	Get(dest interface{}, arg interface{}) error
	// GetContext using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	GetContext(ctx context.Context, dest interface{}, arg interface{}) error
	// MustExec execs a NamedStmt, panicing on error
	// Any named placeholder parameters are replaced with fields from arg.
	MustExec(arg interface{}) sql.Result
	// MustExecContext execs a NamedStmt, panicing on error
	// Any named placeholder parameters are replaced with fields from arg.
	MustExecContext(ctx context.Context, arg interface{}) sql.Result
	// QueryRow this NamedStmt.  Because of limitations with QueryRow, this is
	// an alias for QueryRow.
	// Any named placeholder parameters are replaced with fields from arg.
	QueryRow(arg interface{}) Row
	// QueryRowContext this NamedStmt.  Because of limitations with QueryRow, this is
	// an alias for QueryRow.
	// Any named placeholder parameters are replaced with fields from arg.
	QueryRowContext(ctx context.Context, arg interface{}) Row
	// Query using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	Query(arg interface{}) (Rows, error)
	// QueryContext using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	QueryContext(ctx context.Context, arg interface{}) (Rows, error)
	// Select using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	Select(dest interface{}, arg interface{}) error
	// SelectContext using this NamedStmt
	// Any named placeholder parameters are replaced with fields from arg.
	SelectContext(ctx context.Context, dest interface{}, arg interface{}) error
	// Unsafe returns an unsafe version of the NamedStmt
	Unsafe() *sqlx.NamedStmt
	// Safe returns the underlying named statement
	Safe() *sqlx.NamedStmt
}

// Implements named statement interface
type customNamedStmt struct {
	testErr   error
	namedStmt *sqlx.NamedStmt
}

// pushTestError
func (c *customNamedStmt) pushTestError(err error) {
	c.testErr = err
}

// popTestError
func (c *customNamedStmt) popTestError() error {
	err := c.testErr
	c.testErr = nil
	return err
}

// Close
func (c *customNamedStmt) Close() error {
	return c.namedStmt.Close()
}

// Exec
func (c *customNamedStmt) Exec(arg interface{}) (sql.Result, error) {
	return c.namedStmt.Exec(arg)
}

// ExecContext
func (c *customNamedStmt) ExecContext(ctx context.Context, arg interface{}) (sql.Result, error) {
	return c.namedStmt.ExecContext(ctx, arg)
}

// Get
func (c *customNamedStmt) Get(dest interface{}, arg interface{}) error {
	return c.namedStmt.Get(dest, arg)
}

// GetContext
func (c *customNamedStmt) GetContext(ctx context.Context, dest interface{}, arg interface{}) error {
	return c.namedStmt.GetContext(ctx, dest, arg)
}

// MustExec
func (c *customNamedStmt) MustExec(arg interface{}) sql.Result {
	return c.namedStmt.MustExec(arg)
}

// MustExecContext
func (c *customNamedStmt) MustExecContext(ctx context.Context, arg interface{}) sql.Result {
	return c.namedStmt.MustExecContext(ctx, arg)
}

// QueryRow
func (c *customNamedStmt) QueryRow(arg interface{}) Row {
	return &customRow{row: c.namedStmt.QueryRowx(arg)}
}

// QueryRowContext
func (c *customNamedStmt) QueryRowContext(ctx context.Context, arg interface{}) Row {
	return &customRow{row: c.namedStmt.QueryRowxContext(ctx, arg)}
}

// Query
func (c *customNamedStmt) Query(arg interface{}) (Rows, error) {
	if rows, err := c.query(arg); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// query
func (c *customNamedStmt) query(arg interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.namedStmt.Queryx(arg)
}

// QueryContext
func (c *customNamedStmt) QueryContext(ctx context.Context, arg interface{}) (Rows, error) {
	if rows, err := c.queryContext(ctx, arg); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// queryContext
func (c *customNamedStmt) queryContext(ctx context.Context, arg interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.namedStmt.QueryxContext(ctx, arg)
}

// Select
func (c *customNamedStmt) Select(dest interface{}, arg interface{}) error {
	return c.namedStmt.Select(dest, arg)
}

// SelectContext
func (c *customNamedStmt) SelectContext(ctx context.Context, dest interface{}, arg interface{}) error {
	return c.namedStmt.SelectContext(ctx, dest, arg)
}

// Unsafe
func (c *customNamedStmt) Unsafe() *sqlx.NamedStmt {
	return c.namedStmt.Unsafe()
}

// Safe
func (c *customNamedStmt) Safe() *sqlx.NamedStmt {
	return c.namedStmt
}
