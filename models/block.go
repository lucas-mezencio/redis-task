package models

import (
	"encoding/json"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"redis-task/database"
)

//type Tree struct {
//	Block    Block
//	Children []Block
//}

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
	pattern := "*:*"
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(result)

	return []Block{}
}

func getKeys(pattern string) ([]string, error) {
	db := database.ConnectWithDB()
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetBlockById(key string) Block {
	db := database.ConnectWithDB()
	result, err := db.Get(database.CTX, key).Result()
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
