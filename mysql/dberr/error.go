package dberr

import (
	"errors"
)

//database error
var (
	RecordNotFoundError     = errors.New("database record not found")
	RecordConflictError     = errors.New("database record conflict")
	RecordCreateFailedError = errors.New("database record create failed")
	RecordFetchFailedError  = errors.New("database record fetch failed")
	RecordUpdateFailedError = errors.New("database record update failed")
	RecordDeleteFailedError = errors.New("database record delete failed")
	DbError                 = errors.New("database error")
)

var (
	DuplicateEntryErrNumber uint16 = 1062
)
