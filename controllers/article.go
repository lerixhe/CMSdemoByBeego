package controllers

import (
	"CMSdemoByBeego/models"
	"fmt"
	"math"
	"path"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticleList() {
	o := orm.NewOrm()
	//创建查询器
	qs := o.QueryTable("article")
	var articles []models.Article
	//qs.All(&articles) //select * from article

	//分页实现
	count, err := qs.Count()
	if err != nil {
		fmt.Println("获取记录数错误：", err)
		return
	}

	//定义页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		//若未获取到页码，设置默认页码1
		pageIndex = 1
	}
	//定义每页大小
	pageSize := 3
	//得出开始位置
	start := pageSize * (pageIndex - 1)
	//得出总页数
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//参数1：限制获取的条数，参数2，偏移量，即开始位置
	qs.Limit(pageSize, start).All(&articles)

	//定义页码按钮启用状态
	enablelast, enablenext := true, true
	if pageIndex == 1 {
		enablelast = false
	}
	if pageIndex == pageCount {
		enablenext = false
	}
	this.Data["EnableNext"] = enablenext
	this.Data["EnableLast"] = enablelast
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles
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

	var filename string
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		fmt.Println("文件上传失败:", err)
	} else {
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
		filename = time.Now().Format("20060102150405") + ext

		//保存文件到某路径下，程序默认当前在项目的根目录，故注意相对路径
		err = this.SaveToFile("uploadname", "./static/img/"+filename)
		if err != nil {
			fmt.Println("文件保存失败：", err)
			return
		}
		defer f.Close()
	}

	//插入数据库
	o := orm.NewOrm()
	article := models.Article{}
	article.Title = name
	article.Content = content
	if filename != "" {
		article.Img = "./static/img/" + filename
	}
	_, err = o.Insert(&article)
	if err != nil {
		fmt.Println("插入错误:", err)
		return
	}
	this.Redirect("/ShowArticle", 302)
}
func (this *ArticleController) ShowContent() {
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("获取ID失败：", err)
		return
	}
	content := models.Article{Id: id}
	o := orm.NewOrm()
	err = o.Read(&content)
	if err != nil {
		fmt.Println("查询数据失败：", err)
		return
	}
	//阅读量+1并写回数据库
	content.Count++
	o.Update(&content)
	this.Data["content"] = content
	this.TplName = "content.html"
}
func (this *ArticleController) HandleDelete() {
	/*思路
	1.被点击的url传值
	2.执行对应的删除操作
	*/
	this.TplName = ""
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("获取ID失败：", err)
		return
	}
	article := models.Article{Id: id}
	o := orm.NewOrm()
	_, err = o.Delete(&article)
	if err != nil {
		fmt.Println("删除数据失败：", err)
		return
	}
	//this.TplName = "ShowArticle.html"
	this.Redirect("ShowArticle.html", 302)
}

func (c *ArticleController) ShowUpdate() {
	/*思路
	1. 获取数据，填充数据
	2. 更新数据，更新数据库，返回列表页
	*/
	c.TplName = "update.html"
	id, err := c.GetInt("id")
	if err != nil {
		fmt.Println("id获取失败：", err)
		return
	}
	article := models.Article{Id: id}
	o := orm.NewOrm()
	err = o.ReadForUpdate(&article)
	if err != nil {
		fmt.Println("读取失败：", err)
		return
	}
	c.Data["article"] = article
}

// HandleUpdate 处理更新
func (this *ArticleController) HandleUpdate() {
	this.TplName = "update.html"
	//取得post数据，使用getfile取得文件，注意设置enctype
	name := this.GetString("articleName")
	content := this.GetString("content")
	oldimagepath := this.GetString("oldimagepath")

	var filename string
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("id获取失败：", err)
		return
	}
	article := models.Article{Id: id, Title: name, Content: content, Img: oldimagepath}
	this.Data["article"] = article
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
	} else {
		/*保存之前先做校验处理:
		1.校验文件类型
		2.校验文件大小
		3.防止重名，重新命名
		*/
		ext := path.Ext(h.Filename)
		//fmt.Println(ext)
		if ext != ".jpg" && ext != ".png" && ext != "jpeg" {
			fmt.Println(err)
			this.Data["errmsg"] = "文件类型错误"
			return
		}

		if h.Size > 5000000 {
			fmt.Println(err)
			this.Data["errmsg"] = "文件超出大小"
			return
		}
		filename = time.Now().Format("20060102150405") + ext

		//保存文件到某路径下，程序默认当前在项目的根目录，故注意相对路径
		err = this.SaveToFile("uploadname", "./static/img/"+filename)
		if err != nil {
			fmt.Println("文件保存失败：", err)
			this.Data["errmsg"] = "文件保存失败"
			return
		}
		defer f.Close()
	}

	//若上传了新文件，则使用新文件路径，否则使用旧路径不变
	if filename != "" {
		article.Img = "./static/img/" + filename
	}

	//更新数据库
	o := orm.NewOrm()
	_, err = o.Update(&article, "title", "content", "img", "create_time", "update_time")
	if err != nil {
		fmt.Println("更新错误:", err)
		this.Data["errmsg"] = "更新失败"
		return
	}
	this.Redirect("/ShowArticle", 302)
}

func (c *ArticleController) ShowAddType() {
	c.TplName = "addType.html"
	var types []models.ArticleType
	o := orm.NewOrm()
	o.QueryTable("article_type").All(&types)
	c.Data["types"] = types
}
func (c *ArticleController) HandleAddType() {
	var articleType models.ArticleType
	if articleType.TypeName = c.GetString("typeName"); articleType.TypeName == "" {
		fmt.Println("类型不能为空")
		c.Redirect("/AddArticleType", 302)
		return
	}
	o := orm.NewOrm()
	o.Insert(&articleType)
	c.Redirect("/AddArticleType", 302)
}
func (c *ArticleController) HandleDeleteType() {
	id, err := c.GetInt("id")
	if err != nil {
		fmt.Println("获取ID失败：", err)
		return
	}
	articleType := models.ArticleType{Id: id}
	o := orm.NewOrm()
	o.Delete(&articleType)
	c.Redirect("/AddArticleType", 302)

}
