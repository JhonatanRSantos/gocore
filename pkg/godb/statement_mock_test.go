package godb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_StmtMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   StmtMock
		assert func(t *testing.T, mock StmtMock)
	}{
		{
			name: "Should run CallbackClose",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackClose = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, mock.Close(), "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackClose should be called once")
			},
		},
		{
			name: "Should run default CallbackClose",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to close")
				mock.Error = testErr
				assert.ErrorIs(t, mock.Close(), testErr, "Error should be the same as the Error field")
			},
		},
		{
			name: "Should run CallbackExec",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackExec = func(args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.Exec()
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Equal(t, 1, callCount, "CallbackExec should be called once")
			},
		},
		{
			name: "Should run default CallbackExec",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to exec")
				mock.Error = testErr
				result, err := mock.Exec()
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Empty(t, result, "Result should be empty")
			},
		},
		{
			name: "Should run CallbackExecContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackExecContext = func(ctx context.Context, args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.ExecContext(context.Background())
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Equal(t, 1, callCount, "CallbackExecContext should be called once")
			},
		},
		{
			name: "Should run default CallbackExecContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to exec context")
				mock.Error = testErr
				result, err := mock.ExecContext(context.Background())
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Empty(t, result, "Result should be empty")
			},
		},
		{
			name: "Should run CallbackGet",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				data := 0
				callCount := 0
				mock.CallbackGet = func(dest interface{}, args ...interface{}) error {
					callCount++
					*(dest.(*int)) = 1
					return nil
				}
				err := mock.Get(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data, "Data should be 1")
			},
		},
		{
			name: "Should run default CallbackGet",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to get")
				mock.Error = testErr
				data := 0
				err := mock.Get(&data)
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.Equal(t, 0, data, "Data should be 0")
			},
		},
		{
			name: "Should run CallbackGetContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				data := 0
				callCount := 0
				mock.CallbackGetContext = func(ctx context.Context, dest interface{}, args ...interface{}) error {
					callCount++
					*(dest.(*int)) = 1
					return nil
				}
				err := mock.GetContext(context.Background(), &data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data, "Data should be 1")
			},
		},
		{
			name: "Should run default CallbackGetContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to get context")
				mock.Error = testErr
				data := 0
				err := mock.GetContext(context.Background(), &data)
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.Equal(t, 0, data, "Data should be 0")
			},
		},
		{
			name: "Should run CallbackMustExec",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackMustExec = func(args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExec()
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Equal(t, 1, callCount, "CallbackMustExec should be called once")
			},
		},
		{
			name: "Should run default CallbackMustExec",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				result := mock.MustExec()
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
			},
		},
		{
			name: "Should run CallbackMustExecContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackMustExecContext = func(ctx context.Context, args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExecContext(context.Background())
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
				assert.Equal(t, 1, callCount, "CallbackMustExecContext should be called once")
			},
		},
		{
			name: "Should run default CallbackMustExecContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				result := mock.MustExecContext(context.Background())
				assert.IsType(t, &ResultMock{}, result, "Result should be a ResultMock")
			},
		},
		{
			name: "Should run CallbackQueryRow",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackQueryRow = func(args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := mock.QueryRow()
				assert.IsType(t, &RowMock{}, row, "Row should be a RowMock")
				assert.Equal(t, 1, callCount, "CallbackQueryRow should be called once")
			},
		},
		{
			name: "Should run default CallbackQueryRow",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				row := mock.QueryRow()
				assert.IsType(t, &RowMock{}, row, "Row should be a RowMock")
			},
		},
		{
			name: "Should run CallbackQueryRowContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackQueryRowContext = func(ctx context.Context, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := mock.QueryRowContext(context.Background())
				assert.IsType(t, &RowMock{}, row, "Row should be a RowMock")
				assert.Equal(t, 1, callCount, "CallbackQueryRowContext should be called once")
			},
		},
		{
			name: "Should run default CallbackQueryRowContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				row := mock.QueryRowContext(context.Background())
				assert.IsType(t, &RowMock{}, row, "Row should be a RowMock")
			},
		},
		{
			name: "Should run CallbackQuery",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackQuery = func(args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := mock.Query()
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, &RowsMock{}, rows, "Rows should be a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackQuery should be called once")
			},
		},
		{
			name: "Should run default CallbackQuery",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to query")
				mock.Error = testErr
				rows, err := mock.Query()
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.IsType(t, &RowsMock{}, rows, "Rows should be a RowsMock")
				assert.Empty(t, rows, "Rows should be empty")
			},
		},
		{
			name: "Should run CallbackQueryContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackQueryContext = func(ctx context.Context, args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := mock.QueryContext(context.Background())
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, &RowsMock{}, rows, "Rows should be a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackQueryContext should be called once")
			},
		},
		{
			name: "Should run default CallbackQueryContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to query context")
				mock.Error = testErr
				rows, err := mock.QueryContext(context.Background())
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.IsType(t, &RowsMock{}, rows, "Rows should be a RowsMock")
				assert.Empty(t, rows, "Rows should be empty")
			},
		},
		{
			name: "Should run CallbackSelect",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				data := 0
				callCount := 0
				mock.CallbackSelect = func(dest interface{}, args ...interface{}) error {
					callCount++
					*(dest.(*int)) = 1
					return nil
				}
				err := mock.Select(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data, "Data should be 1")
			},
		},
		{
			name: "Should run default CallbackSelect",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to select")
				mock.Error = testErr
				data := 0
				err := mock.Select(&data)
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.Equal(t, 0, data, "Data should be 0")
			},
		},
		{
			name: "Should run CallbackSelectContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				data := 0
				callCount := 0
				mock.CallbackSelectContext = func(ctx context.Context, dest interface{}, args ...interface{}) error {
					callCount++
					*(dest.(*int)) = 1
					return nil
				}
				err := mock.SelectContext(context.Background(), &data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data, "Data should be 1")
			},
		},
		{
			name: "Should run default CallbackSelectContext",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				testErr := errors.New("failed to select context")
				mock.Error = testErr
				data := 0
				err := mock.SelectContext(context.Background(), &data)
				assert.ErrorIs(t, err, testErr, "Error should be the same as the Error field")
				assert.Equal(t, 0, data, "Data should be 0")
			},
		},
		{
			name: "Should run CallbackUnsafe",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackUnsafe = func() *sqlx.Stmt {
					callCount++
					return &sqlx.Stmt{}
				}
				stmt := mock.Unsafe()
				assert.IsType(t, &sqlx.Stmt{}, stmt, "Stmt should be a sqlx.Stmt")
				assert.Equal(t, 1, callCount, "CallbackUnsafe should be called once")
			},
		},
		{
			name: "Should run default CallbackUnsafe",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				stmt := mock.Unsafe()
				assert.IsType(t, &sqlx.Stmt{}, stmt, "Stmt should be a sqlx.Stmt")
			},
		},
		{
			name: "Should run CallbackSafe",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				callCount := 0
				mock.CallbackSafe = func() *sqlx.Stmt {
					callCount++
					return &sqlx.Stmt{}
				}
				stmt := mock.Safe()
				assert.IsType(t, &sqlx.Stmt{}, stmt, "Stmt should be a sqlx.Stmt")
				assert.Equal(t, 1, callCount, "CallbackSafe should be called once")
			},
		},
		{
			name: "Should run default CallbackSafe",
			mock: StmtMock{},
			assert: func(t *testing.T, mock StmtMock) {
				stmt := mock.Safe()
				assert.IsType(t, &sqlx.Stmt{}, stmt, "Stmt should be a sqlx.Stmt")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
