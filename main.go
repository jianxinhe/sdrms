package main

import (
	_ "github.com/jianxinhe/sdrms/routers"
	_ "github.com/jianxinhe/sdrms/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
