package main

import (
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
	// TODO: do not access cache properties directly
	c.JSON(http.StatusOK, cache.Movies)
}

func routePrivateSeries(c *gin.Context) {
	// TODO: do not access cache properties directly
	c.JSON(http.StatusOK, cache.Series)
}
