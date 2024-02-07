package godb

import (
	"context"
	"database/sql"
)

// ConnMock defines a new connection mock
type ConnMock struct {
	Error                   error
	CallbackClose           func() error
	CallbackExecContext     func(ctx context.Context, query string, args ...any) (sql.Result, error)
	CallbackPingContext     func(ctx context.Context) error
	CallbackRaw             func(f func(driverConn any) error) (err error)
	CallbackBeginTx         func(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	CallbackGetContext      func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	CallbackPrepareContext  func(ctx context.Context, query string) (Stmt, error)
	CallbackQueryRowContext func(ctx context.Context, query string, args ...interface{}) Row
	CallbackQueryContext    func(ctx context.Context, query string, args ...interface{}) (Rows, error)
	CallbackRebind          func(query string) string
	CallbackSelectContext   func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Close closes the connection mock.
// If CallbackClose is set, it will be called and its result will be returned.
// Otherwise, it returns the Error field of the ConnMock struct.
func (cm *ConnMock) Close() error {
	if cm.CallbackClose != nil {
		return cm.CallbackClose()
	}
	return cm.Error
}

// ExecContext executes a query on the connection mock.
// If CallbackExecContext is set, it will be called with the provided arguments and its result will be returned.
// Otherwise, it returns a ResultMock{} and the Error field of the ConnMock struct.
func (cm *ConnMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if cm.CallbackExecContext != nil {
		return cm.CallbackExecContext(ctx, query, args...)
	}
	return &ResultMock{}, cm.Error
}

// PingContext pings the connection mock.
// If CallbackPingContext is set, it will be called with the provided context and its result will be returned.
// Otherwise, it returns the Error field of the ConnMock struct.
func (cm *ConnMock) PingContext(ctx context.Context) error {
	if cm.CallbackPingContext != nil {
		return cm.CallbackPingContext(ctx)
	}
	return cm.Error
}

// Raw executes a raw function on the connection mock.
// If CallbackRaw is set, it will be called with the provided function and its result will be returned.
// Otherwise, it returns the Error field of the ConnMock struct.
func (cm *ConnMock) Raw(f func(driverConn any) error) (err error) {
	if cm.CallbackRaw != nil {
		return cm.CallbackRaw(f)
	}
	return cm.Error
}

// BeginTx begins a transaction on the connection mock.
// If CallbackBeginTx is set, it will be called with the provided context and options, and its result will be returned.
// Otherwise, it returns a TxMock{} and the Error field of the ConnMock struct.
func (cm *ConnMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	if cm.CallbackBeginTx != nil {
		return cm.CallbackBeginTx(ctx, opts)
	}
	return &TxMock{}, cm.Error
}

// GetContext executes a query and retrieves the result on the connection mock.
// If CallbackGetContext is set, it will be called with the provided context, destination, query, and arguments,
// and its result will be returned.
// Otherwise, it returns the Error field of the ConnMock struct.
func (cm *ConnMock) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if cm.CallbackGetContext != nil {
		return cm.CallbackGetContext(ctx, dest, query, args...)
	}
	return cm.Error
}

// PrepareContext prepares a statement on the connection mock.
// If CallbackPrepareContext is set, it will be called with the provided context and query,
// and its result will be returned.
// Otherwise, it returns a StatementMock{} and the Error field of the ConnMock struct.
func (cm *ConnMock) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if cm.CallbackPrepareContext != nil {
		return cm.CallbackPrepareContext(ctx, query)
	}
	return &StmtMock{}, cm.Error
}

// QueryRowContext executes a query and returns a single row on the connection mock.
// If CallbackQueryRowContext is set, it will be called with the provided context, query, and arguments,
// and its result will be returned.
// Otherwise, it returns a RowMock{}.
func (cm *ConnMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	if cm.CallbackQueryRowContext != nil {
		return cm.CallbackQueryRowContext(ctx, query, args...)
	}
	return &RowMock{}
}

// QueryContext executes a query and returns the result on the connection mock.
// If CallbackQueryContext is set, it will be called with the provided context, query, and arguments,
// and its result will be returned.
// Otherwise, it returns a RowsMock{} and the Error field of the ConnMock struct.
func (cm *ConnMock) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if cm.CallbackQueryContext != nil {
		return cm.CallbackQueryContext(ctx, query, args...)
	}
	return &RowsMock{}, cm.Error
}

// Rebind rebinds a query on the connection mock.
// If CallbackRebind is set, it will be called with the provided query,
// and its result will be returned.
// Otherwise, it returns an empty string.
func (cm *ConnMock) Rebind(query string) string {
	if cm.CallbackRebind != nil {
		return cm.CallbackRebind(query)
	}
	return ""
}

// SelectContext executes a query and maps the result to a destination on the connection mock.
// If CallbackSelectContext is set, it will be called with the provided context, destination, query, and arguments,
// and its result will be returned.
// Otherwise, it returns the Error field of the ConnMock struct.
func (cm *ConnMock) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if cm.CallbackSelectContext != nil {
		return cm.CallbackSelectContext(ctx, dest, query, args...)
	}
	return cm.Error
}
