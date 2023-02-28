package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JwtPayload JWT Token 附带的数据
type JwtPayload struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT Token
func GenToken(uid uint, secretKey string, d time.Duration) (string, *JwtPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	payload := &JwtPayload{
		uid,
		jwt.RegisteredClaims{
			// 唯一
			ID:        tokenID.String(),
			Issuer:    "grbac",
			Subject:   "",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
		},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := claims.SignedString([]byte(secretKey))

	return tokenString, payload, err
}

// 验证、解析 JWT Token
func ParseToken(tokenStr string, secretKey string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtPayload{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtPayload); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
