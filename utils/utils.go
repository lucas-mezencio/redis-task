// Package utils holds the utility functions of the application
package utils

import (
	"redis-task/database"
	"strings"
)

// GetKeys returns the list of keys given a redis pattern of the command "KEYS"
//			KEYS <pattern>
// if the consult returns an error returns nil
//
// returns a string slice with the found keys
func GetKeys(pattern string) []string {
	db := database.ConnectWithDB()
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		return nil
	}
	return result
}

// GetIndividualBlockId gets the id of a block given its key
//			composite key:	thisBlock:parentBlock
// returns the individual id of a block
func GetIndividualBlockId(compositeKey string) string {
	return strings.Split(compositeKey, ":")[0]
}

// UpdatedBlockId updates a block key given the new parent id (key)
//
// returns the new modified id of a block
func UpdatedBlockId(key, parentKey string) string {
	blockKey := GetIndividualBlockId(key)
	return blockKey + ":" + parentKey
}
