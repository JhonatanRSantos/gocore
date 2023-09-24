package godb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DBType             uint
	DBConnectionParams map[string]string
)

const (
	MySQLDB DBType = iota + 1
	SQLiteDB
	PostgresDB
)

var (
	defaultTimeout = time.Second * 5

	validDBTypes = map[DBType]string{
		MySQLDB:    "mysql",
		PostgresDB: "postgres",
		SQLiteDB:   "sqlite3",
	}

	// https://github.com/go-sql-driver/mysql
	MySQLDefaultParams = DBConnectionParams{
		"timeout":      defaultTimeout.String(),
		"readTimeout":  defaultTimeout.String(),
		"writeTimeout": defaultTimeout.String(),
	}

	// https://www.postgresql.org/docs/15/libpq-connect.html#LIBPQ-PARAMKEYWORDS
	// https://www.postgresql.org/docs/15/runtime-config-client.html
	// https://www.postgresql.org/docs/current/runtime-config-client.html
	PostgresDefaultParams = DBConnectionParams{
		"connect_timeout":   fmt.Sprintf("%d", defaultTimeout.Milliseconds()),
		"statement_timeout": fmt.Sprintf("%d", defaultTimeout.Milliseconds()),
		"sslmode":           "disable",
	}

	// https://pkg.go.dev/github.com/mattn/go-sqlite3
	SQLiteDefaultParams = DBConnectionParams{
		"cache": "private",
		"mode":  "memory",
	}
)

// New create a database connection
func New(config DBConfig) (*sqlx.DB, error) {
	if !config.DatabaseType.isValid() {
		return nil, ErrInvalidDBType
	}

	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = defaultTimeout
	}

	ctx := context.Background()
	dsn := config.dsn()
	dbType := config.DatabaseType.String()

	ctx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	return connect(ctx, cancel, dbType, dsn)
}

// connect Open a new databse connection
func connect(ctx context.Context, cancel context.CancelFunc, dbType, dsn string) (*sqlx.DB, error) {
	defer cancel()
	var (
		db  *sqlx.DB
		err error
	)

	go func() {
		db, err = sqlx.ConnectContext(ctx, dbType, dsn)
	}()

	<-ctx.Done()
	switch {
	case err != nil:
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrConnectionTimeoutExceeded
		}
		return nil, err
	case db != nil:
		return db, nil
	default:
		return nil, ErrConnectionTimeoutExceeded
	}
}

// isValid check if DBType is valid
func (dbt DBType) isValid() bool {
	for _, currentDB := range validDBTypes {
		if dbt.String() == currentDB {
			return true
		}
	}
	return false
}

// String return DBType as string
func (dbt DBType) String() string {
	return validDBTypes[dbt]
}

type DBConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	Database         string
	DatabaseType     DBType
	ConnectTimeout   time.Duration
	ConnectionParams DBConnectionParams
}

// dsn return data source name
func (dbc DBConfig) dsn() string {
	var connParams string

	if len(dbc.ConnectionParams) > 0 {
		connParams = fmt.Sprintf("?%s", prepareConnectionParams(dbc.ConnectionParams))
	}

	switch dbc.DatabaseType {
	case PostgresDB:
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s%s",
			dbc.User,
			dbc.Password,
			dbc.Host,
			dbc.Port,
			dbc.Database,
			connParams,
		)
	case MySQLDB:
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s%s",
			dbc.User,
			dbc.Password,
			dbc.Host,
			dbc.Port,
			dbc.Database,
			connParams,
		)
	case SQLiteDB:
		if connParams != "" {
			connParams = fmt.Sprintf("&%s", connParams[1:])
		}

		return fmt.Sprintf(
			"file:%s.db?_auth&_auth_user=%s&_auth_pass=%s%s",
			dbc.Database,
			dbc.User,
			dbc.Password,
			connParams,
		)
	default:
		return ""
	}
}

// prepareConnectionParams return the dsn params
func prepareConnectionParams(connectionParams DBConnectionParams) string {
	if len(connectionParams) == 0 {
		return ""
	}

	params := ""
	for key, value := range connectionParams {
		if params == "" {
			params = fmt.Sprintf("%s=%s", key, value)
		} else {
			params += fmt.Sprintf("&%s=%s", key, value)
		}
	}
	return params
}
