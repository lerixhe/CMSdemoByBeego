package routers

import (
	"CMSdemoByBeego/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList")
	beego.Router("/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
}
