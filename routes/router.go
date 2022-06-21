package routes

import (
	"github.com/gin-gonic/gin"
	"redis-task/handlers"
)

func HandleRoutes() {
	r := gin.Default()

	r.NoRoute(handlers.NoRouteHandler)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
