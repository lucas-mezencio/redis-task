package models

import (
	"encoding/json"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"redis-task/database"
)

const defaultPattern string = "*:*"

type Block struct {
	ID       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	ParentID string           `json:"parentID,omitempty"`
	Centroid geojson.Geometry `json:"centroid,omitempty"`
	Value    float64          `json:"value,omitempty"`
}

func (b Block) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Block) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, b)
}

func GetAllBlocks() []Block {
	db := database.ConnectWithDB()
	keys := getKeys(defaultPattern)
	if keys == nil {
		return []Block{}
	}
	result, err := db.MGet(database.CTX, keys...).Result()
	if err != nil {
		return []Block{}
	}

	var blocks []Block
	for _, item := range result {
		var block Block
		err := block.UnmarshalBinary([]byte(fmt.Sprint(item)))
		if err != nil {
			return []Block{}
		}
		blocks = append(blocks, block)
	}

	return blocks
}

func getKeys(pattern string) []string {
	db := database.ConnectWithDB()
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		return nil
	}
	return result
}

func GetBlockById(key string) Block {
	db := database.ConnectWithDB()

	blockKey := getKeys(key + ":*")
	if len(blockKey) != 1 {
		return Block{}
	}
	result, err := db.Get(database.CTX, blockKey[0]).Result()
	if err != nil {
		return Block{}
	}
	var block Block
	if err := block.UnmarshalBinary([]byte(result)); err != nil {
		fmt.Println(err.Error(), err)
		return Block{}
	}
	return block
}
