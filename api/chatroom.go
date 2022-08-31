package api

import (
	"github.com/GPorter-t/gin-plugin-chatroom/model"
	"github.com/GPorter-t/gin-plugin-chatroom/model/response"
	"github.com/GPorter-t/gin-plugin-chatroom/service"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatRoomApi struct{}

func (r *ChatRoomApi) CreateChatRoom(c *gin.Context) {
	var messageReq model.MessageReq
	if err := c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chatId, err := service.ServiceGroupApp.CreateChat(messageReq.Sender, messageReq.Message, messageReq.RoomName, messageReq.Recipients)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(chatId, c)
}

func (r *ChatRoomApi) SendMessage(c *gin.Context) {
	var messageReq model.MessageReq
	if err := c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	chatId, err := service.ServiceGroupApp.SendMessage(messageReq.Sender, messageReq.Message, messageReq.ChatId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(chatId, c)
}

func (r *ChatRoomApi) JoinChat(c *gin.Context) {
	var messageReq model.MessageReq
	if err := c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	stat, err := service.ServiceGroupApp.JoinChat(messageReq.ChatId, messageReq.Sender)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !stat {
		response.Fail(c)
		return
	}
	response.Ok(c)
}

func (r *ChatRoomApi) LeaveChat(c *gin.Context) {
	var messageReq model.MessageReq
	if err := c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	stat, err := service.ServiceGroupApp.LeaveChat(messageReq.ChatId, messageReq.Sender)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !stat {
		response.Fail(c)
		return
	}
	response.Ok(c)
}

func (r *ChatRoomApi) DissolveChat(c *gin.Context) {
	var messageReq model.MessageReq
	if err := c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	stat, err := service.ServiceGroupApp.DissolveChat(messageReq.ChatId, messageReq.Sender)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !stat {
		response.Fail(c)
		return
	}
	response.Ok(c)
}

func (r *ChatRoomApi) GetHistory(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	count, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var messageReq model.MessageReq
	if err = c.ShouldBindJSON(&messageReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	messages, err := service.ServiceGroupApp.History(messageReq.ChatId, messageReq.Sender, page, count)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(messages, c)
}
