package main

import (
	"fmt"
	"redis-task/database/scripts"
	"redis-task/models"
)

func main() {
	scripts.FlushDatabase()
	scripts.PopulateDatabase(nil)
	err := models.DeleteBlockById("F3")
	if err != nil {
		fmt.Println(err)
	}
	//routes.HandleRoutes()
}
