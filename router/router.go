// 2、router 定义了“别人怎么叫我干”
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (registry *Registry) Register(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	v1Router := r.Group("api/v1")
	v1Router.POST("/keygen", registry.KeygenHandler())
	v1Router.POST("/sign", registry.SignTxHandler())
	v1Router.GET("/metrics", registry.PrometheusHandler())
}
