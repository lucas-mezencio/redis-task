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
		c.JSON(http.StatusBadRequest, err)
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
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, block)
}

func UpdateBlockByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var newBlock models.Block

	err := c.ShouldBindJSON(&newBlock)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = models.UpdateBlock(id, newBlock)
	if err == models.ErrBlockNotExists {
		c.JSON(http.StatusNotFound, gin.H{
			"data":  newBlock,
			"error": err.Error(),
		})
		return
	} else if err == models.ErrInvalidParentId {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":  newBlock,
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, newBlock)
}

func DeleteBlockByIdHandler(c *gin.Context) {
	id := c.Param("id")
	err := models.DeleteBlockById(id)
	if err == models.ErrBlockNotExists {
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err == models.ErrInvalidParentId || err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
