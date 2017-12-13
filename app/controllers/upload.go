package controllers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
	"github.com/sinxsoft/web-editor/app/libs"
	"github.com/sinxsoft/web-editor/app/models"
)

const (
	DOC_RESULT_JSON = "{\"state\": \"%s\",\"msg\":\"%s\"}"

	RESULT_JSON = "{\"state\": \"SUCCESS\",\"original\": \"%s\",\"size\": \"%d\",\"title\": \"%s\",\"type\": \"%s\",\"url\": \"%s\"}"
)

type UploadController struct {
	BaseController
}

func (this *UploadController) Object() {
	w := this.Ctx.ResponseWriter
	this.Ctx.Request.ParseForm()
	//file := this.Ctx.Request.Form["file"]
	file := this.Ctx.Input.Param(":splat")
	//从oss拿到对象
	b := this.GetOss(file)
	w.Write(b)
}

func (this *UploadController) Document() {
	w := this.Ctx.ResponseWriter
	this.Ctx.Request.ParseForm()
	id := this.Ctx.Input.Param(":splat")
	//从本地存储拿到html文件
	b, error := ioutil.ReadFile(beego.AppConfig.String("document.path") + id + ".file")
	if error != nil {
		w.Write([]byte(error.Error()))
		return
	}
	w.Write(b)

}

/**
 *所有上传的主入口
 */
func (this *UploadController) Controller() {
	//beego.ReadFromRequest(this.Controller)
	//this.Ctx.Request

	r := this.Ctx.Request
	w := this.Ctx.ResponseWriter

	this.Ctx.Request.ParseForm()
	action := this.Ctx.Request.Form["action"]

	if action != nil {

		act := action[0]
		if "uploadimage" == act && "POST" == r.Method {
			//上传图片、文件
			//file, _ := r.MultipartForm.File[""][0].Open()
			w.Header().Add("Content-Type", "text/json;charset=utf-8")
			w.WriteHeader(200)
			file, fileHeader, err := r.FormFile("upfile")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			defer file.Close()

			fileExtName := path.Ext(fileHeader.Filename)

			uuid := uuid.NewV4()
			fileName := uuid.String() + fileExtName

			f, err := os.Create(fileName)
			//defer f.Close()
			io.Copy(f, file)
			f.Close()

			//fmt.Fprintf(w, "上传文件的大小为: %d", file.(Sizer).Size())
			//开始上传oss
			//fileInfo, _ := f.Stat()

			url := this.UploadOss(fileName, fileName)
			//fmt.Fprintf(w, "上传oss文件的大小为: %d ,url:"+url, file.(Sizer).Size())
			//json := "{\"state\": \"SUCCESS\",\"original\": \"%s\",\"size\": \"%d\",\"title\": \"%s\",\"type\": \"%s\",\"url\": \"%s\"}"
			json := fmt.Sprintf(RESULT_JSON, fileHeader.Filename, fileHeader.Size, fileHeader.Filename, "", url)

			io.WriteString(w, json)

			os.Remove(fileName)

		} else if "config" == act && "GET" == r.Method {
			w.Header().Add("Content-Type", "text/json;charset=utf-8")
			w.WriteHeader(200)
			//var a interface{}
			jsonString := beego.AppConfig.String("config")
			//json.Unmarshal([]byte(jsonString), a)
			io.WriteString(w, jsonString)
		} else if "document" == act && "POST" == r.Method {
			//var a interface{}
			//id := this.Ctx.Request.Form["id"]
			name := this.Ctx.Request.Form["nameX"]
			docID := this.Ctx.Request.Form["docID"]
			//content := this.Ctx.Request.Form["myEditor"]
			content := this.Ctx.Request.Form["editorValue"]
			description := this.Ctx.Request.Form["description"]
			id := this.Ctx.Request.Form["id"]
			if content == nil {
				jString := fmt.Sprintf(DOC_RESULT_JSON, "fail", "content 为空！")
				io.WriteString(w, jString)
				return
			}

			//判断是否存在文件，存在就改名
			fName := ".html.file"
			path := beego.AppConfig.String("document.path") + docID[0] + fName

			bool, _ := libs.Exists(path)
			if bool {
				time := time.Now().Format("2006/01/02 15:04:05")
				time = strings.Replace(time, "/", "_", -1)
				time = strings.Replace(time, ":", "_", -1)
				e := os.Rename(path, path+"_"+time)
				if e != nil {
					models.HandleError(e)
					jString := fmt.Sprintf(DOC_RESULT_JSON, "fail", e)
					io.WriteString(w, jString)
					return
				}
			}
			data := []byte(content[0])
			fl, err := os.OpenFile(beego.AppConfig.String("document.path")+docID[0]+fName,
				os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				models.HandleError(err)
				jString := fmt.Sprintf(DOC_RESULT_JSON, "fail", err)
				io.WriteString(w, jString)
				return
			}
			defer fl.Close()
			_, err = fl.Write(data)

			// err := ioutil.WriteFile(beego.AppConfig.String("document.path")+docID[0],
			// 	[]byte(content[0]), os.ModeAppend)
			if err != nil {
				models.HandleError(err)
				jString := fmt.Sprintf(DOC_RESULT_JSON, "fail", err)
				io.WriteString(w, jString)
				return
			}

			var error error
			//存储入库
			idInt, _ := strconv.Atoi(id[0])
			var c *models.Content
			if idInt >= 0 {
				c, _ = models.ContentGetById(idInt)
				//c.DocId = docID[0]
				//c.CreateTime = time.Now().Format("2006/01/02 15:04:05")
				c.Description = description[0]
				c.Name = name[0]
				//c.Type = "1"
				//c.Url = beego.AppConfig.String("document.url") + c.DocId
				//jsonString := beego.AppConfig.String("config")
				//json.Unmarshal([]byte(jsonString), a)
				error = models.ContentUpdate(c)
			} else {
				c = new(models.Content)
				c.DocId = docID[0]
				c.CreateTime = time.Now().Format("2006/01/02 15:04:05")
				c.Description = description[0]
				c.Name = name[0]
				c.Type = "1"
				c.Url = beego.AppConfig.String("document.url") + c.DocId + ".html"

				c.UserName = this.Data["loginUserName"].(string)
				c.UserId = this.Data["loginUserId"].(int)
				//jsonString := beego.AppConfig.String("config")
				//json.Unmarshal([]byte(jsonString), a)
				_, error = models.ContentAdd(c)

			}
			if error != nil {
				models.HandleError(error)
				jString := fmt.Sprintf(DOC_RESULT_JSON, "fail", error)
				io.WriteString(w, jString)
				return
			}
			// result := "<script>window.location.href = '/';</script>"
			// w.Header().Add("Content-Type", "text/text;charset=utf-8")
			// w.WriteHeader(200)
			// io.WriteString(w, result)
			this.Data["message"] = "发布成功！"
			this.Data["redirect"] = "/"
			this.Data["url"] = c.Url
			this.display("error/result")
			return
		}
	} else {
		edit := this.Ctx.Request.Form["edit"]
		id := this.Ctx.Request.Form["id"]
		if edit != nil && edit[0] == "true" && id != nil {
			i, _ := strconv.Atoi(id[0])

			c, _ := models.ContentGetById(i)
			this.Data["id"] = c.Id
			this.Data["docID"] = c.DocId
			this.Data["name"] = c.Name
			this.Data["description"] = c.Description

			b, error := ioutil.ReadFile(beego.AppConfig.String("document.path") + c.DocId + ".file")
			if error != nil {
				this.Data["content"] = error.Error()
			} else {
				this.Data["content"] = string(b)
			}

		} else {
			this.Data["id"] = -1
			this.Data["docID"] = uuid.NewV1().String()
			this.Data["content"] = ""
			this.Data["name"] = ""
			this.Data["description"] = ""
		}

		this.TplName = "main/add.html"
	}
}

func (this *UploadController) GetOss(objectKey string) []byte {
	accessKey := beego.AppConfig.String("oss.key")
	accessKeySecret := beego.AppConfig.String("oss.secret")
	endPoint := beego.AppConfig.String("oss.clientEndPoint")
	client, err := oss.New(endPoint, accessKey, accessKeySecret)

	if err != nil {
		models.HandleError(err)
	}
	bucket, err := client.Bucket(beego.AppConfig.String("oss.bucketName"))

	body, err := bucket.GetObject(objectKey)
	if err != nil {
		models.HandleError(err)
	}
	data, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		models.HandleError(err)
	}
	return data
}

func (this *UploadController) UploadOss(objectKey string, localFile string) string {
	accessKey := beego.AppConfig.String("oss.key")
	accessKeySecret := beego.AppConfig.String("oss.secret")
	endPoint := beego.AppConfig.String("oss.clientEndPoint")
	client, err := oss.New(endPoint, accessKey, accessKeySecret)
	if err != nil {
		models.HandleError(err)
	}

	bucket, err := client.Bucket(beego.AppConfig.String("oss.bucketName"))

	// 场景4：上传本地文件，不需要打开文件。
	err = bucket.UploadFile(objectKey, localFile, 100*1024, oss.Routines(3))
	if err != nil {
		models.HandleError(err)
	}

	//baseURL := beego.AppConfig.String("oss.endPoint")

	getURL := beego.AppConfig.String("get.url")

	fmt.Println(client)

	fmt.Println(err)

	return getURL + objectKey
}
