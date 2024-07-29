package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	//Environment string `mapstructure:"ENVIRONMENT"`
	//Host        string `mapstructure:"HOST"`

	DBUsername      string `mapstructure:"DB_USERNAME"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBName          string `mapstructure:"DB_NAME"`
	TezosAPIURL     string `mapstructure:"TEZOS_API_URL"`
	PollingInterval string `mapstructure:"POLLING"`
	Port            string `mapstructure:"PORT"`
	RateLimit       int    `mapstructure:"RATE_LIMIT"`
	MigrationPath   string `mapstructure:"MIGRATION_PATH"`
	//DBNameTest    string `mapstructure:"DB_DBNAME_TEST"`
	DBRecreate bool `mapstructure:"DB_RECREATE"`
}

func LoadConfig() (config *Config) {
	rootPath, _ := GetRootPath()

	viper.SetConfigFile(filepath.Join(rootPath, ".env"))

	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	return
}

// GetRootPath
func GetRootPath() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	currDir := filepath.Dir(b)
	for {
		if PathExists(filepath.Join(currDir, "go.mod")) {
			return currDir, nil
		}

		parentDir := filepath.Dir(currDir)
		if parentDir == currDir {
			break
		}
		currDir = parentDir
	}
	return "", os.ErrNotExist
}

// PathExists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
