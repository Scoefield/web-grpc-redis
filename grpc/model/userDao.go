package model

import (
	"github.com/garyburd/redigo/redis"
	"practicProject/myTest/redisDemo/web-grpc-redis/grpc/config"
)

var	pool *redis.Pool


func GetRedisConn() redis.Conn {
	return pool.Get()
}

func init() {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisAddress)
		},
		TestOnBorrow:    nil,
		MaxIdle:         8,
		MaxActive:       0,
		IdleTimeout:     100,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}




