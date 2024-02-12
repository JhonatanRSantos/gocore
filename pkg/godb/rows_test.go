package godb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Rows(t *testing.T) {
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
			name: "should run ColumnTypes",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}
				columns, err := rows.ColumnTypes()
				assert.NoError(t, err, "column types failed")
				assert.Len(t, columns, 3, "column types length mismatch")
				assert.NoError(t, rows.Close(), "close failed")
			},
		},
		{
			name: "should run Columns",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}
				columns, err := rows.Columns()
				assert.NoError(t, err, "columns failed")
				assert.Len(t, columns, 3, "columns length mismatch")
				assert.NoError(t, rows.Close(), "close failed")
			},
		},
		{
			name: "should run Err",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}
				assert.NoError(t, rows.Close(), "close failed")
				assert.NoError(t, rows.Err(), "err failed")
			},
		},
		{
			name: "Should run NextResultSet",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}
				assert.False(t, rows.NextResultSet(), "next result set failed")
				assert.NoError(t, rows.Close(), "close failed")
			},
		},
		{
			name: "should run Scan",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}

				var (
					id   int
					age  int
					name string
				)

				assert.True(t, rows.Next(), "next failed")
				assert.NoError(t, rows.Scan(&id, &name, &age), "scan failed")
				assert.NoError(t, rows.Close(), "close failed")
			},
		},
		{
			name: "should run MapScan",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}

				data := make(map[string]interface{})
				assert.True(t, rows.Next(), "next failed")
				assert.NoError(t, rows.MapScan(data), "map scan failed")
				assert.NoError(t, rows.Close(), "close failed")
				assert.Len(t, data, 3, "map scan length mismatch")
			},
		},
		{
			name: "should run SliceScan",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}

				assert.True(t, rows.Next(), "next failed")
				data, err := rows.SliceScan()
				assert.NoError(t, err, "slice scan failed")
				assert.Len(t, data, 3, "slice scan length mismatch")
				assert.NoError(t, rows.Close(), "close failed")
			},
		},
		{
			name: "should run StructScan",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{"name": "John Wick"})
				if err != nil {
					assert.FailNow(t, err.Error(), "named query failed")
				}

				var data customData
				assert.True(t, rows.Next(), "next failed")
				assert.NoError(t, rows.StructScan(&data), "struct scan failed")
				assert.NoError(t, rows.Close(), "close failed")
				assert.Equal(t, "John Wick", data.Name, "struct scan name mismatch")
				assert.Equal(t, 30, data.Age, "struct scan age mismatch")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}
