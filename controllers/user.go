package controllers

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego"
)

type RegController struct {
	beego.Controller
}

func (c *RegController) ShowReg() {
	c.TplName = "register.html"
}

// HandleReg 处理用户注册请求
/*
1. 取得用户数据
2. 处理用户数据
3. 写入数据库
4. 跳转登录

*/
func (c *RegController) HandleReg() {

	//1. 取得用户数据
	name:=c.GetString("userName")
	password:=c.GetString("password")
	fmt.Println(name,password)

	//2. 处理用户数据
	if name==""||password==""{
		fmt.Println("用户名或密码不能为空！")
		c.TplName="register.html"
	}

	//3. 插入数据库
	o:=orm.NewOrm()
	
	o.Insert(interface{})

}
