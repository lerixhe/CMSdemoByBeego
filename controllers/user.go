package controllers

import (
	"CMSdemoByBeego/models"
	"fmt"
	"time"

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

	//尝试从cookies拿用户名
	name := c.Ctx.GetCookie("username")
	if name != "" {
		//默认勾选
		c.Data["checked"] = "checked"
	}
	c.Data["name"] = name
	c.TplName = "login.html"

}

// HandleLogin 处理登录
func (c *LoginController) HandleLogin() {
	c.TplName = "login.html"

	name := c.GetString("userName")
	password := c.GetString("password")
	remember := c.GetString("remember")
	fmt.Println(name, password, remember)

	//处理用户名
	if name == "" || password == "" {
		fmt.Println("用户名或密码不能为空")
		c.Data["errmsg"] = "用户名或密码不能为空"
		return
	}

	//处理复选框——记住用户名，不必非得登录成功才存储
	if remember == "on" {
		c.Ctx.SetCookie("username", name, 3600*time.Second)
		c.Data["checked"] = "checked"
	} else {
		c.Ctx.SetCookie("username", name, -1)
		c.Data["checked"] = ""
	}
	o := orm.NewOrm()
	user := models.User{Name: name}
	err := o.Read(&user, "name")
	if err != nil {
		fmt.Println("读取错误")
		c.Data["errmsg"] = "发生错误"
		return
	}

	if user.Passwd != password {
		fmt.Println("密码错误")
		c.Data["errmsg"] = "登录失败，密码错误"
		return
	}
	c.Data["errmsg"] = "登录成功"
	//登录成功后，存储到session
	c.SetSession("username", name)
	c.Redirect("/Article/ShowArticle", 302)
}

type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) HandleLogout() {
	//点击注销按钮，删除session
	c.DelSession("username")
	c.Redirect("/", 302)
}
