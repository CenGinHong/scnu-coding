package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

var Hello = helloApi{}

type helloApi struct{}

// Index is a demonstration route handler for output "Hello World!".
func (a *helloApi) Index(r *ghttp.Request) {
	fmt.Println(r.RequestURI)
	router := r.Server.GetRouterArray()
	fmt.Println(router)
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.GetRemoteIp())
	fmt.Println(r.Router)
	type demo struct {
		Id int
	}
	var req *demo
	err := r.Parse(&req)
	if err != nil {
		return
	}
	println(req)
}
func (a *helloApi) Index1(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}
	glog.Info(123)
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			glog.Info(456)
			return
		}
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}

}
