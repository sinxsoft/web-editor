package models

import (
	"github.com/astaxie/beego/orm"
	"fmt"
)

type ShortUri struct {
	ShortUri string `orm:"pk;column(short_uri);"`
	LongUri string
	Description string
	Status string
	UpdateTime string
	CreateUser string
}

//type ShortUriExt struct {
//	ShortUri
//	CreateTime time.Time
//}

func (su *ShortUri) TableName() string  {
	return TableName("short_uri")
}

func SearchShortUriList(pageSize,pageNo int,search string) ([]ShortUri) {
	o := orm.NewOrm()
	qs := o.QueryTable(TableName("short_uri"))
	if search !=""{
		qs=qs.Filter("ShortUri",search)
	}
	var shortUri []ShortUri
	cnt, err :=  qs.Limit(pageSize, (pageNo-1)*pageSize).All(&shortUri)
	if err == nil {
		fmt.Println("count", cnt)
	}
	return shortUri
}

func GetShortUriRecordNum(search string) int64 {

	o := orm.NewOrm()
	qs := o.QueryTable(TableName("short_uri"))
	if search !=""{
		qs=qs.Filter("ShortUri",search)
	}
	var cnt []ShortUri
	num, err :=  qs.All(&cnt)
	if err == nil {
		return num
	}else{
		return 0
	}
}

func (c *ShortUri) Update(fields ...string) error {

	//error := c.Update("Status")
	if _, err := orm.NewOrm().Update(c, fields...); err != nil {
		return err
	}
	return nil
}

func ShortUriAdd(shortUri *ShortUri) (int64, error) {
	return orm.NewOrm().Insert(shortUri)
}

func ShortUriGetAllList() ([]*ShortUri, error) {

	list := make([]*ShortUri, 0)

	_, err := orm.NewOrm().QueryTable(TableName("short_uri")).All(&list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

//另外一个方法
func ShortUriGetAllListExt() ([]ShortUri, error) {

	list := make([]ShortUri, 0)

	_, err := orm.NewOrm().QueryTable(TableName("short_uri")).All(&list)

	if err != nil {
		return nil, err
	}
	return list, nil

}

func ShortUriGetByShort(su string) (*ShortUri, error) {
	c := new(ShortUri)
	err := orm.NewOrm().QueryTable(TableName("short_uri")).Filter("short_uri", su).One(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ShortUriUpdate(shortUri *ShortUri) (int64, error) {
	return orm.NewOrm().Update(shortUri,"LongUri","Status","UpdateTime","Description")
}
