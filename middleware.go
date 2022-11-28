package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationHeader struct {
	Basic string `header:"Authorization"`
}

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := AuthorizationHeader{}

		err := c.ShouldBindHeader(&authorizationHeader)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authorizationHeader.Basic, "Basic ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		basicAuthString, err := base64.StdEncoding.DecodeString(authorizationHeader.Basic[6:])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		basicAuthParts := strings.Split(string(basicAuthString), ":")

		username := basicAuthParts[0]
		password := basicAuthParts[1]

		if user, found := cache.Users[username]; found {
			// check if basic auth string is already cached for user
			if user.basicAuth != "" && user.basicAuth == authorizationHeader.Basic {
				c.Set("userId", user.Id)
				c.Next()
				return
			}
			// check if credentials are invalid
			if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			} else {
				// credentials are valid, save basic auth string in cache to speed up auth
				user.basicAuth = authorizationHeader.Basic
				cache.Users[username] = user
				c.Set("userId", user.Id)
				c.Next()
				return
			}
		} else {
			// credentials are invalid
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
