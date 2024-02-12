package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// NamedStmtMock is a mock implementation of sqlx.NamedStmt.
// It provides callback functions for each method to allow custom behavior during testing.
type NamedStmtMock struct {
	Error                   error
	CallbackClose           func() error
	CallbackExec            func(arg interface{}) (sql.Result, error)
	CallbackExecContext     func(ctx context.Context, arg interface{}) (sql.Result, error)
	CallbackGet             func(dest interface{}, arg interface{}) error
	CallbackGetContext      func(ctx context.Context, dest interface{}, arg interface{}) error
	CallbackMustExec        func(arg interface{}) sql.Result
	CallbackMustExecContext func(ctx context.Context, arg interface{}) sql.Result
	CallbackQueryRow        func(arg interface{}) Row
	CallbackQueryRowContext func(ctx context.Context, arg interface{}) Row
	CallbackQuery           func(arg interface{}) (Rows, error)
	CallbackQueryContext    func(ctx context.Context, arg interface{}) (Rows, error)
	CallbackSelect          func(dest interface{}, arg interface{}) error
	CallbackSelectContext   func(ctx context.Context, dest interface{}, arg interface{}) error
	CallbackUnsafe          func() *sqlx.NamedStmt
	CallbackSafe            func() *sqlx.NamedStmt
}

// Close calls the callback function CallbackClose if it is set.
// Otherwise, it returns the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) Close() error {
	if nsm.CallbackClose != nil {
		return nsm.CallbackClose()
	}
	return nsm.Error
}

// Exec calls the callback function CallbackExec if it is set.
// Otherwise, it returns a ResultMock and the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) Exec(arg interface{}) (sql.Result, error) {
	if nsm.CallbackExec != nil {
		return nsm.CallbackExec(arg)
	}
	return &ResultMock{}, nsm.Error
}

// ExecContext calls the callback function CallbackExecContext if it is set.
// Otherwise, it returns a ResultMock and the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) ExecContext(ctx context.Context, arg interface{}) (sql.Result, error) {
	if nsm.CallbackExecContext != nil {
		return nsm.CallbackExecContext(ctx, arg)
	}
	return &ResultMock{}, nsm.Error
}

// Get calls the callback function CallbackGet if it is set.
// Otherwise, it returns the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) Get(dest interface{}, arg interface{}) error {
	if nsm.CallbackGet != nil {
		return nsm.CallbackGet(dest, arg)
	}
	return nsm.Error
}

// GetContext calls the callback function CallbackGetContext if it is set.
// Otherwise, it returns the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) GetContext(ctx context.Context, dest interface{}, arg interface{}) error {
	if nsm.CallbackGetContext != nil {
		return nsm.CallbackGetContext(ctx, dest, arg)
	}
	return nsm.Error
}

// MustExec calls the callback function CallbackMustExec if it is set.
// Otherwise, it returns a ResultMock.
func (nsm *NamedStmtMock) MustExec(arg interface{}) sql.Result {
	if nsm.CallbackMustExec != nil {
		return nsm.CallbackMustExec(arg)
	}
	return &ResultMock{}
}

// MustExecContext calls the callback function CallbackMustExecContext if it is set.
// Otherwise, it returns a ResultMock.
func (nsm *NamedStmtMock) MustExecContext(ctx context.Context, arg interface{}) sql.Result {
	if nsm.CallbackMustExecContext != nil {
		return nsm.CallbackMustExecContext(ctx, arg)
	}
	return &ResultMock{}
}

// QueryRow calls the callback function CallbackQueryRow if it is set.
// Otherwise, it returns a RowMock.
func (nsm *NamedStmtMock) QueryRow(arg interface{}) Row {
	if nsm.CallbackQueryRow != nil {
		return nsm.CallbackQueryRow(arg)
	}
	return &RowMock{}
}

// QueryRowContext calls the callback function CallbackQueryRowContext if it is set.
// Otherwise, it returns a RowMock.
func (nsm *NamedStmtMock) QueryRowContext(ctx context.Context, arg interface{}) Row {
	if nsm.CallbackQueryRowContext != nil {
		return nsm.CallbackQueryRowContext(ctx, arg)
	}
	return &RowMock{}
}

// Query calls the callback function CallbackQuery if it is set.
// Otherwise, it returns a RowsMock and the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) Query(arg interface{}) (Rows, error) {
	if nsm.CallbackQuery != nil {
		return nsm.CallbackQuery(arg)
	}
	return &RowsMock{}, nsm.Error
}

// QueryContext calls the callback function CallbackQueryContext if it is set.
// Otherwise, it returns a RowsMock and the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) QueryContext(ctx context.Context, arg interface{}) (Rows, error) {
	if nsm.CallbackQueryContext != nil {
		return nsm.CallbackQueryContext(ctx, arg)
	}
	return &RowsMock{}, nsm.Error
}

// Select calls the callback function CallbackSelect if it is set.
// Otherwise, it returns the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) Select(dest interface{}, arg interface{}) error {
	if nsm.CallbackSelect != nil {
		return nsm.CallbackSelect(dest, arg)
	}
	return nsm.Error
}

// SelectContext calls the callback function CallbackSelectContext if it is set.
// Otherwise, it returns the Error field of the NamedStmtMock.
func (nsm *NamedStmtMock) SelectContext(ctx context.Context, dest interface{}, arg interface{}) error {
	if nsm.CallbackSelectContext != nil {
		return nsm.CallbackSelectContext(ctx, dest, arg)
	}
	return nsm.Error
}

// Unsafe calls the callback function CallbackUnsafe if it is set.
// Otherwise, it returns a new sqlx.NamedStmt.
func (nsm *NamedStmtMock) Unsafe() *sqlx.NamedStmt {
	if nsm.CallbackUnsafe != nil {
		return nsm.CallbackUnsafe()
	}
	return &sqlx.NamedStmt{}
}

// Safe calls the callback function CallbackSafe if it is set.
// Otherwise, it returns a new sqlx.NamedStmt.
func (nsm *NamedStmtMock) Safe() *sqlx.NamedStmt {
	if nsm.CallbackSafe != nil {
		return nsm.CallbackSafe()
	}
	return &sqlx.NamedStmt{}
}
