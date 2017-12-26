package database

import (
	"github.com/astaxie/beego/orm"
)


type UserConversations struct{
	RecordID  int `orm:"column(recordid);auto;pk"`
	UserID    int `orm:"column(user_id)"`
	Cid       int `orm:"column(cid)"`
}

func init()  {
	orm.RegisterModel(new(UserConversations))
}


