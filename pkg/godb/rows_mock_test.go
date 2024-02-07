package godb

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RowMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   RowsMock
		assert func(t *testing.T, mock RowsMock)
	}{
		{
			name: "Shoudld run CallbackClose",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackClose = func() error {
					callCount++
					return nil
				}
				err := mock.Close()
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, callCount, "CallbackClose should be called once")
			},
		},
		{
			name: "Shoudld default CallbackClose",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("failed to close rows")
				mock.Error = testErr
				err := mock.Close()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
		{
			name: "Shoudld run CallbackColumnTypes",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackColumnTypes = func() ([]*sql.ColumnType, error) {
					callCount++
					return []*sql.ColumnType{}, nil
				}
				columns, err := mock.ColumnTypes()
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, []*sql.ColumnType{}, columns, "Columns should be []*sql.ColumnType{}")
				assert.Equal(t, 1, callCount, "CallbackColumnTypes should be called once")
			},
		},
		{
			name: "Shoudld default CallbackColumnTypes",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("failed to get column types")
				mock.Error = testErr
				columns, err := mock.ColumnTypes()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.IsType(t, []*sql.ColumnType{}, columns, "Columns should be []*sql.ColumnType{}")
			},
		},
		{
			name: "Shoudld run CallbackColumns",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackColumns = func() ([]string, error) {
					callCount++
					return []string{}, nil
				}
				columns, err := mock.Columns()
				assert.NoError(t, err, "Error should be nil")
				assert.IsType(t, []string{}, columns, "Columns should be []string{}")
				assert.Equal(t, 1, callCount, "CallbackColumns should be called once")
			},
		},
		{
			name: "Shoudld default CallbackColumns",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("failed to get columns")
				mock.Error = testErr
				columns, err := mock.Columns()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.IsType(t, []string{}, columns, "Columns should be []string{}")
			},
		},
		{
			name: "Shoudld run CallbackErr",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("test error")
				mock.Error = testErr
				callCount := 0
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
			name: "Shoudld default CallbackErr",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("test error")
				mock.Error = testErr
				err := mock.Err()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
			},
		},
		{
			name: "Shoudld run CallbackNext",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackNext = func() bool {
					callCount++
					return true
				}
				assert.True(t, mock.Next(), "Next should be true")
				assert.Equal(t, 1, callCount, "CallbackNext should be called once")
			},
		},
		{
			name: "Shoudld default CallbackNext",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				assert.False(t, mock.Next(), "Next should be false")
			},
		},
		{
			name: "Shoudld run CallbackNextResultSet",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackNextResultSet = func() bool {
					callCount++
					return true
				}
				assert.True(t, mock.NextResultSet(), "NextResultSet should be true")
				assert.Equal(t, 1, callCount, "CallbackNextResultSet should be called once")
			},
		},
		{
			name: "Shoudld default CallbackNextResultSet",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				assert.False(t, mock.NextResultSet(), "NextResultSet should be false")
			},
		},
		{
			name: "Shoudld run CallbackScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				data := 0
				callCount := 0
				mock.CallbackScan = func(dest ...interface{}) error {
					callCount++
					*dest[0].(*int) = 1
					return nil
				}
				err := mock.Scan(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data, "Data should be 1")
				assert.Equal(t, 1, callCount, "CallbackScan should be called once")
			},
		},
		{
			name: "Shoudld default CallbackScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				data := 0
				testErr := errors.New("failed to scan")
				mock.Error = testErr
				err := mock.Scan(&data)
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Equal(t, 0, data, "Data should be 0")
			},
		},
		{
			name: "Shoudld run CallbackMapScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				data := map[string]interface{}{}
				callCount := 0
				mock.CallbackMapScan = func(dest map[string]interface{}) error {
					callCount++
					dest["value"] = 1
					return nil
				}
				err := mock.MapScan(data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data["value"], "Data should be 1")
				assert.Equal(t, 1, callCount, "CallbackMapScan should be called once")
			},
		},
		{
			name: "Shoudld default CallbackMapScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				data := map[string]interface{}{}
				testErr := errors.New("failed to map scan")
				mock.Error = testErr
				err := mock.MapScan(data)
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Empty(t, data, "Data should be empty")
			},
		},
		{
			name: "Shoudld run CallbackSliceScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				callCount := 0
				mock.CallbackSliceScan = func() ([]interface{}, error) {
					callCount++
					return []interface{}{1}, nil
				}
				data, err := mock.SliceScan()
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data[0], "Data should be 1")
				assert.Equal(t, 1, callCount, "CallbackSliceScan should be called once")
			},
		},
		{
			name: "Shoudld default CallbackSliceScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				testErr := errors.New("failed to slice scan")
				mock.Error = testErr
				data, err := mock.SliceScan()
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Empty(t, data, "Data should be empty")
			},
		},
		{
			name: "Shoudld run CallbackStructScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				type dataStruct struct {
					Value int
				}
				data := dataStruct{}
				callCount := 0
				mock.CallbackStructScan = func(dest interface{}) error {
					callCount++
					dest.(*dataStruct).Value = 1
					return nil
				}
				err := mock.StructScan(&data)
				assert.NoError(t, err, "Error should be nil")
				assert.Equal(t, 1, data.Value, "Data should be 1")
				assert.Equal(t, 1, callCount, "CallbackStructScan should be called once")
			},
		},
		{
			name: "Shoudld default CallbackStructScan",
			mock: RowsMock{},
			assert: func(t *testing.T, mock RowsMock) {
				type dataStruct struct {
					Value int
				}
				data := dataStruct{}
				testErr := errors.New("failed to struct scan")
				mock.Error = testErr
				err := mock.StructScan(&data)
				assert.ErrorIs(t, err, testErr, "Error should be testErr")
				assert.Zero(t, data.Value, "Data should be zero")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
