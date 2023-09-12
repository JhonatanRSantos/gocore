package pgdb

import "errors"

var (
	ErrFailedToPingDatabase         = errors.New("failed to ping database")
	ErrFailedToCreateConnectionPool = errors.New("failed to create a new connection pool")
)
