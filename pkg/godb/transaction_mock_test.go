package godb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_TxMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   TxMock
		assert func(t *testing.T, mock TxMock)
	}{
		{
			name: "Should run CallbackCommit",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackCommit = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, mock.Commit(), "Should not return error")
				assert.Equal(t, 1, callCount, "Should call CallbackCommit once")
			},
		},
		{
			name: "Should run default CallbackCommit",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to commit")
				mock.Error = testErr
				mock.CallbackCommit = nil
				assert.ErrorIs(t, mock.Commit(), testErr, "Should return error")
			},
		},
		{
			name: "Should run CallbackExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackExec = func(query string, args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.Exec("query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackExec once")
			},
		},
		{
			name: "Should run default CallbackExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to exec")
				mock.Error = testErr
				mock.CallbackExec = nil
				result, err := mock.Exec("query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackExecContext = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.ExecContext(context.Background(), "query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackExecContext once")
			},
		},
		{
			name: "Should run default CallbackExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to exec context")
				mock.Error = testErr
				mock.CallbackExecContext = nil
				result, err := mock.ExecContext(context.Background(), "query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackRollback",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackRollback = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, mock.Rollback(), "Should not return error")
				assert.Equal(t, 1, callCount, "Should call CallbackRollback once")
			},
		},
		{
			name: "Should run default CallbackRollback",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to rollback")
				mock.Error = testErr
				mock.CallbackRollback = nil
				assert.ErrorIs(t, mock.Rollback(), testErr, "Should return error")
			},
		},
		{
			name: "Should run CallbackBindNamed",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackBindNamed = func(query string, arg interface{}) (string, []interface{}, error) {
					callCount++
					return "query", []interface{}{"arg"}, nil
				}
				bquery, args, err := mock.BindNamed("query", "arg")
				assert.NoError(t, err, "Should not return error")
				assert.Equal(t, "query", bquery, "Should return query")
				assert.Equal(t, []interface{}{"arg"}, args, "Should return args")
				assert.Equal(t, 1, callCount, "Should call CallbackBindNamed once")
			},
		},
		{
			name: "Should run default CallbackBindNamed",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to bind named")
				mock.Error = testErr
				mock.CallbackBindNamed = nil
				bquery, args, err := mock.BindNamed("query", "arg")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.Empty(t, bquery, "Should return empty query")
				assert.Empty(t, args, "Should return empty args")
			},
		},
		{
			name: "Should run CallbackDriverName",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackDriverName = func() string {
					callCount++
					return "driver"
				}
				assert.Equal(t, "driver", mock.DriverName(), "Should return driver name")
				assert.Equal(t, 1, callCount, "Should call CallbackDriverName once")
			},
		},
		{
			name: "Should run default CallbackDriverName",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackDriverName = nil
				assert.Empty(t, mock.DriverName(), "Should return empty driver name")
			},
		},
		{
			name: "Should run CallbackGet",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				data := 0
				callCount := 0
				mock.CallbackGet = func(dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = 1
					return nil
				}
				assert.NoError(t, mock.Get(&data, "query"), "Should not return error")
				assert.Equal(t, 1, data, "Should return 1")
				assert.Equal(t, 1, callCount, "Should call CallbackGet once")
			},
		},
		{
			name: "Should run default CallbackGet",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to get")
				mock.Error = testErr
				mock.CallbackGet = nil
				data := 0
				assert.ErrorIs(t, mock.Get(&data, "query"), testErr, "Should return error")
				assert.Equal(t, 0, data, "Should return 0")
			},
		},
		{
			name: "Should run CallbackGetContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				data := 0
				callCount := 0
				mock.CallbackGetContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = 1
					return nil
				}
				assert.NoError(t, mock.GetContext(context.Background(), &data, "query"), "Should not return error")
				assert.Equal(t, 1, data, "Should return 1")
				assert.Equal(t, 1, callCount, "Should call CallbackGetContext once")
			},
		},
		{
			name: "Should run default CallbackGetContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to get context")
				mock.Error = testErr
				mock.CallbackGetContext = nil
				data := 0
				assert.ErrorIs(t, mock.GetContext(context.Background(), &data, "query"), testErr, "Should return error")
				assert.Equal(t, 0, data, "Should return 0")
			},
		},
		{
			name: "Should run CallbackMustExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackMustExec = func(query string, args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExec("query")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackMustExec once")
			},
		},
		{
			name: "Should run default CallbackMustExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to must exec")
				mock.Error = testErr
				mock.CallbackMustExec = nil
				result := mock.MustExec("query")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackMustExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackMustExecContext = func(ctx context.Context, query string, args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{}
				}
				result := mock.MustExecContext(context.Background(), "query")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackMustExecContext once")
			},
		},
		{
			name: "Should run default CallbackMustExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to must exec context")
				mock.Error = testErr
				mock.CallbackMustExecContext = nil
				result := mock.MustExecContext(context.Background(), "query")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackNamedExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackNamedExec = func(query string, arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.NamedExec("query", "arg")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackNamedExec once")
			},
		},
		{
			name: "Should run default CallbackNamedExec",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to named exec")
				mock.Error = testErr
				mock.CallbackNamedExec = nil
				result, err := mock.NamedExec("query", "arg")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackNamedExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackNamedExecContext = func(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{}, nil
				}
				result, err := mock.NamedExecContext(context.Background(), "query", "arg")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
				assert.Equal(t, 1, callCount, "Should call CallbackNamedExecContext once")
			},
		},
		{
			name: "Should run default CallbackNamedExecContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to named exec context")
				mock.Error = testErr
				mock.CallbackNamedExecContext = nil
				result, err := mock.NamedExecContext(context.Background(), "query", "arg")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &ResultMock{}, result, "Should return sql.Result")
				assert.Empty(t, result, "Should return empty result")
			},
		},
		{
			name: "Should run CallbackNamedQuery",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackNamedQuery = func(query string, arg interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := mock.NamedQuery("query", "arg")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
				assert.Equal(t, 1, callCount, "Should call CallbackNamedQuery once")
			},
		},
		{
			name: "Should run default CallbackNamedQuery",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to named query")
				mock.Error = testErr
				mock.CallbackNamedQuery = nil
				rows, err := mock.NamedQuery("query", "arg")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
			},
		},
		{
			name: "Should run CallbackNamedStmt",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackNamedStmt = func(stmt NamedStmt) NamedStmt {
					callCount++
					return &NamedStmtMock{}
				}
				stmt := mock.NamedStmt(&NamedStmtMock{})
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmtMock")
				assert.Empty(t, stmt, "Should return empty NamedStmtMock")
				assert.Equal(t, 1, callCount, "Should call CallbackNamedStmt once")
			},
		},
		{
			name: "Should run default CallbackNamedStmt",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackNamedStmt = nil
				stmt := mock.NamedStmt(&NamedStmtMock{})
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmtMock")
				assert.Empty(t, stmt, "Should return empty NamedStmtMock")
			},
		},
		{
			name: "Should run CallbackNamedStmtContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackNamedStmtContext = func(ctx context.Context, stmt NamedStmt) NamedStmt {
					callCount++
					return &NamedStmtMock{}
				}
				stmt := mock.NamedStmtContext(context.Background(), &NamedStmtMock{})
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmtMock")
				assert.Empty(t, stmt, "Should return empty NamedStmtMock")
				assert.Equal(t, 1, callCount, "Should call CallbackNamedStmtContext once")
			},
		},
		{
			name: "Should run default CallbackNamedStmtContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackNamedStmtContext = nil
				stmt := mock.NamedStmtContext(context.Background(), &NamedStmtMock{})
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmtMock")
				assert.Empty(t, stmt, "Should return empty NamedStmtMock")
			},
		},
		{
			name: "Should run CallbackPrepareNamed",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackPrepareNamed = func(query string) (NamedStmt, error) {
					callCount++
					return &NamedStmtMock{}, nil
				}
				stmt, err := mock.PrepareNamed("query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmt")
				assert.Empty(t, stmt, "Should return empty NamedStmt")
				assert.Equal(t, 1, callCount, "Should call CallbackPrepareNamed once")
			},
		},
		{
			name: "Should run default CallbackPrepareNamed",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to prepare named")
				mock.Error = testErr
				mock.CallbackPrepareNamed = nil
				stmt, err := mock.PrepareNamed("query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmt")
				assert.Empty(t, stmt, "Should return empty NamedStmt")
			},
		},
		{
			name: "Should run CallbackPrepareNamedContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackPrepareNamedContext = func(ctx context.Context, query string) (NamedStmt, error) {
					callCount++
					return &NamedStmtMock{}, nil
				}
				stmt, err := mock.PrepareNamedContext(context.Background(), "query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmt")
				assert.Empty(t, stmt, "Should return empty NamedStmt")
				assert.Equal(t, 1, callCount, "Should call CallbackPrepareNamedContext once")
			},
		},
		{
			name: "Should run default CallbackPrepareNamedContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to prepare named context")
				mock.Error = testErr
				mock.CallbackPrepareNamedContext = nil
				stmt, err := mock.PrepareNamedContext(context.Background(), "query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "Should return NamedStmt")
				assert.Empty(t, stmt, "Should return empty NamedStmt")
			},
		},
		{
			name: "Should run CallbackPrepare",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackPrepare = func(query string) (Stmt, error) {
					callCount++
					return &StmtMock{}, nil
				}
				stmt, err := mock.Prepare("query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
				assert.Equal(t, 1, callCount, "Should call CallbackPrepare once")
			},
		},
		{
			name: "Should run default CallbackPrepare",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to prepare")
				mock.Error = testErr
				mock.CallbackPrepare = nil
				stmt, err := mock.Prepare("query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
			},
		},
		{
			name: "Should run CallbackPrepareContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackPrepareContext = func(ctx context.Context, query string) (Stmt, error) {
					callCount++
					return &StmtMock{}, nil
				}
				stmt, err := mock.PrepareContext(context.Background(), "query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
				assert.Equal(t, 1, callCount, "Should call CallbackPrepareContext once")
			},
		},
		{
			name: "Should run default CallbackPrepareContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to prepare context")
				mock.Error = testErr
				mock.CallbackPrepareContext = nil
				stmt, err := mock.PrepareContext(context.Background(), "query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
			},
		},
		{
			name: "Should run CallbackQueryRow",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackQueryRow = func(query string, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := mock.QueryRow("query")
				assert.IsType(t, &RowMock{}, row, "Should return Row")
				assert.Empty(t, row, "Should return empty Row")
				assert.Equal(t, 1, callCount, "Should call CallbackQueryRow once")
			},
		},
		{
			name: "Should run default CallbackQueryRow",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackQueryRow = nil
				row := mock.QueryRow("query")
				assert.IsType(t, &RowMock{}, row, "Should return Row")
				assert.Empty(t, row, "Should return empty Row")
			},
		},
		{
			name: "Should run CallbackQueryRowContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackQueryRowContext = func(ctx context.Context, query string, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := mock.QueryRowContext(context.Background(), "query")
				assert.IsType(t, &RowMock{}, row, "Should return Row")
				assert.Empty(t, row, "Should return empty Row")
				assert.Equal(t, 1, callCount, "Should call CallbackQueryRowContext once")
			},
		},
		{
			name: "Should run default CallbackQueryRowContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackQueryRowContext = nil
				row := mock.QueryRowContext(context.Background(), "query")
				assert.IsType(t, &RowMock{}, row, "Should return Row")
				assert.Empty(t, row, "Should return empty Row")
			},
		},
		{
			name: "Should run CallbackQuery",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackQuery = func(query string, args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := mock.Query("query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
				assert.Equal(t, 1, callCount, "Should call CallbackQuery once")
			},
		},
		{
			name: "Should run default CallbackQuery",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to query")
				mock.Error = testErr
				mock.CallbackQuery = nil
				rows, err := mock.Query("query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
			},
		},
		{
			name: "Should run CallbackQueryContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackQueryContext = func(ctx context.Context, query string, args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := mock.QueryContext(context.Background(), "query")
				assert.NoError(t, err, "Should not return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
				assert.Equal(t, 1, callCount, "Should call CallbackQueryContext once")
			},
		},
		{
			name: "Should run default CallbackQueryContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to query context")
				mock.Error = testErr
				mock.CallbackQueryContext = nil
				rows, err := mock.QueryContext(context.Background(), "query")
				assert.ErrorIs(t, err, testErr, "Should return error")
				assert.IsType(t, &RowsMock{}, rows, "Should return sql.Rows")
				assert.Empty(t, rows, "Should return empty rows")
			},
		},
		{
			name: "Should run CallbackRebind",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackRebind = func(query string) string {
					callCount++
					return "query"
				}
				assert.Equal(t, "query", mock.Rebind("query"), "Should return query")
				assert.Equal(t, 1, callCount, "Should call CallbackRebind once")
			},
		},
		{
			name: "Should run default CallbackRebind",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackRebind = nil
				assert.Empty(t, mock.Rebind("query"), "Should return empty query")
			},
		},
		{
			name: "Should run CallbackSelect",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				data := 0
				callCount := 0
				mock.CallbackSelect = func(dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = 1
					return nil
				}
				assert.NoError(t, mock.Select(&data, "query"), "Should not return error")
				assert.Equal(t, 1, data, "Should return 1")
				assert.Equal(t, 1, callCount, "Should call CallbackSelect once")
			},
		},
		{
			name: "Should run default CallbackSelect",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to select")
				mock.Error = testErr
				mock.CallbackSelect = nil
				data := 0
				assert.ErrorIs(t, mock.Select(&data, "query"), testErr, "Should return error")
				assert.Equal(t, 0, data, "Should return 0")
			},
		},
		{
			name: "Should run CallbackSelectContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				data := 0
				mock.CallbackSelectContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = 1
					return nil
				}
				assert.NoError(t, mock.SelectContext(context.Background(), &data, "query"), "Should not return error")
				assert.Equal(t, 1, data, "Should return 1")
				assert.Equal(t, 1, callCount, "Should call CallbackSelectContext once")
			},
		},
		{
			name: "Should run default CallbackSelectContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				testErr := errors.New("failed to select context")
				mock.Error = testErr
				mock.CallbackSelectContext = nil
				data := 0
				assert.ErrorIs(t, mock.SelectContext(context.Background(), &data, "query"), testErr, "Should return error")
				assert.Equal(t, 0, data, "Should return 0")
			},
		},
		{
			name: "Should run CallbackStmt",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackStmt = func(stmt interface{}) Stmt {
					callCount++
					return &StmtMock{}
				}
				stmt := mock.Stmt(StmtMock{})
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
				assert.Equal(t, 1, callCount, "Should call CallbackStmt once")
			},
		},
		{
			name: "Should run default CallbackStmt",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackStmt = nil
				stmt := mock.Stmt(StmtMock{})
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
			},
		},
		{
			name: "Should run CallbackStmtContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackStmtContext = func(ctx context.Context, stmt interface{}) Stmt {
					callCount++
					return &StmtMock{}
				}
				stmt := mock.StmtContext(context.Background(), StmtMock{})
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
				assert.Equal(t, 1, callCount, "Should call CallbackStmtContext once")
			},
		},
		{
			name: "Should run default CallbackStmtContext",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackStmtContext = nil
				stmt := mock.StmtContext(context.Background(), StmtMock{})
				assert.IsType(t, &StmtMock{}, stmt, "Should return Stmt")
				assert.Empty(t, stmt, "Should return empty Stmt")
			},
		},
		{
			name: "Should run CallbackUnsafe",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackUnsafe = func() *sqlx.Tx {
					callCount++
					return &sqlx.Tx{}
				}
				assert.IsType(t, &sqlx.Tx{}, mock.Unsafe(), "Should return sqlx.Tx")
				assert.Equal(t, 1, callCount, "Should call CallbackUnsafe once")
			},
		},
		{
			name: "Should run default CallbackUnsafe",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackUnsafe = nil
				assert.Empty(t, mock.Unsafe(), "Should return empty sqlx.Tx")
			},
		},
		{
			name: "Should run CallbackSafe",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				callCount := 0
				mock.CallbackSafe = func() *sqlx.Tx {
					callCount++
					return &sqlx.Tx{}
				}
				assert.IsType(t, &sqlx.Tx{}, mock.Safe(), "Should return sqlx.Tx")
				assert.Equal(t, 1, callCount, "Should call CallbackSafe once")
			},
		},
		{
			name: "Should run default CallbackSafe",
			mock: TxMock{},
			assert: func(t *testing.T, mock TxMock) {
				mock.CallbackSafe = nil
				assert.Empty(t, mock.Safe(), "Should return empty sqlx.Tx")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
