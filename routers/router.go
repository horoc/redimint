package routers

import (
	"github.com/chenzhou9513/DecentralizedRedis/routers/handlers"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {

	r := gin.Default()

	//r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	//r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	//
	//r.GET("/auth", api.GetAuth)
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.POST("/upload", api.UploadImage)

	db := r.Group("/db")
	{
		//获取标签列表
		db.GET("/query", handlers.QueryCommand)
		db.POST("/execute", handlers.ExecuteCommand)
		db.POST("/execute_async", handlers.ExecuteCommandAsync)

		////获取指定文章
		//apiv1.GET("/articles/:id", v1.GetArticle)
		////新建文章
		//apiv1.POST("/articles", v1.AddArticle)
		////更新指定文章
		//apiv1.PUT("/articles/:id", v1.EditArticle)
		////删除指定文章
		//apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		////生成文章海报
		//apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	chain := r.Group("/chain")
	{
		chain.GET("/transaction", handlers.GetTransactionByHash)
		chain.GET("/block", handlers.GetBlockByHeight)
		chain.GET("/state", handlers.GetChainState)
	}

	return r
}
