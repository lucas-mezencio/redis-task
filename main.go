package main

import (
	"redis-task/database/scripts"
	"redis-task/routes"
)

func main() {
	scripts.FlushDatabase()
	scripts.PopulateDatabase()
	routes.HandleRoutes()
}
