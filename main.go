package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	hostname string = "0.0.0.0"
	port     int    = 8080
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Minute,
	}))

	routes(r)

	if err := r.Run(fmt.Sprintf("%s:%d", hostname, port)); err != nil {
		log.Fatal(err)
	}

}
