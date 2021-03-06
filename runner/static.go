package runner

import "github.com/gin-gonic/gin"

// StaticContent site
func StaticContent(r *gin.Engine) {
	r.Static("/js", "./templates/js")
	r.LoadHTMLGlob("./templates/*.html")
	r.GET("/logs/:container_id", func(c *gin.Context) {
		// we need the docker container id
		containerID := c.Param("container_id")
		c.HTML(200, "index.html", gin.H{
			"ContainerID": containerID,
		})
	})
}
