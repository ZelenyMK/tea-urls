package main

import (
	"log"
	"net/http"

	"github.com/ZelenyMK/tea-urls/internal/config"
	"github.com/ZelenyMK/tea-urls/internal/handler"
	"github.com/ZelenyMK/tea-urls/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	conn, err := storage.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.CloseDB(conn)

	redisClient := storage.NewRedisClient()
	defer redisClient.Close()

	hostURL := config.Scheme + config.Host + "/" + config.Port

	linkHandler := handler.NewLinkHandler(
		hostURL, conn,
	)

	router := gin.Default()

	router.LoadHTMLGlob("web/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.POST("/links", linkHandler.Create)
	router.GET("/:alias", linkHandler.Redirect)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	log.Println("server listening on " + hostURL)
}
