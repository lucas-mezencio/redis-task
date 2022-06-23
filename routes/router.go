package routes

import (
	"github.com/gin-gonic/gin"
	"redis-task/handlers"
)

func HandleRoutes() {
	r := gin.Default()

	r.GET("/blocks", handlers.GetAllBlocksHandler)
	r.GET("/blocks/:id", handlers.GetBlockByIdHandler)
	r.POST("/blocks", handlers.CreateBlockHandler)
	r.PUT("/blocks:id", handlers.UpdateBlockByIdHandler)
	r.DELETE("/blocks/:id", handlers.DeleteBlockByIdHandler)

	r.GET("/tree/:id", handlers.GetTreeBellowId)

	r.NoRoute(handlers.NoRouteHandler)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
