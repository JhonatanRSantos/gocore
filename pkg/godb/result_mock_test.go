package godb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ResultMock(t *testing.T) {
	tests := []struct {
		name   string
		mock   ResultMock
		assert func(t *testing.T, mock ResultMock)
	}{
		{
			name: "Should run CallbackLastInsertId",
			mock: ResultMock{},
			assert: func(t *testing.T, mock ResultMock) {
				callCount := 0
				mock.CallbackLastInsertId = func() (int64, error) {
					callCount++
					return int64(1), nil
				}
				id, err := mock.LastInsertId()
				assert.NoError(t, err, "should not return an error")
				assert.Equal(t, int64(1), id, "should return the expected id")
				assert.Equal(t, 1, callCount, "should call the callback function")
			},
		},
		{
			name: "Should run default CallbackLastInsertId",
			mock: ResultMock{},
			assert: func(t *testing.T, mock ResultMock) {
				testErr := errors.New("failed to retrieve last insert id")
				mock.Err = testErr
				id, err := mock.LastInsertId()
				assert.ErrorIs(t, err, testErr, "should return the expected error")
				assert.Equal(t, int64(0), id, "should return the default id")
			},
		},
		{
			name: "Should run CallbackRowsAffected",
			mock: ResultMock{},
			assert: func(t *testing.T, mock ResultMock) {
				callCount := 0
				mock.CallbackRowsAffected = func() (int64, error) {
					callCount++
					return int64(1), nil
				}
				affected, err := mock.RowsAffected()
				assert.NoError(t, err, "should not return an error")
				assert.Equal(t, int64(1), affected, "should return the expected affected")
				assert.Equal(t, 1, callCount, "should call the callback function")
			},
		},
		{
			name: "Should run default CallbackRowsAffected",
			mock: ResultMock{},
			assert: func(t *testing.T, mock ResultMock) {
				testErr := errors.New("failed to retrieve rows affected")
				mock.Err = testErr
				affected, err := mock.RowsAffected()
				assert.ErrorIs(t, err, testErr, "should return the expected error")
				assert.Equal(t, int64(0), affected, "should return the default affected")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t, tt.mock)
		})
	}
}
