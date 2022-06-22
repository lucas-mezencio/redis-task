package models

import (
	"reflect"
	"strings"
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
		keysChildren = getKeys("*:" + id)
		tree.Block = Block{}
	} else {
		tree.Block = GetBlockById(id)
		if reflect.DeepEqual(tree.Block, Block{}) {
			return Tree{}
		}
		blockId = getIndividualBlockId(tree.Block.ID)
		keysChildren = getKeys("*:" + blockId)
	}

	for _, keyChild := range keysChildren {
		tree.Children = append(tree.Children, GetTreeById(getIndividualBlockId(keyChild)))
	}
	return tree
}

func getIndividualBlockId(compositeKey string) string {
	return strings.Split(compositeKey, ":")[0]
}
