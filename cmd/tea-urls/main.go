package main

import (
	"log"

	"github.com/ZelenyMK/tea-urls/internal/config"
	"github.com/ZelenyMK/tea-urls/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	log.Println(config.Scheme + config.Host + "/" + config.Port)
	log.Println(config.Scheme + config.Host + "/" + config.Port)
	log.Println(config.Scheme + config.Host + "/" + config.Port)
	linkHandler := handler.NewLinkHandler(
		config.Scheme + config.Host + "/" + config.Port,
	)

	router := gin.Default()

	router.POST("/links", linkHandler.Create)
	router.GET("/:alias", linkHandler.Redirect)

	log.Println("server listening on http://localhost:8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
