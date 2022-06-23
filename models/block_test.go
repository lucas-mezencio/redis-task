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
	db.FlushAll(database.CTX)
}

func TestGetAllBlocks(t *testing.T) {
	t.Run("get array of existing blocks", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		got := GetAllBlocks()
		assert.Equal(t, []Block{mockBlock}, got)
	})
}

func TestGetBlockById(t *testing.T) {
	t.Run("get existent block", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		got := GetBlockById("C3")
		assert.Equal(t, mockBlock, got)
	})
	t.Run("get inexistent block", func(t *testing.T) {
		got := GetBlockById("C3")
		assert.Equal(t, Block{}, got)
	})
}

func TestCreateBlock(t *testing.T) {
	t.Run("insert existent key", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()

		err := CreateBlock(mockBlock)
		assert.Error(t, err)
	})
	t.Run("insert new block", func(t *testing.T) {
		unmockBlock()
		err := CreateBlock(mockBlock)
		if err != nil {
			t.Error(err)
		}
		gotBlock := GetBlockById("C3")
		assert.Equal(t, mockBlock, gotBlock)
	})
}

func TestUpdateBlock(t *testing.T) {
	updatedMock := mockBlock
	t.Run("valid update block", func(t *testing.T) {
		mockBlockOnDB()
		err := UpdateBlock("C3", updatedMock)
		if err != nil {
			t.Error(err)
		}
		gotBlock := GetBlockById("C3")
		assert.Equal(t, updatedMock, gotBlock)
		unmockBlock()
	})
	t.Run("invalid key", func(t *testing.T) {
		unmockBlock()
		err := UpdateBlock("C3", updatedMock)
		assert.Error(t, err)
	})
	t.Run("indvalid parent id", func(t *testing.T) {
		updatedMock.ParentID = "asdf"
		mockBlockOnDB()
		err := UpdateBlock("C3", updatedMock)
		assert.Error(t, err)
		unmockBlock()
	})
}

func TestDeleteBlock(t *testing.T) {
	t.Run("existent block", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		err := DeleteBlockById("C3")
		if err != nil {
			t.Error(err)
		}
		gotBlock := GetBlockById("C3")
		assert.Equal(t, Block{}, gotBlock)
	})

	t.Run("nonexistent block", func(t *testing.T) {
		unmockBlock()
		err := DeleteBlockById("C3")
		assert.Error(t, err)
	})
}
