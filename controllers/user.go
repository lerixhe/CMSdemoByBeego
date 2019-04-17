package controllers

import (
	"CMSdemoByBeego/models"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
	name := c.GetString("userName")
	password := c.GetString("password")
	fmt.Println(name, password)

	//2. 处理用户数据
	if name == "" || password == "" {
		fmt.Println("用户名或密码不能为空！")
		c.TplName = "register.html"
		return
	}

	//3. 插入数据库
	o := orm.NewOrm()
	user := models.User{}
	user.Name = name
	user.Passwd = password
	o.Insert(&user)
	c.Ctx.WriteString("注册成功")

}

type LoginController struct {
	beego.Controller
}

func (c *LoginController) ShowLogin() {
	c.TplName = "login.html"
}

func (c *LoginController) HandleLogin() {
	name := c.GetString("userName")
	password := c.GetString("password")
	if name == "" || password == "" {
		fmt.Println("用户名或密码不能为空")
		c.TplName = "login.html"
		return
	}
	fmt.Println(name, password)
	o := orm.NewOrm()
	user := models.User{Name: name}
	err := o.Read(&user, "name")
	if err != nil {
		fmt.Println("读取错误")
		c.TplName = "login.html"
		return
	}

	if user.Passwd == password {
		fmt.Println("密码正确")
		c.Ctx.WriteString("登录成功")
	} else {
		fmt.Println("密码错误")
		c.TplName = "login.html"
	}

}
