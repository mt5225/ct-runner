package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mt5225/ct-runner/runner"
)

func main() {
	r := gin.Default()
	runner.StaticContent(r)
	runner.ContainerRunner(r)
	r.Run(":5000")
}
