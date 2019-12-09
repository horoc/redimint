package routers

import (
	"github.com/chenzhou9513/DecentralizedRedis/routers/handlers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/login", handlers.Login)

	db := r.Group("/db")
	//db.Use(middleware.Authorization())
	{
		db.GET("/benchmark", handlers.BenchMarkTest)
		db.GET("/query", handlers.QueryCommand)
		db.GET("/query_private", handlers.QueryPrivateCommand)
		db.POST("/execute", handlers.ExecuteCommand)
		db.POST("/restore", handlers.RestoreLocalDatabase)
	}

	chain := r.Group("/chain")
	//db.Use(middleware.Authorization())
	{
		chain.GET("/transaction", handlers.GetTransactionByHash)
		chain.GET("/transaction_list", handlers.GetCommittedTxList)
		chain.GET("/block", handlers.GetBlockByHeight)
		chain.GET("/state", handlers.GetChainState)
		chain.GET("/info", handlers.GetChainInfo)
		chain.GET("/voting_validators", handlers.GetVotingValidator)
		chain.POST("/update_validators", handlers.UpdateValidators)

	}

	return r
}
