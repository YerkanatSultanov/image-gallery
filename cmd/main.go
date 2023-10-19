package main

import (
	"image-gallery/db"
	"image-gallery/internal/user"
	"image-gallery/router"
	"log"
)

func main() {
	dbConn, err := db.NewDataBase()
	if err != nil {
		log.Fatalf("could not initialize database connection %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRep)
	userHandler := user.NewHandler(userService)

	router.InitRouters(userHandler)
	err = router.Start("0.0.0.0:8080")
	if err != nil {
		return
	}
}
