package godb

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Row defines a new row
type Row interface {
	// sqlx

	// ColumnTypes returns the underlying sql.Rows.ColumnTypes(), or the deferred error
	ColumnTypes() ([]*sql.ColumnType, error)
	// Columns returns the underlying sql.Rows.Columns(), or the deferred error usually
	// returned by Row.Scan()
	Columns() ([]string, error)
	// Err returns the error encountered while scanning
	Err() error
	// MapScan using this Rows
	MapScan(dest map[string]interface{}) error
	// Scan is a fixed implementation of sql.Row.Scan, which does not discard the
	// underlying error from the internal rows object if it exists
	Scan(dest ...interface{}) error
	// SliceScan using this Rows
	SliceScan() ([]interface{}, error)
	// StructScan a single Row into dest
	StructScan(dest interface{}) error
}

// customRow implements the Row interface
type customRow struct {
	row *sqlx.Row
}

// ColumnTypes
func (c *customRow) ColumnTypes() ([]*sql.ColumnType, error) {
	return c.row.ColumnTypes()
}

// Columns
func (c *customRow) Columns() ([]string, error) {
	return c.row.Columns()
}

// Err
func (c *customRow) Err() error {
	return c.row.Err()
}

// Scan
func (c *customRow) Scan(dest ...interface{}) error {
	return c.row.Scan(dest...)
}

// MapScan
func (c *customRow) MapScan(dest map[string]interface{}) error {
	return c.row.MapScan(dest)
}

// SliceScan
func (c *customRow) SliceScan() ([]interface{}, error) {
	return c.row.SliceScan()
}

// StructScan
func (c *customRow) StructScan(dest interface{}) error {
	return c.row.StructScan(dest)
}
