package utils

import (
	"fmt"
	"sync"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/guid"
)

// @Author: 陈健航
// @Date: 2021/5/1 22:28
// @Description:

type MyMutex interface {
	Lock()
	UnLock()
}

func NewMyMutex() (m MyMutex) {
	if g.Cfg().GetBool("server.IsMultiple") {
		m = &RedisMutex{redisLock.NewMutex(guid.S(), nil)}
	} else {
		m = &SyncMutex{&sync.Mutex{}}
	}
	return m
}

type SyncMutex struct {
	*sync.Mutex
}

func (receiver *SyncMutex) Lock() {
	receiver.Mutex.Lock()
}

func (receiver *SyncMutex) UnLock() {
	receiver.Mutex.Unlock()
}

type RedisMutex struct {
	*redsync.Mutex
}

func (receiver *RedisMutex) Lock() {
	_ = receiver.Mutex.Lock()

}

func (receiver *RedisMutex) UnLock() {
	_, _ = receiver.Mutex.Unlock()
}

var redisLock = newRedisLock()

func newRedisLock() (r redsync.Redsync) {
	redisConfig := g.Cfg().GetString("redis.default")
	split := gstr.Split(redisConfig, ":")
	host := split[0]
	split = gstr.Split(split[1], ",")
	port, pass := split[0], split[2]
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
	})
	pool := goredis.NewPool(client)
	return *redsync.New(pool)
}
