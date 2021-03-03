package main

import (
	"github.com/ct-runner/v0.0.1/container"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func main() {
	Router = gin.Default()
	api := Router.Group("/api")
	{
		api.GET("/run", func(ctx *gin.Context) {
			cmd := new(container.Command)
			cID, err := cmd.Run()
			if err != nil {
				ctx.JSON(200, gin.H{
					"message": "fail to run container",
				})
			} else {
				ctx.JSON(200, gin.H{
					"message": cID,
				})
			}

		})
	}
	Router.Run(":5000")
}
