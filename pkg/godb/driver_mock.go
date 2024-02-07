package godb

import (
	"database/sql/driver"
)

// DriverRowsMock defines a new driver rows mock
// It provides callback functions for each method to allow custom behavior during testing.
type DriverRowsMock struct {
	Error           error
	CallbackColumns func() []string
	CallbackClose   func() error
	CallbackNext    func(dest []driver.Value) error
}

// Columns returns the column names of the rows.
// If CallbackColumns is defined, it calls the callback function.
// Otherwise, it returns an empty slice.
func (drm *DriverRowsMock) Columns() []string {
	if drm.CallbackColumns != nil {
		return drm.CallbackColumns()
	}
	return []string{}
}

// Close closes the rows.
// If CallbackClose is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (drm *DriverRowsMock) Close() error {
	if drm.CallbackClose != nil {
		return drm.CallbackClose()
	}
	return drm.Error
}

// Next retrieves the next row of the result set.
// If CallbackNext is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (drm *DriverRowsMock) Next(dest []driver.Value) error {
	if drm.CallbackNext != nil {
		return drm.CallbackNext(dest)
	}
	return drm.Error
}

// DriverResultMock defines a new driver result mock
type DriverResultMock struct {
	Error                error
	CallbackLastInsertId func() (int64, error)
	CallbackRowsAffected func() (int64, error)
}

// LastInsertId returns the last inserted row id.
// If CallbackLastInsertId is defined, it calls the callback function.
// Otherwise, it returns 0 and the error defined in the struct.
func (drm *DriverResultMock) LastInsertId() (int64, error) {
	if drm.CallbackLastInsertId != nil {
		return drm.CallbackLastInsertId()
	}
	return 0, drm.Error
}

// RowsAffected returns the number of rows affected by the query.
// If CallbackRowsAffected is defined, it calls the callback function.
// Otherwise, it returns 0 and the error defined in the struct.
func (drm *DriverResultMock) RowsAffected() (int64, error) {
	if drm.CallbackRowsAffected != nil {
		return drm.CallbackRowsAffected()
	}
	return 0, drm.Error
}

// DriverTxMock defines a new driver transaction mock
type DriverTxMock struct {
	Error            error
	CallbackCommit   func() error
	CallbackRollback func() error
}

// Commit commits the transaction.
// If CallbackCommit is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (dtxm *DriverTxMock) Commit() error {
	if dtxm.CallbackCommit != nil {
		return dtxm.CallbackCommit()
	}
	return dtxm.Error
}

// Rollback rolls back the transaction.
// If CallbackRollback is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (dtxm *DriverTxMock) Rollback() error {
	if dtxm.CallbackRollback != nil {
		return dtxm.CallbackRollback()
	}
	return dtxm.Error
}

// DriverStmtMock defines a new driver statement mock
type DriverStmtMock struct {
	Error            error
	CallbackClose    func() error
	CallbackNumInput func() int
	CallbackExec     func(args []driver.Value) (driver.Result, error)
	CallbackQuery    func(args []driver.Value) (driver.Rows, error)
}

// Close closes the statement.
// If CallbackClose is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (dsm *DriverStmtMock) Close() error {
	if dsm.CallbackClose != nil {
		return dsm.CallbackClose()
	}
	return dsm.Error
}

// NumInput returns the number of placeholder parameters in the statement.
// If CallbackNumInput is defined, it calls the callback function.
// Otherwise, it returns 0.
func (dsm *DriverStmtMock) NumInput() int {
	if dsm.CallbackNumInput != nil {
		return dsm.CallbackNumInput()
	}
	return 0
}

// Exec executes a query that doesn't return rows.
// If CallbackExec is defined, it calls the callback function.
// Otherwise, it returns a new DriverResultMock and the error defined in the struct.
func (dsm *DriverStmtMock) Exec(args []driver.Value) (driver.Result, error) {
	if dsm.CallbackExec != nil {
		return dsm.CallbackExec(args)
	}
	return &DriverResultMock{}, dsm.Error
}

// Query executes a query that may return rows.
// If CallbackQuery is defined, it calls the callback function.
// Otherwise, it returns a new DriverRowsMock and the error defined in the struct.
func (dsm *DriverStmtMock) Query(args []driver.Value) (driver.Rows, error) {
	if dsm.CallbackQuery != nil {
		return dsm.CallbackQuery(args)
	}
	return &DriverRowsMock{}, dsm.Error
}

// DriverConnMock defines a new driver connection mock
type DriverConnMock struct {
	Error           error
	CallbackPrepare func(query string) (driver.Stmt, error)
	CallbackClose   func() error
	CallbackBegin   func() (driver.Tx, error)
}

// Prepare prepares a statement for execution.
// If CallbackPrepare is defined, it calls the callback function.
// Otherwise, it returns a new DriverStmtMock and the error defined in the struct.
func (dcm *DriverConnMock) Prepare(query string) (driver.Stmt, error) {
	if dcm.CallbackPrepare != nil {
		return dcm.CallbackPrepare(query)
	}
	return &DriverStmtMock{}, dcm.Error
}

// Close closes the connection.
// If CallbackClose is defined, it calls the callback function.
// Otherwise, it returns the error defined in the struct.
func (dcm *DriverConnMock) Close() error {
	if dcm.CallbackClose != nil {
		return dcm.CallbackClose()
	}
	return dcm.Error
}

// Begin starts a new transaction.
// If CallbackBegin is defined, it calls the callback function.
// Otherwise, it returns a new DriverTxMock and the error defined in the struct.
func (dcm *DriverConnMock) Begin() (driver.Tx, error) {
	if dcm.CallbackBegin != nil {
		return dcm.CallbackBegin()
	}
	return &DriverTxMock{}, dcm.Error
}

// DriverMock defines a new driver mock
type DriverMock struct {
	Error        error
	CallbackOpen func(name string) (driver.Conn, error)
}

// Open opens a new connection.
// If CallbackOpen is defined, it calls the callback function.
// Otherwise, it returns a new DriverConnMock and the error defined in the struct.
func (dm *DriverMock) Open(name string) (driver.Conn, error) {
	if dm.CallbackOpen != nil {
		return dm.CallbackOpen(name)
	}
	return &DriverConnMock{}, dm.Error
}
