package routers

import (
	_ "github.com/2022AA/bytes-linked/backend/docs"
	"github.com/2022AA/bytes-linked/backend/pkg/util"
	handler "github.com/2022AA/bytes-linked/backend/routers/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiv1 := r.Group("/api/v1")
	userApi := apiv1.Group("/user")
	fileApi := apiv1.Group("/file")

	// v1版本API
	userApi.POST("/signup", handler.DoSignupHandler)
	userApi.POST("/signin", handler.DoSignInHandler)
	userApi.GET("/exists", handler.UserExistsHandler)
	// todo test
	userApi.POST("/transaction", handler.DoFileTransaction)
	userApi.POST("/avatar", handler.UserAvatarUpdateHandler)
	userApi.GET("/file", handler.UserGetFile)
	userApi.Use(util.JWT())
	{
		// 用户相关
		userApi.GET("/info", handler.UserInfoHandler)
		userApi.GET("/query", handler.UserQueryHandler)
		userApi.POST("/update", handler.UserUpdateHandler)
		userApi.POST("/delete", handler.UserDeleteHandler)
		userApi.GET("/secret-download", handler.SecretFileDownloadHandler)
	}

	// 文件操作
	fileApi.Use( /*util.JWT()*/ )
	{
		// 默认的文件操作相关
		fileApi.POST("/commit", handler.CommitFileInfoHandler)
		fileApi.POST("/public", handler.PublicFileHandler)
		fileApi.POST("/upload", handler.UploadFileHandler)
		fileApi.POST("/like", handler.LikeFileHandler)
		fileApi.GET("/list_own", handler.ListOwnFileInfoHandler)
		fileApi.GET("/list_public", handler.ListPublicFileInfoHandler)
	}
	return r
}
