package signature

import (
	"crypto/sha256"
	"fmt"
)

// CreateSignature 生成一个签名
// data 要签名的数据
// secret 密钥
func CreateSignature(data string, secret string) string {
	hash := sha256.New()
	hash.Write([]byte(data + secret))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
