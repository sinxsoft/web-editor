package controllers

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/satori/go.uuid"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/sinxsoft/web-editor/app/models"
	"github.com/dchest/captcha"
)

const EachPageNum = 10

type MainController struct {
	BaseController
}

func (this *MainController) Home() {
	this.redirect("/index?currentPage=1&search=")
}

func (this *MainController) Index() {
	this.Data["pageTitle"] = "webpage编辑器"

	// 已经新建好的好的页面
	entries, _ := models.ContentGetAllListExt()

	this.Data["entries"] = entries

	this.display()
}


func (this *MainController) IndexPager() {
	beego.ReadFromRequest(&this.Controller)
	this.Data["pageTitle"] = "webpage编辑器"
	currentPage := this.Ctx.Request.FormValue("currentPage")
	search := this.Ctx.Request.FormValue("search")

	currPage := 1

	if currentPage != ""{
		currPage,_ = strconv.Atoi(currentPage)
	}


	//已经新建好的好的页面
	entries := models.SearchContentList(EachPageNum,currPage,search)

	this.Data["entries"] = entries
	num:=models.GetContentRecordNum(search)
	res := libs.Paginator(currPage, EachPageNum, num)
	this.Data["paginator"] = res
	this.Data["totals"] = num
	this.Data["search"] = search
	this.display("main/index")
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

func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/index?currentPage=1&search=")
	}
	//fmt.Println("asdfadsf")
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		captchaId := this.GetString("captchaId")
		captchaValue := this.GetString("captcha")

		errorMsg := ""
		if !captcha.VerifyString(captchaId, captchaValue) {
			errorMsg = "验证码错误！"
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			this.redirect(beego.URLFor("MainController.Login"))
			return
		}
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)


			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误！"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用！"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix()
				models.UserUpdate(user)

				uuid,_:=uuid.NewV1()
				token := "webeditor-" + uuid.String()

				userExt := models.GenUserExt()
				userExt.U = *user
				userExt.IP = this.getClientIp()

				second := 7 * 86400
				if remember == "yes" || remember == "1" {
					userExt.RememberMe = true
					this.Ctx.SetCookie("token", token, second) //秒
					this.Ctx.SetCookie("username", username, second)
					libs.SaveToken(token, userExt, second) //秒7天过期
				} else {
					this.Ctx.SetCookie("token", token)
					this.Ctx.SetCookie("username", username)
					libs.SaveToken(token, userExt, second/7) //秒,一天过期
				}

				fmt.Println("login:ok:" + beego.URLFor("MainController.IndexPager"))
				this.redirect(beego.URLFor("MainController.IndexPager"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			fmt.Println("login:fail:" + beego.URLFor("MainController.IndexPager"))
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}

	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	this.Data["CaptchaId"] = d.CaptchaId

	this.TplName = "main/login.html"
}

// 退出登录
func (this *MainController) Logout() {

	libs.DeleteToken(this.Ctx.GetCookie("token"))
	this.Ctx.SetCookie("token", "")
	this.Ctx.SetCookie("username", "")
	this.redirect(beego.URLFor("MainController.Login"))
}

// 获取系统时间
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
