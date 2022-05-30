package models

import (
	"crypto/ecdsa"
	"fmt"

	logger_v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
	"github.com/2022AA/bytes-linked/backend/pkg/util"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type SecretType int

const SecretType_Public = 1
const SecretType_Private = 2

// UserFile : 用户文件表结构体
type UserSecretFile struct {
	ID         int        `gorm:"primary_key" json:"-"`
	UserName   string     `gorm:"column:user_name;not null" json:"Username"`
	FileHash   string     `gorm:"column:file_sha1;not null" json:"FileHash"`
	FileName   string     `gorm:"column:file_name;not null" json:"FileName"`
	FileSize   int64      `gorm:"column:file_size;not null" json:"FileSize"`
	CreateTime string     `gorm:"column:create_time;" json:"CreateTime"`
	LastUpdate string     `gorm:"column:last_update;" json:"LastUpdate"`
	Status     int64      `gorm:"column:status;not null" json:"status"`
	Type       SecretType `gorm:"column:type;not null" json:"type"`
	Data       string     `gorm:"column:data;not null" json:"data"`
}

func GenerateUserSecretFile(username string) ([]*UserSecretFile, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)[2:]
	logger_v2.Debugf("GenerateUserAccount | "+
		"userName: %s, private key: %s, private_key_length: %d",
		username, privateKeyStr, len(privateKeyBytes),
	)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("GenerateUserAccount | " +
			"cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		logger_v2.Errorln(err)
		return nil, err
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyStr := hexutil.Encode(publicKeyBytes)[4:]
	logger_v2.Debugf("GenerateUserAccount | userName: %s, publick key: %s",
		username, publicKeyStr) // publicKey in hex without "0x"

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	logger_v2.Debugf("GenerateUserAccount | userName: %d, address: %s",
		username, address)

	timeNowStr := util.TimeNowString()
	retFiles := []*UserSecretFile{
		{
			UserName:   username,
			FileHash:   util.Sha1(privateKeyBytes),
			FileName:   "private.pem",
			FileSize:   int64(len(privateKeyBytes)),
			CreateTime: timeNowStr,
			LastUpdate: timeNowStr,
			Status:     0,
			Type:       SecretType_Private,
			Data:       privateKeyStr,
		},
		{
			UserName:   username,
			FileHash:   util.Sha1(publicKeyBytes),
			FileName:   "public.pem",
			FileSize:   int64(len(publicKeyBytes)),
			CreateTime: timeNowStr,
			LastUpdate: timeNowStr,
			Status:     0,
			Type:       SecretType_Public,
			Data:       publicKeyStr,
		},
	}
	return retFiles, nil
}

func QuerySecretFile(username string, secretType SecretType) (*UserSecretFile, error) {
	var secretFile UserSecretFile
	err := db.Table("user_secret_file").Where(
		"user_name = ? and status = 0 and type = ?",
		username, secretType,
	).First(&secretFile).Error
	if err != nil {
		return nil, err
	}
	return &secretFile, nil
}

// CheckFileExist : 检查全局文件表中文件元信息是否存在
//func CheckUserFileExist(username, filename, filehash string) bool {
//	var t UserFile
//	err := db.Select("id").Where("status = 0").Where("file_sha1 = ? and file_name = ? and user_name = ?", filehash, filename, username).First(&t).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		logging.Error("Failed to get file, err: " + err.Error())
//		return false
//	}
//	if t.ID > 0 {
//		return true
//	}
//	return false
//}
//
//// OnUserFileUploadFinished : 更新用户文件表
//func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) error {
//	userFile := UserFile{
//		UserName:   username,
//		FileHash:   filehash,
//		FileName:   filename,
//		FileSize:   filesize,
//		UploadAt:   time.Now().Format("2006-01-02 15:04:05"),
//		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
//	}
//	if err := db.Create(&userFile).Error; err != nil {
//		logging.Error("Failed to insert, err:" + err.Error())
//		return err
//	}
//	return nil
//}
//
//// QueryUserFileMetas : 批量获取用户文件信息
//func QueryUserFileMetas(username, filenameFilter string, reverse bool, offset, limit int) ([]UserFile, error) {
//	var userFiles []UserFile
//	orderSQL := "last_update DESC"
//	if reverse {
//		orderSQL = "last_update ASC"
//	}
//	err := db.Table("user_file").Where("status = 0").
//		Where("user_name = ? AND file_name LIKE ?", username, "%"+filenameFilter+"%").
//		UserOrder(orderSQL).Offset(offset).Limit(limit).Find(&userFiles).Error
//	if err != nil {
//		logging.Error(fmt.Sprintf("Failed to get userFileMetas for user:%s, err:%s", username, err.Error()))
//		return nil, err
//	}
//	return userFiles, nil
//}
//
//// DeleteUserFile : 删除文件
//func DeleteUserFile(username, filehash, filename string) error {
//	sql := fmt.Sprintf("DELETE FROM user_file where user_name = '%s' and file_sha1= '%s' and file_name = '%s' limit 1", username, filehash, filename)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error(fmt.Sprintf("Failed to delete userFile with fileHash:%s for user:%s, err:%s",
//			filehash, username, rawDb.Error.Error()))
//		return rawDb.Error
//	}
//	return nil
//}
//
//func DeleteUserFileByUsername(username string) error {
//	sql := fmt.Sprintf("DELETE FROM user_file where user_name = '%s'", username)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error(fmt.Sprintf("Failed to delete userFile with user:%s, err:%s",
//			username, rawDb.Error.Error()))
//		return rawDb.Error
//	}
//	return nil
//}
//
//// RenameUserFileName : 文件重命名
//func RenameUserFileName(username, filehash, oldfilename, filename string) bool {
//	sql := fmt.Sprintf("UPDATE user_file set file_name='%s', last_update='%s' where status = 0 and user_name = '%s' and file_sha1= '%s' and file_name = '%s' limit 1",
//		filename, time.Now().Format("2006-01-02 15:04:05"), username, filehash, oldfilename)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error(fmt.Sprintf("Failed to rename userFile for user:%s, err:%s", username, rawDb.Error.Error()))
//		return false
//	}
//	return true
//}
//
//// QueryUserFileMeta : 获取用户单个文件信息
//func QueryUserFileMeta(username, filehash, filename string) (*UserFile, error) {
//	var ufile UserFile
//	err := db.Where("status = 0").Where("user_name = ? and file_sha1 = ? and file_name = ?", username, filehash, filename).First(&ufile).Error
//	if err != nil && !gorm.IsRecordNotFoundError(err) {
//		logging.Error(fmt.Sprintf("Failed to query userFileMeta for user:%s, err:%s", username, err.Error()))
//		return nil, err
//	}
//	return &ufile, nil
//}
//
//// 文件共享部分
//// QueryShareFileMetas : 批量获取共享文件信息
//func QueryShareFileMetas(shareStatus int, usernameFilter, filenameFilter string, reverse bool, offset, limit int) ([]UserFile, error) {
//	var userFiles []UserFile
//	orderSQL := "last_update DESC"
//	if reverse {
//		orderSQL = "last_update ASC"
//	}
//	err := db.Table("user_file").Where("status = 0 and share_status = ?", shareStatus).
//		Where("user_name LIKE ? AND file_name LIKE ?", "%"+usernameFilter+"%", "%"+filenameFilter+"%").
//		UserOrder(orderSQL).Offset(offset).Limit(limit).Find(&userFiles).Error
//	if err != nil {
//		logging.Error("Failed to get shareFileMetas, err: ", err)
//		return nil, err
//	}
//	return userFiles, nil
//}
//
//// ShareFile: 共享文件
//func ShareFile(username, filehash string) bool {
//	sql := fmt.Sprintf("UPDATE user_file set share_status=%d where status = 0 and user_name = '%s' and file_sha1= '%s' limit 1",
//		1, username, filehash)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error("Failed to share file, err: ", rawDb.Error)
//		return false
//	}
//	return true
//}
//
//// CancelShareFile: 取消共享文件
//func CancelShareFile(username, filehash string) bool {
//	sql := fmt.Sprintf("UPDATE user_file set share_status=%d where status = 0 and user_name = '%s' and file_sha1= '%s' limit 1",
//		0, username, filehash)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error("Failed to cancel share file, err: ", rawDb.Error)
//		return false
//	}
//	return true
//}
//
//// AuditShareFile: 审核文件
//func AuditShareFile(username, filehash string) bool {
//	sql := fmt.Sprintf("UPDATE user_file set share_status=%d where status = 0 and user_name = '%s' and file_sha1= '%s' limit 1",
//		2, username, filehash)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error("Failed to audit share file, err: ", rawDb.Error)
//		return false
//	}
//	return true
//}
//
//// CancelAuditShareFile: 取消审核文件
//func CancelAuditShareFile(username, filehash string) bool {
//	sql := fmt.Sprintf("UPDATE user_file set share_status=%d where status = 0 and user_name = '%s' and file_sha1= '%s' limit 1",
//		1, username, filehash)
//	rawDb := db.Raw(sql)
//	rawDb.Scan(&UserFile{})
//	if rawDb.Error != nil {
//		logging.Error("Failed to cancel audit share file, err: ", rawDb.Error)
//		return false
//	}
//	return true
//}
