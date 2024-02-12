package godb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DBConnectionParams map[string]string

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

// NewDB Returns a new database connection
func NewDB(config DBConfig) (DB, error) {
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

	if dbx, err := connect(ctx, cancel, dbType, dsn); err != nil {
		return &customDB{}, err
	} else {
		return &customDB{db: dbx}, nil
	}
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

type DBType uint

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

// Define database interface
type DB interface {
	popTestError() error
	pushTestError(err error)

	// sql

	// Close closes the database and prevents new queries from starting.
	// Close then waits for all queries that have started processing on the server
	// to finish.
	//
	// It is rare to Close a godb.DB, as the godb.DB handle is meant to be
	// long-lived and shared between many goroutines.
	Close() error
	// Driver returns the database's underlying driver.
	Driver() driver.Driver
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	//
	// Exec uses context.Background internally; to specify the context, use
	// ExecContext.
	Exec(query string, args ...any) (sql.Result, error)
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	//
	// Ping uses context.Background internally; to specify the context, use
	// PingContext.
	Ping() error
	// PingContext verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error
	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's idle time.
	SetConnMaxIdleTime(d time.Duration)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's age.
	SetConnMaxLifetime(d time.Duration)
	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	//
	// The default max idle connections is currently 2. This may change in
	// a future release.
	SetMaxIdleConns(n int)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	SetMaxOpenConns(n int)
	// Stats returns database statistics.
	Stats() sql.DBStats

	// sqlx

	// BeginTx begins a transaction and returns an godb.Tx instead of an *sql.Tx.
	//
	// The provided context is used until the transaction is committed or rolled
	// back. If the context is canceled, the sql package will roll back the
	// transaction. Tx.Commit will return an error if the context provided to
	// BeginxContext is canceled.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	// Begin begins a transaction and returns an godb.Tx instead of an *sql.Tx.
	Begin() (Tx, error)
	// BindNamed binds a query using the DB driver's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	// Conn returns an godb.Conn instead of an *sql.Conn.
	Conn(ctx context.Context) (Conn, error)
	// DriverName returns the driverName passed to the Open function for this DB.
	DriverName() string
	// Get using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error
	// GetContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// MapperFunc sets a new mapper for this db using the default sqlx struct tag
	// and the provided mapper function.
	MapperFunc(mf func(string) string)
	// MustBegin starts a transaction, and panics on error.  Returns an godb.Tx instead
	// of an *sql.Tx.
	MustBegin() Tx
	// MustBeginTx starts a transaction, and panics on error.  Returns an godb.Tx instead
	// of an *sql.Tx.
	//
	// The provided context is used until the transaction is committed or rolled
	// back. If the context is canceled, the sql package will roll back the
	// transaction. Tx.Commit will return an error if the context provided to
	// MustBeginContext is canceled.
	MustBeginTx(ctx context.Context, opts *sql.TxOptions) Tx
	// MustExec (panic) runs MustExec using this database.
	// Any placeholder parameters are replaced with supplied args.
	MustExec(query string, args ...interface{}) sql.Result
	// MustExecContext (panic) runs MustExec using this database.
	// Any placeholder parameters are replaced with supplied args.
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	// NamedExec using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)
	// NamedExecContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	// NamedQuery using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (Rows, error)
	// NamedQueryContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (Rows, error)
	// PrepareNamed returns an godb.NamedStmt
	PrepareNamed(query string) (NamedStmt, error)
	// PrepareNamedContext returns an godb.NamedStmt
	PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error)
	// Prepare returns an godb.Stmt instead of a sql.Stmt
	Prepare(query string) (Stmt, error)
	// PrepareContext returns an godb.Stmt instead of a sql.Stmt.
	//
	// The provided context is used for the preparation of the statement, not for
	// the execution of the statement.
	PrepareContext(ctx context.Context, query string) (Stmt, error)
	// QueryRow queries the database and returns an Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRow(query string, args ...interface{}) Row
	// QueryRowContext queries the database and returns an Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	// Query queries the database and returns an Rows.
	// Any placeholder parameters are replaced with supplied args.
	Query(query string, args ...interface{}) (Rows, error)
	// QueryContext queries the database and returns an Rows.
	// Any placeholder parameters are replaced with supplied args.
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string
	// Select using this DB.
	// Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error
	// SelectContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// Unsafe returns a version of DB which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	// godb.Stmt and godb.Tx which are created from this DB will inherit its
	// safety behavior.
	Unsafe() *sqlx.DB
	// Safe returns the underlying DB
	Safe() *sqlx.DB
}

// Define all database configs
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

// Inplements database interface
type customDB struct {
	db      *sqlx.DB
	testErr error
}

// pushTestError
func (cdb *customDB) pushTestError(err error) {
	cdb.testErr = err
}

// popTestError
func (cdb *customDB) popTestError() error {
	err := cdb.testErr
	cdb.testErr = nil
	return err
}

// Close
func (cdb *customDB) Close() error {
	return cdb.db.Close()
}

// Driver
func (cdb *customDB) Driver() driver.Driver {
	return cdb.db.Driver()
}

// Exec
func (cdb *customDB) Exec(query string, args ...any) (sql.Result, error) {
	return cdb.db.Exec(query, args...)
}

// ExecContext
func (cdb *customDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return cdb.db.ExecContext(ctx, query, args...)
}

// Ping
func (cdb *customDB) Ping() error {
	return cdb.db.Ping()
}

// PingContext
func (cdb *customDB) PingContext(ctx context.Context) error {
	return cdb.db.PingContext(ctx)
}

// SetConnMaxIdleTime
func (cdb *customDB) SetConnMaxIdleTime(duration time.Duration) {
	cdb.db.SetConnMaxIdleTime(duration)
}

// SetConnMaxLifetime
func (cdb *customDB) SetConnMaxLifetime(duration time.Duration) {
	cdb.db.SetConnMaxLifetime(duration)
}

// SetMaxIdleConns
func (cdb *customDB) SetMaxIdleConns(n int) {
	cdb.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns
func (cdb *customDB) SetMaxOpenConns(n int) {
	cdb.db.SetMaxOpenConns(n)
}

// Stats
func (cdb *customDB) Stats() sql.DBStats {
	return cdb.db.Stats()
}

// BeginTx
func (cdb *customDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	if tx, err := cdb.beginTx(ctx, opts); err != nil {
		return &customTx{}, err
	} else {
		return &customTx{tx: tx}, nil
	}
}

// beginTx
func (cdb *customDB) beginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	if cdb.testErr != nil {
		return &sqlx.Tx{}, cdb.popTestError()
	}
	return cdb.db.BeginTxx(ctx, opts)
}

// Begin
func (cdb *customDB) Begin() (Tx, error) {
	if tx, err := cdb.begin(); err != nil {
		return &customTx{}, err
	} else {
		return &customTx{tx: tx}, nil
	}
}

// begin
func (cdb *customDB) begin() (*sqlx.Tx, error) {
	if cdb.testErr != nil {
		return &sqlx.Tx{}, cdb.popTestError()
	}
	return cdb.db.Beginx()
}

// BindNamed
func (cdb *customDB) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return cdb.db.BindNamed(query, arg)
}

// Conn
func (cdb *customDB) Conn(ctx context.Context) (Conn, error) {
	if conn, err := cdb.conn(ctx); err != nil {
		return &customConn{}, err
	} else {
		return &customConn{conn: conn}, nil
	}
}

// conn
func (cdb *customDB) conn(ctx context.Context) (*sqlx.Conn, error) {
	if cdb.testErr != nil {
		return &sqlx.Conn{}, cdb.popTestError()
	}
	return cdb.db.Connx(ctx)
}

// DriverName
func (cdb *customDB) DriverName() string {
	return cdb.db.DriverName()
}

// Get
func (cdb *customDB) Get(dest interface{}, query string, args ...interface{}) error {
	return cdb.db.Get(dest, query, args...)
}

// GetContext
func (cdb *customDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return cdb.db.GetContext(ctx, dest, query, args...)
}

// MapperFunc
func (cdb *customDB) MapperFunc(mf func(string) string) {
	cdb.db.MapperFunc(mf)
}

// MustBegin
func (cdb *customDB) MustBegin() Tx {
	return &customTx{tx: cdb.db.MustBegin()}
}

// MustBeginTx
func (cdb *customDB) MustBeginTx(ctx context.Context, opts *sql.TxOptions) Tx {
	return &customTx{tx: cdb.db.MustBeginTx(ctx, opts)}
}

// MustExec
func (cdb *customDB) MustExec(query string, args ...interface{}) sql.Result {
	return cdb.db.MustExec(query, args...)
}

// MustExecContext
func (cdb *customDB) MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
	return cdb.db.MustExecContext(ctx, query, args...)
}

// NamedExec
func (cdb *customDB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return cdb.db.NamedExec(query, arg)
}

// NamedExecContext
func (cdb *customDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return cdb.db.NamedExecContext(ctx, query, arg)
}

// NamedQuery
func (cdb *customDB) NamedQuery(query string, arg interface{}) (Rows, error) {
	if rows, err := cdb.namedQuery(query, arg); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// namedQuery
func (cdb *customDB) namedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	if cdb.testErr != nil {
		return &sqlx.Rows{}, cdb.popTestError()
	}
	return cdb.db.NamedQuery(query, arg)
}

// NamedQueryContext
func (cdb *customDB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (Rows, error) {
	if rows, err := cdb.namedQueryContext(ctx, query, arg); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// namedQueryContext
func (cdb *customDB) namedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	if cdb.testErr != nil {
		return &sqlx.Rows{}, cdb.popTestError()
	}
	return cdb.db.NamedQueryContext(ctx, query, arg)
}

// PrepareNamed
func (cdb *customDB) PrepareNamed(query string) (NamedStmt, error) {
	if namedStmt, err := cdb.prepareNamed(query); err != nil {
		return &customNamedStmt{}, err
	} else {
		return &customNamedStmt{namedStmt: namedStmt}, nil
	}
}

// prepareNamed
func (cdb *customDB) prepareNamed(query string) (*sqlx.NamedStmt, error) {
	if cdb.testErr != nil {
		return &sqlx.NamedStmt{}, cdb.popTestError()
	}
	return cdb.db.PrepareNamed(query)
}

// PrepareNamedContext
func (cdb *customDB) PrepareNamedContext(ctx context.Context, query string) (NamedStmt, error) {
	if namedStmt, err := cdb.prepareNamedContext(ctx, query); err != nil {
		return &customNamedStmt{}, err
	} else {
		return &customNamedStmt{namedStmt: namedStmt}, nil
	}
}

// prepareNamedContext
func (cdb *customDB) prepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if cdb.testErr != nil {
		return &sqlx.NamedStmt{}, cdb.popTestError()
	}
	return cdb.db.PrepareNamedContext(ctx, query)
}

// Prepare
func (cdb *customDB) Prepare(query string) (Stmt, error) {
	if stmt, err := cdb.prepare(query); err != nil {
		return &customStmt{}, err
	} else {
		return &customStmt{stmt: stmt}, nil
	}
}

// prepare
func (cdb *customDB) prepare(query string) (*sqlx.Stmt, error) {
	if cdb.testErr != nil {
		return &sqlx.Stmt{}, cdb.popTestError()
	}
	return cdb.db.Preparex(query)
}

// PrepareContext
func (cdb *customDB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	if stmt, err := cdb.prepareContext(ctx, query); err != nil {
		return &customStmt{}, err
	} else {
		return &customStmt{stmt: stmt}, nil
	}
}

// prepareContext
func (cdb *customDB) prepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	if cdb.testErr != nil {
		return &sqlx.Stmt{}, cdb.popTestError()
	}
	return cdb.db.PreparexContext(ctx, query)
}

// QueryRow
func (cdb *customDB) QueryRow(query string, args ...interface{}) Row {
	return &customRow{row: cdb.db.QueryRowx(query, args...)}
}

// QueryRowContext
func (cdb *customDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return &customRow{row: cdb.db.QueryRowxContext(ctx, query, args...)}
}

// Query
func (cdb *customDB) Query(query string, args ...interface{}) (Rows, error) {
	if rows, err := cdb.query(query, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// query
func (cdb *customDB) query(query string, args ...interface{}) (*sqlx.Rows, error) {
	if cdb.testErr != nil {
		return &sqlx.Rows{}, cdb.popTestError()
	}
	return cdb.db.Queryx(query, args...)
}

// QueryContext
func (cdb *customDB) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	if rows, err := cdb.queryContext(ctx, query, args...); err != nil {
		return &customRows{}, err
	} else {
		return &customRows{rows: rows}, nil
	}
}

// queryContext
func (cdb *customDB) queryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if cdb.testErr != nil {
		return &sqlx.Rows{}, cdb.popTestError()
	}
	return cdb.db.QueryxContext(ctx, query, args...)
}

// Rebind
func (cdb *customDB) Rebind(query string) string {
	return cdb.db.Rebind(query)
}

// Select
func (cdb *customDB) Select(dest interface{}, query string, args ...interface{}) error {
	return cdb.db.Select(dest, query, args...)
}

// SelectContext
func (cdb *customDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return cdb.db.SelectContext(ctx, dest, query, args...)
}

// Unsafe
func (cdb *customDB) Unsafe() *sqlx.DB {
	return cdb.db.Unsafe()
}

// Safe
func (cdb *customDB) Safe() *sqlx.DB {
	return cdb.db
}
