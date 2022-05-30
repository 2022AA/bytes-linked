package models

import (
	"fmt"
	"time"

	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	"github.com/jinzhu/gorm"
)

// File : 文件表结构体
type File struct {
	ID            int    `gorm:"primary_key" json:"file_id"`
	FileSHA1      string `gorm:"column:file_sha1;not null" json:"file_sha1,omitempty"`
	FileName      string `gorm:"column:file_name;not null" json:"file_name"`
	FileSize      int64  `gorm:"column:file_size;not null" json:"file_size,omitempty"`
	FileAddr      string `gorm:"column:file_addr;not null" json:"file_addr"`
	ImgAddr       string `gorm:"column:img_addr;not null" json:"img_addr"`
	OwnerUID      int    `gorm:"column:owner_uid;not null" json:"owner_uid,omitempty"`
	CreatorUID    int    `gorm:"column:creator_uid;not null" json:"creator_uid,omitempty"`
	CreateAt      string `gorm:"column:create_at;" json:"create_at,omitempty"`
	UpdateAt      string `gorm:"column:update_at;" json:"update_at,omitempty"`
	Status        int64  `gorm:"column:status;" json:"status,omitempty"`
	TransactionId string `gorm:"column:transaction_id;" json:"TransactionId"`
}

type FileRead struct {
	ID            int    `gorm:"primary_key" json:"file_id"`
	FileSHA1      string `gorm:"column:file_sha1;not null" json:"file_sha1,omitempty"`
	FileName      string `gorm:"column:file_name;not null" json:"file_name"`
	FileSize      int64  `gorm:"column:file_size;not null" json:"file_size,omitempty"`
	FileAddr      string `gorm:"column:file_addr;not null" json:"file_addr"`
	ImgAddr       string `gorm:"column:img_addr;not null" json:"img_addr"`
	OwnerUID      int    `gorm:"column:owner_uid;not null" json:"owner_uid,omitempty"`
	CreatorUID    int    `gorm:"column:creator_uid;not null" json:"creator_uid,omitempty"`
	Username      string `gorm:"column:username;not null" json:"username,omitempty"`
	AvatarUrl     string `gorm:"column:avatar_url;" json:"avatar_url,omitempty"`
	LikeCnt       int    `gorm:"column:like_cnt;not null" json:"like_cnt"`
	ArTag         bool   `gorm:"column:ar_tag;not null" json:"ar_tag"`
	CreateAt      string `gorm:"column:create_at;" json:"create_at,omitempty"`
	UpdateAt      string `gorm:"column:update_at;" json:"update_at,omitempty"`
	Status        int64  `gorm:"column:status;" json:"status,omitempty"`
	TransactionId string `gorm:"column:transaction_id;" json:"TransactionId"`
}

/*
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `img_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '缩略文件存储位置',
  `owner_uid` int(11) DEFAULT '0' COMMENT '文件拥有者',
  `creator_uid` int(11) DEFAULT '0' COMMENT '文件创作者',
  `like_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
  `create_at` datetime default NOW() COMMENT '创建日期',
  `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(2公开发布/1可用/0禁用/-1已删除等状态)',
*/

func CommitFileInfo(file File) (int, error) {
	err := db.Table("file").Create(&file).Error
	return file.ID, err
}

func GetFile(ownerUid int, fileId int) (*FileRead, error) {
	var fileRead FileRead
	err := db.Table("file").Where("id = ? and owner_uid = ?", fileId, ownerUid).First(&fileRead).Error
	if err != nil {
		return nil, err
	}
	return &fileRead, nil
}

func ListOwnFileInfo(owner int, page, limit int) ([]FileRead, error) {
	files := []FileRead{}
	err := db.Raw("SELECT f.id AS id, f.file_sha1 AS file_sha1, f.file_name AS file_name, f.file_addr AS file_addr, f.img_addr AS img_addr, f.like_cnt AS like_cnt, f.ar_tag AS ar_tag, f.owner_uid AS owner_uid, u.user_name AS username, u.avatar_url AS avatar_url FROM file AS f JOIN user AS u ON f.owner_uid = u.id WHERE f.status > 0 AND f.owner_uid = ? ORDER BY f.like_cnt DESC", owner).Find(&files).Error
	return files, err
}

func ListPublicFileInfo(page, limit int) ([]FileRead, error) {
	files := []FileRead{}
	err := db.Raw("SELECT f.id AS id, f.file_sha1 AS file_sha1, f.file_name AS file_name, f.file_addr AS file_addr, f.img_addr AS img_addr, f.like_cnt AS like_cnt, f.ar_tag AS ar_tag, f.owner_uid AS owner_uid, u.user_name AS username, u.avatar_url AS avatar_url FROM file AS f JOIN user AS u ON f.owner_uid = u.id WHERE f.status = 2 ORDER BY f.like_cnt DESC").Find(&files).Error
	return files, err
}

func PublicFile(uid, fileID int) error {
	err := db.Table("file").Where("owner_uid = ? AND id = ? AND status > 0", uid, fileID).Update("status", "2").Error
	return err
}

func LikeFile(uid, fileID, likeNum int) (int, string) {
	fileCnt := 0
	db.Table("file").Select("id").Where("id = ? AND owner_uid != ? AND status = 2", fileID, uid).Count(&fileCnt)
	if fileCnt == 0 {
		return -1, "不能给自己的点赞"
	}
	userCnt := 0
	db.Table("user").Select("id").Where("id = ? AND balance >= ?", uid, likeNum).Count(&userCnt)
	if userCnt == 0 {
		return -1, "点赞余额不足"
	}
	if db.Table("user").Where("id = ?", uid).Update("balance", gorm.Expr("balance - ?", likeNum)).Error != nil {
		return -1, "扣除点赞数失败"
	}
	if db.Table("file").Where("id = ?", fileID).Update("like_cnt", gorm.Expr("like_cnt + ?", likeNum)).Error != nil {
		return -1, "添加点赞数失败"
	}
	wrap := struct{ Balance int }{}
	db.Table("user").Select("balance").Where("id = ?", uid).First(&wrap)
	return wrap.Balance, "点赞成功"
}

// CheckFileExist : 检查全局文件表中文件元信息是否存在
func CheckFileExist(filehash string) bool {
	var t File
	err := db.Select("id").Where("file_sha1 = ?", filehash).First(&t).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error("Failed to get file, err: " + err.Error())
		return false
	}
	if t.ID > 0 {
		return true
	}
	return false
}

// OnFileUploadFinished : 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) error {
	tableFile := File{
		FileSHA1: filehash,
		FileName: filename,
		FileSize: filesize,
		FileAddr: fileaddr,
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := db.Create(&tableFile).Error; err != nil {
		logging.Error("Failed to insert, err: " + err.Error())
		return err
	}
	return nil
}

// RenameFileName : 文件重命名
func RenameFileName(filehash, oldfilename, filename string) bool {
	sql := fmt.Sprintf("UPDATE file set file_name='%s', last_update='%s' where status = 0 and file_sha1= '%s' and file_name = '%s' limit 1",
		filename, time.Now().Format("2006-01-02 15:04:05"), filehash, oldfilename)
	rawDb := db.Raw(sql)
	rawDb.Scan(&UserFile{})
	if rawDb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to rename userFile for err:%s", rawDb.Error.Error()))
		return false
	}
	return true
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash, filename string) (*File, error) {
	var tableFile File
	err := db.Where("status = 0").Where("file_sha1 = ? and file_name = ?", filehash, filename).First(&tableFile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("Failed to get fileMeta with hash: %s , err: %s", filehash, err.Error()))
		return &tableFile, err
	}
	return &tableFile, nil
}

// GetFileMetaList : 从mysql批量获取文件元信息
func GetFileMetaList(limit int) ([]File, error) {
	var tableFiles []File
	err := db.Table("file").Where("status = 0").Limit(limit).Find(&tableFiles).Error
	if err != nil {
		logging.Error("Failed to get fileMetaList, err: " + err.Error())
		return nil, err
	}
	return tableFiles, nil
}

// DeleteFile : 删除文件
func DeleteFile(filehash, filename string) error {
	sql := fmt.Sprintf("DELETE FROM file where file_sha1= '%s' and file_name = '%s' limit 1", filehash, filename)
	rawDb := db.Raw(sql)
	rawDb.Scan(&File{})
	if rawDb.Error != nil {
		logging.Error(fmt.Sprintf("Failed to delete file with fileHash: %s, fileName: %s, err: %s",
			filehash, filename, rawDb.Error.Error()))
		return rawDb.Error
	}
	return nil
}
