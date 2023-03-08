package dberr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/errs"
	"gorm.io/gorm"
)

// HandleMySQLError 处理 mysql error 返回一个 errs.AppError 的指针
func HandleMySQLError(err error) *errs.AppError {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NotFound.AsException(err)
		}

		var mysqlErr *mysql.MySQLError

		// 重复键错误, 具体哪个字段重复，由具体业务判断
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errs.ErrAlreadyExists.AsException(err)
		}

		return errs.ServerError.AsException(err)
	}

	return errs.Success
}
