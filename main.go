package main

import (
	"log"
	"product-cmp-api/api"
	"product-cmp-api/storage"
)

func main() {

	store, err := storage.NewStore()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postgress connected successfully")

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", store)
	server.Run()
}
