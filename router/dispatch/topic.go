package dispatch

import (
	"fmt"
	"uno/service/topic"

	"github.com/gin-gonic/gin"
)

// TopicCreate is the entry for topic creation.
func TopicCreate(ctx *gin.Context) {
	entry := &topic.Entry{}
	if ctx.BindJSON(entry) != nil {
		ctx.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	topic.Create(entry)

	fmt.Println(entry)
}

func TopicGet(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	entry, err := topic.Get(id)
	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, entry)
}
