package main

import (
	"context"
	"golang-practise-project/routers"
	"golang-practise-project/utils/database"
	"log"
)

func main() {
	ctx := context.Background()
	initDatabase(ctx)
	initRouter(ctx)
}

func initDatabase(ctx context.Context) {
	var err error
	err = database.InitDatabase()
	if err != nil {
		log.Fatal("Error connecting database")
	}
}

// initRouter initialise the router
func initRouter(ctx context.Context) {
	router := routers.GetRouter(ctx)
	err := router.Run()
	if err != nil {
		log.Fatal("Error occcured")
	}
}
