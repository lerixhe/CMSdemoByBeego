package main

import (
	_ "CMSdemoByBeego/models"
	_ "CMSdemoByBeego/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

