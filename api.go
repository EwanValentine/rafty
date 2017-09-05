package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func API(rafty *Rafty) {
	handlers := &Handlers{rafty}
	r := gin.Default()
	r.GET("/api", handlers.Index)
	log.Println("Running API on :9000")
	r.Run(":9000")
}

type Handlers struct {
	rafty *Rafty
}

func (handlers *Handlers) Index(c *gin.Context) {
	c.JSON(200, &handlers.rafty)
}
