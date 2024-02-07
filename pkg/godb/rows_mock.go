package godb

import (
	"database/sql"
)

// RowsMock is a mock implementation of the sql.Rows interface
type RowsMock struct {
	Error                 error
	CallbackClose         func() error
	CallbackColumnTypes   func() ([]*sql.ColumnType, error)
	CallbackColumns       func() ([]string, error)
	CallbackErr           func() error
	CallbackNext          func() bool
	CallbackNextResultSet func() bool
	CallbackScan          func(dest ...any) error
	CallbackMapScan       func(dest map[string]interface{}) error
	CallbackSliceScan     func() ([]interface{}, error)
	CallbackStructScan    func(dest interface{}) error
}

// Close closes the RowsMock.
// If CallbackClose is set, it will be called and its return value will be returned.
// Otherwise, it returns the Error field of RowsMock.
func (rm *RowsMock) Close() error {
	if rm.CallbackClose != nil {
		return rm.CallbackClose()
	}
	return rm.Error
}

// ColumnTypes returns the column types of the RowsMock.
// If CallbackColumnTypes is set, it will be called and its return value will be returned.
// Otherwise, it returns an empty slice of *sql.ColumnType and the Error field of RowsMock.
func (rm *RowsMock) ColumnTypes() ([]*sql.ColumnType, error) {
	if rm.CallbackColumnTypes != nil {
		return rm.CallbackColumnTypes()
	}
	return []*sql.ColumnType{}, rm.Error
}

// Columns returns the column names of the RowsMock.
// If CallbackColumns is set, it will be called and its return value will be returned.
// Otherwise, it returns an empty slice of strings and the Error field of RowsMock.
func (rm *RowsMock) Columns() ([]string, error) {
	if rm.CallbackColumns != nil {
		return rm.CallbackColumns()
	}
	return []string{}, rm.Error
}

// Err returns the error of the RowsMock.
// If CallbackErr is set, it will be called and its return value will be returned.
// Otherwise, it returns the Error field of RowsMock.
func (rm *RowsMock) Err() error {
	if rm.CallbackErr != nil {
		return rm.CallbackErr()
	}
	return rm.Error
}

// Next returns true if there is another row in the RowsMock, false otherwise.
// If CallbackNext is set, it will be called and its return value will be returned.
// Otherwise, it returns false.
func (rm *RowsMock) Next() bool {
	if rm.CallbackNext != nil {
		return rm.CallbackNext()
	}
	return false
}

// NextResultSet returns true if there is another result set in the RowsMock, false otherwise.
// If CallbackNextResultSet is set, it will be called and its return value will be returned.
// Otherwise, it returns false.
func (rm *RowsMock) NextResultSet() bool {
	if rm.CallbackNextResultSet != nil {
		return rm.CallbackNextResultSet()
	}
	return false
}

// Scan scans the current row of the RowsMock into the provided destination values.
// If CallbackScan is set, it will be called with the destination values and its return value will be returned.
// Otherwise, it returns the Error field of RowsMock.
func (rm *RowsMock) Scan(dest ...any) error {
	if rm.CallbackScan != nil {
		return rm.CallbackScan(dest...)
	}
	return rm.Error
}

// MapScan scans the current row of the RowsMock into the provided map.
// If CallbackMapScan is set, it will be called with the map and its return value will be returned.
// Otherwise, it returns the Error field of RowsMock.
func (rm *RowsMock) MapScan(dest map[string]interface{}) error {
	if rm.CallbackMapScan != nil {
		return rm.CallbackMapScan(dest)
	}
	return rm.Error
}

// SliceScan scans the current row of the RowsMock into a slice of interface{} values.
// If CallbackSliceScan is set, it will be called and its return value will be returned.
// Otherwise, it returns an empty slice of interface{} and the Error field of RowsMock.
func (rm *RowsMock) SliceScan() ([]interface{}, error) {
	if rm.CallbackSliceScan != nil {
		return rm.CallbackSliceScan()
	}
	return []interface{}{}, rm.Error
}

// StructScan scans the current row of the RowsMock into the provided struct.
// If CallbackStructScan is set, it will be called with the struct and its return value will be returned.
// Otherwise, it returns the Error field of RowsMock.
func (rm *RowsMock) StructScan(dest interface{}) error {
	if rm.CallbackStructScan != nil {
		return rm.CallbackStructScan(dest)
	}
	return rm.Error
}
