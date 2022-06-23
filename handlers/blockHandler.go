// Package handlers provides handlers for the application
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redis-task/models"
	"reflect"
)

// GetAllBlocksHandler handles the route to get all blocks on db
//			GET /blocks
// Returns a list of blocks as JSON or an empty list if the db is empty and OK Status
func GetAllBlocksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllBlocks())
}

// GetBlockByIdHandler handles the route to get a single block by a given id
//			GET /blocks/id
// If the block is not found return null body and NotFound Status
// Returns the block data as JSON and OK Status
func GetBlockByIdHandler(c *gin.Context) {
	blockId := c.Param("id")
	block := models.GetBlockById(blockId)

	if reflect.DeepEqual(block, models.Block{}) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, block)
}

// CreateBlockHandler handles the route to create a block
//			POST /blocks
// If the application can't parse the body of request it returns a BadRequest Status
// If the block already exists returns a JSON with the given data and the error - BadRequest Status
// If the block have an invalid parent id returns a JSON with given data and the error - BadRequest Status
// Returns the given data and Created Status
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

// UpdateBlockByIdHandler handles the update route to update a block
//			PUT /blocks/id
// If the application can't parse the body of request it returns a BadRequest Status
// If the block don't exist returns a JSON with given data and the error - NotFound Status
// If the block have an invalid parent id returns a JSON with given data and the error - BadRequest Status
// Returns the updated block data as JSON and Ok Status
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
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, newBlock)
}

// DeleteBlockByIdHandler handles the route to delete a block
//			DELETE /blocks/id
// If the block don't exist returns a NotFound Status
// If the block hava an invalid parent id returns a BadRequestStatus
// Returns an Ok Status
func DeleteBlockByIdHandler(c *gin.Context) {
	id := c.Param("id")
	err := models.DeleteBlockById(id)
	if err == models.ErrBlockNotExists {
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err == models.ErrInvalidParentId {
		c.JSON(http.StatusBadRequest, err)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}
