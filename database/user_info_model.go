package database

import (
	"fmt"
	"strconv"
	"db-server/proto/dbproto"
	"github.com/astaxie/beego/orm"
)

type User struct{
	UserID     int    `orm:"column(user_id);auto;pk"`
	Username   string `orm:"column(username);unique"`
	Password   string `orm:"column(password)"`
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
	return "user_info_tab"
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

func deleteAccount(args SqlArgList) (err error){
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	user := User{ UserID:userId}
	_, err = mysqlORM.Delete(&user)
	return
}

func changeUserState(args SqlArgList) (err error) {
	userId, err := strconv.Atoi(args[0])
	state , err := strconv.Atoi(args[1])
	mysqlORM := orm.NewOrm()
	user := User{
		UserID : userId,
		Enabled : state,
	}
	_, err = mysqlORM.Update(&user, "enabled")
	return
}

func getUserAllInfo(args SqlArgList) (rows []*dbproto.OneRow, err error) {
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	var userInfoList []interface{}
	var r orm.RawSeter
	r = mysqlORM.Raw("select i.username, i.createdate, s.shell, u.nickname,u.level,u.avatar,m.membership,m.expiry_date,u.exp " +
		"from user_info_tab u left join user_shell_tab s on s.user_id=u.user_id left join membership_tab m on m.user_id=u.user_id left join user_tab i on i.user_id=u.user_id where u.user_id=?", userId)
	var basicInfoList []orm.ParamsList
	_, err = r.ValuesList(&basicInfoList)
	if err != nil {
		return
	}
	userInfoList = basicInfoList[0]
	//query user extra info
	var extraInfoList []orm.ParamsList
	_, err = mysqlORM.QueryTable("user_extra_info_t").Filter("user_id", userId).ValuesList(&extraInfoList)
	if err != nil {
		return
	}
	userInfoList = append(userInfoList, extraInfoList[0][1:]...)
	row := serializeRowDate(userInfoList)
	rows = append(rows, row)
	return
}

func changePassword(args SqlArgList) (err error)  {
	userId, err := strconv.Atoi(args[0])
	passwd := args[1]
	mysqlORM := orm.NewOrm()
	user := User{
		UserID : userId,
		Password : passwd,
	}
	_, err = mysqlORM.Update(&user, "password")
	if err != nil {
		log.Errorf("user(%d) change password error : %s", userId, err)
	}
	return
}

func updateUserInfo(args SqlArgList) (err error) {
	userId, _   := strconv.Atoi(args[0])
	nickname, avatar := args[1], args[2]
	mysqlORM := orm.NewOrm()
	user := UserBasicInfo{
		UserID: userId,
		Nickname: nickname,
		Avatar: avatar,
	}
	_, err = mysqlORM.Update(&user, "nickname", "avatar")
	log.Errorf("update user info error : %s", err)
	return
}

func updateExtraInfo(args SqlArgList) (err error) {
	userId, err := strconv.Atoi(args[0])
	gender, err := strconv.Atoi(args[2])
	birth , err := strconv.Atoi(args[3])
	phone,  address, company, signature := args[1],  args[4], args[5], args[6]
	mysqlORM := orm.NewOrm()
	user := UserExtraInfo{UserID:userId}
	created, id, err := mysqlORM.ReadOrCreate(&user, "user_id")
	if err == nil {
		if created {
			log.Infoln("New Insert an object. Id:", id)
		} else {
			user = UserExtraInfo{ userId, phone, gender, birth, address, company, signature }
			_, err = mysqlORM.Update(&user)
			if err != nil {
				log.Errorf("update user(%d) extra info error : %s", userId, err)
			}
		}
	}else{
		log.Errorf("read or create user(%d) extra info error : %s", userId, err)
	}
	return
}

func getUserSyncEmail(args SqlArgList) (rows []*dbproto.OneRow, err error) {
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	mysqlORM.Using("yomail_email_db")
	var userEmails []*UserEmails
	_, err = mysqlORM.QueryTable("user_email_tab").Filter("user_id", userId).All(&userEmails, "email")
	if err != nil  {
		return
	}
	var emailList []interface{}
	for _, user := range userEmails{
		emailList = append(emailList, user.Email)
	}
	row := serializeRowDate(emailList)
	rows = append(rows, row)
	return
}

func getUserPrivacy(args SqlArgList) (rows []*dbproto.OneRow, err error) {
	userId, err := strconv.Atoi(args[0])
	mysqlORM := orm.NewOrm()
	userPrivacy := UserPrivacy{UserID: userId}
	err = mysqlORM.Read(&userPrivacy)
	if err != nil{
		log.Debugf("ID: %d, ERR: %v\n", userId, err)
		err = fmt.Errorf("ID: %d, ERR: %v\n", userId, err)
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
	RegisterQueryCB("get.all.info", getUserAllInfo)
	RegisterQueryCB("get.user.privacy",  getUserPrivacy)
	RegisterQueryCB("get.user.emails", getUserSyncEmail)
	RegisterRunSqlCB("delete.user",  deleteAccount)
	RegisterRunSqlCB("change.pwd",  changePassword)
	RegisterRunSqlCB("user.update", updateUserInfo)
	RegisterRunSqlCB("change.user.state", changeUserState)
	RegisterRunSqlCB("update.extra.info", updateExtraInfo)
}