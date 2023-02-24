package initializer

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var AppConf *appConfig

type appConfig struct {
	Wait *sync.WaitGroup
}

// InitConfig 加载 .env 配置文件
func InitConfig(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

}
