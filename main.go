package main

import (
	"log"
	"os"
	"shortier/db"
	"shortier/handlers"

	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	pool := db.DbStart(databaseURL)

	db := db.NewDB(pool)

	handler := handlers.NewBaseHandler(db)

	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.POST("/create", handler.InsertLink)
		v1.GET("/:id", handler.GetLink)
	}

	r.Run()
}
