package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/pashagolub/pgxmock/v2"
)

var (
	ErrMissingConnections = errors.New("missing connections")
)

type ClientConnectionMock struct {
	pgxmock.PgxConnIface
}

func (ccm *ClientConnectionMock) Release() {}

type DBClientMock struct {
	lock        sync.Mutex
	connPool    pgxmock.PgxPoolIface
	connections []ClientConnectionMock
}

func NewClient() (*DBClientMock, error) {
	connPool, err := pgxmock.NewPool(
		pgxmock.MonitorPingsOption(true),
		pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual),
	)
	if err != nil {
		return nil, err
	}

	return &DBClientMock{
		connPool:    connPool,
		connections: []ClientConnectionMock{},
	}, nil
}

// Close Closes the connection pool
func (dbcm *DBClientMock) Close() {
	dbcm.connPool.Close()
}

// GetConnection Get a connection from pool
func (dbcm *DBClientMock) GetConnection(ctx context.Context) (ClientConnectionMock, error) {
	dbcm.lock.Lock()
	defer dbcm.lock.Unlock()

	if len(dbcm.connections) == 0 {
		return ClientConnectionMock{}, ErrMissingConnections
	}

	conn := dbcm.connections[0]
	dbcm.connections = dbcm.connections[1:]

	return conn, nil
}

// AddConnection Adds a new pre configured connection to the pool
func (dbcm *DBClientMock) AddConnection(conn ClientConnectionMock) {
	dbcm.lock.Lock()
	defer dbcm.lock.Unlock()

	dbcm.connections = append(dbcm.connections, conn)
}
