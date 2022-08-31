package chatroom

import (
	"chatroom/global"
	"chatroom/initialize"
	"chatroom/router"

	"github.com/gin-gonic/gin"
)

type chatRoomPlugin struct{}

func CreateChatRoomPlugin(Addr, Password string, DB int) *chatRoomPlugin {
	global.GVA_CONFIG.Redis.Addr = Addr
	global.GVA_CONFIG.Redis.Password = Password
	global.GVA_CONFIG.Redis.DB = DB
	initialize.Redis()
	return &chatRoomPlugin{}
}

func (*chatRoomPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitChatRoomRouter(group)
}

func (*chatRoomPlugin) RouterPath() string {
	return "chatroom"
}
