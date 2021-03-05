package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mt5225/ct-runner/container"
	"github.com/mt5225/ct-runner/ssesvr"
)

const terraformImage = "radut/terraform-ansible:latest"

func main() {
	r := gin.Default()
	webClient(r)
	stream := new(ssesvr.RespStream)
	stream.SseServer(r)

	api := r.Group("/api")
	{
		api.GET("/run", func(ctx *gin.Context) {
			cmd := new(container.Command)
			cmd.Image = terraformImage
			cmd.Env = make([]string, 0)
			cmd.Commands = []string{"sh", "-c", `watch -n 3 uptime`}
			cID, err := cmd.Run(stream)
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
	r.Run(":5000")
}

func webClient(r *gin.Engine) {
	r.Static("/js", "./templates/js")
	r.LoadHTMLGlob("./templates/*.html")
	r.GET("/log", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
}
