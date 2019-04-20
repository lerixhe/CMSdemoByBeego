package controllers

import (
	"CMSdemoByBeego/models"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// RegController 路由注册页面
type RegController struct {
	beego.Controller
}

// ShowReg 处理注册页面get请求
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
	//c.Ctx.WriteString("注册成功")
	c.Redirect("/", 302)

}

//LoginController 登录路由
type LoginController struct {
	beego.Controller
}

//ShowLogin 路由登录界面的get请求
func (c *LoginController) ShowLogin() {
	c.TplName = "login.html"
}

// HandleLogin 处理登录
func (c *LoginController) HandleLogin() {
	c.TplName = "login.html"
	name := c.GetString("userName")
	password := c.GetString("password")
	if name == "" || password == "" {
		fmt.Println("用户名或密码不能为空")
		c.Data["errmsg"] = "用户名或密码不能为空"
		return
	}
	fmt.Println(name, password)
	o := orm.NewOrm()
	user := models.User{Name: name}
	err := o.Read(&user, "name")
	if err != nil {
		fmt.Println("读取错误")
		c.Data["errmsg"] = "发生错误"
		return
	}

	if user.Passwd == password {
		fmt.Println("密码正确")
		c.Data["errmsg"] = "登录成功"
		c.Redirect("/ShowArticle", 302)
	} else {
		fmt.Println("密码错误")
		c.Data["errmsg"] = "登录失败，密码错误"
	}
}
