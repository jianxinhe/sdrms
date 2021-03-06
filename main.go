package main

import (
	_ "github.com/jianxinhe/sdrms/routers"
	_ "github.com/jianxinhe/sdrms/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	// 读取配置文件
	imgPath := beego.AppConfig.String("img_root_path")
	beego.SetStaticPath("/coursesuite", imgPath)
	beego.Run()
}
