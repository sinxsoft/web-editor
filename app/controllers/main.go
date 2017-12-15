package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/satori/go.uuid"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/sinxsoft/web-editor/app/models"
)

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	this.Data["pageTitle"] = "webpage编辑器"

	// 已经新建好的好的页面
	entries, _ := models.ContentGetAllListExt()

	this.Data["entries"] = entries

	this.display()
}

// 个人信息
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user, _ := models.UserGetById(this.userId)

	if this.isPost() {
		flash := beego.NewFlash()
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("密码长度必须大于6位")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("两次输入的密码不一致")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		flash.Success("修改成功！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "个人信息"
	this.Data["user"] = user
	this.display()
}

// 登录
// func (this *MainController) Login() {
// 	if this.userId > 0 {
// 		this.redirect("/")
// 	}
// 	//fmt.Println("asdfadsf")
// 	beego.ReadFromRequest(&this.Controller)
// 	if this.isPost() {
// 		flash := beego.NewFlash()

// 		username := strings.TrimSpace(this.GetString("username"))
// 		password := strings.TrimSpace(this.GetString("password"))
// 		remember := this.GetString("remember")
// 		if username != "" && password != "" {
// 			user, err := models.UserGetByName(username)
// 			fmt.Println("11")
// 			errorMsg := ""
// 			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
// 				errorMsg = "帐号或密码错误"
// 			} else if user.Status == -1 {
// 				errorMsg = "该帐号已禁用"
// 			} else {
// 				user.LastIp = this.getClientIp()
// 				user.LastLogin = time.Now().Unix()
// 				models.UserUpdate(user)
// 				fmt.Println("33")
// 				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
// 				if remember == "yes" {
// 					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
// 				} else {
// 					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
// 				}
// 				fmt.Println("44")
// 				fmt.Println("login:ok:" + beego.URLFor("MainController.Index"))
// 				this.redirect(beego.URLFor("MainController.Index"))
// 			}
// 			flash.Error(errorMsg)
// 			flash.Store(&this.Controller)
// 			fmt.Println("login:fail:" + beego.URLFor("MainController.Index"))
// 			this.redirect(beego.URLFor("MainController.Login"))
// 		}
// 	}

// 	this.TplName = "main/login.html"
// }

func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/")
	}
	//fmt.Println("asdfadsf")
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)
			fmt.Println("11")
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix()
				models.UserUpdate(user)
				fmt.Println("33")
				token := "webeditor-" + uuid.NewV1().String()

				userExt := models.GenUserExt()
				userExt.U = *user
				userExt.IP = this.getClientIp()

				second := 7 * 86400
				if remember == "yes" {
					userExt.RememberMe = true
					this.Ctx.SetCookie("token", token, second) //秒
					this.Ctx.SetCookie("username", username, second)
					libs.SaveToken(token, userExt, second) //秒7天过期
				} else {
					this.Ctx.SetCookie("token", token)
					this.Ctx.SetCookie("username", username)
					libs.SaveToken(token, userExt, second/7) //秒,一天过期
				}

				//fmt.Println("444444")
				fmt.Println("login:ok:" + beego.URLFor("MainController.Index"))
				this.redirect(beego.URLFor("MainController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			fmt.Println("login:fail:" + beego.URLFor("MainController.Index"))
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}

	this.TplName = "main/login.html"
}

// 退出登录
func (this *MainController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("MainController.Login"))
}

// 获取系统时间
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
