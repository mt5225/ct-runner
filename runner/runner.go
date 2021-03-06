package runner

import (
	"context"

	"github.com/gin-gonic/gin"
)

// docker image to run terraform
const terraformImage = "radut/terraform-ansible:latest"

// ContainerRunner the container runner
func ContainerRunner(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/run", func(ctx *gin.Context) {
			c := new(Command)
			c.Image = terraformImage
			c.Env = make([]string, 0)
			c.Commands = []string{"sh", "-c", `sleep 60&&uptime`}
			c.ContainerInstance = new(Container)
			c.ContainerInstance.Context = context.Background()

			err := c.Create() // create container instance
			if err != nil {
				ctx.JSON(200, gin.H{
					"message": "fail to create container",
				})
			} else {
				ctx.JSON(200, gin.H{
					"container_id": c.ContainerInstance.ID,
				})
				c.Run(r)       // run command
				c.Streaming(r) //straming result
			}
		})
	}
}
