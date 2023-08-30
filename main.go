package main

import (
	"net/http"

	http_listener "github.com/gerbsec/D2/listeners"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/Listeners/:name", func(c *gin.Context) {
		name := c.Param("name")
		info, err := http_listener.GetListenerInfoByName(name)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusOK, info)
	})

	r.GET("/Listeners", func(c *gin.Context) {
		c.String(http.StatusOK, http_listener.GetAllListenersInfo())
	})

	r.POST("/Listener", func(c *gin.Context) {
		var l http_listener.HttpListener
		if err := c.BindJSON(&l); err != nil {
			c.String(http.StatusBadRequest, "Request body is not a valid listener")
			return
		}
		err := http_listener.CreateListener(l.Name, l.BindPort)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusCreated, "Listener created")
	})

	r.DELETE("/StopListener/:name", func(c *gin.Context) {
		name := c.Param("name")
		http_listener.StopListenerByName(name)
		c.String(http.StatusOK, "Listener stopped")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
