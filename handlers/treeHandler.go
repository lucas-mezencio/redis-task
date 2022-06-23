package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/models"
	"reflect"
)

// GetTreeBellowId handles the route to get a tree based on a given id
//			GET /tree/id
// If the root block doesn't exist return a NotFound Status
//
// Returns the tree with the given id block as root and OK Status
func GetTreeBellowId(c *gin.Context) {
	id := c.Param("id")
	tree := models.GetTreeById(id)
	if reflect.DeepEqual(tree, models.Tree{}) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, tree)
}
