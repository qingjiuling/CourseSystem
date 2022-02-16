package db_op

import (
	"github.com/garyburd/redigo/redis"
)

var RedisDb *redis.Pool //创建redis连接池

func RedisInit() {
	pool := &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			c, err := redis.Dial("tcp", "180.184.70.161:6379")
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", "bytedance22"); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}
	RedisDb = pool
}

func RedisClose() {
	RedisDb.Close()
}
