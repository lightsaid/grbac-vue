package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lightsaid/grbac/errs"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/models"
)

const (
	// Session 存储redis的基础可以，一个完整的key是：
	// fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, SessionBaseKey, user.ID)
	SessionBaseKey = "session#"
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

	// 生成签名
	signed := helper.CreateSignature(req.Email)

	user := models.User{
		Name:       req.Name,
		Email:      req.Email,
		VerifyCode: signed,
	}

	if initializer.App.Conf.RunMode == "release" {
		// TODO: 验证邮箱host是否能访问
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

	href := fmt.Sprintf("%s/%s", initializer.App.Conf.ActivateEmailURL, signed)

	// 创建邮件模版
	content, err := user.SetActivateEmailMessage(href)
	if err != nil {
		helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
		return
	}

	payload := models.RegisterMailerPaylod{
		Email:   user.Email,
		Content: content,
	}

	initializer.App.SubPubRabbitMQ.PublishPubSubCtx(c, payload.String())

	// 发送邮件
	// sender := mailer.NewGmailSender(
	// 	initializer.App.Conf.MailSenderName,
	// 	initializer.App.Conf.MailSenderAddress,
	// 	initializer.App.Conf.MailSenderPassword,
	// )

	// err = sender.SendEmail(
	// 	"账户激活",
	// 	content,
	// 	[]string{user.Email},
	// 	nil,
	// 	nil,
	// 	nil,
	// )
	// if err != nil {
	// 	helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
	// 	return
	// }

	helper.ToResponse(c, nil, "注册成功，请注意查收邮件激活用户")
}

func ActivateUser(c *gin.Context) {
	var req models.ActivateUserRequest
	if ok := helper.BindRequestUri(c, &req); !ok {
		return
	}
	var user models.User
	if err := initializer.DB.Where("verify_code = ?", req.VerifyCode).First(&user).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e.AsMessage("激活失败，查询用户错误"))
		return
	}

	if user.ActivatedAt != nil {
		helper.ToResponse(c, "你已经激活了")
		return
	}

	result := initializer.DB.Model(&user).Where("verify_code = ?", req.VerifyCode).Update("activated_at", time.Now())
	if result.Error != nil {
		helper.ToErrResponse(c, errs.InternalServerError.AsException(result.Error, "激活失败"))
		return
	}

	helper.ToResponse(c, "激活成功")
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}

	// 查询用户
	var user models.User
	if err := initializer.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	// 匹配密码
	if err := user.ComparePassword(req.Password); err != nil {
		helper.ToErrResponse(c, errs.BadRequest.AsException(err, "账户或密码不匹配"))
		return
	}

	// 生成access_token 和 refresh_token
	accessToken, _, err := helper.GenToken(
		user.ID, initializer.App.Conf.TokenSecret, initializer.App.Conf.AccessTokenDuration)

	refreshToken, payload, err2 := helper.GenToken(
		user.ID, initializer.App.Conf.TokenSecret, initializer.App.Conf.RefreshTokenDuration)

	if err != nil || err2 != nil {
		helper.ToErrResponse(c, errs.BadRequest.AsException(err, "生成Token失败"))
		return
	}

	// 设置 session
	session := models.Session{
		TID:          payload.ID,
		UID:          user.ID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(initializer.App.Conf.RefreshTokenDuration),
		UserAgent:    c.Request.UserAgent(),
		ClientIP:     c.ClientIP(),
	}
	key := fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, SessionBaseKey, user.ID)

	// 存储 session
	err = session.Save(initializer.RedisPool, key)
	if err != nil {
		helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
		return
	}

	rsp := models.LoginResponse{User: user, AccessToken: accessToken, RefreshToken: refreshToken}
	helper.ToResponse(c, rsp)
}

func Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	payload, err := helper.ParseToken(req.RefreshToken, initializer.App.Conf.TokenSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			helper.ToErrResponse(c, errs.BadRequest.AsException(err, "token 过期"))
			return
		}
		fmt.Println("err2: ", err)
		helper.ToErrResponse(c, errs.BadRequest.AsException(err))
		return
	}
	var session models.Session
	key := fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, SessionBaseKey, payload.UserID)
	err = session.Get(initializer.RedisPool, key)
	if err != nil {
		if err == models.ErrSessionNotFound {
			helper.ToErrResponse(c, errs.NotFound)
			return
		}
		helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
		return
	}

	// 查找成功，refreshToken 在redis中还没过期

	if payload.UserID != session.UID || payload.ID != session.TID || req.RefreshToken != session.RefreshToken {
		helper.ToErrResponse(c, errs.Unauthorized)
		return
	}

	// 生成 refresh TOken
	refreshToken, _, err := helper.GenToken(session.UID, initializer.App.Conf.TokenSecret, initializer.App.Conf.RefreshTokenDuration)
	if err != nil {
		helper.ToResponse(c, errs.InternalServerError.AsException(err))
		return
	}

	helper.ToResponse(c, refreshToken, "成功")
}
