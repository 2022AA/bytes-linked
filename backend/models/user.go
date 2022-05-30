package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/2022AA/bytes-linked/backend/pkg/util"
	"github.com/2022AA/bytes-linked/backend/transaction"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/gorm"

	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
)

type UserOrder struct {
	ID            int    `gorm:"primary_key" json:"id"`
	Type          int    `gorm:"column:type;not null" json:"type"`
	OrderSN       string `gorm:"column:order_sn;not null" json:"order_sn"`
	Content       string `gorm:"column:content;not null" json:"content"`
	CreateTime    string `gorm:"column:create_time;not null" json:"create_time"`
	Uid           int    `gorm:"column:uid;not null" json:"uid"`
	ToId          int    `gorm:"column:to_id;not null" json:"to_id"`
	TransactionId string `gorm:"column:transaction_id;not null" json:"TransactionId"`
	//Username   string `json:"username"`
	//ToUsername string `json:"to_username"`
}

type InviteCode struct {
	ID         int    `gorm:"primary_key" json:"-"`
	Code       string `gorm:"column:code;not null"`
	Status     int    `gorm:"column:status;not null"`
	CreateTime string `gorm:"column:create_time;not null"`
	LastUpdate string `gorm:"column:last_update;not null"`
}

// User : 用户表model
type User struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Username     string `gorm:"column:user_name;not null" json:"Username"`
	Userpwd      string `gorm:"column:user_pwd;not null" json:"-"`
	SignupAt     string `gorm:"column:signup_at;" json:"SignupAt"`
	LastActiveAt string `gorm:"column:last_active;" json:"LastActive"`
	Status       int    `gorm:"column:status;" json:"-"`
	InviteCode   string `gorm:"column:invite_code;" json:"InviteCode"`
	Phone        string `gorm:"column:phone;" json:"Phone"`
	AvatarUrl    string `gorm:"column:avatar_url;" json:"AvatarUrl"`
	Balance      int    `gorm:"column:balance;" json:"Balance"`
}

func CreateInviteCode(code string) error {
	inviteCode := InviteCode{
		Code:       code,
		Status:     0,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := db.Table("invite_code").Create(&inviteCode)
	if err != nil {
		v2.Warnln(err)
	}
	return nil
}

func QueryInviteCode(code string) (*InviteCode, error) {
	var inviteCode InviteCode
	err := db.Table("invite_code").Where("code = ?", code).First(&inviteCode).Error
	if err != nil {
		return nil, err
	}
	return &inviteCode, nil
}

func IsValid(code string) (bool, error) {
	var inviteCode InviteCode
	err := db.Table("invite_code").Where("code = ?", code).First(&inviteCode).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Error("Failed to check invite code exist, err:" + err.Error())
			return false, err
		} else {
			return false, fmt.Errorf("invalid invite code, code: %s", code)
		}
	}
	return true, nil

}

// UserSignup : 通过用户名及密码完成user表的注册操作
func UserSignup(
	username string, passwd string, inviteCode string,
	phone string, avatarUrl string,
) bool {
	user := User{
		Username:     username,
		Userpwd:      passwd,
		SignupAt:     time.Now().Format("2006-01-02 15:04:05"),
		LastActiveAt: time.Now().Format("2006-01-02 15:04:05"),
		Phone:        phone,
		AvatarUrl:    avatarUrl,
	}

	secretFiles, err := GenerateUserSecretFile(username)
	if err != nil {
		v2.Warnf("UserSignup | failed to gen user secret file, err: %v, username: %s", err, username)
		return false
	}

	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			v2.Warnln(err)
			tx.Rollback()
		}
	}()
	//// 先更新激活码状态
	////var inviteCode InviteCode
	//
	//fmt.Println(inviteCode, "111111111111")
	//err = tx.Table("invite_code").Select("id").Where("code = ? ", inviteCode).First(&InviteCode{}).Error
	//if err != nil {
	//	v2.Warnln(err)
	//	tx.Rollback()
	//	return false
	//}
	userInviteCode := InviteCode{
		Code:       util.GenInvite(rand.Uint64()),
		Status:     0,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = tx.Table("invite_code").Create(&userInviteCode).Error
	if err != nil {
		logging.Error("Failed to create user, err:" + err.Error())
		tx.Rollback()
		return false
	}
	user.InviteCode = userInviteCode.Code

	// 创建用户
	if err = tx.Create(&user).Error; err != nil {
		logging.Error("Failed to create user, err:" + err.Error())
		tx.Rollback()
		return false
	}
	// 写入公私钥
	for _, secretFile := range secretFiles {
		err = tx.Create(secretFile).Error
		if err != nil {
			logging.Error("Failed to create secret file, err:" + err.Error())
			tx.Rollback()
			return false
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return false
	}

	return true
}

// UserSignin : 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	var user User
	err := db.Where("status = 0").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		logging.Error("Failed to signIn user, err:" + err.Error())
		return false
	}
	if user.Userpwd == encpwd {
		sql := fmt.Sprintf("UPDATE user set last_active='%s', balance=%d where user_name = '%s' limit 1",
			time.Now().Format("2006-01-02 15:04:05"), user.Balance+2, username)
		rawDb := db.Raw(sql)
		rawDb.Scan(&UserFile{})
		if rawDb.Error != nil {
			logging.Error(fmt.Sprintf("Failed to update user lastActive, err:%s", rawDb.Error.Error()))
		}
		return true
	}
	return false
}

// GetUserInfo : 查询用户信息
func GetUserInfo(username string) (User, error) {
	var user User
	err := db.Table("user").Where("status = 0").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		logging.Error("Failed to get userInfo, err:" + err.Error())
		return user, err
	}
	return user, nil
}

// DeleteUser : 删除用户
func DeleteUser(username string) error {
	sql := fmt.Sprintf("DELETE FROM user where user_name= '%s' limit 1", username)
	rawDb := db.Raw(sql)
	rawDb.Scan(&File{})
	if rawDb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to delete user with user_name: %s, err: %s",
			username, rawDb.Error.Error()))
		return rawDb.Error
	}
	return nil
}

// UserExist : 查询用户是否存在
func UserExist(username string) (bool, error) {
	var user User
	err := db.Table("user").Select("id").Where("status = 0").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Error("Failed to check user exist, err:" + err.Error())
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}

// QueryUserInfoList : 批量获取用户信息
func QueryUserInfoList(offset, limit int) ([]User, error) {
	var users []User
	err := db.Table("user").Where("status = 0").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to get userInfoList, err:%s", err.Error()))
		return nil, err
	}
	return users, nil
}

func BatchQueryUserInfoList(uids []int) ([]User, error) {
	var users []User
	err := db.Table("user").Where("id in ?").Find(&users).Error
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to get userInfoList, err:%s", err.Error()))
		return nil, err
	}
	return users, nil
}

// UpdateUserPassword : 修改用户密码
func UpdateUserPassword(username string, passwd string) (bool, error) {
	var user User
	err := db.Where("status = 0").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		logging.Error("Failed to update user password, err:" + err.Error())
		return false, err
	}
	if err = db.Model(&user).Update("user_pwd", passwd).Error; err != nil {
		logging.Error("Failed to update, err:" + err.Error())
		return false, err
	}
	return true, nil
}

var transferLock sync.Mutex

func FileTransaction(uid int, toUid int, transactionFileId int, privateKey string) (bool, error) {
	if uid == toUid {
		return true, nil
	}
	transferLock.Lock()
	defer transferLock.Unlock()

	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			v2.Warnln(err)
			tx.Rollback()
		}
	}()

	now := time.Now()
	// 进行转移
	fmt.Println(uid, toUid, transactionFileId)
	rowCnt := tx.Model(&File{}).Where("id = ? and owner_uid = ?", transactionFileId, uid).
		Update("owner_uid", toUid).RowsAffected
	fmt.Println("file transaction update >>>>>>>>>...")
	if rowCnt <= 0 {
		tx.Rollback()
		return false, nil
	}

	// 生成流水
	order := UserOrder{
		Type:       1,
		OrderSN:    bson.NewObjectId().Hex(),
		Content:    "",
		CreateTime: util.TimeString(now),
		Uid:        uid,
		ToId:       toUid,
	}
	err := tx.Table("user_order").Create(&order).Error
	fmt.Println("order create >>>>>>>>>...")
	if err != nil {
		logging.Error("Failed to create order, err:" + err.Error())
		tx.Rollback()
		return false, err
	}

	// 生成区块
	contentBytes, _ := json.Marshal(&order)
	fmt.Println("transaction start >>>>>>>>>...")
	transactionId, err := transaction.Transaction(
		fmt.Sprintf("%d", transactionFileId), string(contentBytes), privateKey,
	)
	if err != nil {
		logging.Error("Failed to create Transaction, err:" + err.Error())
		tx.Rollback()
		return false, err
	}

	fmt.Println("transaction end >>>>>>>>>...")
	cnt := tx.Model(&UserOrder{}).Where("order_sn = ?", order.OrderSN).
		Update("transaction_id", transactionId).RowsAffected
	if cnt <= 0 {
		logging.Error("Failed to update transaction id to user_order, transaction_id: %s", transactionId)
	}
	fmt.Println("order transaction_id inserted >>>>>>>>>...")

	cnt = tx.Model(&File{}).Where("id = ?", transactionFileId).Update("transaction_id", transactionId).RowsAffected
	if cnt <= 0 {
		logging.Error("Failed to update transaction id to file, transaction_id: %s", transactionId)
	}
	fmt.Println("file transaction_id insertd >>>>>>>>>...")

	tx.Commit()

	return true, nil
}

func GetOrders(uid int) ([]*UserOrder, error) {
	var orders []UserOrder
	err := db.Table("user_order").Where("uid = ?", uid).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateUserAvatar : 修改用头像
func UpdateUserAvatar(username string, avatarUrl string) (*User, error) {
	var user User
	cnt := db.Table("user").Where("user_name = ?", username).Update("avatar_url", avatarUrl).RowsAffected
	if cnt <= 0 {
		err := fmt.Errorf("failed to update user avatar, username: %s", username)
		return nil, err
	}
	err := db.Table("user").Where("user_name = ?", username).First(&user).Error
	if err != nil {
		v2.Errorf("failed to get user info")
		return nil, err
	}
	return &user, nil
}
