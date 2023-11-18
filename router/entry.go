package router

import (
	"net/http"
	"time"
	"uno/router/dispatch"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

		go time.Sleep(10 * time.Second)
	})

	r.Any("/v1/topic", dispatch.TopicCreate)
	r.Any("/v1/workflow", dispatch.WorkflowCreate)

	return r
}
