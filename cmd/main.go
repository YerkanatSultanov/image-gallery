package main

import (
	"image-gallery/internal/user/db"
	"image-gallery/internal/user/repo"
	"image-gallery/internal/user/service"
	"log"
)

func main() {
	dbConn, err := db.NewDataBase()
	if err != nil {
		log.Fatalf("could not initialize database connection %s", err)
	}

	userRep := repo.NewRepository(dbConn.GetDB())
	userService := service.NewService(userRep)
	userHandler := service.NewHandler(userService)

	service.InitRouters(userHandler)
	err = service.Start("0.0.0.0:8080")
	if err != nil {
		return
	}
}
