package config

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCollectConfig_Empty(t *testing.T) {
	unsetEnv()

	cfg, err := CollectConfig([]string{})
	if err != nil {
		t.Errorf("enexpected error, got: %v, want: nil.", err.Error())
	}

	if cfg.AccessSecret != "" {
		t.Errorf("enexpected config value for AccessSecret, got: '%s', want: ''.", cfg.AccessSecret)
	}
	if cfg.AccessExpiration != time.Minute*15 {
		t.Errorf("enexpected config value for AccessExpiration, got: '%s', want: '%s'.", cfg.AccessExpiration, time.Minute*15)
	}
	if cfg.ExtHostname != "" {
		t.Errorf("enexpected config value for ExtHostname, got: '%s', want: ''.", cfg.ExtHostname)
	}
	if cfg.HttpPort != ":5000" {
		t.Errorf("enexpected config value for ExtHostname, got: '%s', want: ':5000'.", cfg.HttpPort)
	}
	if cfg.LoggerConfig != "<root>=INFO" {
		t.Errorf("enexpected config value for LoggerConfig, got: '%s', want: '<root>=INFO'.", cfg.LoggerConfig)
	}
	if cfg.PostgresDsn != "" {
		t.Errorf("enexpected config value for PostgresDsn, got: '%s', want: ''.", cfg.PostgresDsn)
	}
	if cfg.PrimaryNS != "" {
		t.Errorf("enexpected config value for PrimaryNS, got: '%s', want: ''.", cfg.PrimaryNS)
	}
	if cfg.RedisDNSAddress != "" {
		t.Errorf("enexpected config value for RedisDNSAddress, got: '%s', want: ''.", cfg.RedisDNSAddress)
	}
	if cfg.RedisDNSDB != 0 {
		t.Errorf("enexpected config value for RedisDNSDB, got: %d, want: 0.", cfg.RedisDNSDB)
	}
	if cfg.RedisDNSPassword != "" {
		t.Errorf("enexpected config value for RedisDNSPassword, got: '%s', want: ''.", cfg.RedisDNSPassword)
	}
	if cfg.RedisWebappAddress != "" {
		t.Errorf("enexpected config value for RedisWebappAddress, got: '%s', want: ''.", cfg.RedisWebappAddress)
	}
	if cfg.RedisWebappDB != 0 {
		t.Errorf("enexpected config value for RedisWebappDB, got: %d, want: 0.", cfg.RedisWebappDB)
	}
	if cfg.RedisWebappPassword != "" {
		t.Errorf("enexpected config value for RedisDNSPassword, got: '%s', want: ''.", cfg.RedisWebappPassword)
	}
	if cfg.RefreshSecret != "" {
		t.Errorf("enexpected config value for RefreshSecret, got: '%s', want: ''.", cfg.RefreshSecret)
	}
	if cfg.RefreshExpiration != time.Hour*24*7 {
		t.Errorf("enexpected config value for RefreshExpiration, got: '%s', want: '%s'.", cfg.RefreshExpiration, time.Hour*24*7)
	}
}

func TestCollectConfig_EmptyRequireAll(t *testing.T) {
	unsetEnv()

	requiredEnvVars := []string{
		"ACCESS_SECRET",
		"EXT_HOSTNAME",
		"POSTGRES_DSN",
		"PRIMARY_NS",
		"REDIS_DNS_ADDRESS",
		"REDIS_WEBAPP_ADDRESS",
		"REFRESH_SECRET",
	}

	cfg, err := CollectConfig(requiredEnvVars)
	if err == nil {
		t.Errorf("expected error, got: nil, want: err.")
	}

	if cfg != nil {
		t.Errorf("expected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_InvalidAccessExpiration(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"ACCESS_EXP": "astring",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})
	if err == nil {
		t.Errorf("expected error, got: nil, want: error.")
	}
	if cfg != nil {
		t.Errorf("expected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_InvalidHttPort(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"HTTP_PORT": "astring",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})
	if err == nil {
		t.Errorf("expected error, got: nil, want: error.")
	}
	if cfg != nil {
		t.Errorf("expected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_InvalidRedisDNSDB(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"REDIS_DNS_DB": "astring",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})
	if err == nil {
		t.Errorf("expected error, got: nil, want: error.")
	}
	if cfg != nil {
		t.Errorf("expected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_InvalidRedisWebappDB(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"REDIS_WEBAPP_DB": "astring",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})
	if err == nil {
		t.Errorf("expected error, got: nil, want: error.")
	}
	if cfg != nil {
		t.Errorf("expected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_InvalidRefreshExpiration(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"REFRESH_EXP": "astring",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})

	if err == nil {
		t.Errorf("expected error, got: nil, want: error.")
	}
	if cfg != nil {
		t.Errorf("unexpected config, got: %v, want: nil.", cfg)
	}
}

func TestCollectConfig_Loaded(t *testing.T) {
	unsetEnv()

	setEnvVars := map[string]string{
		"ACCESS_SECRET":         "secret1",
		"ACCESS_EXP":            "300",
		"EXT_HOSTNAME":          "www.bubbles.com",
		"HTTP_PORT":             "9876",
		"LOG_LEVEL":             "trace",
		"POSTGRES_DSN":          "postgresql://test:test@127.0.0.1:5432/test",
		"PRIMARY_NS":            "ns1.ptzo.gdn.",
		"REDIS_DNS_ADDRESS":     "localhost:6379",
		"REDIS_DNS_DB":          "8",
		"REDIS_DNS_PASSWORD":    "P@ssw0rd!",
		"REDIS_WEBAPP_ADDRESS":  "redis.exmple.com:5296",
		"REDIS_WEBAPP_DB":       "5",
		"REDIS_WEBAPP_PASSWORD": "something something",
		"REFRESH_SECRET":        "secret2",
		"REFRESH_EXP":           "123456789",
	}

	for k, v := range setEnvVars {
		err := os.Setenv(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := CollectConfig([]string{})
	if err != nil {
		t.Errorf("enexpected error, got: %v, want: nil.", err.Error())
	}

	if cfg.AccessSecret != "secret1" {
		t.Errorf("enexpected config value for AccessSecret, got: '%s', want: 'secret1'.", cfg.AccessSecret)
	}
	if cfg.AccessExpiration != time.Second*300 {
		t.Errorf("enexpected config value for AccessExpiration, got: '%s', want: '%s'.", cfg.AccessExpiration, time.Second*300)
	}
	if cfg.ExtHostname != "www.bubbles.com" {
		t.Errorf("enexpected config value for ExtHostname, got: '%s', want: 'www.bubbles.com'.", cfg.ExtHostname)
	}
	if cfg.HttpPort != ":9876" {
		t.Errorf("enexpected config value for ExtHostname, got: '%s', want: ':9876'.", cfg.ExtHostname)
	}
	if cfg.LoggerConfig != "<root>=TRACE" {
		t.Errorf("enexpected config value for LoggerConfig, got: '%s', want: '<root>=TRACE'.", cfg.LoggerConfig)
	}
	if cfg.PostgresDsn != "postgresql://test:test@127.0.0.1:5432/test" {
		t.Errorf("enexpected config value for PostgresDsn, got: '%s', want: 'postgresql://test:test@127.0.0.1:5432/test'.", cfg.PostgresDsn)
	}
	if cfg.PrimaryNS != "ns1.ptzo.gdn." {
		t.Errorf("enexpected config value for PrimaryNS, got: '%s', want: 'ns1.ptzo.gdn.'.", cfg.PrimaryNS)
	}
	if cfg.RedisDNSAddress != "localhost:6379" {
		t.Errorf("enexpected config value for RedisDNSAddress, got: '%s', want: 'localhost:6379'.", cfg.RedisDNSAddress)
	}
	if cfg.RedisDNSDB != 8 {
		t.Errorf("enexpected config value for RedisDNSDB, got: %d, want: 8.", cfg.RedisDNSDB)
	}
	if cfg.RedisDNSPassword != "P@ssw0rd!" {
		t.Errorf("enexpected config value for RedisDNSPassword, got: '%s', want: 'P@ssw0rd!'.", cfg.RedisDNSPassword)
	}
	if cfg.RedisWebappAddress != "redis.exmple.com:5296" {
		t.Errorf("enexpected config value for RedisWebappAddress, got: '%s', want: 'redis.exmple.com:5296'.", cfg.RedisWebappAddress)
	}
	if cfg.RedisWebappDB != 5 {
		t.Errorf("enexpected config value for RedisWebappDB, got: %d, want: 5.", cfg.RedisWebappDB)
	}
	if cfg.RedisWebappPassword != "something something" {
		t.Errorf("enexpected config value for RedisDNSPassword, got: '%s', want: 'something something'.", cfg.RedisWebappPassword)
	}
	if cfg.RefreshSecret != "secret2" {
		t.Errorf("enexpected config value for RefreshSecret, got: '%s', want: 'secret2'.", cfg.RefreshSecret)
	}
	if cfg.RefreshExpiration != time.Second*123456789 {
		t.Errorf("enexpected config value for RefreshExpiration, got: '%s', want: '%s'.", cfg.RefreshExpiration, time.Second*123456789)
	}
}

func unsetEnv() {
	envVars := []string{
		"ACCESS_SECRET",
		"ACCESS_EXP",
		"EXT_HOSTNAME",
		"HTTP_PORT",
		"LOG_LEVEL",
		"POSTGRES_DSN",
		"PRIMARY_NS",
		"REDIS_DNS_ADDRESS",
		"REDIS_DNS_DB",
		"REDIS_DNS_PASSWORD",
		"REDIS_WEBAPP_ADDRESS",
		"REDIS_WEBAPP_DB",
		"REDIS_WEBAPP_PASSWORD",
		"REFRESH_SECRET",
		"REFRESH_EXP",
	}

	for _, ev := range envVars {
		err := os.Unsetenv(ev)
		if err != nil {
			log.Fatal(err)
		}
	}
}
