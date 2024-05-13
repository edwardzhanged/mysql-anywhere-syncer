package syncer

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConnectOptions struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type Redis struct {
	connectOptions *RedisConnectOptions

	client *redis.Client
}

var RedisInstance *Redis

func NewRedis(connectOptions *RedisConnectOptions) {
	RedisInstance = &Redis{connectOptions: connectOptions}
}

func (r *Redis) Connect() error {
	redisAddr := fmt.Sprintf("%s:%d", r.connectOptions.Host, r.connectOptions.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: r.connectOptions.Password,
		DB:       r.connectOptions.DB,
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Ping() error {
	_, err := r.client.Ping(r.client.Context()).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
