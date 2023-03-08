package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lightsaid/grbac-vue/dddgo/app"
	"github.com/lightsaid/grbac-vue/dddgo/domain"
	"github.com/lightsaid/grbac-vue/dddgo/dto"
	"github.com/lightsaid/grbac-vue/dddgo/entity"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/dberr"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/errs"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/signature"
	"github.com/lightsaid/grbac-vue/dddgo/pkg/token"
)

const (
	SessionBaseKey = "session#"
)

type AuthService interface {
	Register(ctx context.Context, req dto.NewUserRequest, secret string) (*dto.UserResponse, *errs.AppError)
	Activate(ctx context.Context, verifyCode string) (*dto.UserResponse, *errs.AppError)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	RefreshToken(ctx context.Context, req *dto.RefreshRequest, payload *token.JwtPayload) (string, *errs.AppError)
}

type authService struct {
	repo     domain.UserRepository
	authRepo domain.AuthRepository
	cache    domain.SessionRepository
}

func NewAuthService(
	repo domain.UserRepository,
	auth domain.AuthRepository,
	cache domain.SessionRepository,
) AuthService {
	return &authService{repo: repo, authRepo: auth, cache: cache}
}

func (a *authService) Register(ctx context.Context, req dto.NewUserRequest, secret string) (*dto.UserResponse, *errs.AppError) {
	var user entity.User
	user.Name = req.Name
	user.Email = req.Email
	signed := signature.CreateSignature(req.Email, secret)

	user.VerifyCode = &signed

	user.SetPassword(req.Password)

	err := a.repo.Create(&user)

	if err != nil {
		return nil, dberr.HandleMySQLError(err)
	}

	href := fmt.Sprintf("%s/%s", app.App.Config.ActivateEmailURL, signed)

	// 创建邮件模版
	content, err := user.SetActivateEmailMessage(href)
	if err != nil {
		return nil, errs.ServerError.AsException(err, "激活邮件发送失败")
	}

	payload := dto.RegisterMailerPaylod{
		Email:   user.Email,
		Content: content,
	}

	app.App.SubPubRabbitMQ.PublishPubSubCtx(ctx, payload.String())

	return user.ToDto(), nil
}

func (a *authService) Activate(ctx context.Context, verifyCode string) (*dto.UserResponse, *errs.AppError) {
	user, err := a.authRepo.ActivatedAccount(verifyCode)
	if err != nil {
		e := dberr.HandleMySQLError(err)
		return nil, e
	}
	return user.ToDto(), nil
}

func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	// 查询用户
	user, err := a.repo.FindByEmail(req.Email)
	if err != nil {
		e := dberr.HandleMySQLError(err)
		return nil, e
	}

	// 匹配密码
	if err := user.ComparePassword(req.Password); err != nil {
		return nil, errs.BadRequest.AsException(err, "账户或密码不匹配")
	}

	// 生成access_token 和 refresh_token
	accessToken, _, err := token.GenToken(
		user.ID, app.App.Config.TokenSecret, app.App.Config.AccessTokenDuration)

	refreshToken, payload, err2 := token.GenToken(
		user.ID, app.App.Config.TokenSecret, app.App.Config.RefreshTokenDuration)

	if err != nil || err2 != nil {
		return nil, errs.BadRequest.AsException(err, "生成Token失败")
	}

	// 设置 session
	session := entity.Session{
		TID:          payload.ID,
		UID:          user.ID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(app.App.Config.RefreshTokenDuration),
		UserAgent:    req.UserAgent,
		ClientIP:     req.ClientIP,
	}
	key := fmt.Sprintf("%s%s%d", app.App.Config.RedisPrefixKey, SessionBaseKey, user.ID)

	// 存储 session
	err = a.cache.Save(session, key)
	if err != nil {
		return nil, errs.ServerError.AsException(err)
	}

	rsp := dto.LoginResponse{User: user.ToDto(), AccessToken: accessToken, RefreshToken: refreshToken}
	return &rsp, nil
}

func (a *authService) RefreshToken(ctx context.Context, req *dto.RefreshRequest, payload *token.JwtPayload) (string, *errs.AppError) {
	key := fmt.Sprintf("%s%s%d", app.App.Config.RedisPrefixKey, SessionBaseKey, payload.UserID)

	session, err := a.cache.Get(key)
	if errors.Is(err, domain.ErrSessionNotFound) {
		return "", errs.NotFound
	}

	// 查找成功，refreshToken 在redis中还没过期
	if payload.UserID != session.UID || payload.ID != session.TID || req.RefreshToken != session.RefreshToken {
		return "", errs.Unauthorized
	}

	// 生成 refresh Token
	refreshToken, _, err := token.GenToken(session.UID, app.App.Config.TokenSecret, app.App.Config.RefreshTokenDuration)
	if err != nil {
		return "", errs.ServerError
	}

	// 覆盖session
	session.RefreshToken = refreshToken
	err = a.cache.Save(*session, key)
	if err != nil {
		log.Println("refresh token save session err: ", err)
	}

	return refreshToken, nil
}
