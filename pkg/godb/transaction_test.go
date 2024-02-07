package godb

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_Tx(t *testing.T) {
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
			name: "should run Commit",
			assert: func(t *testing.T, db DB) {
				var (
					tx                Tx
					err               error
					result            sql.Result
					totalAffectedRows int64
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if result, err = tx.Exec("INSERT INTO custom_table (name, age) VALUES ('John Wick', 30);"); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute query")
				}

				if totalAffectedRows, err = result.RowsAffected(); err != nil {
					assert.FailNow(t, err.Error(), "failed to get affected rows")
				}

				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run ExecContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx                Tx
					err               error
					result            sql.Result
					totalAffectedRows int64
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if result, err = tx.ExecContext(context.Background(), "INSERT INTO custom_table (name, age) VALUES ('John Wick', 30);"); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute query")
				}

				if totalAffectedRows, err = result.RowsAffected(); err != nil {
					assert.FailNow(t, err.Error(), "failed to get affected rows")
				}

				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "should run Rollback",
			assert: func(t *testing.T, db DB) {
				var (
					tx                Tx
					err               error
					result            sql.Result
					totalAffectedRows int64
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if result, err = tx.Exec("INSERT INTO custom_table (name, age) VALUES ('Winston', 30);"); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute query")
				}

				if totalAffectedRows, err = result.RowsAffected(); err != nil {
					assert.FailNow(t, err.Error(), "failed to get affected rows")
				}

				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")

				if err := tx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to rollback transaction")
				}
			},
		},
		{
			name: "Should run BindNamed",
			assert: func(t *testing.T, db DB) {
				var (
					tx     Tx
					err    error
					query  = "SELECT * FROM custom_table WHERE name = :name"
					bquery string
					params []interface{}
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if bquery, params, err = tx.BindNamed(query, map[string]interface{}{"name": "John Wick"}); err != nil {
					assert.FailNow(t, err.Error(), "failed to bind named query")
				}

				assert.Equal(t, "SELECT * FROM custom_table WHERE name = ?", bquery, "unexpected binded query")
				assert.Equal(t, 1, len(params), "unexpected params length")
				assert.Equal(t, "John Wick", params[0], "unexpected param value")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run DriverName",
			assert: func(t *testing.T, db DB) {
				var (
					tx  Tx
					err error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.Equal(t, "mysql", tx.DriverName(), "unexpected driver name")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run Get",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					data customData
					err  error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.Get(&data, "SELECT * FROM custom_table LIMIT 1"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run GetContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					data customData
					err  error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.GetContext(context.Background(), &data, "SELECT * FROM custom_table LIMIT 1"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run MustExec",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "INSERT INTO custom_table (name, age) VALUES ('Ariel', 30);"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NotPanics(t, func() {
					tx.MustExec(query)
				}, "must exec should not panic")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run MustExecContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "INSERT INTO custom_table (name, age) VALUES ('King', 30);"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NotPanics(t, func() {
					tx.MustExecContext(context.Background(), query)
				}, "must exec should not panic")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run NamedExec",
			assert: func(t *testing.T, db DB) {
				var (
					tx                Tx
					err               error
					query             = "INSERT INTO custom_table (name, age) VALUES (:name, :age);"
					result            sql.Result
					totalAffectedRows int64
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if result, err = tx.NamedExec(query, customData{Name: "Ariel", Age: 30}); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute named query")
				}

				if totalAffectedRows, err = result.RowsAffected(); err != nil {
					assert.FailNow(t, err.Error(), "failed to get affected rows")
				}

				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run NamedExecContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx                Tx
					err               error
					query             = "INSERT INTO custom_table (name, age) VALUES (:name, :age);"
					result            sql.Result
					totalAffectedRows int64
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if result, err = tx.NamedExecContext(context.Background(), query, customData{Name: "Saitama", Age: 30}); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute named query")
				}

				if totalAffectedRows, err = result.RowsAffected(); err != nil {
					assert.FailNow(t, err.Error(), "failed to get affected rows")
				}

				assert.Equal(t, int64(1), totalAffectedRows, "affected rows should be 1")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Shoduld run NamedQuery",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					rows  Rows
					query = "SELECT * FROM custom_table WHERE name = :name"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if rows, err = tx.NamedQuery(query, map[string]interface{}{"name": "John Wick"}); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute named query")
				}

				assert.True(t, rows.Next(), "no rows returned")
				assert.NoError(t, rows.StructScan(&data), "failed to scan data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, rows.Close(), "failed to close rows")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Shoudl failt to run NamedQuery",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "SELECT * FROM custom_table WHERE name = :name"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to execute named query")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.NamedQuery(query, map[string]interface{}{"name": "John Wick"})
				assert.ErrorIs(t, err, testErr, "should fail to execute named query")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run NamedStmt",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  NamedStmt
					query = "SELECT * FROM custom_table WHERE name = :name"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = db.PrepareNamedContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				namedStmtTx := tx.NamedStmt(&customNamedStmt{namedStmt: stmt.Safe()})

				assert.NoError(t, namedStmtTx.Get(&data, map[string]interface{}{"name": "John Wick"}), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, namedStmtTx.Close(), "failed to close named statement tx")
				assert.NoError(t, stmt.Close(), "failed to close statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run NamedStmtContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  NamedStmt
					query = "SELECT * FROM custom_table WHERE name = :name"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = db.PrepareNamedContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				namedStmtTx := tx.NamedStmtContext(context.Background(), &customNamedStmt{namedStmt: stmt.Safe()})

				assert.NoError(t, namedStmtTx.Get(&data, map[string]interface{}{"name": "John Wick"}), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, namedStmtTx.Close(), "failed to close named statement tx")
				assert.NoError(t, stmt.Close(), "failed to close named statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run PrepareNamed",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  NamedStmt
					query = "SELECT * FROM custom_table WHERE name = :name"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = tx.PrepareNamed(query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				assert.NoError(t, stmt.Get(&data, map[string]interface{}{"name": "John Wick"}), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmt.Close(), "failed to close named statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run PrepareNamed",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "SELECT * FROM custom_table WHERE name = :name"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to prepare named statement")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.PrepareNamed(query)
				assert.ErrorIs(t, err, testErr, "should fail to prepare named statement")

				if err := tx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run PrepareNamedContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  NamedStmt
					query = "SELECT * FROM custom_table WHERE name = :name"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = tx.PrepareNamedContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				assert.NoError(t, stmt.Get(&data, map[string]interface{}{"name": "John Wick"}), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmt.Close(), "failed to close named statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run PrepareNamedContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "SELECT * FROM custom_table WHERE name = :name"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to prepare named statement")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.PrepareNamedContext(context.Background(), query)
				assert.ErrorIs(t, err, testErr, "should fail to prepare named statement")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run Prepare",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  Stmt
					query = "SELECT * FROM custom_table WHERE name = ?"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = tx.Prepare(query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare statement")
				}

				assert.NoError(t, stmt.Get(&data, "John Wick"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmt.Close(), "failed to close statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run Prepare",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "SELECT * FROM custom_table WHERE name = ?"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to prepare statement")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.Prepare(query)
				assert.ErrorIs(t, err, testErr, "should fail to prepare statement")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run PrepareContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					stmt  Stmt
					query = "SELECT * FROM custom_table WHERE name = ?"
					data  customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = tx.PrepareContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare statement")
				}

				assert.NoError(t, stmt.Get(&data, "John Wick"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmt.Close(), "failed to close statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run PrepareContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					query = "SELECT * FROM custom_table WHERE name = ?"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to prepare statement")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.PrepareContext(context.Background(), query)
				assert.ErrorIs(t, err, testErr, "should fail to prepare statement")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run QueryRow",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					err  error
					data customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.QueryRow("SELECT * FROM custom_table WHERE name = ?", "John Wick").Scan(&data.ID, &data.Name, &data.Age), "failed to scan data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run QueryRowContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					err  error
					data customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.QueryRowContext(context.Background(), "SELECT * FROM custom_table WHERE name = ?", "John Wick").Scan(&data.ID, &data.Name, &data.Age), "failed to scan data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run Query",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					err  error
					rows Rows
					data customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if rows, err = tx.Query("SELECT * FROM custom_table WHERE name = ?", "John Wick"); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute query")
				}

				assert.True(t, rows.Next(), "no rows returned")
				assert.NoError(t, rows.StructScan(&data), "failed to scan data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, rows.Close(), "failed to close rows")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run Query",
			assert: func(t *testing.T, db DB) {
				var (
					tx  Tx
					err error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to execute query")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.Query("SELECT * FROM custom_table WHERE name = ?", "John Wick")
				assert.ErrorIs(t, err, testErr, "should fail to execute query")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run QueryContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx   Tx
					err  error
					rows Rows
					data customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if rows, err = tx.QueryContext(context.Background(), "SELECT * FROM custom_table WHERE name = ?", "John Wick"); err != nil {
					assert.FailNow(t, err.Error(), "failed to execute query")
				}

				assert.True(t, rows.Next(), "no rows returned")
				assert.NoError(t, rows.StructScan(&data), "failed to scan data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, rows.Close(), "failed to close rows")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should fail to run QueryContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx  Tx
					err error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				testErr := errors.New("failed to execute query")
				cTx := customTx{tx: tx.Safe()}
				cTx.pushTestError(testErr)

				_, err = cTx.QueryContext(context.Background(), "SELECT * FROM custom_table WHERE name = ?", "John Wick")
				assert.ErrorIs(t, err, testErr, "should fail to execute query")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "should run Rebind",
			assert: func(t *testing.T, db DB) {
				var (
					tx     Tx
					err    error
					query  = "SELECT * FROM custom_table WHERE name = ?"
					bquery string
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				cTx := customTx{tx: tx.Safe()}

				bquery = cTx.Rebind(query)
				assert.Equal(t, "SELECT * FROM custom_table WHERE name = ?", bquery, "unexpected rebinded query")

				if err := cTx.Rollback(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run Select",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					data  []customData
					query = "SELECT * FROM custom_table WHERE name = ?"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.Select(&data, query, "John Wick"), "failed to select data")
				assert.Equal(t, 1, data[0].ID, "unexpected data id")
				assert.Equal(t, "John Wick", data[0].Name, "unexpected data name")
				assert.Equal(t, 30, data[0].Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run SelectContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx    Tx
					err   error
					data  []customData
					query = "SELECT * FROM custom_table WHERE name = ?"
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.NoError(t, tx.SelectContext(context.Background(), &data, query, "John Wick"), "failed to select data")
				assert.Equal(t, 1, data[0].ID, "unexpected data id")
				assert.Equal(t, "John Wick", data[0].Name, "unexpected data name")
				assert.Equal(t, 30, data[0].Age, "unexpected data age")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run Stmt",
			assert: func(t *testing.T, db DB) {
				var (
					tx     Tx
					err    error
					stmt   Stmt
					stmtTx Stmt
					query  = "SELECT * FROM custom_table WHERE name = ?"
					data   customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = db.PrepareContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare statement")
				}
				stmtTx = tx.Stmt(stmt.Safe())

				assert.NoError(t, stmtTx.Get(&data, "John Wick"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmtTx.Close(), "failed to close statement tx")
				assert.NoError(t, stmt.Close(), "failed to close statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "Should run StmtContext",
			assert: func(t *testing.T, db DB) {
				var (
					tx     Tx
					err    error
					stmt   Stmt
					stmtTx Stmt
					query  = "SELECT * FROM custom_table WHERE name = ?"
					data   customData
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				if stmt, err = db.PrepareContext(context.Background(), query); err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare statement")
				}
				stmtTx = tx.StmtContext(context.Background(), stmt.Safe())

				assert.NoError(t, stmtTx.Get(&data, "John Wick"), "failed to get data")
				assert.Equal(t, 1, data.ID, "unexpected data id")
				assert.Equal(t, "John Wick", data.Name, "unexpected data name")
				assert.Equal(t, 30, data.Age, "unexpected data age")
				assert.NoError(t, stmtTx.Close(), "failed to close statement tx")
				assert.NoError(t, stmt.Close(), "failed to close statement")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
		{
			name: "should run Unsafe",
			assert: func(t *testing.T, db DB) {
				var (
					tx  Tx
					err error
				)

				if tx, err = db.Begin(); err != nil {
					assert.FailNow(t, err.Error(), "failed to begin transaction")
				}

				assert.IsType(t, &sqlx.Tx{}, tx.Unsafe(), "unexpected unsafe type")

				if err := tx.Commit(); err != nil {
					assert.FailNow(t, err.Error(), "failed to commit transaction")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}
