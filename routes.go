package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {

	v1 := r.Group("/api/v1")

	public := v1.Group("public")
	{
		public.GET("/", routePublicHome)
	}

	private := v1.Group("private", AuthorizationMiddleware())
	{
		private.GET("/", routePrivateHome)
		private.GET("/movies", routePrivateMovies)
		private.GET("/series", routePrivateSeries)
	}
}

func routePublicHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func routePrivateHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"authorized": true,
	})
}

func routePrivateMovies(c *gin.Context) {

	movies, err := getAllMovies(connectionPool)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, movies)

}

func routePrivateSeries(c *gin.Context) {

	movies, err := getAllSeries(connectionPool)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, movies)

}
