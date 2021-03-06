package runner

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ContainerRunner the container runner
func ContainerRunner(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/run", func(ctx *gin.Context) {
			c, err := ReqToCommand(ctx.Request)
			if err != nil {
				ctx.JSON(500, gin.H{
					"message": "fail to ready message",
				})
				return
			}

			err = c.Create() // create container instance
			if err != nil {
				log.Println(err)
				ctx.JSON(200, gin.H{
					"message": "fail to create container ",
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
