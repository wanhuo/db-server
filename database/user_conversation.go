package database

import (
	"fmt"
	"strconv"
	"db-server/proto/dbproto"
	"github.com/astaxie/beego/orm"
)


type UserConversations struct{
	RecordID  int `orm:"column(recordid);auto;pk"`
	UserID    int `orm:"column(user_id)"`
	Cid       int `orm:"column(cid)"`
}

func (u *UserConversations) TableName() string {
	return "user_conversation_tab"
}

func userCidList(args SqlArgList) (rows []*dbproto.OneRow, error error) {
	userId, _ := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_msg_db")
	var cidList orm.ParamsList
	num, err := mysqlORM.QueryTable("user_conversation_tab").Filter("user_id", userId).ValuesFlat(&cidList, "cid")
	fmt.Printf("Returned Rows Num: %d", num)
	if err != nil {
		error = err
		return
	}
	row := serializeRowDate(cidList)
	rows = append(rows, row)
	return
}

func init()  {
	orm.RegisterModel(new(UserConversations))
	RegisterQueryCB("get.user.cidlist", userCidList)
}


