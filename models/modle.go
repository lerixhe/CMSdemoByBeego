package models

import (
	"github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
)

type User struct{
	Id int
	Name string
	Passwd string
}

func init(){
	orm.RegisterDataBase("default", "mysql", "mysql:123456@tcp(94.191.18.219:3306)/CMSdb?charset=utf8")
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
}