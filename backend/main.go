package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/2022AA/bytes-linked/backend/models"
	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	"github.com/2022AA/bytes-linked/backend/pkg/util"
	"github.com/2022AA/bytes-linked/backend/routers"
	"github.com/2022AA/bytes-linked/backend/setting"
	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化配置文件设置
	setting.Setup()
	// 初始化临时存储文件夹
	//handler.Setup()
	// 初始化数据库连接设置
	models.Setup()
	// 初始化日志设置
	logging.Setup()
	// 初始化jwt验证设置
	util.Setup()
}

// @title 文件存储系统API
// @version 1.2
// @description 功能：用户管理 + 文件上传下载 + 文件共享
func main() {
	rand.Seed(time.Now().UnixNano())
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("[error] ListenAndServe. Err: %s", err)
	}
}
