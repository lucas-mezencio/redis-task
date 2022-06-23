package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NoRouteHandler handlers the default and the root (/) route
//
// Returns a NotFound Status
func NoRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, nil)
}
