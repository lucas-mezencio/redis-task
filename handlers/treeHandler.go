package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/models"
	"reflect"
)

func GetTreeBellowId(c *gin.Context) {
	id := c.Param("id")
	tree := models.GetTreeById(id)
	if reflect.DeepEqual(tree, models.Tree{}) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, tree)
}
