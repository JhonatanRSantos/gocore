package godb

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DriverRowsMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   DriverRowsMock
		assert func(t *testing.T, mock DriverRowsMock)
	}{
		{
			name: "Should run CallbackColumns",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				callCount := 0
				mock.CallbackColumns = func() []string {
					callCount++
					return []string{"column1", "column2"}
				}
				columns := mock.Columns()
				assert.Equal(t, 1, callCount, "Expected CallbackColumns to be called once")
				assert.Equal(t, []string{"column1", "column2"}, columns, "Expected columns to be [column1, column2]")
			},
		},
		{
			name: "Should run default CallbackColumns",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				columns := mock.Columns()
				assert.Equal(t, []string{}, columns, "Expected columns to be []")
			},
		},
		{
			name: "Should run CallbackClose",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				callCount := 0
				mock.CallbackClose = func() error {
					callCount++
					return nil
				}
				err := mock.Close()
				assert.Equal(t, 1, callCount, "Expected CallbackClose to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackClose",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				err := mock.Close()
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run CallbackNext",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				data := make([]driver.Value, 2)
				callCount := 0
				mock.CallbackNext = func(dest []driver.Value) error {
					callCount++
					dest[0] = "value1"
					dest[1] = "value2"
					return nil
				}
				err := mock.Next(data)
				assert.Equal(t, 1, callCount, "Expected CallbackNext to be called once")
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, []driver.Value{"value1", "value2"}, data, "Expected data to be [value1, value2]")
			},
		},
		{
			name: "Should run default CallbackNext",
			mock: DriverRowsMock{},
			assert: func(t *testing.T, mock DriverRowsMock) {
				err := mock.Next([]driver.Value{})
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}

func Test_DriverResultMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   DriverResultMock
		assert func(t *testing.T, mock DriverResultMock)
	}{
		{
			name: "Should run CallbackLastInsertId",
			mock: DriverResultMock{},
			assert: func(t *testing.T, mock DriverResultMock) {
				callCount := 0
				mock.CallbackLastInsertId = func() (int64, error) {
					callCount++
					return 1, nil
				}
				id, err := mock.LastInsertId()
				assert.Equal(t, 1, callCount, "Expected CallbackLastInsertId to be called once")
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, int64(1), id, "Expected id to be 1")
			},
		},
		{
			name: "Should run default CallbackLastInsertId",
			mock: DriverResultMock{},
			assert: func(t *testing.T, mock DriverResultMock) {
				id, err := mock.LastInsertId()
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, int64(0), id, "Expected id to be 0")
			},
		},
		{
			name: "Should run CallbackRowsAffected",
			mock: DriverResultMock{},
			assert: func(t *testing.T, mock DriverResultMock) {
				callCount := 0
				mock.CallbackRowsAffected = func() (int64, error) {
					callCount++
					return 2, nil
				}
				rows, err := mock.RowsAffected()
				assert.Equal(t, 1, callCount, "Expected CallbackRowsAffected to be called once")
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, int64(2), rows, "Expected rows to be 2")
			},
		},
		{
			name: "Should run default CallbackRowsAffected",
			mock: DriverResultMock{},
			assert: func(t *testing.T, mock DriverResultMock) {
				rows, err := mock.RowsAffected()
				assert.NoError(t, err, "Expected err to be nil")
				assert.Equal(t, int64(0), rows, "Expected rows to be 0")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}

func Test_DriverTxMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   DriverTxMock
		assert func(t *testing.T, mock DriverTxMock)
	}{
		{
			name: "Should run CallbackCommit",
			mock: DriverTxMock{},
			assert: func(t *testing.T, mock DriverTxMock) {
				callCount := 0
				mock.CallbackCommit = func() error {
					callCount++
					return nil
				}
				err := mock.Commit()
				assert.Equal(t, 1, callCount, "Expected CallbackCommit to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackCommit",
			mock: DriverTxMock{},
			assert: func(t *testing.T, mock DriverTxMock) {
				err := mock.Commit()
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run CallbackRollback",
			mock: DriverTxMock{},
			assert: func(t *testing.T, mock DriverTxMock) {
				callCount := 0
				mock.CallbackRollback = func() error {
					callCount++
					return nil
				}
				err := mock.Rollback()
				assert.Equal(t, 1, callCount, "Expected CallbackRollback to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackRollback",
			mock: DriverTxMock{},
			assert: func(t *testing.T, mock DriverTxMock) {
				err := mock.Rollback()
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}

func Test_DriverStmtMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   DriverStmtMock
		assert func(t *testing.T, mock DriverStmtMock)
	}{
		{
			name: "Should run CallbackClose",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
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
			name: "Should default CallbackClose",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				testErr := errors.New("failed to close statement")
				mock.Error = testErr
				err := mock.Close()
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
			},
		},
		{
			name: "Should run CallbackNumInput",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				callCount := 0
				mock.CallbackNumInput = func() int {
					callCount++
					return 2
				}
				input := mock.NumInput()
				assert.Equal(t, 1, callCount, "Expected CallbackNumInput to be called once")
				assert.Equal(t, 2, input, "Expected input to be 2")
			},
		},
		{
			name: "Should default CallbackNumInput",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				input := mock.NumInput()
				assert.Equal(t, 0, input, "Expected input to be 0")
			},
		},
		{
			name: "Should run CallbackExec",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				callCount := 0
				mock.CallbackExec = func(args []driver.Value) (driver.Result, error) {
					callCount++
					return &DriverResultMock{}, nil
				}
				result, err := mock.Exec([]driver.Value{})
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &DriverResultMock{}, result, "Expected result to be of type *DriverResultMock")
				assert.Empty(t, result, "Expected result to be empty")
				assert.Equal(t, 1, callCount, "Expected CallbackExec to be called once")
			},
		},
		{
			name: "Should default CallbackExec",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				testErr := errors.New("failed to execute statement")
				mock.Error = testErr
				result, err := mock.Exec([]driver.Value{})
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
				assert.IsType(t, &DriverResultMock{}, result, "Expected result to be of type *DriverResultMock")
				assert.Empty(t, result, "Expected result to be empty")
			},
		},
		{
			name: "Should run CallbackQuery",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				callCount := 0
				mock.CallbackQuery = func(args []driver.Value) (driver.Rows, error) {
					callCount++
					return &DriverRowsMock{}, nil
				}
				rows, err := mock.Query([]driver.Value{})
				assert.NoError(t, err, "Expected err to be nil")
				assert.IsType(t, &DriverRowsMock{}, rows, "Expected rows to be of type *DriverRowsMock")
				assert.Empty(t, rows, "Expected rows to be empty")
				assert.Equal(t, 1, callCount, "Expected CallbackQuery to be called once")
			},
		},
		{
			name: "Should default CallbackQuery",
			mock: DriverStmtMock{},
			assert: func(t *testing.T, mock DriverStmtMock) {
				testErr := errors.New("failed to query statement")
				mock.Error = testErr
				rows, err := mock.Query([]driver.Value{})
				assert.ErrorIs(t, err, testErr, "Expected err to be testErr")
				assert.IsType(t, &DriverRowsMock{}, rows, "Expected rows to be of type *DriverRowsMock")
				assert.Empty(t, rows, "Expected rows to be empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}

func Test_DriverConnMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   DriverConnMock
		assert func(t *testing.T, mock DriverConnMock)
	}{
		{
			name: "Should run CallbackPrepare",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				callCount := 0
				mock.CallbackPrepare = func(query string) (driver.Stmt, error) {
					callCount++
					return nil, nil
				}
				_, err := mock.Prepare("SELECT * FROM table")
				assert.Equal(t, 1, callCount, "Expected CallbackPrepare to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackPrepare",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				_, err := mock.Prepare("SELECT * FROM table")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run CallbackClose",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				callCount := 0
				mock.CallbackClose = func() error {
					callCount++
					return nil
				}
				err := mock.Close()
				assert.Equal(t, 1, callCount, "Expected CallbackClose to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackClose",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				err := mock.Close()
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run CallbackBegin",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				callCount := 0
				mock.CallbackBegin = func() (driver.Tx, error) {
					callCount++
					return nil, nil
				}
				_, err := mock.Begin()
				assert.Equal(t, 1, callCount, "Expected CallbackBegin to be called once")
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
		{
			name: "Should run default CallbackBegin",
			mock: DriverConnMock{},
			assert: func(t *testing.T, mock DriverConnMock) {
				_, err := mock.Begin()
				assert.NoError(t, err, "Expected err to be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
