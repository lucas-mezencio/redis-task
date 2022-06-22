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

func CreateBlockHandler(c *gin.Context) {
	var block models.Block
	err := c.ShouldBindJSON(&block)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = models.CreateBlock(block)
	if err == models.ErrBlockAlreadyExists || err == models.ErrInvalidParentId {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":  block,
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, block)
		return
	}

	c.JSON(http.StatusCreated, block)
}
