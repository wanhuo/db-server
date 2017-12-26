package database

import (
	"fmt"
	"strconv"
	"github.com/astaxie/beego/orm"
	"db-server/proto/dbproto"
)

type User struct{
	UserID     int    `orm:"column(user_id);auto;pk"`
	Username   string `orm:"column(username);unique"`
	Password   string `orm:"column(Password)"`
	Enabled    int	  `orm:"column(enabled);null"`
	CreateDate int	  `orm:"column(createdate)"`
	Ip         string `orm:"column(ip);null"`
}

func (u *User) TableName() string {
	return "user_tab"
}

type UserBasicInfo struct {
	UserID    int    `orm:"column(user_id);pk"`
	Nickname  string `orm:"column(nickname);null"`
	Level     int    `orm:"column(level);"`
	Avatar    string `orm:"column(avatar);null"`
	Exp       int	 `orm:"column(exp);default(0)"`
}

func (u *UserBasicInfo) TableName() string {
	return "user_info_t"
}

type UserExtraInfo struct {
	UserID    int    `orm:"column(user_id);pk"`
	Phone     string `orm:"column(phone);unique"`
	Gender    int    `orm:"column(gender);null"`
	Birth     int    `orm:"column(birth);null"`
	Address   string `orm:"column(address);null"`
	Company   string `orm:"column(company);null"`
	Signature string `orm:"column(signature);null"`
}
func (u *UserExtraInfo) TableName() string {
	return "user_extra_info_t"
}

type UserEmails struct{
	RecordID     int    `orm:"column(record_id);auto;pk"`
	UserID       int    `orm:"column(user_id);index"`
	Email        string `orm:"column(email);index"`
	AccountInfo  byte   `orm:"column(account_info)"`
	Salt         string `orm:"column(salt)"`
	ExpiryDate   int    `orm:"column(expiry_date)"`
}
func (u *UserEmails) TableName() string {
	return "user_email_tab"
}

type UserPrivacy struct {
	UserID          int    `orm:"column(user_id);pk"`
	ReportVisible   int    `orm:"column(report_visible)"`
	UsernameVisible int    `orm:"column(username_visible)"`
}
func (u *UserPrivacy) TableName() string {
	return "user_privacy_tab"
}

func serializeRowDate(params []interface{}) (row *dbproto.OneRow){
	var fieldList []*dbproto.OneField
	for _, value := range params{
		var valueStr string
		switch value.(type) {
		case int:
			valueStr = strconv.Itoa(value.(int))
		default:
			valueStr = value.(string)
		}
		field := dbproto.OneField{
			Value: []byte(valueStr),
		}
		fieldList = append(fieldList, &field)
	}
	row = &dbproto.OneRow{
		OneField:fieldList,
	}
	return
}

func GetUserPrivacy(args sqlArgs) (rows []*dbproto.OneRow, err error) {
	userId, errMsg := strconv.Atoi(args[0].(string))
	o := orm.NewOrm()
	userPrivacy := UserPrivacy{UserID: userId}
	errMsg = nil
	errMsg = o.Read(&userPrivacy)
	if errMsg != nil{
		log.Debugf("ID: %d, ERR: %v\n", userId, errMsg)
		err = fmt.Errorf("ID: %d, ERR: %v\n", userId, errMsg)
	}else{
		log.Debugf("%d user privacy: report visible is %d, username visible is %d",
			userId, userPrivacy.ReportVisible, userPrivacy.UsernameVisible)
		privacyList := []interface{}{
			userPrivacy.ReportVisible,
			userPrivacy.UsernameVisible,
		}
		row := serializeRowDate(privacyList)
		rows = append(rows, row)
	}
	return
}

func init() {
	orm.RegisterModel(new(User), new(UserBasicInfo), new(UserExtraInfo), new(UserEmails), new(UserPrivacy))
	RegisterQueryCB("get.user.privacy", GetUserPrivacy)
}



