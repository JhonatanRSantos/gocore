package godb

import (
	"context"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_NamedStmt(t *testing.T) {
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
			name: "Should run Close",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				_, _ = customNamedStmt.Exec(map[string]interface{}{"name": "John Doe"})
				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run Exec",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("INSERT INTO custom_table (name, age) VALUES (:name, :age)")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				result, err := customNamedStmt.Exec(map[string]interface{}{"name": "Winston", "age": 30})
				assert.NoError(t, err, "failed to execute named statement")
				assert.NotNil(t, result, "failed to execute named statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, totalAffectedRows, int64(1), "should have 1 affected row")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run ExecContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("INSERT INTO custom_table (name, age) VALUES (:name, :age)")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				result, err := customNamedStmt.ExecContext(context.Background(), map[string]interface{}{"name": "Lucas", "age": 30})
				assert.NoError(t, err, "failed to execute named statement")
				assert.NotNil(t, result, "failed to execute named statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, totalAffectedRows, int64(1), "should have 1 affected row")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run Get",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				var data customData
				err = customNamedStmt.Get(&data, map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to get named statement")
				assert.Equal(t, data.Name, "John Wick", "should have the same name")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run GetContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				var data customData
				err = customNamedStmt.GetContext(context.Background(), &data, map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to get named statement")
				assert.Equal(t, data.Name, "John Wick", "should have the same name")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run MustExec",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("INSERT INTO custom_table (name, age) VALUES (:name, :age)")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				result := customNamedStmt.MustExec(map[string]interface{}{"name": "Charon", "age": 30})
				assert.NotNil(t, result, "failed to execute named statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, totalAffectedRows, int64(1), "should have 1 affected row")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run MustExecContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("INSERT INTO custom_table (name, age) VALUES (:name, :age)")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				result := customNamedStmt.MustExecContext(context.Background(), map[string]interface{}{"name": "Ares", "age": 30})
				assert.NotNil(t, result, "failed to execute named statement")

				totalAffectedRows, err := result.RowsAffected()
				assert.NoError(t, err, "failed to get affected rows")
				assert.Equal(t, totalAffectedRows, int64(1), "should have 1 affected row")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run QueryRow",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				row := customNamedStmt.QueryRow(map[string]interface{}{"name": "John Wick"})
				assert.NotNil(t, row, "failed to query named statement")

				var data customData
				assert.NoError(t, row.Scan(&data.ID, &data.Name, &data.Age), "failed to scan row")
				assert.Equal(t, data.Name, "John Wick", "should have the same name")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run QueryRowContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				row := customNamedStmt.QueryRowContext(context.Background(), map[string]interface{}{"name": "John Wick"})
				assert.NotNil(t, row, "failed to query named statement")

				var data customData
				assert.NoError(t, row.Scan(&data.ID, &data.Name, &data.Age), "failed to scan row")
				assert.Equal(t, data.Name, "John Wick", "should have the same name")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run Query",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				rows, err := customNamedStmt.Query(map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to query named statement")
				assert.NoError(t, rows.Close(), "failed to close rows")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "should fail to run Query",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customError := errors.New("failed to run query context")
				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				customNamedStmt.pushTestError(customError)

				_, err = customNamedStmt.Query(map[string]interface{}{"name": "John Doe"})
				assert.ErrorIs(t, err, customError, "should have custom error")
				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run QueryContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				rows, err := customNamedStmt.QueryContext(context.Background(), map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to query named statement")
				assert.NoError(t, rows.Close(), "failed to close rows")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "should fail to run QueryContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customError := errors.New("failed to run query context")
				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				customNamedStmt.pushTestError(customError)

				_, err = customNamedStmt.QueryContext(context.Background(), map[string]interface{}{"name": "John Doe"})
				assert.ErrorIs(t, err, customError, "should have custom error")
				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run Select",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				var data []customData
				err = customNamedStmt.Select(&data, map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to select named statement")
				assert.NotEmpty(t, data, "should have data")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run SelectContext",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")

				if err != nil {
					assert.FailNow(t, err.Error(), "failed to prepare named statement")
				}

				customNamedStmt := customNamedStmt{namedStmt: namedStmt.Safe()}
				var data []customData
				err = customNamedStmt.SelectContext(context.Background(), &data, map[string]interface{}{"name": "John Wick"})
				assert.NoError(t, err, "failed to select named statement")
				assert.NotEmpty(t, data, "should have data")

				assert.NoError(t, customNamedStmt.Close(), "failed to close named statement")
			},
		},
		{
			name: "Should run Unsafe",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")
				assert.NoError(t, err, "failed to prepare named statement")
				assert.IsType(t, &sqlx.NamedStmt{}, namedStmt.Unsafe(), "Unsafe should return *sqlx.NamedStmt")
			},
		},
		{
			name: "Should run Safe",
			assert: func(t *testing.T, db DB) {
				namedStmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")
				assert.NoError(t, err, "failed to prepare named statement")
				assert.IsType(t, &sqlx.NamedStmt{}, namedStmt.Safe(), "Safe should return *sqlx.NamedStmt")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}
