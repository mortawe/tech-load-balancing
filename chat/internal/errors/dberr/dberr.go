package dberr

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func isErr(err error, code string) bool {
	pgerr := &pgconn.PgError{}
	if !errors.As(err, &pgerr) {
		return false
	}
	return pgerr.Code == code
}

func IsUniqueViolationErr(err error) bool {
	return isErr(err, pgerrcode.UniqueViolation)
}

func IsForeignKeyViolation(err error) bool {
	return isErr(err, pgerrcode.ForeignKeyViolation)
}
