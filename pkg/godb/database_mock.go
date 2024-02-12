package godb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/jmoiron/sqlx"
)

// DBMock is a mock implementation of the sql.DB interface.
type DBMock struct {
	Error                       error
	CallbackClose               func() error
	CallbackDriver              func() driver.Driver
	CallbackExec                func(query string, args ...any) (sql.Result, error)
	CallbackExecContext         func(ctx context.Context, query string, args ...any) (sql.Result, error)
	CallbackPing                func() error
	CallbackPingContext         func(ctx context.Context) error
	CallbackSetConnMaxIdleTime  func(d time.Duration)
	CallbackSetConnMaxLifetime  func(d time.Duration)
	CallbackSetMaxOpenConns     func(n int)
	CallbackSetMaxIdleConns     func(n int)
	CallbackStats               func() sql.DBStats
	CallbackBeginTx             func(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	CallbackBegin               func() (Tx, error)
	CallbackBindNamed           func(query string, arg interface{}) (string, []interface{}, error)
	CallbackConn                func(ctx context.Context) (Conn, error)
	CallbackDriverName          func() string
	CallbackGet                 func(dest interface{}, query string, args ...interface{}) error
	CallbackGetContext          func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	CallbackMapperFunc          func(mf func(string) string)
	CallbackMustBegin           func() Tx
	CallbackMustBeginTx         func(ctx context.Context, opts *sql.TxOptions) Tx
	CallbackMustExec            func(query string, args ...interface{}) sql.Result
	CallbackMustExecContext     func(ctx context.Context, query string, args ...interface{}) sql.Result
	CallbackNamedExec           func(query string, arg interface{}) (sql.Result, error)
	CallbackNamedExecContext    func(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	CallbackNamedQuery          func(query string, arg interface{}) (Rows, error)
	CallbackNamedQueryContext   func(ctx context.Context, query string, arg interface{}) (Rows, error)
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
	CallbackUnsafe              func() *sqlx.DB
	CallbackSafe                func() *sqlx.DB
}

// NewMockDB creates a new instance of DBMock.
func NewMockDB() *DBMock {
	return &DBMock{}
}

// Close calls the callback function CallbackClose if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) Close() error {
	if dbm.CallbackClose != nil {
		return dbm.CallbackClose()
	}
	return dbm.Error
}

// Driver calls the callback function CallbackDriver if it is set.
// Otherwise, it returns a new instance of DriverMock.
func (dbm *DBMock) Driver() driver.Driver {
	if dbm.CallbackDriver != nil {
		return dbm.CallbackDriver()
	}
	return &DriverMock{Error: dbm.Error}
}

// Exec calls the callback function CallbackExec if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) Exec(query string, args ...any) (sql.Result, error) {
	if dbm.CallbackExec != nil {
		return dbm.CallbackExec(query, args...)
	}
	return &ResultMock{}, dbm.Error
}

// ExecContext calls the callback function CallbackExecContext if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if dbm.CallbackExecContext != nil {
		return dbm.CallbackExecContext(ctx, query, args...)
	}
	return &ResultMock{}, dbm.Error
}

// Ping calls the callback function CallbackPing if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) Ping() error {
	if dbm.CallbackPing != nil {
		return dbm.CallbackPing()
	}
	return dbm.Error
}

// PingContext calls the callback function CallbackPingContext if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) PingContext(ctx context.Context) error {
	if dbm.CallbackPingContext != nil {
		return dbm.CallbackPingContext(ctx)
	}
	return dbm.Error
}

// SetConnMaxIdleTime calls the callback function CallbackSetConnMaxIdleTime if it is set.
func (dbm *DBMock) SetConnMaxIdleTime(d time.Duration) {
	if dbm.CallbackSetConnMaxIdleTime != nil {
		dbm.CallbackSetConnMaxIdleTime(d)
	}
}

// SetConnMaxLifetime calls the callback function CallbackSetConnMaxLifetime if it is set.
func (dbm *DBMock) SetConnMaxLifetime(d time.Duration) {
	if dbm.CallbackSetConnMaxLifetime != nil {
		dbm.CallbackSetConnMaxLifetime(d)
	}
}

// SetMaxIdleConns calls the callback function CallbackSetMaxIdleConns if it is set.
func (dbm *DBMock) SetMaxIdleConns(n int) {
	if dbm.CallbackSetMaxIdleConns != nil {
		dbm.CallbackSetMaxIdleConns(n)
	}
}

// SetMaxOpenConns calls the callback function CallbackSetMaxOpenConns if it is set.
func (dbm *DBMock) SetMaxOpenConns(n int) {
	if dbm.CallbackSetMaxOpenConns != nil {
		dbm.CallbackSetMaxOpenConns(n)
	}
}

// Stats calls the callback function CallbackStats if it is set.
func (dbm *DBMock) Stats() sql.DBStats {
	if dbm.CallbackStats != nil {
		return dbm.CallbackStats()
	}
	return sql.DBStats{}
}

// BeginTx calls the callback function CallbackBeginTx if it is set.
// Otherwise, it returns a new instance of TxMock.
func (dbm *DBMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	if dbm.CallbackBeginTx != nil {
		return dbm.CallbackBeginTx(ctx, opts)
	}
	return &TxMock{}, dbm.Error
}

// Begin calls the callback function CallbackBegin if it is set.
// Otherwise, it returns a new instance of TxMock.
func (dbm *DBMock) Begin() (Tx, error) {
	if dbm.CallbackBegin != nil {
		return dbm.CallbackBegin()
	}
	return &TxMock{}, dbm.Error
}

// BindNamed calls the callback function CallbackBindNamed if it is set.
// Otherwise, it returns an empty string and an empty slice of interfaces.
func (dbm *DBMock) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	if dbm.CallbackBindNamed != nil {
		return dbm.CallbackBindNamed(query, arg)
	}
	return "", []interface{}{}, dbm.Error
}

// Conn calls the callback function CallbackConn if it is set.
// Otherwise, it returns a new instance of ConnMock.
func (dbm *DBMock) Conn(ctx context.Context) (Conn, error) {
	if dbm.CallbackConn != nil {
		return dbm.CallbackConn(ctx)
	}
	return &ConnMock{}, dbm.Error
}

// DriverName calls the callback function CallbackDriverName if it is set.
func (dbm *DBMock) DriverName() string {
	if dbm.CallbackDriverName != nil {
		return dbm.CallbackDriverName()
	}
	return ""
}

// Get calls the callback function CallbackGet if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) Get(dest interface{}, query string, args ...interface{}) error {
	if dbm.CallbackGet != nil {
		return dbm.CallbackGet(dest, query, args...)
	}
	return dbm.Error
}

// GetContext calls the callback function CallbackGetContext if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if dbm.CallbackGetContext != nil {
		return dbm.CallbackGetContext(ctx, dest, query, args...)
	}
	return dbm.Error
}

// MapperFunc calls the callback function CallbackMapperFunc if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) MapperFunc(mf func(string) string) {
	if dbm.CallbackMapperFunc != nil {
		dbm.CallbackMapperFunc(mf)
	}
}

// MustBegin calls the callback function CallbackMustBegin if it is set.
// Otherwise, it returns a new instance of TxMock.
func (dbm *DBMock) MustBegin() Tx {
	if dbm.CallbackMustBegin != nil {
		return dbm.CallbackMustBegin()
	}
	return &TxMock{}
}

// MustBeginTx calls the callback function CallbackMustBeginTx if it is set.
// Otherwise, it returns a new instance of TxMock.
func (dbm *DBMock) MustBeginTx(ctx context.Context, opts *sql.TxOptions) Tx {
	if dbm.CallbackMustBeginTx != nil {
		return dbm.CallbackMustBeginTx(ctx, opts)
	}
	return &TxMock{}
}

// MustExec calls the callback function CallbackMustExec if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) MustExec(query string, args ...interface{}) sql.Result {
	if dbm.CallbackMustExec != nil {
		return dbm.CallbackMustExec(query, args...)
	}
	return &ResultMock{}
}

// MustExecContext calls the callback function CallbackMustExecContext if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
	if dbm.CallbackMustExecContext != nil {
		return dbm.CallbackMustExecContext(ctx, query, args...)
	}
	return &ResultMock{}
}

// NamedExec calls the callback function CallbackNamedExec if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	if dbm.CallbackNamedExec != nil {
		return dbm.CallbackNamedExec(query, arg)
	}
	return &ResultMock{}, dbm.Error
}

// NamedExecContext calls the callback function CallbackNamedExecContext if it is set.
// Otherwise, it returns a new instance of ResultMock.
func (dbm *DBMock) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if dbm.CallbackNamedExecContext != nil {
		return dbm.CallbackNamedExecContext(ctx, query, arg)
	}
	return &ResultMock{}, dbm.Error
}

// NamedQuery calls the callback function CallbackNamedQuery if it is set.
// Otherwise, it returns a new instance of RowsMock.
func (dbm *DBMock) NamedQuery(query string, arg interface{}) (Rows, error) {
	if dbm.CallbackNamedQuery != nil {
		return dbm.CallbackNamedQuery(query, arg)
	}
	return &RowsMock{}, dbm.Error
}

// NamedQueryContext calls the callback function CallbackNamedQueryContext if it is set.
// Otherwise, it returns a new instance of RowsMock.
func (dbm *DBMock) NamedQueryContext(ctx context.Context, query string, arg interface{}) (Rows, error) {
	if dbm.CallbackNamedQueryContext != nil {
		return dbm.CallbackNamedQueryContext(ctx, query, arg)
	}
	return &RowsMock{}, dbm.Error
}

// PrepareNamed calls the callback function CallbackPrepareNamed if it is set.
// Otherwise, it returns a new instance of NamedStmtMock.
func (dbm *DBMock) PrepareNamed(query string) (NamedStmt, error) {
	if dbm.CallbackPrepareNamed != nil {
		return dbm.CallbackPrepareNamed(query)
	}
	return &NamedStmtMock{}, dbm.Error
}

// PrepareNamedContext calls the callback function CallbackPrepareNamedContext if it is set.
// Otherwise, it returns a new instance of NamedStmtMock.
func (dbm *DBMock) PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error) {
	if dbm.CallbackPrepareNamedContext != nil {
		return dbm.CallbackPrepareNamedContext(ctx, query)
	}
	return &NamedStmtMock{}, dbm.Error
}

// Prepare calls the callback function CallbackPrepare if it is set.
// Otherwise, it returns a new instance of StatementMock.
func (dbm *DBMock) Prepare(query string) (Stmt, error) {
	if dbm.CallbackPrepare != nil {
		return dbm.CallbackPrepare(query)
	}
	return &StmtMock{}, dbm.Error
}

// PrepareContext calls the callback function CallbackPrepareContext if it is set.
// Otherwise, it returns a new instance of StatementMock.
func (dbm *DBMock) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if dbm.CallbackPrepareContext != nil {
		return dbm.CallbackPrepareContext(ctx, query)
	}
	return &StmtMock{}, dbm.Error
}

// QueryRow calls the callback function CallbackQueryRow if it is set.
// Otherwise, it returns a new instance of RowMock.
func (dbm *DBMock) QueryRow(query string, args ...interface{}) Row {
	if dbm.CallbackQueryRow != nil {
		return dbm.CallbackQueryRow(query, args...)
	}
	return &RowMock{}
}

// QueryRowContext calls the callback function CallbackQueryRowContext if it is set.
// Otherwise, it returns a new instance of RowMock.
func (dbm *DBMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	if dbm.CallbackQueryRowContext != nil {
		return dbm.CallbackQueryRowContext(ctx, query, args...)
	}
	return &RowMock{}
}

// Query calls the callback function CallbackQuery if it is set.
// Otherwise, it returns a new instance of RowsMock.
func (dbm *DBMock) Query(query string, args ...interface{}) (Rows, error) {
	if dbm.CallbackQuery != nil {
		return dbm.CallbackQuery(query, args...)
	}
	return &RowsMock{}, dbm.Error
}

// QueryContext calls the callback function CallbackQueryContext if it is set.
// Otherwise, it returns a new instance of RowsMock.
func (dbm *DBMock) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if dbm.CallbackQueryContext != nil {
		return dbm.CallbackQueryContext(ctx, query, args...)
	}
	return &RowsMock{}, dbm.Error
}

// Rebind calls the callback function CallbackRebind if it is set.
// Otherwise, it returns an empty string.
func (dbm *DBMock) Rebind(query string) string {
	if dbm.CallbackRebind != nil {
		return dbm.CallbackRebind(query)
	}
	return ""
}

// Select calls the callback function CallbackSelect if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) Select(dest interface{}, query string, args ...interface{}) error {
	if dbm.CallbackSelect != nil {
		return dbm.CallbackSelect(dest, query, args...)
	}
	return dbm.Error
}

// SelectContext calls the callback function CallbackSelectContext if it is set.
// Otherwise, it returns the Error field of the DBMock.
func (dbm *DBMock) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if dbm.CallbackSelectContext != nil {
		return dbm.CallbackSelectContext(ctx, dest, query, args...)
	}
	return dbm.Error
}

// Unsafe calls the callback function CallbackUnsafe if it is set.
// Otherwise, it returns a new instance of sqlx.DB.
func (dbm *DBMock) Unsafe() *sqlx.DB {
	if dbm.CallbackUnsafe != nil {
		return dbm.CallbackUnsafe()
	}
	return &sqlx.DB{}
}

// Safe calls the callback function CallbackSafe if it is set.
// Otherwise, it returns a new instance of sqlx.DB.
func (dbm *DBMock) Safe() *sqlx.DB {
	if dbm.CallbackSafe != nil {
		return dbm.CallbackSafe()
	}
	return &sqlx.DB{}
}
