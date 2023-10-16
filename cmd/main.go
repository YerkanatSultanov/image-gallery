package main

import (
	"image-gallery/db"
	"log"
)

func main() {
	_, err := db.NewDataBase()
	if err != nil {
		log.Fatalf("could not initialize database connection %s", err)
	}
}
