package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id     int
	Name   string
	Passwd string
}

type Article struct {
	Id      int
	Title   string    //文章标题
	Content string    //文章内容
	Img     string    //图片路径
	Type    string    //文章分类
	Time    time.Time //发布时间
	Count   int       //阅读量
}

func init() {
	orm.RegisterDataBase("default", "mysql", "mysql:123456@tcp(94.191.18.219:3306)/CMSdb?charset=utf8")
	orm.RegisterModel(new(User), new(Article))
	orm.RunSyncdb("default", false, true)
}
