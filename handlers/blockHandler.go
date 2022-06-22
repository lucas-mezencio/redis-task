package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/models"
	"reflect"
)

func GetAllBlocksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllBlocks())
}

func GetBlockByIdHandler(c *gin.Context) {
	blockId := c.Param("id")
	block := models.GetBlockById(blockId)

	if reflect.DeepEqual(block, models.Block{}) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, block)
}
