package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	PostgresSQL struct {
		//Host string `env:"PSQL_HOST"  env-default:"db"`
		Host     string `yaml:"host" env:"PSQL_HOST"  env-default:"localhost"`
		Port     int    `yaml:"port" env:"PSQL_PORT"  env-default:"5432"`
		Username string `yaml:"username" env:"PSQL_USERNAME"  env-default:"postgres"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-default:"postgres"`
		Database string `yaml:"database" env:"PSQL_DATABASE"  env-default:"postgres"`
		SSLMode  string `yaml:"sslmode" env:"SSL_MODE" env-default:"disable"`
	} `yaml:"database"`
	App struct {
		Port      int    `yaml:"port" env:"APP_PORT" env-default:"8000"`
		JwtSecret string `yaml:"jwtSecret" env:"APP_JWT_SECRET" env-default:""`
	} `yaml:"app"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatalln(err)
		}
		//if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
		//	helpText := "sh1neqd - car service"
		//	help, _ := cleanenv.GetDescription(instance, &helpText)
		//	log.Print(help)
		//	log.Fatal(err)
		//}
	})
	return instance
}
