package main

import (
	"fmt"
	"uno/pkg/setting"
	"uno/pkg/util"
	"uno/router"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func init() {
	util.Setup()
}

func main() {
	gin.SetMode(setting.ServerInstance.RunMode)

	endless.ListenAndServe(fmt.Sprintf(":%d", setting.ServerInstance.HttpPort), router.Setup())
}
