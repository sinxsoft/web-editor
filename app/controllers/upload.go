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

func (this *UploadController) DocList() {
	this.Data["pageTitle"] = "文档清单"

	// 已经新建好的好的页面
	entries, _ := models.DocumentListGetAllList()

	this.Data["entries"] = entries

	this.display("main/documentlist")
}

func (this *UploadController) Document() {
	w := this.Ctx.ResponseWriter
	this.Ctx.Request.ParseForm()
	id := this.Ctx.Input.Param(":splat")

	docID := strings.Replace(id, ".html", "", -1)

	c, er := models.ContentGetByDocID(docID)

	if er != nil {
		w.Write([]byte(er.Error()))
		return
	}

	if c.Status == 1 {
		w.Write([]byte("该文档已经删除或者不存在！"))
		return
	}

	b, error := libs.GetObjectAndDelay(id, 60*60*24)

	if b != nil && error == nil {
		//ok
		fmt.Println("获取data成功：" + id)
	} else {
		//从本地存储拿到html文件
		b, error = ioutil.ReadFile(beego.AppConfig.String("document.path") + id + ".file")
		if error != nil {
			this.Ctx.WriteString(error.Error())
			return
		}
		error = libs.SaveObject(id, b, 60*60*24)
		if error != nil {
			fmt.Println("Save data 失败：" + id + "," + error.Error())
		} else {
			fmt.Println("Save data 成功：" + id)
		}

	}

	this.Ctx.WriteString(string(b))

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

			uuid,_ := uuid.NewV1()
			fileName := uuid.String() + fileExtName

			f, err := os.Create(fileName)
			//defer f.Close()
			io.Copy(f, file)
			f.Close()

			//fmt.Fprintf(w, "上传文件的大小为: %d", file.(Sizer).Size())
			//开始上传oss
			//fileInfo, _ := f.Stat()

			url := this.UploadOss(fileName, fileName, fileHeader.Filename)
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
			content := this.Ctx.Request.Form["contentEditor"]

			//iddd := this.Ctx.Request.Form["iddd"]

			//fmt.Println(iddd)
			fmt.Println(content[0])

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

			//删除掉cacheid
			cacheId :=  c.DocId + ".html"
			libs.DelObject(cacheId)


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
		//判断是否存在文件，存在就改名
		fName := ".html.file"
		edit := this.Ctx.Request.Form["edit"]
		id := this.Ctx.Request.Form["id"]
		if edit != nil && edit[0] == "true" && id != nil {
			i, _ := strconv.Atoi(id[0])

			c, _ := models.ContentGetById(i)
			this.Data["id"] = c.Id
			this.Data["docID"] = c.DocId
			this.Data["name"] = c.Name
			this.Data["description"] = c.Description
			this.Data["edit"] = "true"

			b, error := ioutil.ReadFile(beego.AppConfig.String("document.path") + c.DocId + fName)
			if error != nil {
				this.Data["content"] = error.Error()
			} else {
				this.Data["content"] = string(b)
			}

		} else {
			uuid,_ := uuid.NewV1()

			this.Data["id"] = -1
			this.Data["docID"] = strings.Replace(uuid.String(), "-", "", -1)
			this.Data["content"] = ""
			this.Data["name"] = "请修改"
			this.Data["description"] = "请修改"
			this.Data["edit"] = "false"
		}
		this.Data["pageTitle"] = "添加文章"
		this.display("main/add")
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

func (this *UploadController) UploadOss(objectKey string, localFile string, originalFileName string) string {
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

	//fmt.Println(client)

	//fmt.Println(err)
	dl := new(models.DocumentList)
	dl.DocId = objectKey
	dl.CreateTime = time.Now().Format("2006/01/02 15:04:05")
	dl.OssAllUrl = getURL + objectKey
	dl.UserName = this.Data["loginUserName"].(string)
	dl.UserId = this.Data["loginUserId"].(int)
	dl.OriginalFileName = libs.SubString(originalFileName, 0, 100)
	models.DocumentListAdd(dl)

	return getURL + objectKey
}

type ID struct {
	ID string
}

func (this *UploadController) Del() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)

	jsonString := fmt.Sprintf(DOC_RESULT_JSON, "true", "成功！")
	this.Ctx.Input.AcceptsJSON()
	id := this.GetString("ID")
	intID, _ := strconv.Atoi(id)
	model, errorModel := models.ContentGetById(intID)

	if errorModel != nil {
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", errorModel)
		io.WriteString(w, jsonString)
		return
	}
	if model.Status == 1 {
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", "不存在文档！或已经不可用！")
		io.WriteString(w, jsonString)
		return
	}

	path := beego.AppConfig.String("document.path") + model.DocId + ".html.file"

	bool, _ := libs.Exists(path)
	if bool {
		e := os.Rename(path, path+".deleted")
		if e != nil {
			models.HandleError(e)
			jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", e)
		} else {
			jsonString = fmt.Sprintf(DOC_RESULT_JSON, "true", "操作成功！")

			//model, _ := models.ContentGetById(intID)
			model.Status = 1
			error := models.ContentUpdate(model, "Status")
			if error != nil {
				os.Rename(path+".deleted", path)
				jsonString = fmt.Sprintf(DOC_RESULT_JSON, "false", error)
			}
		}
	}
	//jsonString := beego.AppConfig.String("config")
	//json.Unmarshal([]byte(jsonString), a)
	io.WriteString(w, jsonString)
}

func (this *UploadController) CancelDel() {

	w := this.Ctx.ResponseWriter
	w.Header().Add("Content-Type", "text/json;charset=utf-8")
	w.WriteHeader(200)

	jsonString := fmt.Sprintf(DOC_RESULT_JSON, "true", "恢复成功！")
	this.Ctx.Input.AcceptsJSON()
	id := this.GetString("ID")
	intID, _ := strconv.Atoi(id)
	model, errorModel := models.ContentGetById(intID)

	if errorModel != nil {
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", errorModel)
		io.WriteString(w, jsonString)
		return
	}
	if model.Status == 0 {
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", "该文档不需要恢复！")
		io.WriteString(w, jsonString)
		return
	}

	path := beego.AppConfig.String("document.path") + model.DocId + ".html.file"

	bool, _ := libs.Exists(path + ".deleted")
	if bool {
		e := os.Rename(path+".deleted", path)
		if e != nil {
			models.HandleError(e)
			jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", e)
		} else {
			jsonString = fmt.Sprintf(DOC_RESULT_JSON, "true", "操作成功！")

			//model, _ := models.ContentGetById(intID)
			model.Status = 0
			error := models.ContentUpdate(model, "Status")
			if error != nil {
				os.Rename(path, path+".deleted")
				jsonString = fmt.Sprintf(DOC_RESULT_JSON, "false", error)
			}
		}
	} else {
		jsonString = fmt.Sprintf(DOC_RESULT_JSON, "fail", "文件不存在！")
	}
	//jsonString := beego.AppConfig.String("config")
	//json.Unmarshal([]byte(jsonString), a)
	io.WriteString(w, jsonString)
}
