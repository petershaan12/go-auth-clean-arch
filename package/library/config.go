package library

import (
	"fmt"
	"log"
	"time"

	"github.com/petershaan12/go-auth-clean-arch/resource/constants"
	"github.com/spf13/viper"
)

type Env struct {
	Env            string `yaml:"env"`
	Port           string `yaml:"port"`
	FileServerPort string `yaml:"fileServerPort"`
	Timezone       string `yaml:"timezone"`

	Database struct {
		Host            string `yaml:"host"`
		Database        string `yaml:"database"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifeTime string `yaml:"connMaxLifeTime"`
		ConnMaxIdleTime string `yaml:"connMaxIdleTime"`
	} `yaml:"database"`

	Paseto struct {
		Key                string `yaml:"key"`
		AccessTokenExpiry  string `yaml:"accessTokenExpiry"`
		RefreshTokenExpiry string `yaml:"refreshTokenExpiry"`
	} `yaml:"paseto"`
}

var (
	EnvGlobal Env
)

func ModuleConfig() Env {

	viper.SetConfigFile("config.yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("cannot read configuration")
	}

	err = viper.Unmarshal(&EnvGlobal)
	if err != nil {
		log.Println("environment can't be loaded: ", err.Error())
	}

	return EnvGlobal
}

func MaxIdleConns() int {
	if !viper.IsSet("database.maxIdleConns") {
		return 3
	}
	return viper.GetInt("database.maxIdleConns")
}

func MaxOpenConns() int {
	if !viper.IsSet("database.maxOpenConns") {
		return 15
	}
	return viper.GetInt("database.maxOpenConns")
}

func ConnMaxLifeTime() time.Duration {
	time := viper.GetString("database.connMaxLifeTime")
	return ParseTimeDuration(time, constants.DefaultConnMaxLifeTime)
}

func ConnMaxIdleTime() time.Duration {
	time := viper.GetString("database.connMaxIdleTime")
	return ParseTimeDuration(time, constants.DefaultConnMaxIdleTime)
}

func DBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		EnvGlobal.Database.Username, EnvGlobal.Database.Password, EnvGlobal.Database.Host, EnvGlobal.Database.Database)
}

func ParseTimeDuration(t string, defaultt time.Duration) time.Duration {
	timeDurr, err := time.ParseDuration(t)
	if err != nil {
		return defaultt
	}
	return timeDurr
}

func AccessTokenExpiry() time.Duration {
	expiry := viper.GetString("paseto.accessTokenExpiry")
	return ParseTimeDuration(expiry, 15*time.Minute) // Default 15 minutes
}

func RefreshTokenExpiry() time.Duration {
	expiry := viper.GetString("paseto.refreshTokenExpiry")
	return ParseTimeDuration(expiry, 7*24*time.Hour) // Default 7 days
}
