package routes

import (
	"net/http"

	"github.com/gerbsec/D2/listeners"
	"github.com/gin-gonic/gin"
)

func SetupListenerRoutes(r *gin.Engine) {
	r.GET("/Listeners/:name", getListenerInfoByName)
	r.GET("/Listeners", getAllListenersInfo)
	r.POST("/Listener", createListener)
	r.DELETE("/StopListener/:name", stopListener)
}

func getListenerInfoByName(c *gin.Context) {
	name := c.Param("name")
	info, err := listeners.GetListenerInfoByName(name)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, info)
}

func getAllListenersInfo(c *gin.Context) {
	c.String(http.StatusOK, listeners.GetAllListenersInfo())
}

func createListener(c *gin.Context) {
	var l listeners.HttpListener
	if err := c.BindJSON(&l); err != nil {
		c.String(http.StatusBadRequest, "Request body is not a valid listener")
		return
	}
	err := listeners.CreateListener(l.Name, l.BindPort)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusCreated, "Listener created")
}

func stopListener(c *gin.Context) {
	name := c.Param("name")
	listeners.StopListenerByName(name)
	c.String(http.StatusOK, "Listener stopped")
}
