package service

import (
	"fmt"

	"github.com/GPorter-t/gin-plugin-chatroom/global"
	"github.com/GPorter-t/gin-plugin-chatroom/initialize"

	"testing"
)

func init() {
	global.GVA_CONFIG.Redis.Addr = "127.0.0.1:6379"
	global.GVA_CONFIG.Redis.DB = 0
	global.GVA_CONFIG.Redis.Password = ""
	initialize.Redis()
}

var service = new(ChatRoomService)

func TestCreateChat(t *testing.T) {
	recipients := []string{"a", "b", "c", "d", "e", "f"}
	chatId, err := service.CreateChat("admin", "test create chat", "", recipients)
	fmt.Println(chatId, err)
}

func TestJoinChat(t *testing.T) {
	service.JoinChat("c321102b-b6d7-40fe-9d79-4ad5a558860b", "g")
}

func TestSendMessage(t *testing.T) {
	charId := "c321102b-b6d7-40fe-9d79-4ad5a558860b"
	Id, err := service.SendMessage("admin", "test send message", charId)
	fmt.Println(err, Id)
}

func TestGetHistory(t *testing.T) {
	history, err := service.History("c321102b-b6d7-40fe-9d79-4ad5a558860b", "a", 0, 10)
	fmt.Println(err, history)
}

func TestLeaveChat(t *testing.T) {
	stat, err := service.LeaveChat("c321102b-b6d7-40fe-9d79-4ad5a558860b", "e")
	fmt.Println(err, stat)
}

func TestDissolveChat(t *testing.T) {
	stat, err := service.DissolveChat("e3584a19-2994-4c01-be2e-8a98b75493fa", "e")
	fmt.Println(err, stat)

	stat, err = service.DissolveChat("e3584a19-2994-4c01-be2e-8a98b75493fa", "admin")
	fmt.Println(err, stat)
}
