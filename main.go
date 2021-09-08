package main

import (
	_ "scnu-coding/boot"
	_ "scnu-coding/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
