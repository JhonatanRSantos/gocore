package godb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Conn(t *testing.T) {
	var (
		db  DB
		err error
	)

	if db, err = NewDB(mysqlDefaultConfig); err != nil {
		assert.FailNow(t, "failed to create new db: %s", err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			assert.FailNow(t, err.Error())
		}
	}()

	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetConnMaxIdleTime(ConnMaxIdleTime)

	tests := []struct {
		name   string
		assert func(t *testing.T, conn customConn)
	}{
		{
			name: "Should run Close",
			assert: func(t *testing.T, conn customConn) {
				assert.NoError(t, conn.Close())
			},
		},
		{
			name: "Should run ExecContext",
			assert: func(t *testing.T, conn customConn) {
				_, err := conn.ExecContext(context.Background(), "SELECT 1")
				assert.NoError(t, err)
			},
		},
		{
			name: "Should run PingContext",
			assert: func(t *testing.T, conn customConn) {
				assert.NoError(t, conn.PingContext(context.Background()))
			},
		},
		{
			name: "Should run Raw",
			assert: func(t *testing.T, conn customConn) {
				assert.NoError(t, conn.Raw(func(driverConn any) error {
					return nil
				}))
			},
		},
		{
			name: "Should run BeginTx",
			assert: func(t *testing.T, conn customConn) {
				tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)
				assert.NoError(t, tx.Rollback(), "failed to rollback transaction")
			},
		},
		{
			name: "Should fail to run BeginTx",
			assert: func(t *testing.T, conn customConn) {
				testErr := errors.New("failed to begin transaction")
				conn.pushTestError(testErr)
				_, err := conn.BeginTx(context.Background(), nil)
				assert.Error(t, err)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "Should run GetContext",
			assert: func(t *testing.T, conn customConn) {
				var dest int
				assert.NoError(t, conn.GetContext(context.Background(), &dest, "SELECT 1"))
			},
		},
		{
			name: "Should run PrepareContext",
			assert: func(t *testing.T, conn customConn) {
				stmt, err := conn.PrepareContext(context.Background(), "SELECT 1")
				assert.NoError(t, err)
				stmt.Exec()
				assert.NoError(t, stmt.Close())
			},
		},
		{
			name: "Should fail to run PrepareContext",
			assert: func(t *testing.T, conn customConn) {
				testErr := errors.New("failed to prepare statement")
				conn.pushTestError(testErr)
				_, err := conn.PrepareContext(context.Background(), "SELECT 1")
				assert.Error(t, err)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "Should run QueryRowContext",
			assert: func(t *testing.T, conn customConn) {
				var total int
				row := conn.QueryRowContext(context.Background(), "SELECT 1")
				assert.NotNil(t, row)
				assert.NoError(t, row.Scan(&total))
			},
		},
		{
			name: "Should run QueryContext",
			assert: func(t *testing.T, conn customConn) {
				rows, err := conn.QueryContext(context.Background(), "SELECT 1")
				assert.NoError(t, err)
				assert.NoError(t, rows.Close())
			},
		},
		{
			name: "Should fail to run QueryContext",
			assert: func(t *testing.T, conn customConn) {
				testErr := errors.New("failed to query")
				conn.pushTestError(testErr)
				_, err := conn.QueryContext(context.Background(), "SELECT 1")
				assert.Error(t, err)
				assert.ErrorIs(t, err, testErr)
			},
		},
		{
			name: "Should run Rebind",
			assert: func(t *testing.T, conn customConn) {
				assert.Equal(t, "SELECT 1", conn.Rebind("SELECT 1"))
			},
		},
		{
			name: "Should run SelectContext",
			assert: func(t *testing.T, conn customConn) {
				var dest []int
				assert.NoError(t, conn.SelectContext(context.Background(), &dest, "SELECT 1"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, err := db.Safe().Connx(context.Background())
			if err != nil {
				assert.FailNow(t, "%s failed to get connection: %s", tt.name, err.Error())
			}
			tt.assert(t, customConn{conn: connection})
			if err := connection.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
				assert.FailNow(t, "%s failed to close connection: %s", tt.name, err.Error())
			}
		})
	}
}
