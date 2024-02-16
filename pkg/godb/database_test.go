package godb

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const (
	MaxIdleConns    = 10
	MaxOpenConns    = 10
	ConnMaxLifetime = time.Hour * 1
	ConnMaxIdleTime = time.Hour * 1
)

var (
	mysqlDefaultConfig = DBConfig{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     MySQLDB,
		ConnectionParams: MySQLDefaultParams,
	}
	postgresDefaultConfig = DBConfig{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     PostgresDB,
		ConnectionParams: PostgresDefaultParams,
	}
)

// nolint:dupl
func Test_MySQLDB(t *testing.T) {
	var (
		db           DB
		err          error
		result       sql.Result
		rowsAffected int64

		preDDL = `
			CREATE TABLE IF NOT EXISTS users (
				user_id 			VARCHAR(36) NOT NULL DEFAULT (UUID()),
				name 				VARCHAR(255) NOT NULL,
				country 			VARCHAR(255) NOT NULL,
				default_language 	VARCHAR(255) NOT NULL,
				email 				VARCHAR(255) NOT NULL,
				phones 				JSON NOT NULL,
				password 			VARCHAR(255) NOT NULL,
			
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP NULL DEFAULT NULL,
			
				PRIMARY KEY (user_id),
				UNIQUE KEY (email)
			);
		`

		postDDL = `
			DROP TABLE IF EXISTS users;
		`

		dml = `
			INSERT INTO users (
				name, country, default_language, phones, email, password
			) 
			VALUES 
			('John Wick', 'Belarus', 'English', '[]', 'john.wick@continental.com', 'qwerty');
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

	if _, err = db.Exec(preDDL); err != nil {
		assert.FailNow(t, err.Error())
	}

	if result, err = db.Exec(dml); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected != 1 {
		assert.FailNow(t, "invalid amout of affected rows. Expected 1, but got %d", rowsAffected)
	}
}

// nolint:dupl
func Test_PostgresDB(t *testing.T) {
	var (
		db           DB
		err          error
		result       sql.Result
		rowsAffected int64

		preDDL = `
			CREATE TABLE IF NOT EXISTS users (
				user_id 			uuid DEFAULT gen_random_uuid(),
				name 				VARCHAR (255) NOT NULL,
				country 			VARCHAR (255) NOT NULL,
				default_language 	VARCHAR (255) NOT NULL,
				email 				VARCHAR (255) NOT NULL,
				phones 				TEXT [],
				password 			VARCHAR(255) NOT NULL,
				
				created_at 			TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
				updated_at			TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
				deleted_at			TIMESTAMP DEFAULT NULL,
				
				CONSTRAINT users_pk 			PRIMARY KEY (user_id),
				CONSTRAINT users_unique_email 	UNIQUE (email)	
			);
			CREATE INDEX users_unique_email_idx ON users (email);
		`

		postDDL = `
			DROP TABLE IF EXISTS users;
		`

		dml = `
			INSERT INTO users (
				name, country, default_language, email, password
			) 
			VALUES 
			('John Wick', 'Belarus', 'English', 'john.wick@continental.com', 'qwerty');
		`
	)

	if db, err = NewDB(postgresDefaultConfig); err != nil {
		assert.FailNow(t, err.Error())
	}

	defer func() {
		_, _ = db.Exec(postDDL)
		if err := db.Close(); err != nil {
			assert.FailNow(t, err.Error())
		}
	}()

	if _, err = db.Exec(preDDL); err != nil {
		assert.FailNow(t, err.Error())
	}

	if result, err = db.Exec(dml); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected != 1 {
		assert.FailNow(t, "invalid amout of affected rows. Expected 1, but got %d", rowsAffected)
	}
}

func Test_SQLiteDB(t *testing.T) {
	var (
		db           DB
		err          error
		result       sql.Result
		rowsAffected int64

		preDDL = `
			CREATE TABLE users (
				user_id				VARCHAR(36) PRIMARY KEY,
				name              	VARCHAR(255) NOT NULL,
				country           	VARCHAR(255) NOT NULL,
				default_language  	VARCHAR(255) NOT NULL,
				email             	VARCHAR(255) NOT NULL UNIQUE,
				phones            	TEXT NOT NULL,
				password          	VARCHAR(255) NOT NULL
			);
		`

		postDDL = `
			DROP TABLE users;
		`

		dml = `
			INSERT INTO users (
				user_id, name, country, default_language, phones, email, password
			)
			VALUES
			('6ffbde59-f292-4062-b6d5-45dc7acfbd5b', 'John Wick', 'Belarus', 'English', '[]', 'john.wick@continental.com', 'qwerty');
		`
	)

	if db, err = NewDB(DBConfig{
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     SQLiteDB,
		ConnectionParams: SQLiteDefaultParams,
	}); err != nil {
		assert.FailNow(t, err.Error())
	}

	defer func() {
		_, _ = db.Exec(postDDL)
		if err := db.Close(); err != nil {
			assert.FailNow(t, err.Error())
		}
	}()

	if _, err = db.Exec(preDDL); err != nil {
		assert.FailNow(t, err.Error())
	}

	if result, err = db.Exec(dml); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected, err = result.RowsAffected(); err != nil {
		assert.FailNow(t, err.Error())
	}

	if rowsAffected != 1 {
		assert.FailNow(t, "invalid amout of affected rows. Expected 1, but got %d", rowsAffected)
	}
}

func Test_NewDB(t *testing.T) {
	_, err := NewDB(DBConfig{DatabaseType: 30})
	assert.Error(t, err, "NewDB should return an error")
	assert.ErrorIs(t, err, ErrInvalidDBType, "NewDB should return ErrInvalidDBType")

	_, err = NewDB(DBConfig{
		Host:         mysqlDefaultConfig.Host,
		Port:         "7070",
		DatabaseType: mysqlDefaultConfig.DatabaseType,
	})
	assert.Error(t, err, "NewDB should return an error")
	assert.IsType(t, &net.OpError{}, err, "NewDB should return *net.OpError")

	_, err = NewDB(DBConfig{
		Host:         mysqlDefaultConfig.Host,
		Port:         mysqlDefaultConfig.Port,
		User:         mysqlDefaultConfig.User,
		Password:     mysqlDefaultConfig.Password,
		Database:     mysqlDefaultConfig.Database,
		DatabaseType: mysqlDefaultConfig.DatabaseType,
		ConnectionParams: DBConnectionParams{
			"timeout": (time.Microsecond * 1).String(),
		},
	})
	assert.Error(t, err, "NewDB should return an error")
	assert.ErrorIs(t, err, ErrConnectionTimeoutExceeded, "NewDB should return ErrConnectionTimeoutExceeded")
}

func Test_DB(t *testing.T) {
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
			INSERT INTO custom_table (name, age) VALUES ('John Doe', 30);
		`
	)

	if db, err = NewDB(mysqlDefaultConfig); err != nil {
		assert.FailNow(t, "failed to create new db: %s", err.Error())
	}

	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetConnMaxIdleTime(ConnMaxIdleTime)

	defer func() {
		_, _ = db.Exec(postDDL)
		if err := db.Close(); err != nil {
			assert.FailNow(t, "failed to close db: %s", err.Error())
		}
	}()

	if _, err = db.Exec(preDDL); err != nil {
		assert.FailNow(t, "failed to execute preDDL: %s", err.Error())
	}

	if _, err = db.Exec(dml); err != nil {
		assert.FailNow(t, "failed to execute dml: %s", err.Error())
	}

	tests := []struct {
		name   string
		assert func(t *testing.T, db DB)
	}{
		{
			name: "Should get database driver",
			assert: func(t *testing.T, db DB) {
				driver := db.Driver()

				assert.NotNil(t, driver, "driver should not be nil")
				assert.IsType(t, &mysql.MySQLDriver{}, driver, "driver should be of type *mysql.MySQLDriver")
			},
		},
		{
			name: "should run ExecContext",
			assert: func(t *testing.T, db DB) {
				result, err := db.ExecContext(context.Background(), "SELECT * FROM custom_table")
				assert.NoError(t, err, "ExecContext should not return error")
				assert.NotNil(t, result, "ExecContext result should not be nil")
			},
		},
		{
			name: "Should run Ping",
			assert: func(t *testing.T, db DB) {
				err := db.Ping()
				assert.NoError(t, err, "Ping should not return an error")
			},
		},
		{
			name: "Should run PingContext",
			assert: func(t *testing.T, db DB) {
				err := db.PingContext(context.Background())
				assert.NoError(t, err, "PingContext should not return an error")
			},
		},
		{
			name: "Should run Stats",
			assert: func(t *testing.T, db DB) {
				stats := db.Stats()
				assert.NotNil(t, stats, "Stats should not be nil")
				assert.Equal(t, MaxOpenConns, stats.MaxOpenConnections, "MaxOpenConnections should be equal to %d", MaxOpenConns)
			},
		},
		{
			name: "Should run BeginTx",
			assert: func(t *testing.T, db DB) {
				tx, err := db.BeginTx(context.Background(), nil)
				assert.NoError(t, err, "BeginTx should not return an error")
				assert.NotNil(t, tx, "BeginTx should not return nil")
				assert.NoError(t, tx.Commit(), "Commit should not return an error")
			},
		},
		{
			name: "Should run BeginTx with options",
			assert: func(t *testing.T, db DB) {
				tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
					Isolation: sql.LevelDefault,
					ReadOnly:  false,
				})
				assert.NoError(t, err, "BeginTx should not return an error")
				assert.NotNil(t, tx, "BeginTx should not return nil")
				assert.NoError(t, tx.Commit(), "Commit should not return an error")
			},
		},
		{
			name: "Should fail to run BeginTx",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to begin transaction")
				db.pushTestError(testErr)
				tx, err := db.BeginTx(context.Background(), nil)
				assert.Error(t, err, "BeginTx should return an error")
				assert.Empty(t, tx, "BeginTx should return nil")
				assert.ErrorIs(t, err, testErr, "BeginTx should return testErr")
			},
		},
		{
			name: "Should run Begin",
			assert: func(t *testing.T, db DB) {
				tx, err := db.Begin()
				assert.NoError(t, err, "Begin should not return an error")
				assert.NotNil(t, tx, "Begin should not return nil")
				assert.NoError(t, tx.Commit(), "Commit should not return an error")
			},
		},
		{
			name: "Should fail to run Begin",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to begin transaction")
				db.pushTestError(testErr)
				tx, err := db.Begin()
				assert.Error(t, err, "Begin should return an error")
				assert.Empty(t, tx, "Begin should return nil")
				assert.ErrorIs(t, err, testErr, "Begin should return testErr")
			},
		},
		{
			name: "Should run BindNamed",
			assert: func(t *testing.T, db DB) {
				_, params, err := db.BindNamed("SELECT * FROM custom_table WHERE name = :name", map[string]interface{}{
					"name": "John Doe",
				})
				assert.NoError(t, err, "BindNamed should not return an error")
				assert.NotEmpty(t, params, "BindNamed should not return empty params")
				assert.Len(t, params, 1, "BindNamed should return 1 param")
			},
		},
		{
			name: "Should run Conn",
			assert: func(t *testing.T, db DB) {
				conn, err := db.Conn(context.Background())
				assert.NoError(t, err, "Conn should not return an error")
				assert.NotNil(t, conn, "Conn should not return nil")
				assert.NoError(t, conn.Close(), "Conn close should not return an error")
			},
		},
		{
			name: "should fail to run Conn",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to create connection")
				db.pushTestError(testErr)
				conn, err := db.Conn(context.Background())
				assert.Error(t, err, "Conn should return an error")
				assert.Empty(t, conn, "Conn should return nil")
				assert.ErrorIs(t, err, testErr, "Conn should return testErr")
			},
		},
		{
			name: "Should run DriverName",
			assert: func(t *testing.T, db DB) {
				driverName := db.DriverName()
				assert.NotEmpty(t, driverName, "DriverName should not be empty")
			},
		},
		{
			name: "Should run Get",
			assert: func(t *testing.T, db DB) {
				var data customData
				err := db.Get(&data, "SELECT * FROM custom_table limit 1")
				assert.NoError(t, err, "Get should not return an error")
				assert.NotEmpty(t, data, "Get should not return empty data")
				assert.Equal(t, "John Doe", data.Name, "Get should return John Doe")
			},
		},
		{
			name: "Should run GetContext",
			assert: func(t *testing.T, db DB) {
				var data customData
				err := db.GetContext(context.Background(), &data, "SELECT * FROM custom_table limit 1")
				assert.NoError(t, err, "GetContext should not return an error")
				assert.NotEmpty(t, data, "GetContext should not return empty data")
				assert.Equal(t, "John Doe", data.Name, "Get should return John Doe")
			},
		},
		{
			name: "Should run MapperFunc",
			assert: func(t *testing.T, db DB) {
				db.Safe().Mapper = nil
				assert.Nil(t, db.Safe().Mapper, "Mapper should be nil")
				db.MapperFunc(func(s string) string {
					return s
				})
				assert.NotNil(t, db.Safe().Mapper, "Mapper should not be nil")
			},
		},
		{
			name: "Should run MustBegin",
			assert: func(t *testing.T, db DB) {
				tx := db.MustBegin()
				assert.NotNil(t, tx, "MustBegin should not return nil")
				assert.NoError(t, tx.Commit(), "MustBegin commit should not return an error")
			},
		},
		{
			name: "Should run MustBeginTx",
			assert: func(t *testing.T, db DB) {
				tx := db.MustBeginTx(context.Background(), nil)
				assert.NotNil(t, tx, "MustBeginTx should not return nil")
				assert.NoError(t, tx.Commit(), "MustBeginTx commit should not return an error")
			},
		},
		{
			name: "Should run MustExec",
			assert: func(t *testing.T, db DB) {
				result := db.MustExec("INSERT INTO custom_table (name, age) VALUES ('Dracula', 30)")
				assert.NotNil(t, result, "MustExec should not return nil")

				id, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(2), id, "LastInsertId should return 2")

				totalAffected, err := result.RowsAffected()
				assert.NoError(t, err, "Affected should not return an error")
				assert.Equal(t, int64(1), totalAffected, "MustExec should return 1 affected row")
			},
		},
		{
			name: "Should run MustExecContext",
			assert: func(t *testing.T, db DB) {
				result := db.MustExecContext(context.Background(), "INSERT INTO custom_table (name, age) VALUES ('Werewolf', 30)")
				assert.NotNil(t, result, "MustExecContext should not return nil")

				id, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(3), id, "LastInsertId should return 3")

				totalAffected, err := result.RowsAffected()
				assert.NoError(t, err, "Affected should not return an error")
				assert.Equal(t, int64(1), totalAffected, "MustExecContext should return 1 affected row")
			},
		},
		{
			name: "Should run NamedExec",
			assert: func(t *testing.T, db DB) {
				result, err := db.NamedExec("INSERT INTO custom_table (name, age) VALUES (:name, :age)", customData{
					Name: "John Wick",
					Age:  30,
				})
				assert.NoError(t, err, "NamedExec should not return an error")
				assert.NotNil(t, result, "NamedExec should not return nil")

				id, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(4), id, "LastInsertId should return 4")

				totalAffected, err := result.RowsAffected()
				assert.NoError(t, err, "Affected should not return an error")
				assert.Equal(t, int64(1), totalAffected, "NamedExec should return 1 affected row")
			},
		},
		{
			name: "Should run NamedExecContext",
			assert: func(t *testing.T, db DB) {
				result, err := db.NamedExecContext(context.Background(), "INSERT INTO custom_table (name, age) VALUES (:name, :age)", customData{
					Name: "Winston",
					Age:  30,
				})
				assert.NoError(t, err, "NamedExecContext should not return an error")
				assert.NotNil(t, result, "NamedExecContext should not return nil")

				id, err := result.LastInsertId()
				assert.NoError(t, err, "LastInsertId should not return an error")
				assert.Equal(t, int64(5), id, "LastInsertId should return 5")

				totalAffected, err := result.RowsAffected()
				assert.NoError(t, err, "Affected should not return an error")
				assert.Equal(t, int64(1), totalAffected, "NamedExecContext should return 1 affected row")
			},
		},
		{
			name: "Should run NamedQuery",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", customData{
					Name: "John Doe",
				})
				assert.NoError(t, err, "NamedQuery should not return an error")
				assert.NotNil(t, rows, "NamedQuery should not return nil")

				assert.True(t, rows.Next(), "Next should return true")
				results, err := rows.SliceScan()
				assert.NoError(t, err, "SliceScan should not return an error")
				assert.NotEmpty(t, results, "SliceScan should not return empty results")
				assert.Len(t, results, 3, "SliceScan should return 3 results")
			},
		},
		{
			name: "should fail to run NamedQuery",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to run named query")
				db.pushTestError(testErr)
				rows, err := db.NamedQuery("SELECT * FROM custom_table WHERE name = :name", customData{
					Name: "John Doe",
				})
				assert.Error(t, err, "NamedQuery should return an error")
				assert.Empty(t, rows, "NamedQuery should return nil")
				assert.ErrorIs(t, err, testErr, "NamedQuery should return testErr")
			},
		},
		{
			name: "Should run NamedQueryContext",
			assert: func(t *testing.T, db DB) {
				rows, err := db.NamedQueryContext(context.Background(), "SELECT * FROM custom_table WHERE name = :name", customData{
					Name: "John Doe",
				})
				assert.NoError(t, err, "NamedQueryContext should not return an error")
				assert.NotNil(t, rows, "NamedQueryContext should not return nil")

				assert.True(t, rows.Next(), "Next should return true")
				results, err := rows.SliceScan()
				assert.NoError(t, err, "SliceScan should not return an error")
				assert.NotEmpty(t, results, "SliceScan should not return empty results")
				assert.Len(t, results, 3, "SliceScan should return 3 results")
			},
		},
		{
			name: "Should fail to run NamedQueryContext",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to run named query")
				db.pushTestError(testErr)
				rows, err := db.NamedQueryContext(context.Background(), "SELECT * FROM custom_table WHERE name = :name", customData{
					Name: "John Doe",
				})
				assert.Error(t, err, "NamedQueryContext should return an error")
				assert.Empty(t, rows, "NamedQueryContext should return nil")
				assert.ErrorIs(t, err, testErr, "NamedQueryContext should return testErr")
			},
		},
		{
			name: "Should run PrepareNamed",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")
				assert.NoError(t, err, "PrepareNamed should not return an error")
				assert.NotNil(t, stmt, "PrepareNamed should not return nil")

				var data customData
				err = stmt.Get(&data, map[string]interface{}{"name": "John Doe"})
				assert.NoError(t, err, "Get should not return an error")
				assert.NotEmpty(t, data, "Get should not return empty data")
			},
		},
		{
			name: "Should fail to run PrepareNamed",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to prepare named statement")
				db.pushTestError(testErr)
				stmt, err := db.PrepareNamed("SELECT * FROM custom_table WHERE name = :name")
				assert.Error(t, err, "PrepareNamed should return an error")
				assert.Empty(t, stmt, "PrepareNamed should return nil")
				assert.ErrorIs(t, err, testErr, "PrepareNamed should return testErr")
			},
		},
		{
			name: "Should run PrepareNamedContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.PrepareNamedContext(context.Background(), "SELECT * FROM custom_table WHERE name = :name")
				assert.NoError(t, err, "PrepareNamedContext should not return an error")
				assert.NotNil(t, stmt, "PrepareNamedContext should not return nil")

				var data customData
				err = stmt.GetContext(context.Background(), &data, map[string]interface{}{"name": "John Doe"})
				assert.NoError(t, err, "GetContext should not return an error")
				assert.NotEmpty(t, data, "GetContext should not return empty data")
			},
		},
		{
			name: "Should fail to run PrepareNamedContext",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to prepare named statement")
				db.pushTestError(testErr)
				stmt, err := db.PrepareNamedContext(context.Background(), "SELECT * FROM custom_table WHERE name = :name")
				assert.Error(t, err, "PrepareNamedContext should return an error")
				assert.Empty(t, stmt, "PrepareNamedContext should return nil")
				assert.ErrorIs(t, err, testErr, "PrepareNamedContext should return testErr")
			},
		},
		{
			name: "Should run Prepare",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.Prepare("SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "Prepare should not return an error")
				assert.NotNil(t, stmt, "Prepare should not return nil")

				var data customData
				err = stmt.Get(&data, "John Doe")
				assert.NoError(t, err, "Get should not return an error")
				assert.NotEmpty(t, data, "Get should not return empty data")
			},
		},
		{
			name: "Should fail to run Prepare",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to prepare statement")
				db.pushTestError(testErr)
				stmt, err := db.Prepare("SELECT * FROM custom_table WHERE name = ?")
				assert.Error(t, err, "Prepare should return an error")
				assert.Empty(t, stmt, "Prepare should return nil")
				assert.ErrorIs(t, err, testErr, "Prepare should return testErr")
			},
		},
		{
			name: "Should run PrepareContext",
			assert: func(t *testing.T, db DB) {
				stmt, err := db.PrepareContext(context.Background(), "SELECT * FROM custom_table WHERE name = ?")
				assert.NoError(t, err, "PrepareContext should not return an error")
				assert.NotNil(t, stmt, "PrepareContext should not return nil")

				var data customData
				err = stmt.GetContext(context.Background(), &data, "John Doe")
				assert.NoError(t, err, "GetContext should not return an error")
				assert.NotEmpty(t, data, "GetContext should not return empty data")
			},
		},
		{
			name: "Should fail to run PrepareContext",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to prepare statement")
				db.pushTestError(testErr)
				stmt, err := db.PrepareContext(context.Background(), "SELECT * FROM custom_table WHERE name = ?")
				assert.Error(t, err, "PrepareContext should return an error")
				assert.Empty(t, stmt, "PrepareContext should return nil")
				assert.ErrorIs(t, err, testErr, "PrepareContext should return testErr")
			},
		},
		{
			name: "Should run QueryRow",
			assert: func(t *testing.T, db DB) {
				var data customData
				err := db.QueryRow("SELECT * FROM custom_table").Scan(&data.ID, &data.Name, &data.Age)
				assert.NoError(t, err, "QueryRow should not return an error")
				assert.NotEmpty(t, data, "QueryRow should not return empty data")
			},
		},
		{
			name: "Should run QueryRowContext",
			assert: func(t *testing.T, db DB) {
				var data customData
				err := db.QueryRowContext(context.Background(), "SELECT * FROM custom_table").Scan(&data.ID, &data.Name, &data.Age)
				assert.NoError(t, err, "QueryRowContext should not return an error")
				assert.NotEmpty(t, data, "QueryRowContext should not return empty data")
			},
		},
		{
			name: "Should run Query",
			assert: func(t *testing.T, db DB) {
				rows, err := db.Query("SELECT id, name, age FROM custom_table")
				assert.NoError(t, err, "Query should not return an error")
				assert.NotNil(t, rows, "Query should not return nil")

				var data []customData
				for rows.Next() {
					var d customData
					err = rows.Scan(&d.ID, &d.Name, &d.Age)
					assert.NoError(t, err, "Scan should not return an error")
					data = append(data, d)
				}
				assert.NotEmpty(t, data, "Query should not return empty data")
				assert.Len(t, data, 5, "Query should return 5 rows")
			},
		},
		{
			name: "should fail to run Query",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to run query")
				db.pushTestError(testErr)
				rows, err := db.Query("SELECT id, name, age FROM custom_table")
				assert.Error(t, err, "Query should return an error")
				assert.Empty(t, rows, "Query should return nil")
				assert.ErrorIs(t, err, testErr, "Query should return testErr")
			},
		},
		{
			name: "Should run QueryContext",
			assert: func(t *testing.T, db DB) {
				rows, err := db.QueryContext(context.Background(), "SELECT * FROM custom_table")
				assert.NoError(t, err, "QueryContext should not return an error")
				assert.NotNil(t, rows, "QueryContext should not return nil")

				var data []customData
				for rows.Next() {
					var d customData
					err = rows.Scan(&d.ID, &d.Name, &d.Age)
					assert.NoError(t, err, "Scan should not return an error")
					data = append(data, d)
				}
				assert.NotEmpty(t, data, "QueryContext should not return empty data")
				assert.Len(t, data, 5, "QueryContext should return 5 rows")
			},
		},
		{
			name: "should fail to run QueryContext",
			assert: func(t *testing.T, db DB) {
				testErr := errors.New("failed to run query")
				db.pushTestError(testErr)
				rows, err := db.QueryContext(context.Background(), "SELECT id, name, age FROM custom_table")
				assert.Error(t, err, "QueryContext should return an error")
				assert.Empty(t, rows, "QueryContext should return nil")
				assert.ErrorIs(t, err, testErr, "QueryContext should return testErr")
			},
		},
		{
			name: "Should run Rebind",
			assert: func(t *testing.T, db DB) {
				query := db.Rebind("SELECT * FROM custom_table WHERE name = ?")
				assert.NotEmpty(t, query, "Rebind should not return empty query")

				rows, err := db.Query(query, "John Doe")
				assert.NoError(t, err, "Query should not return an error")
				assert.NotNil(t, rows, "Query should not return nil")

				var data []customData
				for rows.Next() {
					var d customData
					err = rows.Scan(&d.ID, &d.Name, &d.Age)
					assert.NoError(t, err, "Scan should not return an error")
					data = append(data, d)
				}
				assert.NotEmpty(t, data, "Query should not return empty data")
				assert.Len(t, data, 1, "Query should return 1 row")
			},
		},
		{
			name: "Should run Select",
			assert: func(t *testing.T, db DB) {
				var data []customData
				err := db.Select(&data, "SELECT * FROM custom_table")
				assert.NoError(t, err, "Select should not return an error")
				assert.NotEmpty(t, data, "Select should not return empty data")
				assert.Len(t, data, 5, "Select should return 5 rows")
			},
		},
		{
			name: "Should run SelectContext",
			assert: func(t *testing.T, db DB) {
				var data []customData
				err := db.SelectContext(context.Background(), &data, "SELECT * FROM custom_table")
				assert.NoError(t, err, "SelectContext should not return an error")
				assert.NotEmpty(t, data, "SelectContext should not return empty data")
				assert.Len(t, data, 5, "SelectContext should return 5 rows")
			},
		},
		{
			name: "Should run Unsafe",
			assert: func(t *testing.T, db DB) {
				assert.IsType(t, &sqlx.DB{}, db.Unsafe(), "Unsafe should return *sqlx.DB")
			},
		},
		{
			name: "Should run Safe",
			assert: func(t *testing.T, db DB) {
				assert.IsType(t, &sqlx.DB{}, db.Safe(), "Safe should return *sqlx.DB")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, db)
		})
	}
}

func Test_prepareConnectionParams(t *testing.T) {
	params := prepareConnectionParams(DBConnectionParams{})
	assert.Empty(t, params, "prepareConnectionParams should return empty params when no DBConnectionParams provided")
}
