package repository

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func IsDuplicateErr(err error) bool {
	mysqlErr := &mysql.MySQLError{}
	if !errors.As(err, &mysqlErr) {
		return false
	}
	return mysqlErr.Number == 0x426
}
