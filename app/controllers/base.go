package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/sinxsoft/web-editor/app/models"
	"time"
)

var (
	Filters []string
)

func IsEscape(url string) bool {
	if url == "" {
		return false
	}
	for _, v := range Filters {
		if strings.Index(url, v) == 0 {
			return true
		}
	}
	return false
}

const (
	MSG_OK  = 0
	MSG_ERR = -1
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *models.User
	userId         int
	userName       string
	pageSize       int
}

func (this *BaseController) Prepare() {
	this.pageSize = 20
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)

	url := this.Ctx.Request.URL

	// fmt.Println("\r\n")

	// fmt.Println(this.Ctx.Request)

	// fmt.Println("对方IP：" + this.Ctx.Request.RemoteAddr)

	// fmt.Println("\r\n")

	if !IsEscape(url.Path) {
		//fmt.Println("1")
		this.auth()
	}

	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["siteName"] = beego.AppConfig.String("site.name")
	this.Data["curRoute"] = this.controllerName + "." + this.actionName
	this.Data["curController"] = this.controllerName
	this.Data["curAction"] = this.actionName
	this.Data["loginUserId"] = this.userId
	this.Data["loginUserName"] = this.userName
	fmt.Println("5")
}

//登录状态验证

func (this *BaseController) auth() {

	username := this.Ctx.GetCookie("username")
	token := this.Ctx.GetCookie("token")

	if token == "" || username == "" {
		this.redirect(beego.URLFor("MainController.Login"))
	} else {
		t, b := libs.GetToken(token)
		if !b {
			this.redirect(beego.URLFor("MainController.Login"))
		} else {
			if username != t.U.UserName || this.getClientIp() != t.IP {
				fmt.Println("why:" + this.getClientIp() + "," + t.IP)
				fmt.Println("why:" + username + "," + t.U.UserName)
				this.redirect(beego.URLFor("MainController.Login"))
			}
		}
		user := t.U
		this.userId = user.Id
		this.userName = user.UserName
		this.user = &user
	}

}


// 1 success
// 0 fail
func JsonFormat(retcode int, retmsg string, retdata interface{}, stime time.Time) (json map[string]interface{}) {
	cost := time.Now().Sub(stime).Seconds()
	if retcode == 1 {
		json = map[string]interface{}{
			"code": retcode,
			"data": retdata,
			"desc": retmsg,
			"cost": cost,
		}
	} else {
		json = map[string]interface{}{
			"code": retcode,
			"desc": retmsg,
			"cost": cost,
		}
	}

	return json
}


//渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = this.controllerName + "/" + this.actionName + ".html"
	}
	this.Layout = "layout/layout.html"
	this.TplName = tplname
}

// 重定向
func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
	this.StopRun()
}

// 是否POST提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}

// 显示错误信息
func (this *BaseController) showMsg(args ...string) {
	this.Data["message"] = args[0]
	redirect := this.Ctx.Request.Referer()
	if len(args) > 1 {
		redirect = args[1]
	}

	this.Data["redirect"] = redirect
	this.Data["pageTitle"] = "系统提示"
	this.display("error/message")
	this.Render()
	this.StopRun()
}

// 输出json
func (this *BaseController) jsonResult(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["msg"] = msg

	this.jsonResult(out)
}

//获取用户IP地址
func (this *BaseController) getClientIp() string {
	var realIP string
	realIP = this.Ctx.Request.Header.Get("X-Real-Forwarded-For")
	if realIP != "" {
		fmt.Println("X-Real-Forwarded-For" + realIP)
		return realIP
	}
	realIP = this.Ctx.Request.Header.Get("X-Forwarded-For")
	if realIP != "" {
		fmt.Println("X-Forwarded-For" + realIP)
		return realIP
	}
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}
