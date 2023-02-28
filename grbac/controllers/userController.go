package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/lightsaid/grbac/errs"
	"github.com/lightsaid/grbac/helper"
	"github.com/lightsaid/grbac/initializer"
	"github.com/lightsaid/grbac/middleware"
	"github.com/lightsaid/grbac/models"
)

const (
	UserRedisBaseKey = "user#"
	OneHourSecond    = 60 * 60
	OneDaySecond     = 24 * 60 * 60
)

func GetProfile(c *gin.Context) {
	var user models.User
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*helper.JwtPayload)
	conn := initializer.RedisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, UserRedisBaseKey, payload.UserID)
	result, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		// key 不存在情况, 请求数据库， 并做缓存
		if err = initializer.DB.Where("id = ?", payload.UserID).First(&user).Error; err != nil {
			e := helper.HandleMySQLError(c, err)
			helper.ToErrResponse(c, e)
			return
		}
		// 缓存用户信息
		buf, _ := json.Marshal(user)
		if len(buf) > 0 {
			_, err1 := conn.Do("SET", key, string(buf))
			// 随机过期时间，介于1h~24h, 在更新接口需要删除
			d := helper.RandomInt(OneHourSecond, OneDaySecond)
			_, err2 := conn.Do("EXPIRE", key, d)
			if err1 != nil || err2 != nil {
				log.Printf("conn.Do(%q, %s, ...) %s", "SET", key, err1)
				log.Printf("conn.Do(%q, %s, ...) %s, d=%d", "EXPIRE", key, err2, d)
			}
		}
	} else if err != nil {
		helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
		conn.Do("DEL", key)
		return
	} else {
		err = json.Unmarshal(result, &user)
		if err != nil {
			helper.ToErrResponse(c, errs.InternalServerError.AsException(err))
			conn.Do("DEL", key)
			return
		}
	}

	// 响应
	if user.ID > 0 {
		helper.ToResponse(c, user)
	} else {
		helper.ToErrResponse(c, errs.NotFound)
	}
}

func UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	var user models.User
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*helper.JwtPayload)
	if err := initializer.DB.Where("id = ?", payload.UserID).First(&user).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}
	user.Avatar = req.Avatar
	user.Name = req.Name
	if err := initializer.DB.Save(user).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
		return
	}

	// 清理一下缓存
	conn := initializer.RedisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, UserRedisBaseKey, payload.UserID)
	conn.Do("DEL", key)

	helper.ToResponse(c, user)
}

func ListUsers(c *gin.Context) {
	var req models.Pagination
	var users []*models.User
	if ok := helper.BindRequest(c, &req); !ok {
		return
	}
	db := initializer.DB.Offset(req.Offset()).Limit(req.Size)
	if err := db.Find(&users).Error; err != nil {
		e := helper.HandleMySQLError(c, err)
		helper.ToErrResponse(c, e)
	}
	helper.ToResponse(c, users)
}

func Logout(c *gin.Context) {
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*helper.JwtPayload)
	key := fmt.Sprintf("%s%s%d", initializer.App.Conf.RedisPrefixKey, SessionBaseKey, payload.UserID)
	conn := initializer.RedisPool.Get()
	defer conn.Close()
	conn.Do("DEL", key)

	helper.ToResponse(c, nil, "成功")
}
