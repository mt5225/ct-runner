package main

import (
	"github.com/ct-runner/container"
	"github.com/gin-gonic/gin"
)

const terraformImage = "radut/terraform-ansible:latest"

func main() {
	r := gin.Default()
	r.Static("/js", "./templates/js")
	r.LoadHTMLGlob("./templates/*.html")
	api := r.Group("/api")
	{
		api.GET("/run", func(ctx *gin.Context) {
			cmd := new(container.Command)
			cmd.Image = terraformImage
			cmd.Env = make([]string, 0)
			cmd.Commands = []string{"sh", "-c", `for i in {1..10}; do sleep 1&&echo $i; done`}
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

	ses.WebClient(r)
	r.Run(":5000")
}
