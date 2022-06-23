package utils

import (
	"redis-task/database"
	"strings"
)

func GetKeys(pattern string) []string {
	db := database.ConnectWithDB()
	result, err := db.Keys(database.CTX, pattern).Result()
	if err != nil {
		return nil
	}
	return result
}

func GetIndividualBlockId(compositeKey string) string {
	return strings.Split(compositeKey, ":")[0]
}

func UpdatedBlockId(key, parentKey string) string {
	blockKey := GetIndividualBlockId(key)
	return blockKey + ":" + parentKey
}
