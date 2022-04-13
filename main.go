package main

import (
	"fmt"
	"ip_lookup/controllers"
	"ip_lookup/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ip/myip", func(c *gin.Context) {
		ip, err := controllers.GetClientIPHelper(c.Request)
		if err == nil {
			c.String(http.StatusOK, fmt.Sprint(ip))
		} else {
			c.String(http.StatusBadRequest, "")
		}

	})
	r.GET("/ip/:ip_address", func(c *gin.Context) {
		ip := c.Params.ByName("ip_address")
		c.String(http.StatusOK, fmt.Sprint(controllers.SingleIPLookup(ip)))
	})

	// TODO: Complete multi IP post func
	r.POST("/ip", func(c *gin.Context) {
		ippost := models.IPPost{}
		if err := c.BindJSON(&ippost); err != nil {
			return
		}

	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
