package config

import (
	"errors"
	"fmt"
	"github.com/tyrm/supreme-robot/util"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ExtHostname string

	LoggerConfig string

	PostgresDsn string

	RedisDnsAddress     string
	RedisDnsDB          int
	RedisDnsPassword    string
	RedisWebappAddress  string
	RedisWebappDB       int
	RedisWebappPassword string

	Secret string
}

func CollectConfig(requiredVars []string) (*Config, error) {
	var config Config

	// EXT_HOSTNAME
	config.ExtHostname = os.Getenv("EXT_HOSTNAME")
	if config.ExtHostname != "" {
		requiredVars = util.FastPopString(requiredVars, "EXT_HOSTNAME")
	}

	// LOG_LEVEL
	config.LoggerConfig = os.Getenv("LOG_LEVEL")
	if config.LoggerConfig == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(config.LoggerConfig))
	}

	// POSTGRES_DSN
	config.PostgresDsn = os.Getenv("POSTGRES_DSN")
	if config.PostgresDsn != "" {
		requiredVars = util.FastPopString(requiredVars, "POSTGRES_DSN")
	}

	// REDIS_DNS_ADDRESS
	config.RedisDnsAddress = os.Getenv("REDIS_DNS_ADDRESS")
	if config.RedisDnsAddress != "" {
		requiredVars = util.FastPopString(requiredVars, "REDIS_DNS_ADDRESS")
	}

	// REDIS_DNS_DB
	if os.Getenv("REDIS_DNS_DB") == "" {
		config.RedisDnsDB = 0
	} else {
		var err error
		config.RedisDnsDB, err = strconv.Atoi(os.Getenv("REDIS_DNS_DB"))
		if err != nil {
			return nil, err
		}
	}

	// REDIS_DNS_PASSWORD
	config.RedisDnsPassword = os.Getenv("REDIS_DNS_PASSWORD")

	// REDIS_WEBAPP_ADDRESS
	config.RedisWebappAddress = os.Getenv("REDIS_WEBAPP_ADDRESS")
	if config.RedisWebappAddress != "" {
		requiredVars = util.FastPopString(requiredVars, "REDIS_WEBAPP_ADDRESS")
	}

	// REDIS_WEBAPP_DB
	if os.Getenv("REDIS_WEBAPP_DB") == "" {
		config.RedisWebappDB = 0
	} else {
		var err error
		config.RedisWebappDB, err = strconv.Atoi(os.Getenv("REDIS_WEBAPP_DB"))
		if err != nil {
			return nil, err
		}
	}

	// REDIS_WEBAPP_PASSWORD
	config.RedisWebappPassword = os.Getenv("REDIS_WEBAPP_PASSWORD")

	// SECRET
	config.Secret = os.Getenv("SECRET")
	if config.Secret != "" {
		requiredVars = util.FastPopString(requiredVars, "SECRET")
	}

	// Validation
	if len(requiredVars) > 0 {
		return nil, errors.New(fmt.Sprintf("Environment variables missing: %v", requiredVars))
	}

	return &config, nil
}
