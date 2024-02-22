package main

import (
	"log"
	"sfvn-hbduy/api"
	"sfvn-hbduy/common/cache"
	"sfvn-hbduy/service"

	"github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type Config struct {
	Dir  string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port string
}

var config Config

func init() {
	if err := env.Parse(&config); err != nil {
		log.Panic("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Port: viper.GetString(`main.port`),
	}

	config = cfg
}

func main() {
	cache.MCache = cache.NewMemCache()
	defer cache.MCache.Close()

	server := api.NewServer()
	baseUrl := viper.GetString("coingecko.url")
	apiKey := viper.GetString("coingecko.api_key")

	historiesService := service.NewHistories(baseUrl, apiKey)
	api.APIHistoriesHandler(server.Engine, historiesService)
	server.Start(config.Port)
}
