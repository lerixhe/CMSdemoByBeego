package controllers

import (
	"CMSdemoByBeego/models"
	"fmt"
	"path"
	"time"

	"github.com/astaxie/beego/orm"

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
func (this *ArticleController) HandleAddArticle() {
	this.TplName = "add.html"

	//取得post数据，使用getfile取得文件，注意设置enctype
	name := this.GetString("articleName")
	content := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		fmt.Println("文件上传失败")
		return
	}
	defer f.Close()
	/*保存之前先做校验处理:
	1.校验文件类型
	2.校验文件大小
	3.防止重名，重新命名
	*/
	ext := path.Ext(h.Filename)
	fmt.Println(ext)
	if ext != ".jpg" && ext != ".png" && ext != "jpeg" {
		fmt.Println("文件类型错误")
		return
	}

	if h.Size > 5000000 {
		fmt.Println("文件超出大小")
		return
	}
	filename := time.Now().Format("20060102150405") + ext

	//保存文件到某路径下，程序默认当前在项目的根目录，故注意相对路径
	err = this.SaveToFile("uploadname", "./static/img/"+filename)
	if err != nil {
		fmt.Println("文件保存失败：", err)
		return
	}
	fmt.Println(name, content)
	fmt.Println(filename)

	//插入数据库
	o := orm.NewOrm()
	article := models.Article{}
	article.Title = name
	article.Content = content
	article.Img = "./static/img/" + filename
	_, err = o.Insert(&article)
	if err != nil {
		fmt.Println("插入错误:", err)
		return
	}
}
