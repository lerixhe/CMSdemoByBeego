package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id     int
	Name   string
	Passwd string
}

type Article struct {
	Id      int    `orm:"pk;auto"`
	Title   string `orm:"size(20)"`  //文章标题
	Content string `orm:"size(500)"` //文章内容
	Img     string `orm:siza(50)`    //图片路径
	//	Type    string    //文章分类
	CreateTime time.Time `orm:"auto_now;type(datetime)"`    //发布时间
	UpdateTime time.Time `orm:"auto_now_add;type(datetime)` //修改时间
	Count      int       `orm:"default(0)"`                 //阅读量
}

func init() {
	orm.RegisterDataBase("default", "mysql", "mysql:123456@tcp(94.191.18.219:3306)/CMSdb?charset=utf8")
	orm.RegisterModel(new(User), new(Article))
	orm.RunSyncdb("default", false, true)
}
