package config

import (
	"bytes"
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const CONF_DIR = "config/"

//go:embed *.yaml
var configs embed.FS

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	// env := os.Getenv("ENV")
	env, exists := os.LookupEnv("ENV")
	if !exists {
		panic("ENV is not set")
	}
	vp := viper.New()
	// 根据环境变量 ENV 决定要读取的应用启动配置
	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(err)
	}
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewBuffer(configFileStream))
	if err != nil {
		// 加载不到配置，阻挡应用的继续启动
		panic(err)
	}
	vp.UnmarshalKey("app", &App)
	vp.UnmarshalKey("database", &Database)
	vp.UnmarshalKey("redis", &Redis)
}
