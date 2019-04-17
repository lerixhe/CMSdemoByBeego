package routers

import (
	"CMSdemoByBeego/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
}
