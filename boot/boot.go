package boot

import (
	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/frame/g"
	_ "scnu-coding/packed"
)

func init() {
	if g.Cfg().GetBool("server.IsMultiple") {
		redisAdapter := adapter.NewRedis(g.Redis())
		g.DB().GetCache().SetAdapter(redisAdapter)
	}
}
