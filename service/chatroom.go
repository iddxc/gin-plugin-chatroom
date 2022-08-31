package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/GPorter-t/gin-plugin-chatroom/global"
	"github.com/GPorter-t/gin-plugin-chatroom/model"
	"github.com/GPorter-t/gin-plugin-chatroom/utils"

	"github.com/google/uuid"
)

type ChatRoomService struct{}

var ctx = context.Background()

/*
func: CreateChat
desc: 根据 发送者、信息、房间名称、接收者 进行创建会话
return err
*/
func (c *ChatRoomService) CreateChat(sender, message, roomName string, recipients []string) (chatId string, err error) {
	chatId = uuid.NewString()
	recipients = append(recipients, sender)

	pipe := global.GVA_REDIS.Pipeline()
	for _, recipient := range recipients {
		pipe.SAdd(ctx, "chat:"+chatId, recipient)
		pipe.SAdd(ctx, "user:"+recipient, chatId)
	}
	pipe.HSet(ctx, "chat:room", chatId, roomName)
	pipe.HSet(ctx, "chat:creator", chatId, sender)
	pipe.Exec(ctx)
	chatId, err = c.SendMessage(sender, message, chatId)
	return
}

/*
func: SendMessage
desc: 向指定会话id发送消息，必须发送者为会话id中的成员
return: chat_id, err
*/
func (c *ChatRoomService) SendMessage(sender, message, chatId string) (chat_id string, err error) {
	if !c.IsInChat(chatId, sender) {
		return "", errors.New("not in chat")
	}
	stat, identifier := utils.AcquireLock("chat:"+chatId, 10)
	if !stat {
		return "", errors.New("get the lock")
	}
	ts := time.Now()
	packed := model.Message{
		Sender:  sender,
		Message: message,
		TS:      ts,
	}
	packed_data, _ := json.Marshal(packed)

	count, err := global.GVA_REDIS.LPush(ctx, "msgs:"+chatId, packed_data).Result()
	if err != nil {
		return "", err
	}
	fmt.Println(count)
	utils.ReleaseLock("chat:"+chatId, identifier)
	return chatId, nil
}

/*
func: IsInChat
desc: 判断用户是否为会话id的成员
return bool
*/
func (c *ChatRoomService) IsInChat(chatId, userId string) bool {
	return global.GVA_REDIS.SIsMember(ctx, "chat:"+chatId, userId).Val()
}

/*
func: isChat
desc: 判断会话id是否有效
return bool
*/
func (c *ChatRoomService) isChat(chatId string) (bool, error) {
	return global.GVA_REDIS.HExists(ctx, "chat:room", chatId).Result()
}

/*
func: isInUsers
*/
func (c *ChatRoomService) IsInUsers(userId string) bool {
	return global.GVA_REDIS.SIsMember(ctx, "users:", userId).Val()
}

/*
func: JoinUsers
*/
func (c *ChatRoomService) JoinUsers(userId string) bool {
	pipe := global.GVA_REDIS.Pipeline()
	pipe.SAdd(ctx, "users:", userId)
	pipe.Exec(ctx)
	return true
}

/*
func：JoinChat
desc: 用户加入指定会话中
return bool, err
*/
func (c *ChatRoomService) JoinChat(chatId, userId string) (stat bool, err error) {
	stat, err = c.isChat(chatId)
	if !stat {
		return false, fmt.Errorf("join Chat Failed: %v", err)
	}
	pipe := global.GVA_REDIS.Pipeline()
	pipe.SAdd(ctx, "chat:"+chatId, userId)
	pipe.SAdd(ctx, "user:"+userId, chatId)
	pipe.Exec(ctx)
	return true, nil
}

/*
func: LeaveChat
desc: 用户退出指定会话
return bool, err
*/
func (c *ChatRoomService) LeaveChat(chatId, userId string) (stat bool, err error) {
	stat, identifier := utils.AcquireLock("chat:"+chatId, 10)
	if !stat {
		return false, fmt.Errorf("acquire lock Failed: %v", err)
	}

	stat = c.IsInChat(chatId, userId)
	if !stat {
		return false, fmt.Errorf("leave Chat Failed: %v", err)
	}
	pipe := global.GVA_REDIS.Pipeline()
	pipe.SRem(ctx, "chat:"+chatId, userId)
	pipe.SRem(ctx, "user:"+userId, chatId)
	pipe.Exec(ctx)

	members := global.GVA_REDIS.SCard(ctx, "chat:"+chatId).Val()
	if members == 0 {
		pipe.Del(ctx, "msgs:"+chatId)
		pipe.Del(ctx, "ids:"+chatId)
		pipe.Exec(ctx)
	}
	utils.ReleaseLock("chat:"+chatId, identifier)
	return true, nil
}

/*
func: DissolveChat
desc: 解散会话，必须用户为会话的创建者
return bool, err
*/
func (c *ChatRoomService) DissolveChat(chatId, userId string) (stat bool, err error) {
	uId := global.GVA_REDIS.HGet(ctx, "chat:creator", chatId).Val()
	if uId != userId {
		return false, fmt.Errorf("operator is not creator")
	}

	userIds := global.GVA_REDIS.SMembers(ctx, "chat:"+chatId).Val()
	for _, memberId := range userIds {
		t, err := c.LeaveChat(chatId, memberId)
		if err != nil {
			return t, err
		}
	}
	stat, identifier := utils.AcquireLock("chat:"+chatId, 10)
	if !stat {
		return false, fmt.Errorf("lock for chat: %v", chatId)
	}
	pipe := global.GVA_REDIS.Pipeline()
	pipe.HDel(ctx, "chat:creator", chatId)
	pipe.HDel(ctx, "chat:room", chatId)
	pipe.Exec(ctx)

	utils.ReleaseLock("chat:"+chatId, identifier)
	return true, nil
}

/*
func: History
desc: 根据page和count获取会话数据，用户必须为会话成员，
return []Message, error
*/
func (c *ChatRoomService) History(chatId, userId string, page, count int) (messages []model.Message, err error) {
	if !c.IsInChat(chatId, userId) {
		return
	}
	msgs, err := global.GVA_REDIS.LRange(ctx, "msgs:"+chatId, int64(page*count), int64(page*count+count)).Result()
	for _, msg := range msgs {
		message := model.Message{}
		err = json.Unmarshal([]byte(msg), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return
}
