package godb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_DBMock(t *testing.T) {
	tests := []struct {
		name   string
		dbMock *DBMock
		assert func(t *testing.T, dbMock *DBMock)
	}{
		{
			name:   "Should run CallbackClose",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackClose = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, dbMock.Close(), "Close should not return an error")
				assert.Equal(t, 1, callCount, "CallbackClose should be called once")
			},
		},
		{
			name:   "should run default CallbackClose",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackClose = nil
				testErr := errors.New("failed to run CallbackClose")
				dbMock.Error = testErr
				assert.ErrorIs(t, dbMock.Close(), testErr, "Close should return an error")
			},
		},
		{
			name:   "Should run CallbackDriver",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackDriver = func() driver.Driver {
					callCount++
					return &DriverMock{}
				}
				assert.IsType(t, &DriverMock{}, dbMock.Driver(), "Driver should return a DriverMock")
				assert.Equal(t, 1, callCount, "CallbackDriver should be called once")
			},
		},
		{
			name:   "should run default CallbackDriver",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackDriver")
				dbMock.Error = testErr
				dbMock.CallbackDriver = nil

				if driver, ok := dbMock.Driver().(*DriverMock); ok {
					assert.ErrorIs(t, driver.Error, testErr, "mismatch error")
				} else {
					assert.Fail(t, "Driver should return an error")
				}
			},
		},
		{
			name:   "Should run CallbackExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackExec = func(query string, args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}, nil
				}
				result, err := dbMock.Exec("")
				assert.NoError(t, err, "Exec should not return an error")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackExec should be called once")
			},
		},
		{
			name:   "should run default CallbackExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackExec")
				dbMock.Error = testErr
				dbMock.CallbackExec = nil
				_, err := dbMock.Exec("")
				assert.ErrorIs(t, err, testErr, "Exec should return an error")
			},
		},
		{
			name:   "Should run CallbackExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackExecContext = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}, nil
				}
				result, err := dbMock.ExecContext(context.Background(), "")
				assert.NoError(t, err, "ExecContext should not return an error")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackExecContext should be called once")
			},
		},
		{
			name:   "should run default CallbackExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackExecContext")
				dbMock.Error = testErr
				dbMock.CallbackExecContext = nil
				result, err := dbMock.ExecContext(context.Background(), "")
				assert.ErrorIs(t, err, testErr, "ExecContext should return an error")
				assert.IsType(t, &ResultMock{}, result, "ExecContext should return an empty result")
				assert.Empty(t, result, "ExecContext should return an empty result")
			},
		},
		{
			name:   "Should run CallbackPing",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPing = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, dbMock.Ping(), "Ping should not return an error")
				assert.Equal(t, 1, callCount, "CallbackPing should be called once")
			},
		},
		{
			name:   "should run default CallbackPing",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackPing")
				dbMock.Error = testErr
				dbMock.CallbackPing = nil
				assert.ErrorIs(t, dbMock.Ping(), testErr, "Ping should return an error")
			},
		},
		{
			name:   "Should run default CallbackPingContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackPingContext = nil
				testErr := errors.New("failed to run CallbackPingContext")
				dbMock.Error = testErr
				assert.ErrorIs(t, dbMock.PingContext(context.Background()), testErr, "PingContext should return an error")
			},
		},
		{
			name:   "Should run default CallbackPingContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPingContext = func(ctx context.Context) error {
					callCount++
					return nil
				}
				assert.NoError(t, dbMock.PingContext(context.Background()), "PingContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackPingContext should be called once")
			},
		},
		{
			name:   "Should run CallbackSetConnMaxIdleTime",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackSetConnMaxIdleTime = func(d time.Duration) {
					callCount++
				}
				dbMock.SetConnMaxIdleTime(1)
				assert.Equal(t, 1, callCount, "CallbackSetConnMaxIdleTime should be called once")
			},
		},
		{
			name:   "Should run CallbackSetConnMaxLifetime",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackSetConnMaxLifetime = func(d time.Duration) {
					callCount++
				}
				dbMock.SetConnMaxLifetime(1)
				assert.Equal(t, 1, callCount, "CallbackSetConnMaxLifetime should be called once")
			},
		},
		{
			name:   "Should run CallbackSetMaxOpenConns",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackSetMaxOpenConns = func(n int) {
					callCount++
				}
				dbMock.SetMaxOpenConns(1)
				assert.Equal(t, 1, callCount, "CallbackSetMaxOpenConns should be called once")
			},
		},
		{
			name:   "Should run CallbackSetMaxIdleConns",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackSetMaxIdleConns = func(n int) {
					callCount++
				}
				dbMock.SetMaxIdleConns(1)
				assert.Equal(t, 1, callCount, "CallbackSetMaxIdleConns should be called once")
			},
		},
		{
			name:   "Should run CallbackStats",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackStats = func() sql.DBStats {
					callCount++
					return sql.DBStats{MaxOpenConnections: 1024}
				}
				assert.Equal(t, 1024, dbMock.Stats().MaxOpenConnections, "Stats should return MaxOpenConnections = 1024")
				assert.Equal(t, 1, callCount, "CallbackStats should be called once")
			},
		},
		{
			name:   "Should run default CallbackStats",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackStats = nil
				assert.Equal(t, sql.DBStats{}, dbMock.Stats(), "Stats should return an empty sql.DBStats")
			},
		},
		{
			name:   "Should run CallbackBeginTx",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackBeginTx = func(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
					callCount++
					return &customTx{tx: nil}, nil
				}
				tx, err := dbMock.BeginTx(context.Background(), nil)
				assert.NoError(t, err, "BeginTx should not return an error")
				assert.Nil(t, tx.Safe(), "BeginTx should return a customTx")
				assert.Equal(t, 1, callCount, "CallbackBeginTx should be called once")
			},
		},
		{
			name:   "Should run default CallbackBeginTx",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackBeginTx")
				dbMock.Error = testErr
				dbMock.CallbackBeginTx = nil
				_, err := dbMock.BeginTx(context.Background(), nil)
				assert.ErrorIs(t, err, testErr, "BeginTx should return an error")
			},
		},
		{
			name:   "Should run CallbackBegin",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackBegin = func() (Tx, error) {
					callCount++
					return &customTx{tx: nil}, nil
				}
				tx, err := dbMock.Begin()
				assert.NoError(t, err, "Begin should not return an error")
				assert.Nil(t, tx.Safe(), "Begin should return a customTx")
				assert.Equal(t, 1, callCount, "CallbackBegin should be called once")
			},
		},
		{
			name:   "Should run default CallbackBegin",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackBegin")
				dbMock.Error = testErr
				dbMock.CallbackBegin = nil
				_, err := dbMock.Begin()
				assert.ErrorIs(t, err, testErr, "Begin should return an error")
			},
		},
		{
			name:   "Should run CallbackBindNamed",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackBindNamed = func(query string, arg interface{}) (string, []interface{}, error) {
					callCount++
					return "select 1", []interface{}{"select 1"}, nil
				}
				bquery, args, err := dbMock.BindNamed("", nil)
				assert.NoError(t, err, "BindNamed should not return an error")
				assert.Equal(t, "select 1", bquery, "BindNamed should retunr 'select 1'")
				assert.Equal(t, []interface{}{"select 1"}, args, "BindNamed should return ['select 1']")
				assert.Equal(t, 1, callCount, "CallbackBindNamed should be called once")
			},
		},
		{
			name:   "Should run default CallbackBindNamed",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackBindNamed")
				dbMock.Error = testErr
				dbMock.CallbackBindNamed = nil
				bquery, args, err := dbMock.BindNamed("", nil)
				assert.ErrorIs(t, err, testErr, "BindNamed should return an error")
				assert.Empty(t, bquery, "BindNamed should return an empty string")
				assert.Empty(t, args, "BindNamed should return an empty slice")
			},
		},
		{
			name:   "Should run CallbackConn",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackConn = func(ctx context.Context) (Conn, error) {
					callCount++
					return &ConnMock{}, nil
				}
				conn, err := dbMock.Conn(context.Background())
				assert.NoError(t, err, "Conn should not return an error")
				assert.IsType(t, &ConnMock{}, conn, "Conn should return a ConnMock")
				assert.Equal(t, 1, callCount, "CallbackConn should be called once")
			},
		},
		{
			name:   "Should run default CallbackConn",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackConn")
				dbMock.Error = testErr
				dbMock.CallbackConn = nil
				_, err := dbMock.Conn(context.Background())
				assert.ErrorIs(t, err, testErr, "Conn should return an error")
			},
		},
		{
			name:   "Should run CallbackDriverName",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackDriverName = func() string {
					callCount++
					return "test driver"
				}
				assert.Equal(t, "test driver", dbMock.DriverName(), "DriverName should return 'test driver'")
				assert.Equal(t, 1, callCount, "CallbackDriverName should be called once")
			},
		},
		{
			name:   "Should run default CallbackDriverName",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackDriverName")
				dbMock.Error = testErr
				dbMock.CallbackDriverName = nil
				assert.Equal(t, "", dbMock.DriverName(), "DriverName should return an empty string")
				assert.ErrorIs(t, dbMock.Error, testErr, "DriverName should return an error")
			},
		},
		{
			name:   "Should run CallbackGet",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				data := 0
				callCount := 0
				dbMock.CallbackGet = func(dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = callCount + *(dest.(*int))
					return nil
				}
				err := dbMock.Get(&data, "")
				assert.NoError(t, err, "Get should not return an error")
				assert.Equal(t, 1, data, "Get should return 1")
				assert.Equal(t, 1, callCount, "CallbackGet should be called once")
			},
		},
		{
			name:   "Should run default CallbackGet",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackGet")
				dbMock.Error = testErr
				dbMock.CallbackGet = nil
				err := dbMock.Get(nil, "")
				assert.ErrorIs(t, err, testErr, "Get should return an error")
			},
		},
		{
			name:   "Should run CallbackGetContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				data := 0
				callCount := 0
				dbMock.CallbackGetContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = callCount + *(dest.(*int))
					return nil
				}
				err := dbMock.GetContext(context.Background(), &data, "")
				assert.NoError(t, err, "GetContext should not return an error")
				assert.Equal(t, 1, data, "GetContext should return 1")
				assert.Equal(t, 1, callCount, "CallbackGetContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackGetContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackGetContext")
				dbMock.Error = testErr
				dbMock.CallbackGetContext = nil
				err := dbMock.GetContext(context.Background(), nil, "")
				assert.ErrorIs(t, err, testErr, "GetContext should return an error")
			},
		},
		{
			name:   "Should run CallbackMapperFunc",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				mapperFunc := func(data string) string { return fmt.Sprintf("test - %s", data) }
				mapperData := ""
				dbMock.CallbackMapperFunc = func(mf func(string) string) {
					callCount++
					mapperData = mf(fmt.Sprint(callCount))
				}
				dbMock.MapperFunc(mapperFunc)
				assert.Equal(t, "test - 1", mapperData, "MapperFunc should return 'test - 1'")
				assert.Equal(t, 1, callCount, "CallbackMapperFunc should be called once")
			},
		},
		{
			name:   "Should run default CallbackMapperFunc",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				mapperFunc := func(data string) string {
					callCount++
					return fmt.Sprintf("test - %s", data)
				}
				dbMock.CallbackMapperFunc = nil
				dbMock.MapperFunc(mapperFunc)
				assert.Equal(t, 0, callCount, "CallbackMapperFunc should not be called")
			},
		},
		{
			name:   "Should run CallbackMustBegin",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackMustBegin = func() Tx {
					callCount++
					return &customTx{tx: nil}
				}
				tx := dbMock.MustBegin()
				assert.Nil(t, tx.Safe(), "MustBegin should return a customTx")
				assert.Equal(t, 1, callCount, "CallbackMustBegin should be called once")
			},
		},
		{
			name:   "Should run default CallbackMustBegin",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackMustBegin = nil
				tx := dbMock.MustBegin()
				assert.Empty(t, tx, "MustBegin should return an empty Tx")
			},
		},
		{
			name:   "Should run CallbackMustBeginTx",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackMustBeginTx = func(ctx context.Context, opts *sql.TxOptions) Tx {
					callCount++
					return &customTx{tx: nil}
				}
				tx := dbMock.MustBeginTx(context.Background(), nil)
				assert.Nil(t, tx.Safe(), "MustBeginTx should return a customTx")
				assert.Equal(t, 1, callCount, "CallbackMustBeginTx should be called once")
			},
		},
		{
			name:   "Should run default CallbackMustBeginTx",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackMustBeginTx = nil
				tx := dbMock.MustBeginTx(context.Background(), nil)
				assert.Empty(t, tx, "MustBeginTx should return an empty Tx")
			},
		},
		{
			name:   "Should run CallbackMustExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackMustExec = func(query string, args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}
				}
				result := dbMock.MustExec("")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackMustExec should be called once")
			},
		},
		{
			name:   "Should run default CallbackMustExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackMustExec = nil
				result := dbMock.MustExec("")
				assert.Empty(t, result, "MustExec should return an empty result")
			},
		},
		{
			name:   "Should run CallbackMustExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackMustExecContext = func(ctx context.Context, query string, args ...interface{}) sql.Result {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}
				}
				result := dbMock.MustExecContext(context.Background(), "")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackMustExecContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackMustExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackMustExecContext = nil
				result := dbMock.MustExecContext(context.Background(), "")
				assert.Empty(t, result, "MustExecContext should return an empty result")
			},
		},
		{
			name:   "Should run CallbackNamedExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackNamedExec = func(query string, arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}, nil
				}
				result, err := dbMock.NamedExec("", nil)
				assert.NoError(t, err, "NamedExec should not return an error")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackNamedExec should be called once")
			},
		},
		{
			name:   "Should run default CallbackNamedExec",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackNamedExec")
				dbMock.Error = testErr
				dbMock.CallbackNamedExec = nil
				result, err := dbMock.NamedExec("", nil)
				assert.Empty(t, result, "NamedExec should return an empty result")
				assert.ErrorIs(t, err, testErr, "NamedExec should return an error")
			},
		},
		{
			name:   "Should run CallbackNamedExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackNamedExecContext = func(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
					callCount++
					return &ResultMock{CallbackLastInsertId: func() (int64, error) { return 1024, nil }}, nil
				}
				result, err := dbMock.NamedExecContext(context.Background(), "", nil)
				assert.NoError(t, err, "NamedExecContext should not return an error")
				lastID, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(1024), lastID, "LastInsertId should return 1024")
				assert.Equal(t, 1, callCount, "CallbackNamedExecContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackNamedExecContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackNamedExecContext")
				dbMock.Error = testErr
				dbMock.CallbackNamedExecContext = nil
				result, err := dbMock.NamedExecContext(context.Background(), "", nil)
				assert.Empty(t, result, "NamedExecContext should return an empty result")
				assert.ErrorIs(t, err, testErr, "NamedExecContext should return an error")
			},
		},
		{
			name:   "Should run CallbackNamedQuery",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackNamedQuery = func(query string, arg interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := dbMock.NamedQuery("", nil)
				assert.NoError(t, err, "NamedQuery should not return an error")
				assert.IsType(t, &RowsMock{}, rows, "NamedQuery should return a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackNamedQuery should be called once")
			},
		},
		{
			name:   "Should run default CallbackNamedQuery",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackNamedQuery")
				dbMock.Error = testErr
				dbMock.CallbackNamedQuery = nil
				rows, err := dbMock.NamedQuery("", nil)
				assert.Empty(t, rows, "NamedQuery should return an empty Rows")
				assert.ErrorIs(t, err, testErr, "NamedQuery should return an error")
			},
		},
		{
			name:   "Should run CallbackNamedQueryContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackNamedQueryContext = func(ctx context.Context, query string, arg interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := dbMock.NamedQueryContext(context.Background(), "", nil)
				assert.NoError(t, err, "NamedQueryContext should not return an error")
				assert.IsType(t, &RowsMock{}, rows, "NamedQueryContext should return a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackNamedQueryContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackNamedQueryContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackNamedQueryContext")
				dbMock.Error = testErr
				dbMock.CallbackNamedQueryContext = nil
				rows, err := dbMock.NamedQueryContext(context.Background(), "", nil)
				assert.Empty(t, rows, "NamedQueryContext should return an empty Rows")
				assert.ErrorIs(t, err, testErr, "NamedQueryContext should return an error")
			},
		},
		{
			name:   "Should run CallbackPrepareNamed",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPrepareNamed = func(query string) (NamedStmt, error) {
					callCount++
					return &NamedStmtMock{}, nil
				}
				stmt, err := dbMock.PrepareNamed("")
				assert.NoError(t, err, "PrepareNamed should not return an error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "PrepareNamed should return a NamedStmtMock")
				assert.Equal(t, 1, callCount, "CallbackPrepareNamed should be called once")
			},
		},
		{
			name:   "Should run default CallbackPrepareNamed",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackPrepareNamed")
				dbMock.Error = testErr
				dbMock.CallbackPrepareNamed = nil
				stmt, err := dbMock.PrepareNamed("")
				assert.Empty(t, stmt, "PrepareNamed should return an empty NamedStmt")
				assert.ErrorIs(t, err, testErr, "PrepareNamed should return an error")
			},
		},
		{
			name:   "Should run CallbackPrepareNamedContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPrepareNamedContext = func(ctx context.Context, query string) (NamedStmt, error) {
					callCount++
					return &NamedStmtMock{}, nil
				}
				stmt, err := dbMock.PrepareNamedContext(context.Background(), "")
				assert.NoError(t, err, "PrepareNamedContext should not return an error")
				assert.IsType(t, &NamedStmtMock{}, stmt, "PrepareNamedContext should return a NamedStmtMock")
				assert.Equal(t, 1, callCount, "CallbackPrepareNamedContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackPrepareNamedContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackPrepareNamedContext")
				dbMock.Error = testErr
				dbMock.CallbackPrepareNamedContext = nil
				stmt, err := dbMock.PrepareNamedContext(context.Background(), "")
				assert.Empty(t, stmt, "PrepareNamedContext should return an empty NamedStmt")
				assert.ErrorIs(t, err, testErr, "PrepareNamedContext should return an error")
			},
		},
		{
			name:   "Should run CallbackPrepare",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPrepare = func(query string) (Stmt, error) {
					callCount++
					return &StmtMock{}, nil
				}
				stmt, err := dbMock.Prepare("")
				assert.NoError(t, err, "Prepare should not return an error")
				assert.IsType(t, &StmtMock{}, stmt, "Prepare should return a StmtMock")
				assert.Equal(t, 1, callCount, "CallbackPrepare should be called once")
			},
		},
		{
			name:   "Should run default CallbackPrepare",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackPrepare")
				dbMock.Error = testErr
				dbMock.CallbackPrepare = nil
				stmt, err := dbMock.Prepare("")
				assert.Empty(t, stmt, "Prepare should return an empty Stmt")
				assert.ErrorIs(t, err, testErr, "Prepare should return an error")
			},
		},
		{
			name:   "Should run CallbackPrepareContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackPrepareContext = func(ctx context.Context, query string) (Stmt, error) {
					callCount++
					return &StmtMock{}, nil
				}
				stmt, err := dbMock.PrepareContext(context.Background(), "")
				assert.NoError(t, err, "PrepareContext should not return an error")
				assert.IsType(t, &StmtMock{}, stmt, "PrepareContext should return a StmtMock")
				assert.Equal(t, 1, callCount, "CallbackPrepareContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackPrepareContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackPrepareContext")
				dbMock.Error = testErr
				dbMock.CallbackPrepareContext = nil
				stmt, err := dbMock.PrepareContext(context.Background(), "")
				assert.Empty(t, stmt, "PrepareContext should return an empty Stmt")
				assert.ErrorIs(t, err, testErr, "PrepareContext should return an error")
			},
		},
		{
			name:   "Should run CallbackQueryRow",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackQueryRow = func(query string, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := dbMock.QueryRow("")
				assert.IsType(t, &RowMock{}, row, "QueryRow should return a RowMock")
				assert.Equal(t, 1, callCount, "CallbackQueryRow should be called once")
			},
		},
		{
			name:   "Should run default CallbackQueryRow",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackQueryRow = nil
				row := dbMock.QueryRow("")
				assert.Empty(t, row, "QueryRow should return an empty Row")
			},
		},
		{
			name:   "Should run CallbackQueryRowContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackQueryRowContext = func(ctx context.Context, query string, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				row := dbMock.QueryRowContext(context.Background(), "")
				assert.IsType(t, &RowMock{}, row, "QueryRowContext should return a RowMock")
				assert.Equal(t, 1, callCount, "CallbackQueryRowContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackQueryRowContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackQueryRowContext = nil
				row := dbMock.QueryRowContext(context.Background(), "")
				assert.Empty(t, row, "QueryRowContext should return an empty Row")
			},
		},
		{
			name:   "Should run CallbackQuery",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackQuery = func(query string, args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := dbMock.Query("")
				assert.NoError(t, err, "Query should not return an error")
				assert.IsType(t, &RowsMock{}, rows, "Query should return a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackQuery should be called once")
			},
		},
		{
			name:   "Should run default CallbackQuery",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackQuery")
				dbMock.Error = testErr
				dbMock.CallbackQuery = nil
				rows, err := dbMock.Query("")
				assert.Empty(t, rows, "Query should return an empty Rows")
				assert.ErrorIs(t, err, testErr, "Query should return an error")
			},
		},
		{
			name:   "Should run CallbackQueryContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackQueryContext = func(ctx context.Context, query string, args ...interface{}) (Rows, error) {
					callCount++
					return &RowsMock{}, nil
				}
				rows, err := dbMock.QueryContext(context.Background(), "")
				assert.NoError(t, err, "QueryContext should not return an error")
				assert.IsType(t, &RowsMock{}, rows, "QueryContext should return a RowsMock")
				assert.Equal(t, 1, callCount, "CallbackQueryContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackQueryContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackQueryContext")
				dbMock.Error = testErr
				dbMock.CallbackQueryContext = nil
				rows, err := dbMock.QueryContext(context.Background(), "")
				assert.Empty(t, rows, "QueryContext should return an empty Rows")
				assert.ErrorIs(t, err, testErr, "QueryContext should return an error")
			},
		},
		{
			name:   "Should run CallbackRebind",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackRebind = func(query string) string {
					callCount++
					return "select 1"
				}
				assert.Equal(t, "select 1", dbMock.Rebind(""), "Rebind should return 'select 1'")
				assert.Equal(t, 1, callCount, "CallbackRebind should be called once")
			},
		},
		{
			name:   "Should run default CallbackRebind",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackRebind = nil
				assert.Equal(t, "", dbMock.Rebind(""), "Rebind should return an empty string")
			},
		},
		{
			name:   "Should run CallbackSelect",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				data := 0
				callCount := 0
				dbMock.CallbackSelect = func(dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = callCount + *(dest.(*int))
					return nil
				}
				err := dbMock.Select(&data, "")
				assert.NoError(t, err, "Select should not return an error")
				assert.Equal(t, 1, data, "Select should return 1")
				assert.Equal(t, 1, callCount, "CallbackSelect should be called once")
			},
		},
		{
			name:   "Should run default CallbackSelect",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackSelect")
				dbMock.Error = testErr
				dbMock.CallbackSelect = nil
				err := dbMock.Select(nil, "")
				assert.ErrorIs(t, err, testErr, "Select should return an error")
			},
		},
		{
			name:   "Should run CallbackSelectContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				data := 0
				callCount := 0
				dbMock.CallbackSelectContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					*dest.(*int) = callCount + *(dest.(*int))
					return nil
				}
				err := dbMock.SelectContext(context.Background(), &data, "")
				assert.NoError(t, err, "SelectContext should not return an error")
				assert.Equal(t, 1, data, "SelectContext should return 1")
				assert.Equal(t, 1, callCount, "CallbackSelectContext should be called once")
			},
		},
		{
			name:   "Should run default CallbackSelectContext",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				testErr := errors.New("failed to run CallbackSelectContext")
				dbMock.Error = testErr
				dbMock.CallbackSelectContext = nil
				err := dbMock.SelectContext(context.Background(), nil, "")
				assert.ErrorIs(t, err, testErr, "SelectContext should return an error")
			},
		},
		{
			name:   "Should run CallbackUnsafe",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackUnsafe = func() *sqlx.DB {
					callCount++
					return &sqlx.DB{}
				}
				assert.Empty(t, dbMock.Unsafe(), "Unsafe should return an empty *sqlx.DB")
				assert.Equal(t, 1, callCount, "CallbackUnsafe should be called once")
			},
		},
		{
			name:   "Should run default CallbackUnsafe",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackUnsafe = nil
				assert.Empty(t, dbMock.Unsafe(), "Unsafe should return an empty *sqlx.DB")
			},
		},
		{
			name:   "Should run CallbackSafe",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				callCount := 0
				dbMock.CallbackSafe = func() *sqlx.DB {
					callCount++
					return &sqlx.DB{}
				}
				assert.Empty(t, dbMock.Safe(), "Safe should return an empty *sqlx.DB")
				assert.Equal(t, 1, callCount, "CallbackSafe should be called once")
			},
		},
		{
			name:   "Should run default CallbackSafe",
			dbMock: NewMockDB(),
			assert: func(t *testing.T, dbMock *DBMock) {
				dbMock.CallbackSafe = nil
				assert.Empty(t, dbMock.Safe(), "Safe should return an empty *sqlx.DB")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.dbMock)
		})
	}
}
