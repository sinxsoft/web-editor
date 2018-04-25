package controllers

import (
	"github.com/sinxsoft/web-editor/app/models"
	"strconv"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/astaxie/beego"
	"fmt"
	"io"
	"strings"
)

const PRERIX="SHORT_URI_"

type ShortUriController struct {
	BaseController
}

//接收跳转的controller 如： https://t.cn/123   -->  https://www.baidu.com/xxx
func (this *ShortUriController) Redirect(){
	//redis 处理
	id :=this.Ctx.Input.Param(":id")
	//baseUrl := beego.AppConfig.String("base.url") +"/"
	//redis 处理
 	b,e:=libs.GetObjectAndDelay(PRERIX+id,60*60*24)
	if e == nil && b != nil {
		//success
		this.redirect(string(b))
	}else{
		su,error:=models.ShortUriGetByShort(id)
		if error !=nil || su == nil{
			this.Ctx.WriteString("非法短链接")
		}else {
			libs.SaveObject(PRERIX+id, []byte(su.LongUri), 60*60*24)
			this.redirect(su.LongUri)
		}
	}
}

//全部记录，分页
func (this *ShortUriController) Index() {
	this.Data["pageTitle"] = "短链接编辑"
	beego.ReadFromRequest(&this.Controller)

	currentPage := this.Ctx.Request.FormValue("currentPage")
	search := this.Ctx.Request.FormValue("search")
	currPage := 1

	if currentPage != ""{
		currPage,_ = strconv.Atoi(currentPage)
	}
	// 已经新建好的好的页面
	entries:= models.SearchShortUriList(EachPageNum,currPage,search)

	this.Data["entries"] = entries
	num:=models.GetShortUriRecordNum(search)
	res := libs.Paginator(currPage, EachPageNum, num)
	this.Data["paginator"] = res
	this.Data["totals"] = num
	this.Data["search"] = search
	this.display("shorturi/shorturi")
}

func (this *ShortUriController) PutOne() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)
	this.Ctx.Input.AcceptsJSON()

	beego.ReadFromRequest(&this.Controller)
	shortUri := this.Ctx.Request.FormValue("shortUri")
	LongUri := this.Ctx.Request.FormValue("longUri")
	description := this.Ctx.Request.FormValue("description")

	if strings.Count(shortUri,"")-1 <3 ||
		strings.Count(shortUri,"")-1 >10{
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "短链接长度必须在区间[3,10]！")
		io.WriteString(w, jsonString)
		return
	}

	if strings.Count(LongUri,"")-1 <10 ||
		strings.Count(LongUri,"")-1 >300{
		jsonString := fmt.Sprintf(DOC_RESULT_JSON, "false", "长链接长度必须在区间[10,300]！")
		io.WriteString(w, jsonString)
		return
	}

	su,error:= models.ShortUriGetByShort(shortUri)
	msg := "新增成功!"
	if error ==nil && su != nil{
		su.LongUri = LongUri
		su.Status = "01"
		su.Description = description
		models.ShortUriUpdate(su)
		msg = "修改成功!"
	}else{
		su = new(models.ShortUri)
		su.ShortUri = shortUri
		su.LongUri = LongUri
		//su.Status = status
		su.Description = description
		su.CreateUser = this.Data["loginUserName"].(string)
		models.ShortUriAdd(su)
	}

	jsonString := fmt.Sprintf(DOC_RESULT_JSON, "true", msg)
	io.WriteString(w, jsonString)
}


func (this *ShortUriController) ShowOne(){
	beego.ReadFromRequest(&this.Controller)
	shortUri := this.Ctx.Request.FormValue("shortUri")
	su,e:=models.ShortUriGetByShort(shortUri)
	if su == nil || e != nil{
		if "new" == this.Ctx.Request.FormValue("action"){
			this.Data["su"] = new(models.ShortUri)
			this.display("shorturi/add")
		}else{
			this.Ctx.WriteString("非法短链接")
		}

	}else {
		this.Data["su"] = su
		this.display("shorturi/add")
	}
}

