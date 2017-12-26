package database

import (
	"github.com/astaxie/beego/orm"
)

type UserTopUpShell struct{
	RecordID    	int    `orm:"column(record_id);pk;auto"`
	TransactionID   string `orm:"column(transaction_id);index"`
	UserID          int    `orm:"column(user_id);index"`
	Shell           int    `orm:"column(shell)"`
	RawInfo         string `orm:"column(raw_info)"`
	Ip              string `orm:"column(ip)"`
	OpDate          int    `orm:"column(op_date)"`
}


type UserBuyMembership struct{
	AutoID    int    `orm:"column(auto_id);auto;pk"`
	UserId    int	 `orm:"column(userId);"`
	OrderID   string `orm:"column(orderId);index"`
	Shell     int    `orm:"column(shell)"`
	OrderTime int    `orm:"column(orderTime)"`
	GoodsId   string `orm:"column(goodsId)"`
	StartTime int    `orm:"column(startTime)"`
	EndTime   int    `orm:"column(endTime)"`
	Ip        string `orm:"column(ip);null"`
}

func init() {
	orm.RegisterModel(new(UserTopUpShell), new(UserBuyMembership))
}
