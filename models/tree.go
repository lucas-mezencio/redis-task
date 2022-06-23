package models

import (
	"redis-task/utils"
	"reflect"
)

type Tree struct {
	Block    Block  `json:"block,omitempty"`
	Children []Tree `json:"children,omitempty"`
}

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
