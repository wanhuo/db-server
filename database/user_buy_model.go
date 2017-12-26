package database

import (
	//"strconv"
	"github.com/astaxie/beego/orm"
	//"db-server/proto/dbproto"
)

type UserTopUpShell struct{
	RecordID    	int    `orm:"column(record_id);pk;auto"`
	TransactionID   string `orm:"column(transaction_id);index"`
	UserID          int    `orm:"column(user_id);index"`
	Shell           int    `orm:"column(shell)"`
	RawInfo         string `orm:"column(raw_info)"`
	IP              string `orm:"column(ip)"`
	OpDate          int    `orm:"column(op_date)"`
}

func (u *UserTopUpShell) TableName() string {
	return "topup_history_tab"
}

type UserTopUpRecord struct{
	AutoID        int    `orm:"column(auto_id);auto;pk"`
	UserId        int	 `orm:"column(userId);"`
	TransactionID string `orm:"column(transactionId);index"`
	RMB			  int    `orm:"column(rmb)"`
	Shell         int    `orm:"column(shell)"`
	Status        int    `orm:"column(status)"`
	PayType       string `orm:"column(payType)"`
	Time          int    `orm:"column(transactionTime)"`
}
func (u *UserTopUpRecord) TableName() string {
	return "topup_record_t"
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
	IP        string `orm:"column(ip);null"`
}

func (u *UserBuyMembership) TableName() string {
	return "purchase_record_t"
}



func init() {
	orm.RegisterModel(new(UserTopUpShell), new(UserBuyMembership))
}


