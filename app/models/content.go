package models

import (
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
