package handler

import (
	"strconv"

	"github.com/2022AA/bytes-linked/backend/models"
	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
	"github.com/2022AA/bytes-linked/backend/pkg/util"
	"github.com/2022AA/bytes-linked/backend/pkg/util/e"
	"github.com/gin-gonic/gin"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt    = "*#890"
	adminToken = "admin@1234"
)

// @Title DoSignupHandler
// @Tags user
// @Description 处理用户注册请求
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Param inviteCode query string true "邀请码"
// @Param phone query string true "手机号"
// @Param avatarUrl query string true "头像地址"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/signup [Post]
func DoSignupHandler(c *gin.Context) {
	statusCode := 0
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
		})
	}()
	params := struct {
		UserName   string `json:"username,omitempty"`
		Passwd     string `json:"password,omitempty"`
		InviteCode string `json:"inviteCode,omitempty"`
		Phone      string `json:"phone,omitempty"`
		AvatarUrl  string `json:"avatarUrl,omitempty"`
	}{}
	err := c.Bind(&params)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		v2.Warnln(err)
		return
	}
	username := params.UserName
	passwd := params.Passwd
	inviteCode := params.InviteCode

	_, err = models.GetUserInfo(username)
	if err == nil {
		statusCode = e.ERROR_USER_ALREADY_EXIST
		return
	}

	ok, err := models.IsValid(inviteCode)
	if err != nil {
		statusCode = e.ERROR_USER_INVITE_CODE
		return
	}
	if !ok {
		statusCode = e.ERROR_USER_INVITE_CODE
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	// 将用户信息注册到用户表中
	suc := models.UserSignup(username, encPasswd, inviteCode, params.Phone, params.AvatarUrl)
	if !suc {
		statusCode = e.ERROR_USER_SIGNUP_FAIL
		return
	}
}

// @Title DoSignInHandler
// @Tags user
// @Description 登录接口
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {object} util.RespMsg
//@router /api/v1/user/signin [Post]
func DoSignInHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()
	params := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{}
	err := c.Bind(&params)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		v2.Warnln(err)
		return
	}

	username := params.Username
	password := params.Password

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	// 1. 校验用户名及密码
	pwdChecked := models.UserSignin(username, encPasswd)
	if !pwdChecked {
		// 密码错误
		statusCode = e.ERROR_USER_PASSWORD
		return
	}

	// 2. 生成访问凭证(token)
	token, err := util.GenerateToken(username, password)
	if err != nil {
		logging.Error("[DoSignInHandler] GenerateToken. err:", err.Error())
		statusCode = e.ERROR_USER_TOKEN
		return
	}

	// 3. 判断是否管理员
	//user, err := models.GetUserInfo(username)
	//if err != nil {
	//	logging.Error("[DoSignInHandler] GetUserInfo. err:", err.Error())
	//	statusCode = e.ERROR_USER_GET_INFO_FAIL
	//	return
	//}

	// 登录成功，返回用户信息
	responseData = struct {
		Username string `json:"Username"`
		Token    string `json:"Token"`
	}{
		Username: username,
		Token:    token,
	}
}

// @Title UserExistsHandler
// @Tags user
// @Description 下载公私钥
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Param  secrete_type query int true "类型 1: 公钥, 2: 私钥"
// @Success 200 {object} util.RespMsg
//@router /api/v1/user/secret_download [Get]
func SecretFileDownloadHandler(c *gin.Context) {
	statusCode := 0
	defer func() {
		if statusCode != 0 {
			c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
				Code: statusCode,
				Msg:  e.GetMsg(statusCode),
			})
		}
	}()

	username := c.Request.FormValue("username")
	secretTypeStr := c.Request.FormValue("secrete_type")

	secretType, err := strconv.Atoi(secretTypeStr)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}
	if secretType != models.SecretType_Private && secretType != models.SecretType_Public {
		statusCode = e.INVALID_PARAMS
		return
	}

	// 2. 查询用户信息
	ok, err := models.UserExist(username)
	if err != nil {
		// 获取用户信息失败
		logging.Error("[SecretFileDownloadHandler] GetUserInfo. err:", err.Error())
		statusCode = e.ERROR_USER_GET_INFO_FAIL
		return
	}
	if !ok {
		statusCode = e.ERROR_USER_CHECK_EXIST_FAIL
		return
	}

	// 3. 查询密钥信息
	secretFile, err := models.QuerySecretFile(username, models.SecretType(secretType))
	if err != nil {
		logging.Error("[SecretFileDownloadHandler] QuerySecretFile. err:", err.Error())
		statusCode = e.ERROR_USER_SECRET_FILE_QUERY_FAILED
		return
	}
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+secretFile.FileName)
	cnt, err := c.Writer.Write([]byte(secretFile.Data))
	v2.Debugf("SecretFileDownloadHandler | username: %s, write_count: %d", username, cnt)
	return

}

// @Title UserExistsHandler
// @Tags user
// @Description 查询用户是否存在
// @Param username query string true "用户名"
// @Success 200 {object} util.RespMsg
//@router /api/v1/user/exists [Get]
func UserExistsHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	// 1. 解析请求参数
	username := c.Request.FormValue("username")

	// 2. 查询用户是否存在
	exists, err := models.UserExist(username)
	if err != nil {
		logging.Error("[UserExistsHandler] UserExist. err:", err.Error())
		statusCode = e.ERROR_USER_CHECK_EXIST_FAIL
		return
	}

	responseData = map[string]interface{}{
		"exists": exists,
	}
}

// @Title UserInfoHandler
// @Tags user
// @Description 查询用户信息
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/info [Get]
func UserInfoHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	// 1. 解析请求参数
	username := c.Request.FormValue("username")

	// 2. 查询用户信息
	user, err := models.GetUserInfo(username)
	if err != nil {
		// 获取用户信息失败
		logging.Error("[UserInfoHandler] GetUserInfo. err:", err.Error())
		statusCode = e.ERROR_USER_GET_INFO_FAIL
		return
	}

	// 3. 组装并且响应用户数据
	responseData = user
}

// @Title UserQueryHandler
// @Tags user
// @Description 批量查询的用户信息（需要管理员账号）
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Param offset query string true "偏移量"
// @Param limit query string true "数量"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/query [Get]
func UserQueryHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	//username := c.Request.FormValue("username")
	////user, err := models.GetUserInfo(username)
	////if err != nil {
	////	logging.Error(fmt.Sprintf("[UserQueryHandler] GetUserInfo. user:%s, err:%s", username, err.Error()))
	////	statusCode = e.ERROR_USER_GET_INFO_FAIL
	////	return
	//}
	//if !user.IsAdmin {
	//	statusCode = e.ERROR_USER_IS_NOT_ADMIN
	//	return
	//}

	offset, _ := strconv.Atoi(c.Request.FormValue("offset"))
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	users, err := models.QueryUserInfoList(offset, limit)
	if err != nil {
		logging.Error("[UserQueryHandler] QueryUserInfoList. err:", err.Error())
		statusCode = e.ERROR_FILE_QUERY_USERFILEMETA_FAIL
		return
	}

	responseData = users
}

// @Title UserUpdateHandler
// @Tags user
// @Description 修改用户密码
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Param password query string true "初始密码"
// @Param newPassword query string true "新密码"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/update [Post]
func UserUpdateHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	newPassword := c.Request.FormValue("newPassword")

	encPasswd := util.Sha1([]byte(password + pwdSalt))
	pwdChecked := models.UserSignin(username, encPasswd)
	if !pwdChecked {
		// 密码错误
		statusCode = e.ERROR_USER_PASSWORD
		return
	}

	if len(newPassword) < 5 {
		statusCode = e.INVALID_PARAMS
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd = util.Sha1([]byte(newPassword + pwdSalt))
	// 更新到用户表中
	suc, err := models.UpdateUserPassword(username, encPasswd)
	if !suc || err != nil {
		statusCode = e.ERROR_USER_UPDATE_PASSWORD_FAIL
		return
	}
	responseData = suc
}

// @Title UserUpdateHandler
// @Tags user
// @Description 删除用户（需要管理员权限）
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Param deleteUsername query string true "需要删除的用户"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/delete [Post]
func UserDeleteHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	//username := c.Request.FormValue("username")
	deleteUsername := c.Request.FormValue("deleteUsername")
	//user, err := models.GetUserInfo(username)
	//if err != nil {
	//	logging.Error(fmt.Sprintf("[UserQueryHandler] GetUserInfo. user:%s, err:%s", username, err.Error()))
	//	statusCode = e.ERROR_USER_GET_INFO_FAIL
	//	return
	//}
	//if !user.IsAdmin {
	//	statusCode = e.ERROR_USER_IS_NOT_ADMIN
	//	return
	//}

	err := models.DeleteUserFileByUsername(deleteUsername)
	if err != nil {
		statusCode = e.ERROR_USER_DELETE_FAIL
		return
	}

	err = models.DeleteUser(deleteUsername)
	if err != nil {
		statusCode = e.ERROR_USER_DELETE_FAIL
		return
	}
	responseData = "OK"
}

// @Title DoFileTransaction
// @Tags user
// @Description 转移接口
// @Param token query string true "token"
// @Param uid query int true "用户名"
// @Param toUsername query string true "转移的用户名"
// @Param transferFileId query int true "转移的file id"
// @Param privateKey query string true "私钥"
// @Success 200 {object} util.RespMsg
//@router /api/v1/user/transaction [Post]
func DoFileTransaction(c *gin.Context) {

	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()

	param := struct {
		Uid            int    `json:"uid,omitempty"`
		ToUsername     string `json:"toUsername,omitempty"`
		TransferFileId int    `json:"transferFileId,omitempty"`
		PrivateKey     string `json:"privateKey,omitempty"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	toUser, err := models.GetUserInfo(param.ToUsername)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	ok, err := models.FileTransaction(
		param.Uid, toUser.ID, param.TransferFileId, param.PrivateKey,
	)
	if err != nil {
		statusCode = e.ERROR_USER_TRANSACTION_ERROR
		return
	}
	if !ok {
		statusCode = e.ERROR_USER_INVALID_ITEM
		return
	}

	return

}

// @Title UserAvatarUpdateHandler
// @Tags user
// @Description 修改用户头像
// @Param token query string true "token"
// @Param username query string true "用户名"
// @Param avatarUrl query string true "头像地址"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/avatar [POST]
func UserAvatarUpdateHandler(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()
	params := struct {
		Username  string `json:"username,omitempty"`
		AvatarUrl string `json:"avatarUrl,omitempty"`
	}{}
	err := c.Bind(&params)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	// 更新到用户表中
	user, err := models.UpdateUserAvatar(params.Username, params.AvatarUrl)
	if err != nil {
		statusCode = e.ERROR_USER_UPDATE_AVATAR_FAIL
		return
	}
	responseData = user
}

// @Title UserFile
// @Tags user
// @Description 获取用户藏品
// @Param token query string true "token"
// @Param uid query int true "uid"
// @Param fileId query int true "文件id"
// @Success 200 {object} util.RespMsg
// @router /api/v1/user/file [GET]
func UserGetFile(c *gin.Context) {
	statusCode := 0
	var responseData interface{}
	defer func() {
		c.JSON(e.GetHttpCode(statusCode), util.RespMsg{
			Code: statusCode,
			Msg:  e.GetMsg(statusCode),
			Data: responseData,
		})
	}()
	uidStr := c.Request.FormValue("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}
	fileIdStr := c.Request.FormValue("fileId")
	fileId, err := strconv.Atoi(fileIdStr)
	if err != nil {
		statusCode = e.INVALID_PARAMS
		return
	}

	file, err := models.GetFile(uid, fileId)
	if err != nil {
		statusCode = e.ERROR_USER_FILE_NOT_FOUND
		return
	}
	responseData = file
}
