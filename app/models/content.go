package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Content struct {
	Id          int
	DocId       string
	Name        string
	Type        string
	Url         string
	Description string
	Status      int
	CreateTime  string
	UserName    string
	UserId      int
}

func (c *Content) TableName() string {
	return TableName("content")
}

//根据条件返回总数
func GetContentRecordNum(search string) int64 {

	o := orm.NewOrm()
	qs := o.QueryTable(TableName("content"))
	if search != "" {
		qs = qs.Filter("Name", search)
	}
	var cnt []Content
	num, err := qs.All(&cnt)
	if err == nil {
		return num
	} else {
		return 0
	}
}

//总的content数量
func SearchContentList(pageSize, pageNo int, search string) []Content {
	o := orm.NewOrm()
	qs := o.QueryTable(TableName("content"))
	if search != "" {
		qs = qs.Filter("Name", search)
	}
	var content []Content
	cnt, err := qs.Limit(pageSize, (pageNo-1)*pageSize).All(&content)
	if err == nil {
		fmt.Println("count", cnt)
	}
	return content
}

func (c *Content) Update(fields ...string) error {

	//error := c.Update("Status")
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func ContentAdd(content *Content) (int64, error) {
	return orm.NewOrm().Insert(content)
}

func ContentGetAllList() ([]*Content, error) {

	list := make([]*Content, 0)

	_, err := orm.NewOrm().QueryTable(TableName("content")).All(&list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

//另外一个方法
func ContentGetAllListExt() ([]Content, error) {

	list := make([]Content, 0)

	_, err := orm.NewOrm().QueryTable(TableName("content")).All(&list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

func ContentGetById(id int) (*Content, error) {
	c := new(Content)
	err := orm.NewOrm().QueryTable(TableName("content")).Filter("id", id).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ContentGetByName(contentName string) (*Content, error) {
	c := new(Content)

	err := orm.NewOrm().QueryTable(TableName("content")).Filter("Name", contentName).One(c)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func ContentGetByDocID(docID string) (*Content, error) {
	c := new(Content)

	err := orm.NewOrm().QueryTable(TableName("content")).Filter("DocID", docID).One(c)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func ContentUpdate(content *Content, fields ...string) error {

	err := content.Update(fields...)

	return err
}
