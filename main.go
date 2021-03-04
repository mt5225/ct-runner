package main

import (
	"github.com/ct-runner/v0.0.1/container"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

const TerraformImage string = "radut/terraform-ansible:latest"

func main() {
	Router = gin.Default()
	api := Router.Group("/api")
	{
		api.POST("/run", func(ctx *gin.Context) {
			cmd := new(container.Command)
			cmd.Image = TerraformImage
			cmd.Env = make([]string, 0)
			cmd.Commands = []string{"sh", "-c", `terraform version`}
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
