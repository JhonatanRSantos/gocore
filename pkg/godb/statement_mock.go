package godb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// StmtMock is a mock implementation of the sqlx.Stmt interface
type StmtMock struct {
	Error                   error
	CallbackClose           func() error
	CallbackExec            func(args ...any) (sql.Result, error)
	CallbackExecContext     func(ctx context.Context, args ...any) (sql.Result, error)
	CallbackGet             func(dest interface{}, args ...interface{}) error
	CallbackGetContext      func(ctx context.Context, dest interface{}, args ...interface{}) error
	CallbackMustExec        func(args ...interface{}) sql.Result
	CallbackMustExecContext func(ctx context.Context, args ...interface{}) sql.Result
	CallbackQueryRow        func(args ...interface{}) Row
	CallbackQueryRowContext func(ctx context.Context, args ...interface{}) Row
	CallbackQuery           func(args ...interface{}) (Rows, error)
	CallbackQueryContext    func(ctx context.Context, args ...interface{}) (Rows, error)
	CallbackSelect          func(dest interface{}, args ...interface{}) error
	CallbackSelectContext   func(ctx context.Context, dest interface{}, args ...interface{}) error
	CallbackUnsafe          func() *sqlx.Stmt
	CallbackSafe            func() *sqlx.Stmt
}

// Close is a mock implementation of the Close method of sqlx.Stmt interface.
// It calls the CallbackClose function if it is not nil, otherwise it returns the Error field.
func (sm *StmtMock) Close() error {
	if sm.CallbackClose != nil {
		return sm.CallbackClose()
	}
	return sm.Error
}

// Exec is a mock implementation of the Exec method of sqlx.Stmt interface.
// It calls the CallbackExec function if it is not nil, otherwise it returns a ResultMock and the Error field.
func (sm *StmtMock) Exec(args ...any) (sql.Result, error) {
	if sm.CallbackExec != nil {
		return sm.CallbackExec(args...)
	}
	return &ResultMock{}, sm.Error
}

// ExecContext is a mock implementation of the ExecContext method of sqlx.Stmt interface.
// It calls the CallbackExecContext function if it is not nil, otherwise it returns a ResultMock and the Error field.
func (sm *StmtMock) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	if sm.CallbackExecContext != nil {
		return sm.CallbackExecContext(ctx, args...)
	}
	return &ResultMock{}, sm.Error
}

// Get is a mock implementation of the Get method of sqlx.Stmt interface.
// It calls the CallbackGet function if it is not nil, otherwise it returns the Error field.
func (sm *StmtMock) Get(dest interface{}, args ...interface{}) error {
	if sm.CallbackGet != nil {
		return sm.CallbackGet(dest, args...)
	}
	return sm.Error
}

// GetContext is a mock implementation of the GetContext method of sqlx.Stmt interface.
// It calls the CallbackGetContext function if it is not nil, otherwise it returns the Error field.
func (sm *StmtMock) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	if sm.CallbackGetContext != nil {
		return sm.CallbackGetContext(ctx, dest, args...)
	}
	return sm.Error
}

// MustExec is a mock implementation of the MustExec method of sqlx.Stmt interface.
// It calls the CallbackMustExec function if it is not nil, otherwise it returns a ResultMock.
func (sm *StmtMock) MustExec(args ...interface{}) sql.Result {
	if sm.CallbackMustExec != nil {
		return sm.CallbackMustExec(args...)
	}
	return &ResultMock{}
}

// MustExecContext is a mock implementation of the MustExecContext method of sqlx.Stmt interface.
// It calls the CallbackMustExecContext function if it is not nil, otherwise it returns a ResultMock.
func (sm *StmtMock) MustExecContext(ctx context.Context, args ...interface{}) sql.Result {
	if sm.CallbackMustExecContext != nil {
		return sm.CallbackMustExecContext(ctx, args...)
	}
	return &ResultMock{}
}

// QueryRow is a mock implementation of the QueryRow method of sqlx.Stmt interface.
// It calls the CallbackQueryRow function if it is not nil, otherwise it returns a RowMock.
func (sm *StmtMock) QueryRow(args ...interface{}) Row {
	if sm.CallbackQueryRow != nil {
		return sm.CallbackQueryRow(args...)
	}
	return &RowMock{}
}

// QueryRowContext is a mock implementation of the QueryRowContext method of sqlx.Stmt interface.
// It calls the CallbackQueryRowContext function if it is not nil, otherwise it returns a RowMock.
func (sm *StmtMock) QueryRowContext(ctx context.Context, args ...interface{}) Row {
	if sm.CallbackQueryRowContext != nil {
		return sm.CallbackQueryRowContext(ctx, args...)
	}
	return &RowMock{}
}

// Query is a mock implementation of the Query method of sqlx.Stmt interface.
// It calls the CallbackQuery function if it is not nil, otherwise it returns a RowsMock and the Error field.
func (sm *StmtMock) Query(args ...interface{}) (Rows, error) {
	if sm.CallbackQuery != nil {
		return sm.CallbackQuery(args...)
	}
	return &RowsMock{}, sm.Error
}

// QueryContext is a mock implementation of the QueryContext method of sqlx.Stmt interface.
// It calls the CallbackQueryContext function if it is not nil, otherwise it returns a RowsMock and the Error field.
func (sm *StmtMock) QueryContext(ctx context.Context, args ...interface{}) (Rows, error) {
	if sm.CallbackQueryContext != nil {
		return sm.CallbackQueryContext(ctx, args...)
	}
	return &RowsMock{}, sm.Error
}

// Select is a mock implementation of the Select method of sqlx.Stmt interface.
// It calls the CallbackSelect function if it is not nil, otherwise it returns the Error field.
func (sm *StmtMock) Select(dest interface{}, args ...interface{}) error {
	if sm.CallbackSelect != nil {
		return sm.CallbackSelect(dest, args...)
	}
	return sm.Error
}

// SelectContext is a mock implementation of the SelectContext method of sqlx.Stmt interface.
// It calls the CallbackSelectContext function if it is not nil, otherwise it returns the Error field.
func (sm *StmtMock) SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	if sm.CallbackSelectContext != nil {
		return sm.CallbackSelectContext(ctx, dest, args...)
	}
	return sm.Error
}

// Unsafe is a mock implementation of the Unsafe method of sqlx.Stmt interface.
// It calls the CallbackUnsafe function if it is not nil, otherwise it returns a sqlx.Stmt.
func (sm *StmtMock) Unsafe() *sqlx.Stmt {
	if sm.CallbackUnsafe != nil {
		return sm.CallbackUnsafe()
	}
	return &sqlx.Stmt{}
}

// Safe is a mock implementation of the Safe method of sqlx.Stmt interface.
// It calls the CallbackSafe function if it is not nil, otherwise it returns a sqlx.Stmt.
func (sm *StmtMock) Safe() *sqlx.Stmt {
	if sm.CallbackSafe != nil {
		return sm.CallbackSafe()
	}
	return &sqlx.Stmt{}
}
