package router

import (
	// _ "ccian.cc/really/arweave-api/docs"
	mlogger "ccian.cc/really/arweave-api/middleware/logger"
	"ccian.cc/really/arweave-api/pkg/setting"
	"ccian.cc/really/arweave-api/router/v1/images"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRouter(config *setting.ServerConfig, logger *logrus.Entry) *gin.Engine {
	r := gin.New()

	r.Use(mlogger.Logger(), mlogger.Time(), gin.Recovery())

	gin.SetMode(config.RunMode)

	// // docs
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Group := r.Group("/api/v1")
	addImageRoutes(v1Group)

	return r
}

func addImageRoutes(parent *gin.RouterGroup) {
	// images
	group := parent.Group("/images")
	group.POST("", images.UploadImage)
	group.POST("/uri", images.SaveImageUri)
}
