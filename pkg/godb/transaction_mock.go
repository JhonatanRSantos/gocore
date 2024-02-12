package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// TxMock is a mock implementation of the sqlx.Tx interface
type TxMock struct {
	Error                       error
	CallbackCommit              func() error
	CallbackExec                func(query string, args ...any) (sql.Result, error)
	CallbackExecContext         func(ctx context.Context, query string, args ...any) (sql.Result, error)
	CallbackRollback            func() error
	CallbackBindNamed           func(query string, arg interface{}) (string, []interface{}, error)
	CallbackDriverName          func() string
	CallbackGet                 func(dest interface{}, query string, args ...interface{}) error
	CallbackGetContext          func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	CallbackMustExec            func(query string, args ...interface{}) sql.Result
	CallbackMustExecContext     func(ctx context.Context, query string, args ...interface{}) sql.Result
	CallbackNamedExec           func(query string, arg interface{}) (sql.Result, error)
	CallbackNamedExecContext    func(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	CallbackNamedQuery          func(query string, arg interface{}) (Rows, error)
	CallbackNamedStmt           func(stmt NamedStmt) NamedStmt
	CallbackNamedStmtContext    func(ctx context.Context, stmt NamedStmt) NamedStmt
	CallbackPrepareNamed        func(query string) (NamedStmt, error)
	CallbackPrepareNamedContext func(ctx context.Context, query string) (NamedStmt, error)
	CallbackPrepare             func(query string) (Stmt, error)
	CallbackPrepareContext      func(ctx context.Context, query string) (Stmt, error)
	CallbackQueryRow            func(query string, args ...interface{}) Row
	CallbackQueryRowContext     func(ctx context.Context, query string, args ...interface{}) Row
	CallbackQuery               func(query string, args ...interface{}) (Rows, error)
	CallbackQueryContext        func(ctx context.Context, query string, args ...interface{}) (Rows, error)
	CallbackRebind              func(query string) string
	CallbackSelect              func(dest interface{}, query string, args ...interface{}) error
	CallbackSelectContext       func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	CallbackStmt                func(stmt interface{}) Stmt
	CallbackStmtContext         func(ctx context.Context, stmt interface{}) Stmt
	CallbackUnsafe              func() *sqlx.Tx
	CallbackSafe                func() *sqlx.Tx
}

// Commit commits the transaction.
// If a callback function is set for Commit, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) Commit() error {
	if tm.CallbackCommit != nil {
		return tm.CallbackCommit()
	}
	return tm.Error
}

// Exec executes a query without returning any rows.
// If a callback function is set for Exec, it will be called instead of the default implementation.
// Returns a sql.Result and an error if the callback function returns a result and an error or if an error is set in the TxMock struct.
func (tm *TxMock) Exec(query string, args ...any) (sql.Result, error) {
	if tm.CallbackExec != nil {
		return tm.CallbackExec(query, args...)
	}
	return &ResultMock{}, tm.Error
}

// ExecContext executes a query without returning any rows with a context.
// If a callback function is set for ExecContext, it will be called instead of the default implementation.
// Returns a sql.Result and an error if the callback function returns a result and an error or if an error is set in the TxMock struct.
func (tm *TxMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tm.CallbackExecContext != nil {
		return tm.CallbackExecContext(ctx, query, args...)
	}
	return &ResultMock{}, tm.Error
}

// Rollback rolls back the transaction.
// If a callback function is set for Rollback, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) Rollback() error {
	if tm.CallbackRollback != nil {
		return tm.CallbackRollback()
	}
	return tm.Error
}

// BindNamed binds named parameters in a query.
// If a callback function is set for BindNamed, it will be called instead of the default implementation.
// Returns the modified query string, a slice of interface{} containing the bound values, and an error if the callback function returns these values or if an error is set in the TxMock struct.
func (tm *TxMock) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	if tm.CallbackBindNamed != nil {
		return tm.CallbackBindNamed(query, arg)
	}
	return "", []interface{}{}, tm.Error
}

// DriverName returns the name of the driver.
// If a callback function is set for DriverName, it will be called instead of the default implementation.
// Returns the driver name or an empty string if the callback function returns a driver name or if no callback function is set.
func (tm *TxMock) DriverName() string {
	if tm.CallbackDriverName != nil {
		return tm.CallbackDriverName()
	}
	return ""
}

// Get retrieves a single row from the database and stores the result in the provided struct or pointer.
// If a callback function is set for Get, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) Get(dest interface{}, query string, args ...interface{}) error {
	if tm.CallbackGet != nil {
		return tm.CallbackGet(dest, query, args...)
	}
	return tm.Error
}

// GetContext retrieves a single row from the database with a context and stores the result in the provided struct or pointer.
// If a callback function is set for GetContext, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if tm.CallbackGetContext != nil {
		return tm.CallbackGetContext(ctx, dest, query, args...)
	}
	return tm.Error
}

// MustExec executes a query without returning any rows and panics if an error occurs.
// If a callback function is set for MustExec, it will be called instead of the default implementation.
// Returns a sql.Result or a ResultMock if the callback function returns a result or if no callback function is set.
func (tm *TxMock) MustExec(query string, args ...interface{}) sql.Result {
	if tm.CallbackMustExec != nil {
		return tm.CallbackMustExec(query, args...)
	}
	return &ResultMock{}
}

// MustExecContext executes a query without returning any rows with a context and panics if an error occurs.
// If a callback function is set for MustExecContext, it will be called instead of the default implementation.
// Returns a sql.Result or a ResultMock if the callback function returns a result or if no callback function is set.
func (tm *TxMock) MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
	if tm.CallbackMustExecContext != nil {
		return tm.CallbackMustExecContext(ctx, query, args...)
	}
	return &ResultMock{}
}

// NamedExec executes a named query without returning any rows.
// If a callback function is set for NamedExec, it will be called instead of the default implementation.
// Returns a sql.Result and an error if the callback function returns a result and an error or if an error is set in the TxMock struct.
func (tm *TxMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	if tm.CallbackNamedExec != nil {
		return tm.CallbackNamedExec(query, arg)
	}
	return &ResultMock{}, tm.Error
}

// NamedExecContext executes a named query without returning any rows with a context.
// If a callback function is set for NamedExecContext, it will be called instead of the default implementation.
// Returns a sql.Result and an error if the callback function returns a result and an error or if an error is set in the TxMock struct.
func (tm *TxMock) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if tm.CallbackNamedExecContext != nil {
		return tm.CallbackNamedExecContext(ctx, query, arg)
	}
	return &ResultMock{}, tm.Error
}

// NamedQuery executes a named query that returns rows.
// If a callback function is set for NamedQuery, it will be called instead of the default implementation.
// Returns a Rows and an error if the callback function returns rows and an error or if an error is set in the TxMock struct.
func (tm *TxMock) NamedQuery(query string, arg interface{}) (Rows, error) {
	if tm.CallbackNamedQuery != nil {
		return tm.CallbackNamedQuery(query, arg)
	}
	return &RowsMock{}, tm.Error
}

// NamedStmt returns a NamedStmt for the provided NamedStmt.
// If a callback function is set for NamedStmt, it will be called instead of the default implementation.
// Returns a NamedStmt or a NamedStmtMock if the callback function returns a NamedStmt or if no callback function is set.
func (tm *TxMock) NamedStmt(stmt NamedStmt) NamedStmt {
	if tm.CallbackNamedStmt != nil {
		return tm.CallbackNamedStmt(stmt)
	}
	return &NamedStmtMock{}
}

// NamedStmtContext returns a NamedStmt for the provided NamedStmt with a context.
// If a callback function is set for NamedStmtContext, it will be called instead of the default implementation.
// Returns a NamedStmt or a NamedStmtMock if the callback function returns a NamedStmt or if no callback function is set.
func (tm *TxMock) NamedStmtContext(ctx context.Context, stmt NamedStmt) NamedStmt {
	if tm.CallbackNamedStmtContext != nil {
		return tm.CallbackNamedStmtContext(ctx, stmt)
	}
	return &NamedStmtMock{}
}

// PrepareNamed prepares a named statement for execution.
// If a callback function is set for PrepareNamed, it will be called instead of the default implementation.
// Returns a NamedStmt and an error if the callback function returns a NamedStmt and an error or if an error is set in the TxMock struct.
func (tm *TxMock) PrepareNamed(query string) (NamedStmt, error) {
	if tm.CallbackPrepareNamed != nil {
		return tm.CallbackPrepareNamed(query)
	}
	return &NamedStmtMock{}, tm.Error
}

// PrepareNamedContext prepares a named statement for execution with a context.
// If a callback function is set for PrepareNamedContext, it will be called instead of the default implementation.
// Returns a NamedStmt and an error if the callback function returns a NamedStmt and an error or if an error is set in the TxMock struct.
func (tm *TxMock) PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error) {
	if tm.CallbackPrepareNamedContext != nil {
		return tm.CallbackPrepareNamedContext(ctx, query)
	}
	return &NamedStmtMock{}, tm.Error
}

// Prepare prepares a statement for execution.
// If a callback function is set for Prepare, it will be called instead of the default implementation.
// Returns a Stmt and an error if the callback function returns a Stmt and an error or if an error is set in the TxMock struct.
func (tm *TxMock) Prepare(query string) (Stmt, error) {
	if tm.CallbackPrepare != nil {
		return tm.CallbackPrepare(query)
	}
	return &StmtMock{}, tm.Error
}

// PrepareContext prepares a statement for execution with a context.
// If a callback function is set for PrepareContext, it will be called instead of the default implementation.
// Returns a Stmt and an error if the callback function returns a Stmt and an error or if an error is set in the TxMock struct.
func (tm *TxMock) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if tm.CallbackPrepareContext != nil {
		return tm.CallbackPrepareContext(ctx, query)
	}
	return &StmtMock{}, tm.Error
}

// QueryRow executes a query that is expected to return at most one row.
// If a callback function is set for QueryRow, it will be called instead of the default implementation.
// Returns a Row or a RowMock if the callback function returns a Row or if no callback function is set.
func (tm *TxMock) QueryRow(query string, args ...interface{}) Row {
	if tm.CallbackQueryRow != nil {
		return tm.CallbackQueryRow(query, args...)
	}
	return &RowMock{}
}

// QueryRowContext executes a query that is expected to return at most one row with a context.
// If a callback function is set for QueryRowContext, it will be called instead of the default implementation.
// Returns a Row or a RowMock if the callback function returns a Row or if no callback function is set.
func (tm *TxMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	if tm.CallbackQueryRowContext != nil {
		return tm.CallbackQueryRowContext(ctx, query, args...)
	}
	return &RowMock{}
}

// Query executes a query that returns rows.
// If a callback function is set for Query, it will be called instead of the default implementation.
// Returns a Rows and an error if the callback function returns rows and an error or if an error is set in the TxMock struct.
func (tm *TxMock) Query(query string, args ...interface{}) (Rows, error) {
	if tm.CallbackQuery != nil {
		return tm.CallbackQuery(query, args...)
	}
	return &RowsMock{}, tm.Error
}

// QueryContext executes a query that returns rows with a context.
// If a callback function is set for QueryContext, it will be called instead of the default implementation.
// Returns a Rows and an error if the callback function returns rows and an error or if an error is set in the TxMock struct.
func (tm *TxMock) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if tm.CallbackQueryContext != nil {
		return tm.CallbackQueryContext(ctx, query, args...)
	}
	return &RowsMock{}, tm.Error
}

// Rebind returns a query with the placeholders replaced with the appropriate bindvar for the database.
// If a callback function is set for Rebind, it will be called instead of the default implementation.
// Returns the modified query string or an empty string if the callback function returns a modified query or if no callback function is set.
func (tm *TxMock) Rebind(query string) string {
	if tm.CallbackRebind != nil {
		return tm.CallbackRebind(query)
	}
	return ""
}

// Select executes a query that returns multiple rows and stores the result in the provided slice or struct.
// If a callback function is set for Select, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) Select(dest interface{}, query string, args ...interface{}) error {
	if tm.CallbackSelect != nil {
		return tm.CallbackSelect(dest, query, args...)
	}
	return tm.Error
}

// SelectContext executes a query that returns multiple rows with a context and stores the result in the provided slice or struct.
// If a callback function is set for SelectContext, it will be called instead of the default implementation.
// Returns an error if the callback function returns an error or if an error is set in the TxMock struct.
func (tm *TxMock) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if tm.CallbackSelectContext != nil {
		return tm.CallbackSelectContext(ctx, dest, query, args...)
	}
	return tm.Error
}

// Stmt returns a Stmt for the provided Stmt.
// If a callback function is set for Stmt, it will be called instead of the default implementation.
// Returns a Stmt or a StatementMock if the callback function returns a Stmt or if no callback function is set.
func (tm *TxMock) Stmt(stmt interface{}) Stmt {
	if tm.CallbackStmt != nil {
		return tm.CallbackStmt(stmt)
	}
	return &StmtMock{}
}

// StmtContext returns a Stmt for the provided Stmt with a context.
// If a callback function is set for StmtContext, it will be called instead of the default implementation.
// Returns a Stmt or a StatementMock if the callback function returns a Stmt or if no callback function is set.
func (tm *TxMock) StmtContext(ctx context.Context, stmt interface{}) Stmt {
	if tm.CallbackStmtContext != nil {
		return tm.CallbackStmtContext(ctx, stmt)
	}
	return &StmtMock{}
}

// Unsafe returns the underlying *sqlx.Tx.
// If a callback function is set for Unsafe, it will be called instead of the default implementation.
// Returns the underlying *sqlx.Tx or a new *sqlx.Tx if the callback function returns a *sqlx.Tx or if no callback function is set.
func (tm *TxMock) Unsafe() *sqlx.Tx {
	if tm.CallbackUnsafe != nil {
		return tm.CallbackUnsafe()
	}
	return &sqlx.Tx{}
}

// Safe returns the underlying *sqlx.Tx.
// If a callback function is set for Safe, it will be called instead of the default implementation.
// Returns the underlying *sqlx.Tx or a new *sqlx.Tx if the callback function returns a *sqlx.Tx or if no callback function is set.
func (tm *TxMock) Safe() *sqlx.Tx {
	if tm.CallbackSafe != nil {
		return tm.CallbackSafe()
	}
	return &sqlx.Tx{}
}
