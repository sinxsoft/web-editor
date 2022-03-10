package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"strings"
	"web-editor/app/controllers"
	"web-editor/app/libs"
	"web-editor/app/models"
)

const VERSION = "1.0.0"

func main() {

	//自定义load自己的配置，配置文件放到外面
	error := beego.LoadAppConfig("ini", "/a/www/webeditor/app.conf")
	//error := beego.LoadAppConfig("ini", "/Users/henrik/Documents/golang/src/github.com/sinxsoft/web-editor/conf/app.conf")

	if error != nil {
		fmt.Println(error)
	}

	models.Init()

	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})

	// 生产环境不输出debug日志
	if beego.AppConfig.String("runmode") == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}
	beego.AppConfig.Set("version", VERSION)

	// 路由设置
	beego.Router("/", &controllers.MainController{}, "*:Home")
	beego.Router("/index", &controllers.MainController{}, "*:IndexPager")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/gettime", &controllers.MainController{}, "*:GetTime")
	beego.Router("/help", &controllers.HelpController{}, "*:Index")
	// beego.AutoRouter(&controllers.TaskController{})
	beego.Router("/controller", &controllers.UploadController{}, "*:Controller")
	//beego.Router("/jsp/controller.jsp", &controllers.UploadController{}, "*:Controller")
	beego.Router("/object/*", &controllers.UploadController{}, "*:Object")

	beego.Router("/doc/*", &controllers.UploadController{}, "*:Document")
	beego.Router("/doclist", &controllers.UploadController{}, "*:DocList")

	beego.Router("/shorturi", &controllers.ShortUriController{}, "*:Index")
	beego.Router("/showshorturi", &controllers.ShortUriController{}, "*:ShowOne")
	beego.Router("/addshorturi", &controllers.ShortUriController{}, "*:PutOne")
	beego.Router("/t/:id:string", &controllers.ShortUriController{}, "*:Redirect")
	beego.Router("/genshorturi", &controllers.ShortUriController{}, "*:GenShortUri")
	beego.Router("/short/changeStatus", &controllers.ShortUriController{}, "*:ChangeStatus")

	beego.Router("/verify", &controllers.CaptchaController{}, "post:VerifyCaptcha")
	beego.Handler("/captcha/*.png", captcha.Server(240, 80))

	beego.AutoRouter(&controllers.UploadController{})
	//add a test page by henrik
	//beego.AutoRouter(&controllers.CommandController{})

	beego.BConfig.WebConfig.Session.SessionOn = false
	beego.SetStaticPath("/public", beego.AppConfig.String("public.dir"))
	viewsDir := beego.AppConfig.String("views.dir")
	if viewsDir != "" {
		beego.SetViewsPath(viewsDir)
	}

	staticDir := beego.AppConfig.String("static.dir")
	if staticDir != "" {
		beego.SetStaticPath("static", staticDir)
	}

	filters := beego.AppConfig.String("action.noauth.url")
	controllers.Filters = strings.Split(filters, ",")

	//设置自己的store
	captcha.SetCustomStore(new(libs.RedisStore))

	//golang 内存分析/动态追踪
	//golang pprof
	//当你的golang程序在运行过程中消耗了超出你理解的内存时，你就需要搞明白，到底是 程序中哪些代码导致了这些内存消耗。
	// 此时golang编译好的程序对你来说是个黑盒，该 如何搞清其中的内存使用呢？幸好golang已经内置了一些机制来帮助我们进行分析和追 踪。
	//此时，通常我们可以采用golang的pprof来帮助我们分析golang进程的内存使用。
	pprof := beego.AppConfig.String("pprof.port")
	if pprof == "" {
		pprof = "8867"
	}
	go func() {
		http.ListenAndServe("0.0.0.0:"+pprof, nil)
	}()

	beego.Run()
}

//---------------------------------以下是用reflect实现一些类型无关的泛型编程示例
//new object same the type as sample
func New(sample interface{}) interface{} {
	t := reflect.ValueOf(sample).Type()
	v := reflect.New(t).Interface()
	return v
}

//---------------------------------check type of aninterface
func CheckType(val interface{}, kind reflect.Kind) bool {
	v := reflect.ValueOf(val)
	return kind == v.Kind()
}

//---------------------------------if _func is not a functionor para num and type not match,it will cause panic
func Call(_func interface{}, params ...interface{}) (result []interface{}, err error) {
	f := reflect.ValueOf(_func)
	if len(params) != f.Type().NumIn() {
		ss := fmt.Sprintf("The number of params is not adapted.%s", f.String())
		panic(ss)
		return
	}
	var in []reflect.Value
	if len(params) > 0 { //prepare in paras
		in = make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
	}
	out := f.Call(in)
	if len(out) > 0 { //prepare out paras
		result = make([]interface{}, len(out), len(out))
		for i, v := range out {
			result[i] = v.Interface()
		}
	}
	return
}
