package runner

import "github.com/gin-gonic/gin"

// StaticContent
func StaticContent(r *gin.Engine) {
	r.Static("/js", "./templates/js")
	r.LoadHTMLGlob("./templates/*.html")
	r.GET("/log", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
}
