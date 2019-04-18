package controllers

import (
	"github.com/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticleList() {
	this.TplName = "index.html"
}
func (this *ArticleController) ShowAddArticle() {
	this.TplName = "add.html"
}
