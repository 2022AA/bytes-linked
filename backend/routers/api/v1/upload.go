package handler

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/2022AA/bytes-linked/backend/models"
	"github.com/2022AA/bytes-linked/backend/pkg/util"
	"github.com/2022AA/bytes-linked/backend/pkg/util/e"
	"github.com/2022AA/bytes-linked/backend/setting"
	"github.com/2022AA/bytes-linked/backend/store/ipfs"
	"github.com/gin-gonic/gin"
)

const (
	uploadModelPathDir       = "backend/model3d" // -> access
	uploadModelPathDirOutput = "assets"
	uploadImagePathDir       = "backend/model3d/image"
	uploadImagePathDirOutput = "assets/image"

	bucketName      = "userfile"
	fileUploadedErr = "File has been uploaded before. "
)

func Setup() {
	// 目录已存在
	_, err := os.Stat(setting.AppSetting.TempLocalRootDir)
	if os.IsExist(err) {
		return
	}

	os.MkdirAll(uploadModelPathDir, 0744)
	os.MkdirAll(uploadImagePathDir, 0744)

	// 尝试创建目录
	err = os.MkdirAll(setting.AppSetting.TempLocalRootDir, 0744)
	if err != nil {
		log.Fatalf("无法创建临时存储目录，程序将退出")
	}
}

// File : 文件提交请求结构体
type CommitFileReq struct {
	UID      int    `json:"uid"`
	FileSHA1 string `json:"file_sha1"`
	FileName string `json:"file_name"`
	FileAddr string `json:"file_addr"`
	ImgAddr  string `json:"img_addr"`
	FiscoKey string `json:"fisco_key"`
}

// @Title CommitFileInfoHandler
// @Tags file
// @Description 提交上传模型文件
// @Accept json
// @Produce json
// @Param post body token query string true "token"
// @Param post body uid query int true "用户ID"
// @Param post body file_name query string true "文件名"
// @Param post body file_sha1 query string true "文件哈希"
// @Param post body file_addr query string true "文件下载访问地址"
// @Param post body img_addr query string true "缩略图下载访问地址"
// @Param post body fisco_key query string true "区块链密钥"
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/commit [Post]
func CommitFileInfoHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	f := CommitFileReq{}
	err := c.BindJSON(&f)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	file := models.File{
		FileSHA1:   f.FileSHA1,
		FileName:   f.FileName,
		FileAddr:   f.FileAddr,
		ImgAddr:    f.ImgAddr,
		OwnerUID:   f.UID,
		CreatorUID: f.UID,
		CreateAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt:   time.Now().Format("2006-01-02 15:04:05"),
		Status:     1,
	}

	id, err := models.CommitFileInfo(file)
	if err != nil || id == 0 {
		statusCode = e.ERROR_FILE_COMMIT_FAIL
		return
	}

	ok, err := models.FileTransaction(f.UID, f.UID, id, f.FiscoKey)
	if err != nil || !ok {
		statusCode = e.ERROR_FILE_COMMIT_FAIL
		return
	}

	responseData = id
}

type LikeFileReq struct {
	UID     int `json:"uid"`
	FileID  int `json:"file_id"`
	LikeNum int `json:"like_num"`
}

// @Title LikeFileHandler
// @Tags file
// @Description 模型文件点赞
// @Accept json
// @Produce json
// @Param post body token query string true "token"
// @Param post body uid query int true "用户ID"
// @Param post body file_id query int true "文件ID"
// @Param post body like_num query int true "点赞数"
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/like [Post]
func LikeFileHandler(c *gin.Context) {
	statusCode, statusMsg := 0, ""
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  statusMsg,
			Data: responseData,
		})
	}()

	req := LikeFileReq{}
	err := c.BindJSON(&req)
	if err != nil {
		statusCode, statusMsg = e.INVALID_PARAMS, "参数错误"
		return
	}

	if req.FileID == 0 || req.UID == 0 || req.LikeNum < 0 {
		statusCode, statusMsg, responseData = e.ERROR_FILE_COMMIT_FAIL, "参数错误", false
		return
	} else if req.LikeNum == 0 {
		req.LikeNum = 1
	}
	balance, statusMsg := models.LikeFile(req.UID, req.FileID, req.LikeNum)
	if balance < 0 {
		statusCode, responseData = e.ERROR_FILE_COMMIT_FAIL, 0
		return
	}

	responseData = balance
}

// @Title ListOwnFileInfoHandler
// @Tags file
// @Description 列出属于自己的模型文件
// @Accept json
// @Produce json
// @Param uid query int true "用户ID"
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/list_own [Get]
func ListOwnFileInfoHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	uid, _ := strconv.Atoi(c.Request.FormValue("uid"))

	files, err := models.ListOwnFileInfo(uid, 0, 0)
	if err != nil || uid == 0 {
		statusCode = e.ERROR_FILE_GET_FILEMETA_FAIL
		return
	}

	responseData = files
}

// @Title ListPublicFileInfoHandler
// @Tags file
// @Description 列出属于自己的模型文件
// @Accept json
// @Produce json
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/list_public [Get]
func ListPublicFileInfoHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	// uid, _ := strconv.Atoi(c.Request.FormValue("uid"))

	files, err := models.ListPublicFileInfo(0, 0)
	if err != nil {
		fmt.Println(err)
		statusCode = e.ERROR_FILE_GET_FILEMETA_FAIL
		return
	}

	responseData = files
}

// File : 文件提交请求结构体
type PublicFileReq struct {
	UID    int `json:"uid"`
	FileID int `json:"file_id"`
}

// @Title PublicFileHandler
// @Tags file
// @Description 分享模型文件(设成公开)
// @Accept json
// @Produce json
// @Param post body token query string true "token"
// @Param post body uid query int true "用户ID"
// @Param post body file_id query string true "文件ID"
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/public [Post]
func PublicFileHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	f := PublicFileReq{}
	err := c.BindJSON(&f)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	err = models.PublicFile(f.UID, f.FileID)
	if err != nil {
		statusCode = e.ERROR_FILE_SHARE_FAIL
		return
	}
}

type UploadFileRsp struct {
	FileSHA1 string `json:"file_sha1,omitempty"`
	FileAddr string `json:"file_addr"`
}

// @Title UploadFileHandler
// @Tags file
// @Description 上传文件接口(上传到obs/ipfs)
// @Accept json
// @Produce json
// @Param model post body model query string true "模型文件"
// @Param model post body image query string true "缩略图文件"
// @Success 200 {object} util.RespMsg
// @router /api/v1/file/upload [Post]
func UploadFileHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	var rsp UploadFileRsp

	fileHeader, err := c.FormFile("model")
	if err == nil {
		fileSHA1, fileAddr := uploadModelFile(c, fileHeader)
		rsp = UploadFileRsp{FileSHA1: fileSHA1, FileAddr: fileAddr}
	}

	fileHeader, err = c.FormFile("image")
	if err == nil {
		fileAddr := uploadImageFile(c, fileHeader)
		rsp = UploadFileRsp{FileAddr: fileAddr}
	}

	if fileHeader, err = c.FormFile("key"); err == nil {
		responseData = readKeyFile(c, fileHeader)
		return
	}

	if rsp.FileAddr == "" {
		statusCode = e.ERROR_FILE_UPLOAD_FAIL
	}
	responseData = rsp
}

func uploadModelFile(c *gin.Context, fileHeader *multipart.FileHeader) (string, string) {
	fileSHA1, err := uploadIPfs(fileHeader) // genSHA1(fileHeader)
	if err != nil {
		return "", ""
	}
	fileUnzipDir := path.Join(uploadModelPathDir, fileSHA1)
	filePath := fileUnzipDir + ".zip"
	os.Mkdir(fileUnzipDir, 0644)
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		return "", ""
	}
	unzipModel(filePath, fileUnzipDir)
	return fileSHA1, path.Join(uploadModelPathDirOutput, fileSHA1, "scene.gltf")
}

func uploadImageFile(c *gin.Context, fileHeader *multipart.FileHeader) string {
	fileSHA1 := genSHA1(fileHeader)
	fileSplits := strings.Split(fileHeader.Filename, ".")
	filePostfix := fileSplits[len(fileSplits)-1]
	if len(fileSplits) == 1 {
		filePostfix = ".png"
	}
	filePath := path.Join(uploadImagePathDir, fileSHA1+"."+filePostfix)
	err := c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		return ""
	}
	return path.Join(uploadImagePathDirOutput, fileSHA1+"."+filePostfix)
}

func readKeyFile(c *gin.Context, fileHeader *multipart.FileHeader) string {
	f, err := fileHeader.Open()
	if err != nil {
		return ""
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	return string(data)
}

func uploadIPfs(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return ipfs.Client().Add(f)

}

func genSHA1(fh *multipart.FileHeader) string {
	f, err := fh.Open()
	if err != nil {
		return util.Sha1([]byte(strconv.Itoa(rand.Int())))
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return util.Sha1([]byte(strconv.Itoa(rand.Int())))
	}

	return util.Sha1(data)
}

func unzipModel(src, dst string) {
	archive, err := zip.OpenReader(src)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		// fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			// fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
