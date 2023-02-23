package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lightsaid/grbac/errs"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/models"
)

// Register godoc
// @Summary 用户注册
// @Description 注册，成为系统用户
// @Tags         Auth
// @Accept       json
// @Produce      json
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	err := user.SetPassword(req.Password)
	if err != nil {
		helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
		return
	}
	result := initializer.DB.Create(&user)
	if result.Error != nil {
		e := helper.HandleMySQLError(c, result.Error)
		helper.ToErrResponse(c, e)
		return
	}

	// TODO:
	helper.ToResponse(c, user)
}

func ActivateUser(c *gin.Context) {

}

func Login(c *gin.Context) {

}

func Logout(c *gin.Context) {

}

func Refresh(c *gin.Context) {

}

func ForgotPswd(c *gin.Context) {

}

func RestPswd(c *gin.Context) {

}
