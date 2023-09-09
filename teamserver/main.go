package main

import (
	"github.com/gerbsec/G2/teamserver/routes"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	routes.SetupListenerRoutes(r)
	routes.SetupAgentRoutes(r)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
