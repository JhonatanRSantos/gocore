package godb

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

//nolint:dupl
func TestMySQLDB(t *testing.T) {
	var (
		db           *sqlx.DB
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

	if db, err = New(DBConfig{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     MySQLDB,
		ConnectionParams: MySQLDefaultParams,
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

//nolint:dupl
func TestPostgres(t *testing.T) {
	var (
		db           *sqlx.DB
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

	if db, err = New(DBConfig{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     PostgresDB,
		ConnectionParams: PostgresDefaultParams,
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

//nolint:dupl
func TestSQLiteDB(t *testing.T) {
	var (
		db           *sqlx.DB
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

	if db, err = New(DBConfig{
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

func TestConnectTimeout(t *testing.T) {
	_, err := New(DBConfig{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     MySQLDB,
		ConnectionParams: PostgresDefaultParams,
	})
	assert.ErrorIs(t, err, ErrConnectionTimeoutExceeded)
}

func TestInvalidDatabaseType(t *testing.T) {
	var newDB DBType = 100
	_, err := New(DBConfig{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             "admin",
		Password:         "qwerty",
		Database:         "test-db",
		DatabaseType:     newDB,
		ConnectionParams: PostgresDefaultParams,
	})
	assert.ErrorIs(t, err, ErrInvalidDBType)
}
