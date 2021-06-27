package main

import (
	_ "assets/boot"
	_ "assets/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
