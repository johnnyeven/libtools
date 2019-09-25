package trans

import (
	"fmt"
	"runtime/debug"

	"golib/gorm"
)

const (
	internalErrorCode = 1
	dbErrCode         = 2
)

type serverError struct {
	Code int
	Msg  string
}

const serverErrorFmt = "ErrorCode:[%d]; Msg:[%s]"

func (this serverError) Error() string {
	return fmt.Sprintf(serverErrorFmt, this.Code, this.Msg)
}

type Task func(db *gorm.DB) error

func ErrHandler(db *gorm.DB, task Task) (err error) {
	defer func() {
		if e := recover(); e != nil {
			msg := fmt.Sprintf("panic: %s; calltrace:%s", fmt.Sprint(e), string(debug.Stack()))
			err = &serverError{internalErrorCode, msg}
		}
	}()
	return task(db)
}

func ExecTransaction(m_db *gorm.DB, transction ...Task) error {
	exec_db := m_db.Begin()
	if exec_db.Error != nil {
		return &serverError{dbErrCode, exec_db.Error.Error()}
	}
	for _, task := range transction {
		if err := ErrHandler(exec_db, task); err != nil {
			if rberr := exec_db.Rollback().Error; rberr != nil {
				return &serverError{dbErrCode, rberr.Error()}
			}
			//return &serverError{dbErrCode, err.Error()}
			return err
		}
	}

	if err := exec_db.Commit().Error; err != nil {
		exec_db.Rollback()
		return &serverError{dbErrCode, err.Error()}
	}
	return nil
}
