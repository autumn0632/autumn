package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type redisConf struct {
	Url       string
	MaxIdle   int
	MaxActive int
	Password  string
}

var confRedis = redisConf{
	//Url:       "127.0.0.1:6379",
	Url:       "49.235.235.221:6379",
	MaxIdle:   3,
	MaxActive: 5,
	Password:  "123456",
}

func NewRedis() (redis.Conn, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil, err
	}
	//defer c.Close()
	return c, nil
}

func NewRedisPool(c *redisConf) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   c.MaxIdle,
		MaxActive: c.MaxActive,
		Dial: func() (redis.Conn, error) {
			//1. 连接
			cli, err := redis.DialURL(c.Url)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//2. 密码验证
			if _, authErr := cli.Do("AUTH", c.Password); authErr != nil {
				return nil, fmt.Errorf("redis auth password error: %s", authErr)
			}
			return cli, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func CheckKey(c redis.Conn, k string) bool {
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return exist
	}
}

// 写入的值永远不会过期
func Set(c redis.Conn, k, v string) {
	//defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func SetJson(c redis.Conn, k string, data interface{}) error {
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}


func SetSet(c redis.Conn, key string, data ...string) error{
	var dataSlice = []string{}
	for _, v :=  range data {
		dataSlice = append(dataSlice, v)
	}
	_, err := c.Do("sadd", key, dataSlice)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// 过期方法1：写入值时附加过期时间
func SetEx(c redis.Conn, k, v string, ex int) {
	//defer c.Close()
	_, err := c.Do("SET", k, v, "EX", "5")
	if err != nil {
		log.Printf("%v", err)
	}

}

// 过期方法2：对key值设置过期时间
func SetKeyEx(c redis.Conn, k string, ex int) error {
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		fmt.Println("set error", err.Error())
		return err
	}
	return nil
}

func GetStringValue(c redis.Conn, k string) (string, error) {
	//defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		log.Printf("Get Error: %v", err)
		return "", err
	}
	return username, nil
}
func GetJson(c redis.Conn, k string) (string, error) {
	jsonGet, err := redis.String(c.Do("GET", k))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return jsonGet, nil
}

func GetSet(c redis.Conn, k string) ([]string, error) {
	dataSet, err := redis.Strings(c.Do("SMEMBERS", k))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dataSet, nil
}