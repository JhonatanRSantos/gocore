package godb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Row(t *testing.T) {
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
				row := db.QueryRow("SELECT * FROM custom_table")
				columnTypes, err := row.ColumnTypes()
				assert.NoError(t, err)
				assert.Len(t, columnTypes, 3, "column types should have 3 elements")
			},
		},
		{
			name: "should run Columns",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				columns, err := row.Columns()
				assert.NoError(t, err)
				assert.Len(t, columns, 3, "columns should have 3 elements")
			},
		},
		{
			name: "should run Err",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				err := row.Err()
				assert.NoError(t, err)
			},
		},
		{
			name: "should run Scan",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				data := customData{}
				err := row.Scan(&data.ID, &data.Name, &data.Age)
				assert.NoError(t, err)
				assert.Equal(t, 1, data.ID, "id should be 1")
				assert.Equal(t, "John Wick", data.Name, "name should be John Wick")
				assert.Equal(t, 30, data.Age, "age should be 30")
			},
		},
		{
			name: "should run MapScan",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				data := make(map[string]interface{})
				err := row.MapScan(data)
				assert.NoError(t, err)
				assert.Len(t, data, 3, "map should have 3 elements")
			},
		},
		{
			name: "should run SliceScan",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				data, err := row.SliceScan()
				assert.NoError(t, err)
				assert.Len(t, data, 3, "slice should have 3 elements")
			},
		},
		{
			name: "should run StructScan",
			assert: func(t *testing.T, db DB) {
				row := db.QueryRow("SELECT * FROM custom_table")
				data := customData{}
				err := row.StructScan(&data)
				assert.NoError(t, err)
				assert.Equal(t, 1, data.ID, "id should be 1")
				assert.Equal(t, "John Wick", data.Name, "name should be John Wick")
				assert.Equal(t, 30, data.Age, "age should be 30")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}
