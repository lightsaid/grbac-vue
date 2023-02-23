package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

// InitConfig 加载 .env 配置文件
func InitConfig(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}
