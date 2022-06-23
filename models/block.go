// Package models holds the models of the applications and their database operations
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"redis-task/database"
	"redis-task/utils"
	"reflect"
)

// defaultPattern the default pattern to search all ids of blocks
const defaultPattern string = "*:*"

var (
	// ErrBlockAlreadyExists Error thrown when a given block already exists on db
	ErrBlockAlreadyExists = errors.New("this key already exists on database")

	// ErrInvalidParentId  Error thrown when a given block have an invalid parent id
	ErrInvalidParentId = errors.New("invalid parent id or parent doesn't exists")

	// ErrBlockNotExists  Error thrown when a given block don't exist on db
	ErrBlockNotExists = errors.New("key not exists on database")

	//ErrInvalidKey         = errors.New("this key is invalid - follow the pattern it:father")
)

// Block Definition of the base structure of the application
type Block struct {
	ID       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	ParentID string           `json:"parentID,omitempty"`
	Centroid geojson.Geometry `json:"centroid,omitempty"`
	Value    float64          `json:"value,omitempty"`
}

// MarshalBinary implementation of a needed interface for the go-redis lib
func (b Block) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// UnmarshalBinary implementation of a needed interface for the go-redis lib
func (b *Block) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, b)
}

// GetAllBlocks consults the database to search for all blocks
//
// in case of error or empty database returns an empty Block slice
//
// Returns a BlockSlice
func GetAllBlocks() []Block {
	db := database.ConnectWithDB()
	keys := utils.GetKeys(defaultPattern)
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

// GetBlockById consults the database to search for a specific block given an id (key)
//
// in case of error or Block not found on database returns an empty Block
//
// Returns the Block that matches the key parameter
func GetBlockById(key string) Block {
	db := database.ConnectWithDB()

	blockKey := utils.GetKeys(key + ":*")
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

// CreateBlock create/insert a block on database
//
// if the key of the block already exists return a ErrBlockAlreadyExists
//
// if the parent id of block is invalid returns a ErrInvalidParentId
//
// Returns a nil value as error if the block is created
func CreateBlock(block Block) error {
	if existentKeys := utils.GetKeys(block.ID); len(existentKeys) != 0 {
		return ErrBlockAlreadyExists
	}
	return setBlock(block)
}

// UpdateBlock updates a block on database
//
// if the block doesn't exist on database returns a ErrBlockNotExists
//
// if the parent id of block is invalid returns a ErrInvalidParentId
//
// Returns a nil value as error if the block gets updated
func UpdateBlock(key string, block Block) error {
	if checkBlockKey := utils.GetKeys(key + ":*"); len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}
	return setBlock(block)
}

// DeleteBlockById deletes a block on database and transfer its parent's to its children as parent
//
// if the block doesn't exist on database returns a ErrBlockNotExists
//
// if the parent id of block is invalid returns a ErrInvalidParentId
//
// Returns a nil value as error if the block gets deleted
func DeleteBlockById(key string) error {
	db := database.ConnectWithDB()
	checkBlockKey := utils.GetKeys(key + ":*")
	if len(checkBlockKey) != 1 {
		return ErrBlockNotExists
	}

	blockKey := checkBlockKey[0]
	childrenKeys := utils.GetKeys("*:" + utils.GetIndividualBlockId(blockKey))
	if len(childrenKeys) == 0 {
		err := db.Del(database.CTX, blockKey).Err()
		return err
	}

	childrenBlocks, err := getChildren(childrenKeys)
	if err != nil {
		return err
	}

	block := GetBlockById(utils.GetIndividualBlockId(blockKey))

	for _, childBlock := range childrenBlocks {
		childBlock.ID = utils.UpdatedBlockId(childBlock.ID, block.ParentID)
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

// setBlock directly inserts a block on database and treats the validity of the block parent
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

// getChildren gets a slice of Block given the slice of keys
func getChildren(childrenKeys []string) ([]Block, error) {
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
