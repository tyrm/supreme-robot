package redis

import (
	"errors"
	"github.com/tyrm/supreme-robot/startup"
	"time"

	redisCon "github.com/gomodule/redigo/redis"
)

type Client struct {
	redisAddress  string
	redisDB       int
	redisPassword string

	db *redisCon.Pool
}

var ErrCantConnect = errors.New("can't connect to redis")

func NewClient(cfg *startup.StartupConfig) (*Client, error) {
	client := Client{
		redisAddress:  cfg.RedisDnsAddress,
		redisDB:       cfg.RedisDnsDB,
		redisPassword: cfg.RedisDnsPassword,
	}

	// connect to redis
	client.db = &redisCon.Pool{
		Dial: func() (redisCon.Conn, error) {
			var opts []redisCon.DialOption
			if client.redisPassword != "" {
				opts = append(opts, redisCon.DialPassword(client.redisPassword))
			}
			opts = append(opts, redisCon.DialDatabase(client.redisDB))
			opts = append(opts, redisCon.DialConnectTimeout(100*time.Millisecond))
			opts = append(opts, redisCon.DialReadTimeout(1000*time.Millisecond))

			return redisCon.Dial("tcp", client.redisAddress, opts...)
		},
	}

	logger.Tracef("new redis client created: %s(%d)", cfg.RedisDnsAddress, cfg.RedisDnsDB)
	return &client, nil
}
