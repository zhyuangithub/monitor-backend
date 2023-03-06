package uiserver

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	broadcast "github.com/teivah/broadcast"

	"monitor-backend/docs"
	"monitor-backend/utils"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var instance *UiServer

type UiServer struct {
	relay *broadcast.Relay[string]
}

func UiServerInstance() *UiServer {
	if instance == nil {
		instance = new(UiServer)
		//instance.init()
	}
	return instance
}
func (u *UiServer) Init(relay *broadcast.Relay[string]) {
	u.relay = relay
	u.startServer()
}

func (u *UiServer) startServer() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		output := fmt.Sprintf("ui server:%s %s\n", utils.Version, time.Now().Format("15:04:05"))
		c.String(200, output)
	})
	r.GET("/GetNodeDataHistory", getNodeDataHistoryHandler)
	r.GET("/GetNodeInfo", getNodeInfoHandler)
	r.GET("/GetNodeEventLogs", getNodeEventLogsHandler)
	r.GET("/GetNodesInfo", nodesHandler)
	//编译swag命令：swag init
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/ChangeNodeName", changeNodeNameHandler)
	r.Any("/notificationCenter", func(c *gin.Context) {
		socketHandler(c.Writer, c.Request)
	})
	docs.SwaggerInfo.BasePath = fmt.Sprintf(":%s/", os.Getenv("UISERVER_PORT"))

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("UISERVER_PORT"))
	r.Run(addr)
}
