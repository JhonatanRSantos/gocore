package godb

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RowsMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   RowMock
		assert func(t *testing.T, mock RowMock)
	}{
		{
			name: "Should run CallbackColumnTypes",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				callCount := 0
				mock.CallbackColumnTypes = func() ([]*sql.ColumnType, error) {
					callCount++
					return []*sql.ColumnType{{}}, nil
				}
				columns, err := mock.ColumnTypes()
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackColumnTypes should be called once")
				assert.Len(t, columns, 1, "Columns should have length 1")
			},
		},
		{
			name: "Should default CallbackColumnTypes",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to get column types")
				mock.Error = testErr
				columns, err := mock.ColumnTypes()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Len(t, columns, 0, "Columns should have length 0")
			},
		},
		{
			name: "Should run CallbackColumns",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				callCount := 0
				mock.CallbackColumns = func() ([]string, error) {
					callCount++
					return []string{"column"}, nil
				}
				columns, err := mock.Columns()
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackColumns should be called once")
				assert.Len(t, columns, 1, "Columns should have length 1")
			},
		},
		{
			name: "Should default CallbackColumns",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to get columns")
				mock.Error = testErr
				columns, err := mock.Columns()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Len(t, columns, 0, "Columns should have length 0")
			},
		},
		{
			name: "Should run CallbackErr",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				callCount := 0
				testErr := errors.New("failed to get error")
				mock.CallbackErr = func() error {
					callCount++
					return testErr
				}
				err := mock.Err()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Equal(t, 1, callCount, "CallbackErr should be called once")
			},
		},
		{
			name: "Should default CallbackErr",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to get error")
				mock.Error = testErr
				err := mock.Err()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
		{
			name: "Should run CallbackMapScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				data := map[string]interface{}{}
				callCount := 0
				mock.CallbackMapScan = func(dest map[string]interface{}) error {
					callCount++
					dest["column"] = "value"
					return nil
				}
				err := mock.MapScan(data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, "value", data["column"], "Data should have column value")
				assert.Equal(t, 1, callCount, "CallbackMapScan should be called once")
			},
		},
		{
			name: "Should default CallbackMapScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to map scan")
				mock.Error = testErr
				err := mock.MapScan(map[string]interface{}{})
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
		{
			name: "Should run CallbackScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				data := 0
				callCount := 0
				mock.CallbackScan = func(dest ...interface{}) error {
					callCount++
					*dest[0].(*int) = 1
					return nil
				}
				err := mock.Scan(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackScan should be called once")
				assert.Equal(t, 1, data, "Data should be 1")
			},
		},
		{
			name: "Should default CallbackScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to scan")
				mock.Error = testErr
				err := mock.Scan(nil)
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
		{
			name: "Should run CallbackSliceScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				callCount := 0
				mock.CallbackSliceScan = func() ([]interface{}, error) {
					callCount++
					return []interface{}{"value"}, nil
				}
				data, err := mock.SliceScan()
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackSliceScan should be called once")
				assert.Equal(t, "value", data[0], "Data should have value")
			},
		},
		{
			name: "Should default CallbackSliceScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to slice scan")
				mock.Error = testErr
				data, err := mock.SliceScan()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Len(t, data, 0, "Data should have length 0")
			},
		},
		{
			name: "Should run CallbackStructScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				type Data struct {
					Column string
				}
				data := Data{}
				callCount := 0
				mock.CallbackStructScan = func(dest interface{}) error {
					callCount++
					data.Column = "value"
					return nil
				}
				err := mock.StructScan(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, "value", data.Column, "Data should have column value")
				assert.Equal(t, 1, callCount, "CallbackStructScan should be called once")
			},
		},
		{
			name: "Should default CallbackStructScan",
			mock: RowMock{},
			assert: func(t *testing.T, mock RowMock) {
				testErr := errors.New("failed to struct scan")
				mock.Error = testErr
				err := mock.StructScan(nil)
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
