package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Tx defines a new transaction interface
type Tx interface {
	// sql

	// Commit commits the transaction.
	Commit() error
	// Exec executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	//
	// Exec uses context.Background internally; to specify the context, use
	// ExecContext.
	Exec(query string, args ...any) (sql.Result, error)
	// ExecContext executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Rollback aborts the transaction.
	Rollback() error

	// sqlx.Tx

	// BindNamed binds a query within a transaction's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	// DriverName returns the driverName used by the DB which began this transaction.
	DriverName() string
	// Get within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error
	// GetContext within a transaction and context.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// MustExec runs MustExec within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	MustExec(query string, args ...interface{}) sql.Result
	// MustExecContext runs MustExecContext within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	// NamedExec a named query within a transaction.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)
	// NamedExecContext using this Tx.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	// NamedQuery within a transaction.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (Rows, error)
	// NamedStmt returns a version of the prepared statement which runs within a transaction.
	NamedStmt(stmt NamedStmt) NamedStmt
	// NamedStmtContext returns a version of the prepared statement which runs
	// within a transaction.
	NamedStmtContext(ctx context.Context, stmt NamedStmt) NamedStmt
	// PrepareNamed returns an NamedStmt
	PrepareNamed(query string) (NamedStmt, error)
	// PrepareNamedContext returns an NamedStmt
	PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error)
	// Prepare  a statement within a transaction.
	Prepare(query string) (Stmt, error)
	// PrepareContext returns an godb.Stmt instead of a sql.Stmt.
	//
	// The provided context is used for the preparation of the statement, not for
	// the execution of the statement.
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	// QueryRow within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	QueryRow(query string, args ...interface{}) Row
	// QueryRowContext within a transaction and context.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	// Query within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	Query(query string, args ...interface{}) (Rows, error)
	// QueryContext within a transaction and context.
	// Any placeholder parameters are replaced with supplied args.
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	// Rebind a query within a transaction's bindvar type.
	Rebind(query string) string
	// Select within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error
	// SelectContext within a transaction and context.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// Stmt returns a version of the prepared statement which runs within a transaction.  Provided
	// stmt can be either *sql.Stmt or *sqlx.Stmt.
	Stmt(stmt interface{}) Stmt
	// StmtContext returns a version of the prepared statement which runs within a
	// transaction. Provided stmt can be either *sql.Stmt or *sqlx.Stmt.
	StmtContext(ctx context.Context, stmt interface{}) Stmt
	// Unsafe returns a version of Tx which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	Unsafe() *sqlx.Tx
	// Safe returns the underlying Tx
	Safe() *sqlx.Tx
}

// customTx implements the Tx interface
type customTx struct {
	tx      *sqlx.Tx
	testErr error
}

// pushTestError
func (c *customTx) pushTestError(err error) {
	c.testErr = err
}

// popTestError
func (c *customTx) popTestError() error {
	err := c.testErr
	c.testErr = nil
	return err
}

// Commit
func (c *customTx) Commit() error {
	return c.tx.Commit()
}

// Exec
func (c *customTx) Exec(query string, args ...any) (sql.Result, error) {
	return c.tx.Exec(query, args...)
}

// ExecContext
func (c *customTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.tx.ExecContext(ctx, query, args...)
}

// Rollback
func (c *customTx) Rollback() error {
	return c.tx.Rollback()
}

// BindNamed
func (c *customTx) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return c.tx.BindNamed(query, arg)
}

// DriverName
func (c *customTx) DriverName() string {
	return c.tx.DriverName()
}

// Get
func (c *customTx) Get(dest interface{}, query string, args ...interface{}) error {
	return c.tx.Get(dest, query, args...)
}

// GetContext
func (c *customTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.tx.GetContext(ctx, dest, query, args...)
}

// MustExec
func (c *customTx) MustExec(query string, args ...interface{}) sql.Result {
	return c.tx.MustExec(query, args...)
}

// MustExecContext
func (c *customTx) MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
	return c.tx.MustExecContext(ctx, query, args...)
}

// NamedExec
func (c *customTx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return c.tx.NamedExec(query, arg)
}

// NamedExecContext
func (c *customTx) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return c.tx.NamedExecContext(ctx, query, arg)
}

// NamedQuery
func (c *customTx) NamedQuery(query string, arg interface{}) (Rows, error) {
	if rows, err := c.namedQuery(query, arg); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// namedQuery
func (c *customTx) namedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.tx.NamedQuery(query, arg)
}

// NamedStmt
func (c *customTx) NamedStmt(stmt NamedStmt) NamedStmt {
	return &customNamedStmt{namedStmt: c.tx.NamedStmt(stmt.Safe())}
}

// NamedStmtContext
func (c *customTx) NamedStmtContext(ctx context.Context, stmt NamedStmt) NamedStmt {
	return &customNamedStmt{namedStmt: c.tx.NamedStmtContext(ctx, stmt.Safe())}
}

// PrepareNamed
func (c *customTx) PrepareNamed(query string) (NamedStmt, error) {
	if namedStmt, err := c.prepareNamed(query); err != nil {
		return &customNamedStmt{}, err
	} else {
		return &customNamedStmt{namedStmt: namedStmt}, nil
	}
}

// prepareNamed
func (c *customTx) prepareNamed(query string) (*sqlx.NamedStmt, error) {
	if c.testErr != nil {
		return &sqlx.NamedStmt{}, c.popTestError()
	}
	return c.tx.PrepareNamed(query)
}

// PrepareNamedContext
func (c *customTx) PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error) {
	if namedStmt, err := c.prepareNamedContext(ctx, query); err != nil {
		return &customNamedStmt{}, err
	} else {
		return &customNamedStmt{namedStmt: namedStmt}, nil
	}
}

// prepareNamedContext
func (c *customTx) prepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if c.testErr != nil {
		return &sqlx.NamedStmt{}, c.popTestError()
	}
	return c.tx.PrepareNamedContext(ctx, query)
}

// Prepare
func (c *customTx) Prepare(query string) (Stmt, error) {
	if stmt, err := c.prepare(query); err != nil {
		return &customStmt{}, err
	} else {
		return &customStmt{stmt: stmt}, nil
	}
}

// prepare
func (c *customTx) prepare(query string) (*sqlx.Stmt, error) {
	if c.testErr != nil {
		return &sqlx.Stmt{}, c.popTestError()
	}
	return c.tx.Preparex(query)
}

// PrepareContext
func (c *customTx) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if stmt, err := c.prepareContext(ctx, query); err != nil {
		return &customStmt{}, err
	} else {
		return &customStmt{stmt: stmt}, nil
	}
}

// prepareContext
func (c *customTx) prepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	if c.testErr != nil {
		return &sqlx.Stmt{}, c.popTestError()
	}
	return c.tx.PreparexContext(ctx, query)
}

// QueryRow
func (c *customTx) QueryRow(query string, args ...interface{}) Row {
	return &customRow{row: c.tx.QueryRowx(query, args...)}
}

// QueryRowContext
func (c *customTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return &customRow{row: c.tx.QueryRowxContext(ctx, query, args...)}
}

// Query
func (c *customTx) Query(query string, args ...interface{}) (Rows, error) {
	if rows, err := c.query(query, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// query
func (c *customTx) query(query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.tx.Queryx(query, args...)
}

// QueryContext
func (c *customTx) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if rows, err := c.queryContext(ctx, query, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// queryContext
func (c *customTx) queryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if c.testErr != nil {
		return &sqlx.Rows{}, c.popTestError()
	}
	return c.tx.QueryxContext(ctx, query, args...)
}

// Rebind
func (c *customTx) Rebind(query string) string {
	return c.tx.Rebind(query)
}

// Select
func (c *customTx) Select(dest interface{}, query string, args ...interface{}) error {
	return c.tx.Select(dest, query, args...)
}

// SelectContext
func (c *customTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.tx.SelectContext(ctx, dest, query, args...)
}

// Stmt
func (c *customTx) Stmt(stmt interface{}) Stmt {
	return &customStmt{stmt: c.tx.Stmtx(stmt)}
}

// StmtContext
func (c *customTx) StmtContext(ctx context.Context, stmt interface{}) Stmt {
	return &customStmt{stmt: c.tx.StmtxContext(ctx, stmt)}
}

// Unsafe
func (c *customTx) Unsafe() *sqlx.Tx {
	return c.tx.Unsafe()
}

// Safe
func (c *customTx) Safe() *sqlx.Tx {
	return c.tx
}
