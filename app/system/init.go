package system

import (
	"sync"

	"github.com/jinzhu/configor"
)

var (
	configOnce sync.Once
	configInst AppConfig = AppConfig{}
)

func GetConfig() *AppConfig {
	configOnce.Do(func() {
		err := configor.Load(&configInst)
		if err != nil {
			panic(err)
		}
	})
	return &configInst
}

type AppConfig struct {
	Username    string `required:"true" env:"DB_USERNAME"`
	Password    string `required:"true" env:"DB_PASSWORD"`
	Hostname    string `required:"true" env:"DB_HOST"`
	Port        string `required:"true" env:"DB_PORT"`
	Environment string `required:"true" env:"ENVIRONMENT"`
	BucketName  string `required:"true" env:"BUCKET_NAME"`
	PrivateKey  string `required:"true" env:"PRIVATE_KEY"`
	PublicKey   string `required:"true" env:"PUBLIC_KEY"`
}

func Initialize() {
	GetConfig()
}
