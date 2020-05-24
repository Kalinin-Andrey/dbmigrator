package api

import(
	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/dbmigrator/internal/pkg/apperror"
)

// ErrNotFound is error for case when entity not found
var ErrNotFound error = errors.New("Not found")
// ErrBadRequest is error for case when bad request
var ErrBadRequest error = errors.New("Bad request")
// ErrInternal is error for case when smth went wrong
var ErrInternal error = errors.New("Internal error")
// ErrDuplicate error
var ErrDuplicate error = errors.New("Duplicate error")
// ErrUsersFunc error
var ErrUsersFunc error = errors.New("User's func error")
// ErrUsersSQL error
var ErrUsersSQL error = errors.New("User's SQL string error")
// ErrUndefinedTypeOfAction error
var ErrUndefinedTypeOfAction error = errors.New("Undefined type of action")
// ErrNotInitialised error
var ErrNotInitialised error = errors.New("SQL Migrator is not initialised")

// AppErrorConv is a converter from app errors to api errors
func AppErrorConv(err error) (res error) {

	switch {
	case errors.Is(err, apperror.ErrNotFound):
		res = errors.Wrapf(ErrNotFound, "%v", err.Error())
	case errors.Is(err, apperror.ErrBadRequest):
		res = errors.Wrapf(ErrBadRequest, "%v", err.Error())
	case errors.Is(err, apperror.ErrInternal):
		res = errors.Wrapf(ErrInternal, "%v", err.Error())
	case errors.Is(err, apperror.ErrDuplicate):
		res = errors.Wrapf(ErrDuplicate, "%v", err.Error())
	case errors.Is(err, apperror.ErrUsersFunc):
		res = errors.Wrapf(ErrUsersFunc, "%v", err.Error())
	case errors.Is(err, apperror.ErrUsersSQL):
		res = errors.Wrapf(ErrUsersSQL, "%v", err.Error())
	case errors.Is(err, apperror.ErrUndefinedTypeOfAction):
		res = errors.Wrapf(ErrUndefinedTypeOfAction, "%v", err.Error())
	case errors.Is(err, apperror.ErrNotInitialised):
		res = errors.Wrapf(ErrNotInitialised, "%v", err.Error())
	default:
		res = err
	}
	return res
}

