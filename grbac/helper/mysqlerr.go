package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/lightsaid/grbac/errs"
	"gorm.io/gorm"
)

func HandleMySQLError(c *gin.Context, err error) *errs.AppError {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NotFound.AsException(err)
		}

		var mysqlErr *mysql.MySQLError
		// 重复键错误
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errs.AlreadyExist.AsException(err)
		}

		return errs.InternalServerError.AsException(err)
	}

	return errs.StatusOK
}
