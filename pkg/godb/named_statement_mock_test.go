package godb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_NamedStmtMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   NamedStmtMock
		assert func(t *testing.T, mock NamedStmtMock)
	}{
		{
			name: "Should run CallbackClose",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackClose = func() error {
					callCount++
					return nil
				}
				err := mock.Close()
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackClose to be called once")
			},
		},
		{
			name: "Should run default CallbackClose",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				testErr := errors.New("failed to close")
				mock.Error = testErr
				err := mock.Close()
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
			},
		},
		{
			name: "Should run CallbackExec",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackExec = func(arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.Exec(nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackExec to be called once")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
			},
		},
		{
			name: "Should run default CallbackExec",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result, err := mock.Exec(nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackExecContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackExecContext = func(ctx context.Context, arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.ExecContext(context.Background(), nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackExecContext to be called once")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
			},
		},
		{
			name: "Should default CallbackExecContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result, err := mock.ExecContext(context.Background(), nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackGet",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				callCount := 0
				mock.CallbackGet = func(dest interface{}, arg interface{}) error {
					callCount++
					*(dest).(*int) = 1
					return nil
				}
				err := mock.Get(&data, nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackGet to be called once")
				assert.Equal(t, 1, data, "Expected data to be 1")
			},
		},
		{
			name: "Should default CallbackGet",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				testErr := errors.New("failed to get")
				mock.Error = testErr
				err := mock.Get(&data, nil)
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
				assert.Equal(t, 0, data, "Expected data to be 0")
			},
		},
		{
			name: "Should run CallbackGetContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				callCount := 0
				mock.CallbackGetContext = func(ctx context.Context, dest interface{}, arg interface{}) error {
					callCount++
					*(dest).(*int) = 1
					return nil
				}
				err := mock.GetContext(context.Background(), &data, nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackGetContext to be called once")
				assert.Equal(t, 1, data, "Expected data to be 1")
			},
		},
		{
			name: "Should default CallbackGetContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				testErr := errors.New("failed to get")
				mock.Error = testErr
				err := mock.GetContext(context.Background(), &data, nil)
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
				assert.Equal(t, 0, data, "Expected data to be 0")
			},
		},
		{
			name: "Should run CallbackMustExec",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackMustExec = func(arg interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExec(nil)
				assert.Equal(t, 1, callCount, "Expected CallbackMustExec to be called once")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
			},
		},
		{
			name: "Should default CallbackMustExec",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.MustExec(nil)
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackMustExecContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackMustExecContext = func(ctx context.Context, arg interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExecContext(context.Background(), nil)
				assert.Equal(t, 1, callCount, "Expected CallbackMustExecContext to be called once")
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
			},
		},
		{
			name: "Should default CallbackMustExecContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.MustExecContext(context.Background(), nil)
				assert.IsType(t, &ResultMock{}, result, "Expected result to be of type *ResultMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackQueryRow",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackQueryRow = func(arg interface{}) Row {
					callCount++
					return &RowMock{}
				}
				result := mock.QueryRow(nil)
				assert.Equal(t, 1, callCount, "Expected CallbackQueryRow to be called once")
				assert.IsType(t, &RowMock{}, result, "Expected result to be of type *RowMock")
			},
		},
		{
			name: "Should default CallbackQueryRow",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.QueryRow(nil)
				assert.IsType(t, &RowMock{}, result, "Expected result to be of type *RowMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackQueryRowContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackQueryRowContext = func(ctx context.Context, arg interface{}) Row {
					callCount++
					return &RowMock{}
				}
				result := mock.QueryRowContext(context.Background(), nil)
				assert.Equal(t, 1, callCount, "Expected CallbackQueryRowContext to be called once")
				assert.IsType(t, &RowMock{}, result, "Expected result to be of type *RowMock")
			},
		},
		{
			name: "Should default CallbackQueryRowContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.QueryRowContext(context.Background(), nil)
				assert.IsType(t, &RowMock{}, result, "Expected result to be of type *RowMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackQuery",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackQuery = func(arg interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				result, err := mock.Query(nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackQuery to be called once")
				assert.IsType(t, &RowsMock{}, result, "Expected result to be of type *RowsMock")
			},
		},
		{
			name: "Should default CallbackQuery",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result, err := mock.Query(nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &RowsMock{}, result, "Expected result to be of type *RowsMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackQueryContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackQueryContext = func(ctx context.Context, arg interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				result, err := mock.QueryContext(context.Background(), nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackQueryContext to be called once")
				assert.IsType(t, &RowsMock{}, result, "Expected result to be of type *RowsMock")
			},
		},
		{
			name: "Should default CallbackQueryContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result, err := mock.QueryContext(context.Background(), nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &RowsMock{}, result, "Expected result to be of type *RowsMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackSelect",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				callCount := 0
				mock.CallbackSelect = func(dest interface{}, arg interface{}) error {
					callCount++
					*(dest).(*int) = 1
					return nil
				}
				err := mock.Select(&data, nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackSelect to be called once")
				assert.Equal(t, 1, data, "Expected data to be 1")
			},
		},
		{
			name: "Should default CallbackSelect",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				testErr := errors.New("failed to select")
				mock.Error = testErr
				err := mock.Select(nil, nil)
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
			},
		},
		{
			name: "Should run CallbackSelectContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				data := 0
				callCount := 0
				mock.CallbackSelectContext = func(ctx context.Context, dest interface{}, arg interface{}) error {
					callCount++
					*(dest).(*int) = 1
					return nil
				}
				err := mock.SelectContext(context.Background(), &data, nil)
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, 1, callCount, "Expected CallbackSelectContext to be called once")
				assert.Equal(t, 1, data, "Expected data to be 1")
			},
		},
		{
			name: "Should default CallbackSelectContext",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				testErr := errors.New("failed to select")
				mock.Error = testErr
				err := mock.SelectContext(context.Background(), nil, nil)
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
			},
		},
		{
			name: "Should run CallbackUnsafe",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackUnsafe = func() *sqlx.NamedStmt {
					callCount++
					return &sqlx.NamedStmt{}
				}
				result := mock.Unsafe()
				assert.Equal(t, 1, callCount, "Expected CallbackUnsafe to be called once")
				assert.IsType(t, &sqlx.NamedStmt{}, result, "Expected result to be of type *sqlx.NamedStmt")
			},
		},
		{
			name: "Should default CallbackUnsafe",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.Unsafe()
				assert.IsType(t, &sqlx.NamedStmt{}, result, "Expected result to be of type *sqlx.NamedStmt")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackSafe",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				callCount := 0
				mock.CallbackSafe = func() *sqlx.NamedStmt {
					callCount++
					return &sqlx.NamedStmt{}
				}
				result := mock.Safe()
				assert.Equal(t, 1, callCount, "Expected CallbackSafe to be called once")
				assert.IsType(t, &sqlx.NamedStmt{}, result, "Expected result to be of type *sqlx.NamedStmt")
			},
		},
		{
			name: "Should default CallbackSafe",
			mock: NamedStmtMock{},
			assert: func(t *testing.T, mock NamedStmtMock) {
				result := mock.Safe()
				assert.IsType(t, &sqlx.NamedStmt{}, result, "Expected result to be of type *sqlx.NamedStmt")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
