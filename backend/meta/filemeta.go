package meta

import (
	"errors"
	"sort"

	models2 "github.com/2022AA/bytes-linked/backend/models"
	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	"github.com/jinzhu/gorm"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : 新增/更新缓存的文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB : 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) error {
	if models2.CheckFileExist(fmeta.FileSha1) {
		err := errors.New("File has been uploaded before. ")
		logging.Error(err)
		return err
	}
	return models2.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// UpdateUserFileMetaDB : 新增/更新用户文件元信息到mysql中
func UpdateUserFileMetaDB(username, filehash, filename string, filesize int64) error {
	if models2.CheckUserFileExist(username, filename, filehash) {
		err := errors.New("File has been uploaded before. ")
		logging.Error(err)
		return err
	}
	return models2.OnUserFileUploadFinished(username, filehash, filename, filesize)
}

// GetFileMeta : 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB : 从mysql获取文件元信息
func GetFileMetaDB(fileSha1, filename string) (*FileMeta, error) {
	tfile, err := models2.GetFileMeta(fileSha1, filename)
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error("Failed to get filemeta. err:", err)
		return nil, err
	}
	if tfile == nil {
		logging.Error("Failed to get filemeta. err:", err)
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileSHA1,
		FileName: tfile.FileName,
		FileSize: tfile.FileSize,
		Location: tfile.FileAddr,
	}
	return &fmeta, nil
}

// GetLastFileMetas : 获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// GetLastFileMetasDB : 批量从mysql获取文件元信息
func GetLastFileMetasDB(limit int) ([]FileMeta, error) {
	tfiles, err := models2.GetFileMetaList(limit)
	if err != nil {
		logging.Error("Failed to get filemetaList. err:", err)
		return make([]FileMeta, 0), err
	}

	tfilesm := make([]FileMeta, len(tfiles))
	for i := 0; i < len(tfilesm); i++ {
		tfilesm[i] = FileMeta{
			FileSha1: tfiles[i].FileSHA1,
			FileName: tfiles[i].FileName,
			FileSize: tfiles[i].FileSize,
			Location: tfiles[i].FileAddr,
		}
	}
	return tfilesm, nil
}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

// RemoveFileMeta : 删除元信息
func RemoveFileMetaDB(fileSha1 string) {
	// TODO 后续考虑是否真的删除
}
