package routers

import (
	"CMSdemoByBeego/controllers"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

func init() {
	//设置路由过滤，使用正则匹配，在路由之前，判断session，否则重定向
	beego.InsertFilter("/Artiicle/*", beego.BeforeRouter, filterFunc)
	//定义路由：
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList;post:HandleTypeSelected")
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/content", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/Article/DeleteArticle", &controllers.ArticleController{}, "get:HandleDelete")
	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdate;post:HandleUpdate")
	//定义路由：添加文章类型
	beego.Router("/Article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")
	//定义路由：删除文章类型
	beego.Router("/Article/DeleteArticleType", &controllers.ArticleController{}, "get:HandleDeleteType")
	//定义路由：退出路径
	beego.Router("/logout", &controllers.LogoutController{}, "get:HandleLogout")

}

//检查session，此函数由正则匹配的路由，在寻找路由之前执行
func filterFunc(ctx *context.Context) {
	if name := ctx.Input.Session("username"); name == nil {
		ctx.Redirect(302, "/")
		return
	}

}
