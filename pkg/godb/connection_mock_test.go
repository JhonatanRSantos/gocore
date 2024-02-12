package godb

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConnMock(t *testing.T) {
	tests := []struct {
		name     string
		connMock ConnMock
		assert   func(t *testing.T, connMock ConnMock)
	}{
		{
			name:     "Should run CallbackClose",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackClose = func() error {
					callCount++
					return nil
				}
				assert.NoError(t, connMock.Close(), "Close should not return an error")
				assert.Equal(t, 1, callCount, "CallbackClose should be called once")
			},
		},
		{
			name:     "should fail to run CallbackClose",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackClose = nil
				testErr := errors.New("failed to run CallbackClose")
				connMock.Error = testErr
				assert.ErrorIs(t, connMock.Close(), testErr, "Close should return an error")
			},
		},
		{
			name:     "Should run CallbackExecContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackExecContext = func(ctx context.Context, query string, args ...any) (sql.Result, error) {
					callCount++
					return nil, nil
				}
				_, err := connMock.ExecContext(context.Background(), "")
				assert.NoError(t, err, "ExecContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackExecContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackExecContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackExecContext = nil
				testErr := errors.New("failed to run CallbackExecContext")
				connMock.Error = testErr
				_, err := connMock.ExecContext(context.Background(), "")
				assert.ErrorIs(t, err, testErr, "ExecContext should return an error")
			},
		},
		{
			name:     "Should run CallbackPingContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackPingContext = func(ctx context.Context) error {
					callCount++
					return nil
				}
				assert.NoError(t, connMock.PingContext(context.Background()), "PingContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackPingContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackPingContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackPingContext = nil
				testErr := errors.New("failed to run CallbackPingContext")
				connMock.Error = testErr
				assert.ErrorIs(t, connMock.PingContext(context.Background()), testErr, "PingContext should return an error")
			},
		},
		{
			name:     "Should run CallbackRaw",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackRaw = func(f func(driverConn any) error) error {
					callCount++
					return nil
				}
				assert.NoError(t, connMock.Raw(func(driverConn any) error {
					return nil
				}), "Raw should not return an error")
				assert.Equal(t, 1, callCount, "CallbackRaw should be called once")
			},
		},
		{
			name:     "should fail to run CallbackRaw",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackRaw = nil
				testErr := errors.New("failed to run CallbackRaw")
				connMock.Error = testErr
				assert.ErrorIs(t, connMock.Raw(func(driverConn any) error {
					return nil
				}), testErr, "Raw should return an error")
			},
		},
		{
			name:     "Should run CallbackBeginTx",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackBeginTx = func(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
					callCount++
					return nil, nil
				}
				_, err := connMock.BeginTx(context.Background(), nil)
				assert.NoError(t, err, "BeginTx should not return an error")
				assert.Equal(t, 1, callCount, "CallbackBeginTx should be called once")
			},
		},
		{
			name:     "should fail to run CallbackBeginTx",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackBeginTx = nil
				testErr := errors.New("failed to run CallbackBeginTx")
				connMock.Error = testErr
				_, err := connMock.BeginTx(context.Background(), nil)
				assert.ErrorIs(t, err, testErr, "BeginTx should return an error")
			},
		},
		{
			name:     "Should run CallbackGetContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackGetContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					return nil
				}
				assert.NoError(t, connMock.GetContext(context.Background(), nil, "", nil), "GetContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackGetContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackGetContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackGetContext = nil
				testErr := errors.New("failed to run CallbackGetContext")
				connMock.Error = testErr
				assert.ErrorIs(t, connMock.GetContext(context.Background(), nil, "", nil), testErr, "GetContext should return an error")
			},
		},
		{
			name:     "Should run CallbackPrepareContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackPrepareContext = func(ctx context.Context, query string) (Stmt, error) {
					callCount++
					return nil, nil
				}
				_, err := connMock.PrepareContext(context.Background(), "")
				assert.NoError(t, err, "PrepareContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackPrepareContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackPrepareContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackPrepareContext = nil
				testErr := errors.New("failed to run CallbackPrepareContext")
				connMock.Error = testErr
				_, err := connMock.PrepareContext(context.Background(), "")
				assert.ErrorIs(t, err, testErr, "PrepareContext should return an error")
			},
		},
		{
			name:     "Should run CallbackQueryRowContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackQueryRowContext = func(ctx context.Context, query string, args ...interface{}) Row {
					callCount++
					return &RowMock{}
				}
				assert.IsType(t, &RowMock{}, connMock.QueryRowContext(context.Background(), ""), "QueryRowContext should return a Row")
				assert.Equal(t, 1, callCount, "CallbackQueryRowContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackQueryRowContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackQueryRowContext = nil
				testErr := errors.New("failed to run CallbackQueryRowContext")
				connMock.Error = testErr
				assert.IsType(t, &RowMock{}, connMock.QueryRowContext(context.Background(), ""), "QueryRowContext should return a Row")
			},
		},
		{
			name:     "Should run CallbackQueryContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackQueryContext = func(ctx context.Context, query string, args ...interface{}) (Rows, error) {
					callCount++
					return nil, nil
				}
				_, err := connMock.QueryContext(context.Background(), "", nil)
				assert.NoError(t, err, "QueryContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackQueryContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackQueryContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackQueryContext = nil
				testErr := errors.New("failed to run CallbackQueryContext")
				connMock.Error = testErr
				_, err := connMock.QueryContext(context.Background(), "", nil)
				assert.ErrorIs(t, err, testErr, "QueryContext should return an error")
			},
		},
		{
			name:     "Should run CallbackRebind",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackRebind = func(query string) string {
					callCount++
					return strings.ReplaceAll(query, "?", "@")
				}
				assert.Equal(t, "SELECT * FROM @", connMock.Rebind("SELECT * FROM ?"), "Rebind should return a rebinded query")
				assert.Equal(t, 1, callCount, "CallbackRebind should be called once")
			},
		},
		{
			name:     "should fail to run CallbackRebind",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackRebind = nil
				testErr := errors.New("failed to run CallbackRebind")
				connMock.Error = testErr
				assert.Equal(t, "", connMock.Rebind(""), "Rebind should not return an empty string")
			},
		},
		{
			name:     "Should run CallbackSelectContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				callCount := 0
				connMock.CallbackSelectContext = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
					callCount++
					return nil
				}
				assert.NoError(t, connMock.SelectContext(context.Background(), nil, "", nil), "SelectContext should not return an error")
				assert.Equal(t, 1, callCount, "CallbackSelectContext should be called once")
			},
		},
		{
			name:     "should fail to run CallbackSelectContext",
			connMock: ConnMock{},
			assert: func(t *testing.T, connMock ConnMock) {
				connMock.CallbackSelectContext = nil
				testErr := errors.New("failed to run CallbackSelectContext")
				connMock.Error = testErr
				assert.ErrorIs(t, connMock.SelectContext(context.Background(), nil, "", nil), testErr, "SelectContext should return an error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.connMock)
		})
	}
}
