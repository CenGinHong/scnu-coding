package api

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/library/response"
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
	p := r.GetUploadFile("file")
	f := r.GetMultipartFiles("file")
	fmt.Println(f)
	fmt.Println(p)
	response.Succ(r, g.Map{
		"id": 1,
	})
}
func (a *helloApi) Index2(r *ghttp.Request) {
	fmt.Println("receive a stop")
}
