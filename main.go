package main

import (
	"github.com/gerbsec/D2/routes"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Setup routes for listeners
	routes.SetupListenerRoutes(r)

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
