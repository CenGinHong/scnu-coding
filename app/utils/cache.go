package utils

import (
	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
)

// @Author: 陈健航
// @Date: 2021/5/1 22:52
// @Description:

type MyCache struct {
	*gcache.Cache
}

func NewMyCache() (m *MyCache) {
	cache := gcache.New()
	// 多副本部署下用redis
	if g.Cfg().GetBool("server.IsMultiple") {
		redisAdapter := adapter.NewRedis(g.Redis())
		cache.SetAdapter(redisAdapter)
	}
	m = &MyCache{cache}
	return m
}
