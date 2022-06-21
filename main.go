package main

import (
	"redis-task/database/scripts"
)

func main() {
	scripts.FlushDatabase()
	scripts.PopulateDatabase(nil)
	//routes.HandleRoutes()
}
