package e

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	INVALID_PARAMS:                 "请求参数错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",

	ERROR_USER_PASSWORD:                 "用户密码错误",
	ERROR_USER_TOKEN:                    "用户Token错误",
	ERROR_USER_SIGNUP_FAIL:              "用户注册失败",
	ERROR_USER_GET_INFO_FAIL:            "用户信息获取失败",
	ERROR_USER_CHECK_EXIST_FAIL:         "用户是否存在查询失败",
	ERROR_USER_ALREADY_EXIST:            "用户已经存在",
	ERROR_USER_IS_NOT_ADMIN:             "用户需要管理员权限",
	ERROR_USER_QUERY_FAIL:               "用户队列查找失败",
	ERROR_USER_UPDATE_PASSWORD_FAIL:     "用户密码更新失败",
	ERROR_USER_DELETE_FAIL:              "用户删除失败",
	ERROR_USER_ADMIN_TOKEN:              "管理员token错误",
	ERROR_USER_INVITE_CODE:              "无效邀请码",
	ERROR_USER_SECRET_FILE_QUERY_FAILED: "密钥不存在",
	ERROR_USER_PARAMS_ERROR:             "用户名或密码长度不合法",
	ERROR_USER_PHONE_INVALID:            "手机号长度不合法",
	ERROR_USER_TRANSACTION_ERROR:        "转移发生错误，请重试",
	ERROR_USER_INVALID_ITEM:             "转移物品已无效",
	ERROR_USER_UPDATE_AVATAR_FAIL:       "更新头像失败，用户处于无效状态",

	ERROR_FILE_GET_FORM_FAIL:                    "文件参数获取失败",
	ERROR_FILE_GET_FILEMETA_NIL:                 "文件元信息获取为空",
	ERROR_FILE_UPLOAD_TO_CEPH_FAIL:              "文件上传到Ceph失败",
	ERROR_FILE_UPLOAD_LOCAL_FAIL:                "文件上传到本地失败",
	ERROR_FILE_CHANGE_TO_BYTE_FAIL:              "文件转化成字节失败",
	ERROR_FILE_SAVE_FILEMETA_FAIL:               "文件元信息保存失败",
	ERROR_FILE_GET_FILEMETA_FAIL:                "文件元信息获取失败",
	ERROR_FILE_QUERY_USERFILEMETA_FAIL:          "文件用户元信息查询失败",
	ERROR_FILE_DOWNLOAD_ROOTDIR:                 "文件下载类型异常",
	ERROR_FILE_DELETE_FAIL:                      "文件删除失败",
	ERROR_FILE_DELETE_USERFILE_FAIL:             "文件用户元信息删除失败",
	ERROR_FILE_DELETE_USERFILE_FAIL_FOR_NOOWNER: "文件用户元信息删除非法，该用户不是拥有者",
	ERROR_FILE_SAVE_USERFILEMETA_FAIL:           "文件用户元信息保存失败",
	ERROR_FILE_RENAME_FAIL:                      "文件名修改失败",
	ERROR_FILE_SHARE_FAIL:                       "共享文件申请失败",
	ERROR_FILE_CANCEL_SHARE_FAIL:                "共享文件取消申请失败",
	ERROR_FILE_AUDIT_SHARE_FAIL:                 "共享文件审核失败",
	ERROR_FILE_CANCEL_AUDIT_SHARE_FAIL:          "共享文件取消审核失败",
	ERROR_FILE_UNSHARE_DOWNLOAD_FAIL:            "文件未共享下载失败",
	ERROR_FILE_UNPERMISSION_DELETE_FILE:         "文件删除非法，权限不够",
	ERROR_FILE_DELETE_FILE_FAIL:                 "文件删除失败",
	ERROR_FILE_UPLOADED_SAVE_FILEMETA_FAIL:      "文件已上传，请勿重复上传",

	ERROR_FILE_UPLOAD_MPUPLOAD_FAIL:     "文件分块上传失败",
	ERROR_FILE_INVALID_MPUPLOAD_REQUEST: "文件分块上传非法请求",
	ERROR_FILE_CANCEL_MPUPLOAD_FAIL:     "文件分块上传取消失败",
	ERROR_FILE_QUERY_MPUPLOAD_FAIL:      "文件分块上传查询失败",
	ERROR_FILE_COMPLETE_MPUPLOAD_FAIL:   "文件分块上传失败",

	ERROR_FILE_NORECORD_FASTUPLOAD_FAIL: "文件秒传，查不到记录，请使用普通上传",
	ERROR_FILE_UPLOAD_FASTUPLOAD_FAIL:   "文件秒传失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	if code == 0 {
		return MsgFlags[SUCCESS]
	}

	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
