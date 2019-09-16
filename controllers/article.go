package controllers

import (
	"CMSdemoByBeego/models"
	"CMSdemoByBeego/redispool"
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"path"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) ShowArticleList() {
	c.TplName = "index.html"
	o := orm.NewOrm()
	//创建文章表查询器，但不查询
	qs := o.QueryTable("article")
	var articles []models.Article //qs.All(&articles) //select * from article

	//先从redis中读取需要的数据
	fmt.Println("【a】尝试从redis查询数据，若查询成功直接显示")
	articletypes := []models.ArticleType{}
	//1. 从redis连接池中获取1个连接
	conn := redispool.Redisclient.Get()
	defer conn.Close()

	//2. 将类型从redis中取出并打印
	//正常情况下，存进去什么类型，就利用什么类型的回复助手函数。但是自定义类型不支持，需使用字节流存入和取出。
	relbytes, err := redis.Bytes(conn.Do("get", "articletypes"))
	if err != nil {
		fmt.Println("get错误：", err)
	}
	dec := gob.NewDecoder(bytes.NewReader(relbytes))
	dec.Decode(&articletypes)
	if len(articletypes) != 0 {
		fmt.Println("【a】从redis中成功查询到数据", articletypes)
	} else {
		//如果以上操作没有从redis中读出数据，则去数据库中查询，并存入redis
		fmt.Println("【a】redis未读取到缓存，本次去mysql中查询数据")
		o.QueryTable("article_type").All(&articletypes)
		fmt.Println("【a】从mysql中成功查到数据查询数据", articletypes)

		//将文章类型存入redis数据库
		//正常情况下，存进去什么类型，就利用什么类型的回复助手函数。但是自定义类型不支持，需使用字节流存入和取出。
		//首先序列化内容
		/*
			1. 初始化一个buffer类型内存，用来存储编码的结果。（造一个内存卡）
			2. 获取一个编码器对象，并给他刚刚创建的buffer内存（得到一个播放器，把内存卡插进去）
			3. 使用编码器对象的编码方法开始编码，输入参数为要编码的内容，buffer存储编码结果，返回值为是否出错。
		*/
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(articletypes)
		if _, err := conn.Do("set", "articletypes", buffer.String()); err != nil {
			fmt.Println("set错误：", err)
			return
		}
		if _, err := conn.Do("EXPIRE", "articletypes", 120); err != nil {
			fmt.Println("过期时间错误：", err)
			return
		}
		fmt.Println("【a】从将mysql查询dao的数据成功存入redis", articletypes)
	}

	//获取本次查询的页码
	pageIndex, err := c.GetInt("pageIndex")
	if err != nil {
		//若未获取到页码，设置默认页码1
		pageIndex = 1
	}
	//定义每页大小，即本次请求的条数
	pageSize := 6
	//根据以上信息，获取开始查询的位置
	start := pageSize * (pageIndex - 1)

	//使用文章查询器，简单获得记录总数
	count, err := qs.RelatedSel("ArticleType").Count()
	if err != nil {
		fmt.Println("获取记录数错误：", err)
		return
	}
	//根据查询头和查询量，开始查询数据
	//参数1：限制获取的条数，参数2，偏移量，即开始位置
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	//加入文章类型筛选，默认全部,选择类型后，再次筛选。
	selectedtype := c.GetString("select")
	if selectedtype == "" || selectedtype == "全部类型" {
		fmt.Println("本次GET请求全部,未加入select参数,默认全部")
	} else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).Count()
		if err != nil {
			fmt.Println("获取记录数错误：", err)
			return
		}
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).All(&articles)
	}
	//得出总页数
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//定义页码按钮启用状态
	enablelast, enablenext := true, true
	if pageIndex == 1 {
		enablelast = false
	}
	if pageIndex == pageCount {
		enablenext = false
	}
	c.Data["username"] = c.GetSession("username")
	c.Data["typename"] = selectedtype
	c.Data["articletypes"] = articletypes
	c.Data["EnableNext"] = enablenext
	c.Data["EnableLast"] = enablelast
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount
	c.Data["pageIndex"] = pageIndex
	c.Data["articles"] = articles
}
func (c *ArticleController) HandleTypeSelected() {
	selectedtype := c.GetString("select")
	articles := []models.Article{}
	o := orm.NewOrm()
	o.QueryTable("article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", selectedtype).All(&articles)
	c.Data["artciles"] = articles

	//文章类型下拉
	articletypes := []models.ArticleType{}
	o.QueryTable("article_type").All(&articletypes)
	c.Data["articletypes"] = articletypes
	c.Data["username"] = c.GetSession("username")
	c.TplName = "index.html"
}

func (c *ArticleController) ShowAddArticle() {
	//文章类型下拉
	o := orm.NewOrm()
	articletypes := []models.ArticleType{}
	o.QueryTable("article_type").All(&articletypes)
	c.Data["articletypes"] = articletypes
	c.Data["username"] = c.GetSession("username")
	c.TplName = "add.html"
}
func (c *ArticleController) HandleAddArticle() {
	// c.Layout = "layout.html"
	c.TplName = "add.html"

	//取得post数据，使用getfile取得文件，注意设置enctype
	name := c.GetString("articleName")
	content := c.GetString("content")
	//取得上传文件，需判断是否传了文件
	var filename string
	f, h, err := c.GetFile("uploadname")
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

		//保存文件到某路径下，程序默认当前路由的路径，故注意相对路径
		err = c.SaveToFile("uploadname", "../static/img/"+filename)
		if err != nil {
			fmt.Println("文件保存失败：", err)
			return
		}
		defer f.Close()
	}

	o := orm.NewOrm()
	//取得文章类型
	selectedtype := c.GetString("select")
	//利用此类型获取完整对象
	articletype := models.ArticleType{TypeName: selectedtype}
	o.Read(&articletype, "TypeName")
	//已知某个字段，查询所有字段时，如果字段为主键，则可省略，否则必须填列名。

	fmt.Println("aaaaaaaaa:", articletype.Id)
	article := models.Article{Title: name, Content: content, ArticleType: &articletype}
	//根据文件上传情况，判断是否更新路径
	if filename != "" {
		article.Img = "../static/img/" + filename
	}
	//插入数据库

	_, err = o.Insert(&article)
	if err != nil {
		fmt.Println("插入错误:", err)
		return
	}

	c.Redirect("/Article/ShowArticle", 302)
}
func (c *ArticleController) ShowContent() {
	id, err := c.GetInt("id")
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

	/*处理最近浏览,
	1. 首先需确定当前浏览者登录状态,获取浏览者信息
	2. 将浏览者信息插入数据表
	3. 将历史浏览者信息从表中读出，去重，显示*/
	if username := c.GetSession("username"); username != nil {
		user := models.User{Name: username.(string)}
		o.Read(&user, "Name")
		//目的：构造多对多查询器,并执行添加插入方法
		o.QueryM2M(&content, "Users").Add(&user)
	}
	//开始读出历史浏览者信息
	users := []models.User{}
	o.QueryTable("User").Filter("Articles__Article__Id", content.Id).Distinct().All(&users)
	c.Data["users"] = users
	c.Data["content"] = content
	c.Data["username"] = c.GetSession("username")
	c.TplName = "content.html"
}
func (c *ArticleController) HandleDelete() {
	/*思路
	1.被点击的url传值
	2.执行对应的删除操作
	*/
	c.TplName = ""
	id, err := c.GetInt("id")
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
	//c.TplName = "ShowArticle.html"
	c.Redirect("/Article/ShowArticle", 302)
}

func (c *ArticleController) ShowUpdate() {
	/*思路
	1. 获取数据，填充数据
	2. 更新数据，更新数据库，返回列表页
	*/
	// c.Layout = "layout.html"
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
	c.Data["username"] = c.GetSession("username")
}

// HandleUpdate 处理更新
func (c *ArticleController) HandleUpdate() {

	c.TplName = "update.html"
	//取得post数据，使用getfile取得文件，注意设置enctype
	name := c.GetString("articleName")
	content := c.GetString("content")
	oldimagepath := c.GetString("oldimagepath")

	var filename string
	id, err := c.GetInt("id")
	if err != nil {
		fmt.Println("id获取失败：", err)
		return
	}
	article := models.Article{Id: id, Title: name, Content: content, Img: oldimagepath}
	c.Data["article"] = article
	f, h, err := c.GetFile("uploadname")
	if err != nil {
		c.Data["errmsg"] = "文件上传失败"
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
			c.Data["errmsg"] = "文件类型错误"
			return
		}

		if h.Size > 5000000 {
			fmt.Println(err)
			c.Data["errmsg"] = "文件超出大小"
			return
		}
		filename = time.Now().Format("20060102150405") + ext

		//保存文件到某路径下，程序默认当前在项目的根目录，故注意相对路径
		err = c.SaveToFile("uploadname", "./static/img/"+filename)
		if err != nil {
			fmt.Println("文件保存失败：", err)
			c.Data["errmsg"] = "文件保存失败"
			return
		}
		defer f.Close()
	}

	//若上传了新文件，则使用新文件路径，否则使用旧路径不变
	if filename != "" {
		article.Img = "../static/img/" + filename
	}

	//更新数据库
	o := orm.NewOrm()
	_, err = o.Update(&article, "title", "content", "img", "create_time", "update_time")
	if err != nil {
		fmt.Println("更新错误:", err)
		c.Data["errmsg"] = "更新失败"
		return
	}
	c.Redirect("/Article/ShowArticle", 302)
}

func (c *ArticleController) ShowAddType() {
	c.TplName = "addType.html"
	var types []models.ArticleType
	o := orm.NewOrm()
	o.QueryTable("article_type").All(&types)
	c.Data["types"] = types
	c.Data["username"] = c.GetSession("username")
	//刷新页面时更新缓存。
	err := updateRedisDate("set", "articletypes", types, 300)
	if err != nil {
		fmt.Println("更新缓存失败：", err)
	}

}

//处理更新redis的功能函数:将自定义类型变量序列化存储到redis
//handlestr为操作名：如get set等
//key为redis中的key
//cont为需要序列号写入的自定义类型变量，需要传指针类型
//time为更新后的过期时间（秒），-1代表永不过期
func updateRedisDate(handlestr string, key string, cont interface{}, time int) error {
	fmt.Println("【b】准备序列化：", cont)
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(cont)
	if err != nil {
		return err
	}
	//	fmt.Println("【b】准备写入redis", buffer.String())
	conn := redispool.Redisclient.Get()
	_, err = conn.Do(handlestr, key, buffer.String())
	if err != nil {
		return err
	}
	fmt.Println("time")
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	fmt.Println("xier")
	return nil

}
func (c *ArticleController) HandleAddType() {
	var articleType models.ArticleType
	if articleType.TypeName = c.GetString("typeName"); articleType.TypeName == "" {
		fmt.Println("类型不能为空")
		c.Redirect("/Article/AddArticleType", 302)
		return
	}
	fmt.Println("您输入的类型名为：", articleType.Id, articleType.TypeName)
	o := orm.NewOrm()
	_, err := o.Insert(&articleType)
	if err != nil {
		fmt.Println("插入数据失败：", err)
		return
	}
	c.Redirect("/Article/AddArticleType", 302)

	//插入数据库成功后，此处不更新缓存，否则需要再次请求所有类型，刷新页面时更新更合适。

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
	c.Redirect("/Article/AddArticleType", 302)

}
