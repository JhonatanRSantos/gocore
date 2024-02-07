package godb

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_Stmt(t *testing.T) {
	type customData struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
		Age  int    `db:"age"`
	}

	var (
		db  DB
		err error

		preDDL = `
			CREATE TABLE IF NOT EXISTS custom_table (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				age INT NOT NULL
			);
		`

		postDDL = `
			DROP TABLE IF EXISTS custom_table;
		`

		dml = `
			INSERT INTO custom_table (name, age) VALUES ('John Wick', 30);
		`
	)

	if db, err = NewDB(mysqlDefaultConfig); err != nil {
		assert.FailNow(t, err.Error())
	}

	defer func() {
		_, _ = db.Exec(postDDL)
		if err := db.Close(); err != nil {
			assert.FailNow(t, err.Error())
		}
	}()

	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetConnMaxIdleTime(ConnMaxIdleTime)

	if _, err = db.Exec(preDDL); err != nil {
		assert.FailNow(t, err.Error())
	}

	if _, err = db.Exec(dml); err != nil {
		assert.FailNow(t, err.Error())
	}

	tests := []struct {
		name   string
		assert func(t *testing.T, db DB)
	}{
		{
			name: "should run Close",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run Exec",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("INSERT INTO custom_table (name, age) VALUES ('Winston', 30)")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				result, err := customStmt.Exec()
				assert.NoError(t, err, "failed to execute statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run ExecContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("INSERT INTO custom_table (name, age) VALUES ('Lucas', 30)")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				result, err := customStmt.ExecContext(context.Background())
				assert.NoError(t, err, "failed to execute statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run Get",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				customData := customData{}
				err = customStmt.Get(&customData, "John Wick")
				assert.NoError(t, err, "failed to get data")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.Equal(t, 1, customData.ID, "id should be 1")
				assert.Equal(t, "John Wick", customData.Name, "name should be John Wick")
			},
		},
		{
			name: "should run GetContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				customData := customData{}
				err = customStmt.GetContext(context.Background(), &customData, "John Wick")
				assert.NoError(t, err, "failed to get data")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.Equal(t, 1, customData.ID, "id should be 1")
				assert.Equal(t, "John Wick", customData.Name, "name should be John Wick")
			},
		},
		{
			name: "should run MustExec",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("INSERT INTO custom_table (name, age) VALUES ('Charon', 30)")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				result := customStmt.MustExec()
				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run MustExecContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("INSERT INTO custom_table (name, age) VALUES ('Sofia', 30)")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				result := customStmt.MustExecContext(context.Background())
				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run QueryRow",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				row := customStmt.QueryRow("John Wick")
				assert.NotNil(t, row, "row should not be nil")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run QueryRowContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				row := customStmt.QueryRowContext(context.Background(), "John Wick")
				assert.NotNil(t, row, "row should not be nil")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run Query",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				rows, err := customStmt.Query()
				assert.NoError(t, err, "failed to execute statement")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.NotNil(t, rows, "rows should not be nil")
			},
		},
		{
			name: "should fail to run Query",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				testErr := errors.New("failed to execute statement")
				customStmt.pushTestError(testErr)

				_, err = customStmt.Query("John Wick")
				assert.ErrorIs(t, err, testErr, "error should be test error")
			},
		},
		{
			name: "should run QueryContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				rows, err := customStmt.QueryContext(context.Background())
				assert.NoError(t, err, "failed to execute statement")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.NotNil(t, rows, "rows should not be nil")
			},
		},
		{
			name: "should fail to run QueryContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				testErr := errors.New("failed to execute statement")
				customStmt.pushTestError(testErr)

				_, err = customStmt.QueryContext(context.Background(), "John Wick")
				assert.ErrorIs(t, err, testErr, "error should be test error")
			},
		},
		{
			name: "should run Select",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table LIMIT 2")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				var customData []customData
				err = customStmt.Select(&customData)
				assert.NoError(t, err, "failed to get data")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.Len(t, customData, 2, "data length should be 2")
			},
		},
		{
			name: "should run SelectContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table LIMIT 3")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				var customData []customData
				err = customStmt.SelectContext(context.Background(), &customData)
				assert.NoError(t, err, "failed to get data")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
				assert.Len(t, customData, 3, "data length should be 3")
			},
		},
		{
			name: "should run Unsafe",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				unsafeStmt := customStmt.Unsafe()
				assert.IsType(t, &sqlx.Stmt{}, unsafeStmt, "unsafe statement should be sqlx statement")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
		{
			name: "should run Safe",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Safe().Preparex("SELECT * FROM custom_table")
				assert.NoError(t, err, "failed to prepare statement")

				customStmt := &customStmt{stmt: stmt}
				safeStmt := customStmt.Safe()
				assert.IsType(t, &sqlx.Stmt{}, safeStmt, "safe statement should be sqlx statement")
				assert.NoError(t, customStmt.Close(), "failed to close statement")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}
