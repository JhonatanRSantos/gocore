package godb

import "database/sql"

// RowMock is a mock implementation of the sql.Row interface
type RowMock struct {
	Error               error
	CallbackColumnTypes func() ([]*sql.ColumnType, error)
	CallbackColumns     func() ([]string, error)
	CallbackErr         func() error
	CallbackMapScan     func(dest map[string]interface{}) error
	CallbackScan        func(dest ...interface{}) error
	CallbackSliceScan   func() ([]interface{}, error)
	CallbackStructScan  func(dest interface{}) error
}

// ColumnTypes returns the column types of the row.
// If CallbackColumnTypes is defined, it calls the callback function.
// Otherwise, it returns an empty slice of sql.ColumnType and the stored error.
func (rm *RowMock) ColumnTypes() ([]*sql.ColumnType, error) {
	if rm.CallbackColumnTypes != nil {
		return rm.CallbackColumnTypes()
	}
	return []*sql.ColumnType{}, rm.Error
}

// Columns returns the column names of the row.
// If CallbackColumns is defined, it calls the callback function.
// Otherwise, it returns an empty slice of strings and the stored error.
func (rm *RowMock) Columns() ([]string, error) {
	if rm.CallbackColumns != nil {
		return rm.CallbackColumns()
	}
	return []string{}, rm.Error
}

// Err returns the error associated with the row.
// If CallbackErr is defined, it calls the callback function.
// Otherwise, it returns the stored error.
func (rm *RowMock) Err() error {
	if rm.CallbackErr != nil {
		return rm.CallbackErr()
	}
	return rm.Error
}

// MapScan scans the row into a map of column names to values.
// If CallbackMapScan is defined, it calls the callback function.
// Otherwise, it returns the stored error.
func (rm *RowMock) MapScan(dest map[string]interface{}) error {
	if rm.CallbackMapScan != nil {
		return rm.CallbackMapScan(dest)
	}
	return rm.Error
}

// Scan scans the row into the provided destination values.
// If CallbackScan is defined, it calls the callback function.
// Otherwise, it returns the stored error.
func (rm *RowMock) Scan(dest ...interface{}) error {
	if rm.CallbackScan != nil {
		return rm.CallbackScan(dest...)
	}
	return rm.Error
}

// SliceScan scans the row into a slice of interface{} values.
// If CallbackSliceScan is defined, it calls the callback function.
// Otherwise, it returns an empty slice of interface{} and the stored error.
func (rm *RowMock) SliceScan() ([]interface{}, error) {
	if rm.CallbackSliceScan != nil {
		return rm.CallbackSliceScan()
	}
	return []interface{}{}, rm.Error
}

// StructScan scans the row into the provided destination struct.
// If CallbackStructScan is defined, it calls the callback function.
// Otherwise, it returns the stored error.
func (rm *RowMock) StructScan(dest interface{}) error {
	if rm.CallbackStructScan != nil {
		return rm.CallbackStructScan(dest)
	}
	return rm.Error
}
