package database

import (
	"strconv"
	"db-server/proto/dbproto"
	"github.com/astaxie/beego/orm"
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

func userTopup(args SqlArgList) (err error) {
	topupId, payType := args[1], args[5]
	userId, err := strconv.Atoi(args[0])
	rmb,    err := strconv.Atoi(args[2])
	shell,  err := strconv.Atoi(args[3])
	status, err := strconv.Atoi(args[4])
	time,   err := strconv.Atoi(args[6])
	if err != nil {
		log.Errorf("parse topup argument error: %s", err)
		return
	}
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_web_db")
	UserTopup := UserTopUpRecord{0, userId, topupId, rmb, shell, status, payType, time}
	_, err = mysqlORM.Insert(&UserTopup)
	if err != nil {
		log.Errorf("Insert new transaction error: %s", err)
	}
	return
}

func userBuyVip(args SqlArgList) (err error) {
	orderId, goodsId := args[1], args[4]
	userId, err := strconv.Atoi(args[0])
	shell,  err := strconv.Atoi(args[2])
	time,   err := strconv.Atoi(args[3])
	start,  err := strconv.Atoi(args[5])
	end,    err := strconv.Atoi(args[6])
	if err != nil {
		log.Errorf("parse BuyVip argument error: %s", err)
		return
	}
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_web_db")
	UserBuyVip := UserBuyMembership{0,userId, orderId, shell, time, goodsId, start, end, ""}
	_, err = mysqlORM.Insert(&UserBuyVip)
	if err != nil {
		log.Errorf("Insert new transaction error: %s", err)
	}
	return
}

func topupHistory(args SqlArgList) (rows []*dbproto.OneRow, err error) {
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_web_db")
	var historyList []orm.ParamsList
	_, err = mysqlORM.QueryTable("topup_record_t").Filter("userId", userId).OrderBy("-transactionTime").ValuesList(&historyList,
		"transactionId", "rmb", "shell", "status", "payType", "transactionTime")
	if err != nil {
		log.Errorf("query topup history error: %s", err)
		return
	}
	for _, record := range historyList {
		row  := serializeRowDate(record)
		rows = append(rows, row)
	}
	return
}

func userVipHistory(args SqlArgList) (rows []*dbproto.OneRow, err error) {
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_web_db")
	var historyList []orm.ParamsList
	_, err = mysqlORM.QueryTable("purchase_record_t").Filter("userId", userId).OrderBy("-orderTime").ValuesList(&historyList,
		"orderId", "shell", "orderTime", "goodsId", "startTime", "endTime")
	if err != nil {
		log.Errorf("query topup history error: %s", err)
		return
	}
	for _, record := range historyList {
		row  := serializeRowDate(record)
		rows = append(rows, row)
	}
	return
}

func init() {
	orm.RegisterModel(new(UserTopUpShell), new(UserTopUpRecord), new(UserBuyMembership))
	RegisterRunSqlCB("user.buy.shell",  userTopup)
	RegisterRunSqlCB("user.buy.vip",   userBuyVip)
	RegisterQueryCB("topup.history", topupHistory)
	RegisterQueryCB("buy.history", userVipHistory)
}


