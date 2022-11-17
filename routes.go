package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {

	r.GET("/", routeHome)

}

func routeHome(c *gin.Context) {
	c.Status(http.StatusOK)
}
