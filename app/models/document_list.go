package models

import "github.com/astaxie/beego/orm"

type DocumentList struct {
	Id               int
	DocId            string
	OssAllUrl        string
	CreateTime       string
	UserName         string
	UserId           int
	OriginalFileName string
}

func (c *DocumentList) TableName() string {
	return TableName("document_list")
}

func (c *DocumentList) Update(fields ...string) error {

	//error := c.Update("Status")
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func DocumentListAdd(documentList *DocumentList) (int64, error) {
	return orm.NewOrm().Insert(documentList)
}

func DocumentListGetAllList() ([]*DocumentList, error) {

	list := make([]*DocumentList, 0)

	_, err := orm.NewOrm().QueryTable(TableName("document_list")).All(&list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

func DocumentListGetById(id int) (*DocumentList, error) {
	c := new(DocumentList)
	err := orm.NewOrm().QueryTable(TableName("document_list")).Filter("id", id).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func DocumentListGetByDocID(docID string) (*DocumentList, error) {
	c := new(DocumentList)

	err := orm.NewOrm().QueryTable(TableName("document_list")).Filter("DocID", docID).One(c)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func DocumentListUpdate(documentList *DocumentList, fields ...string) error {

	err := documentList.Update(fields...)

	return err
}
