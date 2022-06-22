package models

import (
	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
	"redis-task/database"
	"testing"
)

var (
	c0 = Block{
		ID:       "C0:0",
		Name:     "Cliente A",
		ParentID: "0",
		Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
		Value:    10000,
	}
	f1 = Block{
		ID:       "F1:C0",
		Name:     "FAZENDA 1",
		ParentID: "C0",
		Centroid: *geojson.NewPointGeometry([]float64{-52.9046630859375, -18.132801356084773}),
		Value:    1000,
	}
	//f2 = Block{
	//	ID:       "F2:C0",
	//	Name:     "FAZENDA 2",
	//	ParentID: "C0",
	//	Centroid: *geojson.NewPointGeometry([]float64{54.60205078125, -25.52509317964987}),
	//	Value:    2000,
	//}
)

var treeMock = Tree{
	Block: c0,
	Children: []Tree{
		{
			Block:    f1,
			Children: nil,
		},
		//{
		//	Block:    f2,
		//	Children: nil,
		//},
	},
}

func MockTree(t *testing.T) {
	UnmockTree(t)
	db := database.ConnectWithDB()
	blocks := []Block{c0, f1}
	for _, block := range blocks {
		err := db.Set(database.CTX, block.ID, block, 0).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func UnmockTree(t *testing.T) {
	db := database.ConnectWithDB()
	err := db.FlushAll(database.CTX).Err()
	if err != nil {
		t.Error(err)
	}
}

func TestGetTreeById(t *testing.T) {
	t.Run("mocked tree", func(t *testing.T) {
		MockTree(t)
		defer UnmockTree(t)

		got := GetTreeById("C0")

		assert.Equal(t, treeMock, got)
	})
	t.Run("nonexistent tree", func(t *testing.T) {
		got := GetTreeById("C0")
		assert.Equal(t, Tree{}, got)
		assert.NotEqual(t, treeMock, got)
	})
}
