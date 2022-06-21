package models

import (
	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
	"redis-task/database"
	"testing"
)

var mockBlock = Block{
	ID:       "C3:0",
	Name:     "Bloco teste",
	ParentID: "0",
	Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
	Value:    50000000,
}

func mockBlockOnDB() {
	db := database.ConnectWithDB()
	db.Set(database.CTX, mockBlock.ID, mockBlock, 0)
}

func unmockBlock() {
	db := database.ConnectWithDB()
	db.Del(database.CTX, mockBlock.ID)
}

func TestGetAllBlocks(t *testing.T) {
	GetAllBlocks()
}

func TestGetBlockById(t *testing.T) {
	t.Run("get existent block", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		got := GetBlockById(mockBlock.ID)
		assert.Equal(t, mockBlock, got)
	})
	t.Run("get inexistent block", func(t *testing.T) {
		got := GetBlockById(mockBlock.ID)
		assert.Equal(t, Block{}, got)
	})
}