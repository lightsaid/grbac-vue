package helper

import (
	"crypto/sha256"
	"fmt"

	"github.com/lightsaid/grbac/initializer"
)

// CreateSignature 生成一个签名
func CreateSignature(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data + initializer.App.Conf.SignatureSecret))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
