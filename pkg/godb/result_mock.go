package godb

// ResultMock defines a mock for the sql.Result interface
type ResultMock struct {
	Err                  error
	CallbackLastInsertId func() (int64, error)
	CallbackRowsAffected func() (int64, error)
}

// LastInsertId returns the last inserted row id.
// If a callback function is defined, it will be called to retrieve the value.
// Otherwise, it returns 0 and the error defined in ResultMock.Err.
func (rm *ResultMock) LastInsertId() (int64, error) {
	if rm.CallbackLastInsertId != nil {
		return rm.CallbackLastInsertId()
	}
	return 0, rm.Err
}

// RowsAffected returns the number of rows affected by the query.
// If a callback function is defined, it will be called to retrieve the value.
// Otherwise, it returns 0 and the error defined in ResultMock.Err.
func (rm *ResultMock) RowsAffected() (int64, error) {
	if rm.CallbackRowsAffected != nil {
		return rm.CallbackRowsAffected()
	}
	return 0, rm.Err
}
