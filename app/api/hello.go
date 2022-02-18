package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
)

var Hello = helloApi{}

type helloApi struct{}

type T struct {
	U *ghttp.UploadFile
	I int
}

// Index is a demonstration route handler for output "Hello World!".
func (a *helloApi) Index(r *ghttp.Request) {
	fmt.Println(r.Router.Uri)
}
func (a *helloApi) Index1(r *ghttp.Request) {
	var t T
	err := r.Parse(&t)
	file := r.GetUploadFile("u")
	t.U = file
	if err != nil {
		return
	}
	fmt.Println(t)
}
