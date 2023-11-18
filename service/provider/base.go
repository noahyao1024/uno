package provider

import (
	"uno/service/subscriber"
	"uno/service/template"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Sender string
}

type Response struct {
	MessageID string
}

type Base interface {
	SetOption(*gin.Context, *Option) error
	Send(*gin.Context, *subscriber.Entry, *template.Entry) (*Response, error)
}
