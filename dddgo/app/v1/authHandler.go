package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lightsaid/grbac-vue/dddgo/app"
	"github.com/lightsaid/grbac-vue/dddgo/domain"
	"github.com/lightsaid/grbac-vue/dddgo/dto"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/errs"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/token"
	"github.com/lightsaid/grbac-vue/dddgo/service"
)

const (
	SessionBaseKey = "session#"
)

type AuthHandler struct {
	service service.AuthService
}

// Register godoc
// @Summary 用户注册
// @Description 注册，成为系统用户
// @Tags         Auth
// @Accept       json
// @Produce      json
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.NewUserRequest
	if ok := app.BindRequest(c, &req); !ok {
		return
	}
	_, err := h.service.Register(c.Request.Context(), req, app.App.Config.SignatureSecret)
	if err != nil {
		app.ToErrResponse(c, err)
		return
	}
	app.ToResponse(c, nil, "注册成功，请注意查收邮件激活用户")
}

func (h *AuthHandler) ActivateUserHandler(c *gin.Context) {
	var req dto.ActivateUserRequest
	if ok := app.BindRequestUri(c, &req); !ok {
		return
	}
	_, err := h.service.Activate(c.Request.Context(), req.VerifyCode)
	if err != nil {
		if errors.Is(err, domain.ErrUserIsActived) {
			app.ToErrResponse(c, errs.BadRequest.AsMessage("您已经激活了，无须再次激活"))
			return
		}
		app.ToErrResponse(c, err.AsMessage("激活失败"))
		return
	}
	app.ToResponse(c, nil, "激活成功")
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if ok := app.BindRequest(c, &req); !ok {
		return
	}
	req.UserAgent = c.Request.UserAgent()
	req.ClientIP = c.ClientIP()
	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		app.ToErrResponse(c, err)
		return
	}

	app.ToResponse(c, res)
}

func (h *AuthHandler) RefreshHandler(c *gin.Context) {
	var req dto.RefreshRequest
	if ok := app.BindRequest(c, &req); !ok {
		return
	}
	payload, err := token.ParseToken(req.RefreshToken, app.App.Config.TokenSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			app.ToErrResponse(c, errs.BadRequest.AsException(err, "token 过期"))
			return
		}
		app.ToErrResponse(c, errs.BadRequest.AsException(err))
		return
	}
	refreshToken, e := h.service.RefreshToken(c.Request.Context(), &req, payload)
	if e != nil {
		app.ToErrResponse(c, e)
		return
	}
	app.ToResponse(c, refreshToken)
}
