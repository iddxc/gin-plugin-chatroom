package router

import (
	"github.com/GPorter-t/gin-plugin-chatroom/api"

	"github.com/gin-gonic/gin"
)

type ChatRoomRouter struct{}

func (c *ChatRoomRouter) InitChatRoomRouter(Router *gin.RouterGroup) {
	//chatRoomRouter := Router.Use(middleware.OperationRecord())
	chatRoomRouter := Router
	{
		chatRoomRouter.POST("create", api.ApiGroupApp.CreateChatRoom)
		chatRoomRouter.POST("leave", api.ApiGroupApp.LeaveChat)
		chatRoomRouter.POST("join", api.ApiGroupApp.JoinChat)
		chatRoomRouter.POST("delete", api.ApiGroupApp.DissolveChat)
		chatRoomRouter.POST("send", api.ApiGroupApp.SendMessage)
		chatRoomRouter.POST("history", api.ApiGroupApp.GetHistory)
	}
}
