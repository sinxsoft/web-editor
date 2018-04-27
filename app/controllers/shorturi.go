package controllers

import (
	"github.com/sinxsoft/web-editor/app/models"
	"strconv"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/astaxie/beego"
	"fmt"
	"io"
	"strings"
	"time"
	"math/rand"
)

const PRERIX = "SHORT_URI_"

type ShortUriController struct {
	BaseController
}

//接收跳转的controller 如： https://t.cn/123   -->  https://www.baidu.com/xxx
func (this *ShortUriController) Redirect() {
	//redis 处理
	id := this.Ctx.Input.Param(":id")
	//baseUrl := beego.AppConfig.String("base.url") +"/"
	//redis 处理
	b, e := libs.GetObjectAndDelay(PRERIX+id, 60*60*24)
	if e == nil && b != nil { //存在对象，但是不确保longuri是http（s）
		//success
		this.handleRequest(string(b))

	} else {
		su, error := models.ShortUriGetByShort(id)
		if error != nil || su == nil {
			this.Ctx.WriteString("非法短链接")
		} else {
			if su.Status != "10" {
				su.LongUri = "error:无此链接！"
			}
			libs.SaveObject(PRERIX+id, []byte(su.LongUri), 60*60*24)

			this.handleRequest(su.LongUri)
		}
	}
}

func (this *ShortUriController) handleRequest(longUri string) {
	if strings.Index(strings.ToLower(longUri), "http") < 0 {
		this.writeHtml(longUri)
	} else {
		this.redirect(longUri)
	}
}
func (this *ShortUriController) writeHtml(html string) {
	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(200)
	io.WriteString(w, html)
}

//全部记录，分页
func (this *ShortUriController) Index() {
	this.Data["pageTitle"] = "短链接编辑"
	beego.ReadFromRequest(&this.Controller)

	currentPage := this.Ctx.Request.FormValue("currentPage")
	search := this.Ctx.Request.FormValue("search")
	currPage := 1

	if currentPage != "" {
		currPage, _ = strconv.Atoi(currentPage)
	}
	// 已经新建好的好的页面
	entries := models.SearchShortUriList(EachPageNum, currPage, search)

	this.Data["entries"] = entries
	num := models.GetShortUriRecordNum(search)
	res := libs.Paginator(currPage, EachPageNum, num)
	this.Data["paginator"] = res
	this.Data["totals"] = num
	this.Data["search"] = search
	this.Data["baseurl"] = beego.AppConfig.String("short.baseurl")
	this.display("shorturi/shorturi")
}

func (this *ShortUriController) PutOne() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)
	this.Ctx.Input.AcceptsJSON()

	beego.ReadFromRequest(&this.Controller)
	shortUri := this.Ctx.Request.FormValue("shortUri")
	longUri := this.Ctx.Request.FormValue("longUri")
	description := this.Ctx.Request.FormValue("description")
	status := this.Ctx.Request.FormValue("status")
	action := this.Ctx.Request.FormValue("action")
	if strings.Count(shortUri, "")-1 < 3 ||
		strings.Count(shortUri, "")-1 > 10 {
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "短链接长度必须在区间[3,10]！")
		io.WriteString(w, jsonString)
		return
	}

	if strings.Count(longUri, "")-1 < 10 ||
		strings.Count(longUri, "")-1 > 300 {
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "长链接长度必须在区间[10,300]！")
		io.WriteString(w, jsonString)
		return
	}

	if !strings.HasPrefix(strings.ToLower(longUri), "http://") &&
		!strings.HasPrefix(strings.ToLower(longUri), "https://") {
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "长链接必须使用https://或者http://开头")
		io.WriteString(w, jsonString)
		return
	}
	jsonString := ""
	su, error := models.ShortUriGetByShort(shortUri)
	msg := "新增成功!"
	if error == nil && su != nil && action == "edit" {
		su.LongUri = longUri
		if strings.TrimSpace(status) != "" {
			su.Status = status
		}
		su.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
		su.Description = description
		models.ShortUriUpdate(su)
		msg = "修改成功!"
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "true", msg)
	} else if error != nil && su == nil && action == "new" {
		su = new(models.ShortUri)
		su.ShortUri = shortUri
		su.LongUri = longUri
		su.Status = "10" //default 10
		su.Description = description
		su.CreateUser = this.Data["loginUserName"].(string)
		models.ShortUriAdd(su)
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "true", msg)
	} else if error == nil && su != nil && action == "new" {
		msg = "短链接重复!"
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "false", msg)
	} else {
		msg = "状态异常!"
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "false", msg)
	}

	io.WriteString(w, jsonString)

}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (this *ShortUriController) GenShortUri() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)
	this.Ctx.Input.AcceptsJSON()
	for {
		rand := getRandomString(5)
		su, e := models.ShortUriGetByShort(rand)
		if su == nil && e != nil {
			jsonString := fmt.Sprintf(DOC_RESULT_JSON, "true", rand)
			io.WriteString(w, jsonString)
			break
		}
	}

}

func (this *ShortUriController) ShowOne() {
	beego.ReadFromRequest(&this.Controller)
	shortUri := this.Ctx.Request.FormValue("shortUri")
	action := this.Ctx.Request.FormValue("action")

	if "new" == action {
		this.Data["su"] = new(models.ShortUri)
		this.Data["action"] = "new"
		this.display("shorturi/add")
	} else if "edit" == action && shortUri != "" {
		su, e := models.ShortUriGetByShort(shortUri)
		if su == nil || e != nil {
			this.Ctx.WriteString("非法短链接" + e.Error())
		} else {
			this.Data["su"] = su
			this.Data["action"] = "edit"
			this.display("shorturi/add")
		}
	} else {
		this.Ctx.WriteString("非法短链接")
	}
}

func (this *ShortUriController) ChangeStatus() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)
	this.Ctx.Input.AcceptsJSON()
	shortUri := this.Ctx.Request.FormValue("shortUri")
	status := this.Ctx.Request.FormValue("status")
	su, e := models.ShortUriGetByShort(shortUri)
	if su == nil && e != nil {
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "不存在短链接！")
		io.WriteString(w, jsonString)
	} else {
		if su.Status == status {
			jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "状态无需改变！")
			io.WriteString(w, jsonString)
		} else {
			su.Status = status
			models.ShortUriUpdate(su)
			jsonString := fmt.Sprintf(DOC_RESULT_JSON, "true", "成功！")
			io.WriteString(w, jsonString)

			go func() {
				libs.DelObject(PRERIX + shortUri)
			}()
		}
	}

}
