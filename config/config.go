package config

import (
	"errors"
	"fmt"
	"github.com/tyrm/supreme-robot/util"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config hold config collected from the environment.
type Config struct {
	AccessSecret     string
	AccessExpiration time.Duration

	ExtHostname string

	LoggerConfig string

	PostgresDsn string

	PrimaryNS string

	RedisDnsAddress     string
	RedisDnsDB          int
	RedisDnsPassword    string
	RedisWebappAddress  string
	RedisWebappDB       int
	RedisWebappPassword string

	RefreshSecret     string
	RefreshExpiration time.Duration
}

// CollectConfig will gather configuration from env vars and return a Config object
func CollectConfig(requiredVars []string) (*Config, error) {
	var config Config

	// ACCESS_SECRET
	config.AccessSecret = os.Getenv("ACCESS_SECRET")
	if config.AccessSecret != "" {
		requiredVars = util.FastPopString(requiredVars, "ACCESS_SECRET")
	}

	// ACCESS_EXP
	if os.Getenv("ACCESS_EXP") == "" {
		config.AccessExpiration = time.Minute * 15
	} else {
		exp, err := strconv.Atoi(os.Getenv("ACCESS_EXP"))
		if err != nil {
			return nil, err
		}
		config.AccessExpiration = time.Second * time.Duration(exp)
	}

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

	// PRIMARY_NS
	config.PrimaryNS = os.Getenv("PRIMARY_NS")
	if config.PrimaryNS != "" {
		requiredVars = util.FastPopString(requiredVars, "PRIMARY_NS")
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

	// REFRESH_SECRET
	config.RefreshSecret = os.Getenv("REFRESH_SECRET")
	if config.RefreshSecret != "" {
		requiredVars = util.FastPopString(requiredVars, "REFRESH_SECRET")
	}

	// REFRESH_EXP
	if os.Getenv("REFRESH_EXP") == "" {
		config.RefreshExpiration = time.Hour * 24 * 7
	} else {
		exp, err := strconv.Atoi(os.Getenv("REFRESH_EXP"))
		if err != nil {
			return nil, err
		}
		config.RefreshExpiration = time.Second * time.Duration(exp)
	}

	// Validation
	if len(requiredVars) > 0 {
		return nil, errors.New(fmt.Sprintf("Environment variables missing: %v", requiredVars))
	}

	return &config, nil
}
