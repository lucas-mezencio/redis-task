package models

import (
	"redis-task/utils"
	"reflect"
)

// Tree Definition of the Block tree of the application
type Tree struct {
	Block    Block  `json:"block,omitempty"`
	Children []Tree `json:"children,omitempty"`
}

// GetTreeById gets a Block tree given a block id as root block.
//
// in case of error, invalid (or nonexistent) block or database empty returns an empty tree.
//
// Returns a Tree of Blocks.
func GetTreeById(id string) Tree {
	var tree Tree
	var blockId string
	var keysChildren []string
	if id == "0" {
		keysChildren = utils.GetKeys("*:" + id)
		tree.Block = Block{}
	} else {
		tree.Block = GetBlockById(id)
		if reflect.DeepEqual(tree.Block, Block{}) {
			return Tree{}
		}
		blockId = utils.GetIndividualBlockId(tree.Block.ID)
		keysChildren = utils.GetKeys("*:" + blockId)
	}

	for _, keyChild := range keysChildren {
		tree.Children = append(tree.Children, GetTreeById(utils.GetIndividualBlockId(keyChild)))
	}
	return tree
}
