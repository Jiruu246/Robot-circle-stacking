package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

const GridSize = 3

func main() {
	dataStore := NewDataStore()
	service := NewService(dataStore)
	handler := NewHandler(service)

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/state", handler.GetState)
	r.POST("/command", handler.ProcessCommand)
	r.GET("/export", handler.ExportHistory)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
