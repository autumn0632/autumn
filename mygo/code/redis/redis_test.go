package redis

import (
	"log"
	"testing"
	"time"
)

func TestNewRedis(t *testing.T) {
	c, err := NewRedis()
	if err != nil {
		log.Printf("%v", err)
	}
	defer c.Close()
	c.Do("AUTH", "123456")

	//1. set string
	Set(c, "time", "2011-11-1 10:10:00")
	v, err := GetStringValue(c, "time")
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", v)


	//2. set json
	info := map[string]string{
		"name": "bill",
		"age":  "18",
	}
	SetJson(c, "info", info)
	infoJson, err := GetJson(c, "info")
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", infoJson)

	//3. 集合操作
	setList := []string{"你好", "嘛", "美", "女"}
	err = SetSet(c, "test", setList...)
	if err != nil {
		log.Printf("set error: %v", err)
	}
	setInfo, err := GetSet(c, "test")
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("get set info: %v", setInfo)
}

func TestNewRedisPool(t *testing.T) {
	confRedis := &redisConf{
		Url:       "redis://127.0.0.1:6379",
		MaxIdle:   3,
		MaxActive: 5,
		Password:  "123456",
	}

	p := NewRedisPool(confRedis)
	err := p.TestOnBorrow(p.Get(), time.Now())
	if err != nil {
		log.Printf("%v", err)
	}
	Set(p.Get(), "name", "bill")
	v, err := GetStringValue(p.Get(), "name")
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%v", v)
}
