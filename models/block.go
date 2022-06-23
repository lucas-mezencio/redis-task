package models

import (
	"encoding/json"
	"errors"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"redis-task/database"
	"reflect"
	"strings"
)

const defaultPattern string = "*:*"

var (
	ErrBlockAlreadyExists = errors.New("this key already exists on database")
	ErrInvalidParentId    = errors.New("invalid parent id or parent doesn't exists")
	ErrBlockNotExists     = errors.New("key not exists on database")
	//ErrInvalidKey         = errors.New("this key is invalid - follow the pattern it:father")
)

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

func CreateBlock(block Block) error {
	if existentKeys := getKeys(block.ID); len(existentKeys) != 0 {
		return ErrBlockAlreadyExists
	}
	return setBlock(block)
}

func UpdateBlock(key string, block Block) error {
	if checkBlockKey := getKeys(key + ":*"); len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}
	return setBlock(block)
}

func DeleteBlockById(key string) error {
	db := database.ConnectWithDB()
	checkBlockKey := getKeys(key + ":*")
	if len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}
	blockKey := checkBlockKey[0]
	childrenKeys := getKeys("*:" + getIndividualBlockId(blockKey))
	if len(childrenKeys) == 0 {
		err := db.Del(database.CTX, blockKey).Err()
		return err
	}

	childrenBlocks, err := getChildrenById(childrenKeys)
	if err != nil {
		return err
	}

	block := GetBlockById(getIndividualBlockId(blockKey))
	for _, childBlock := range childrenBlocks {
		childBlock.ID = updatedBlockId(childBlock.ID, block.ParentID)
		childBlock.ParentID = block.ParentID
		err := CreateBlock(childBlock)
		if err != nil {
			return err
		}
	}
	keysToDelete := append(childrenKeys, block.ID)
	err = db.Del(database.CTX, keysToDelete...).Err()

	return err
}

func setBlock(block Block) error {
	if block.ParentID != "0" {
		parentBlock := GetBlockById(block.ParentID)
		if reflect.DeepEqual(parentBlock, Block{}) {
			return ErrInvalidParentId
		}
	}

	db := database.ConnectWithDB()
	err := db.Set(database.CTX, block.ID, block, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func getChildrenById(childrenKeys []string) ([]Block, error) {
	db := database.ConnectWithDB()
	result, err := db.MGet(database.CTX, childrenKeys...).Result()
	if err != nil {
		return nil, err
	}

	var childrenBlocks []Block
	for _, item := range result {
		var childBlock Block
		err := childBlock.UnmarshalBinary([]byte(fmt.Sprint(item)))
		if err != nil {
			return nil, err
		}
		childrenBlocks = append(childrenBlocks, childBlock)
	}
	return childrenBlocks, nil
}

func getKeys(pattern string) []string {
	db := database.ConnectWithDB()
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		return nil
	}
	return result
}

func getIndividualBlockId(compositeKey string) string {
	return strings.Split(compositeKey, ":")[0]
}

func updatedBlockId(key, parentKey string) string {
	blockKey := getIndividualBlockId(key)
	return blockKey + ":" + parentKey
}
