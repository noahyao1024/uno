package main

import (
	"fmt"
	"uno/pkg/database"
	"uno/pkg/setting"
	"uno/router"
	"uno/util"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func init() {
	database.Setup()
	util.Setup()
}

func main() {
	gin.SetMode(setting.ServerInstance.RunMode)

	endless.ListenAndServe(fmt.Sprintf(":%d", setting.ServerInstance.HttpPort), router.Setup())
}
