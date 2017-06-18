package main

import (
	_ "./routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
)

func main() {
	beego.SetStaticPath("/plugins", "plugins")
	beego.SetStaticPath("/upload", "upload")
	
	beego.Run()
}