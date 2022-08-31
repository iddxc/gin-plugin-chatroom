# gin-plugin-chatroom
基于gin的聊天室插件，使用Redis进行缓存。

## 使用

### 载入插件
运行指令`go get github.com/GPorter-t/gin-plugin-chatroom`进行安装。

在 gin 项目的初始化过程中加入 plugin的注册项，如果已经设置好该步骤，则直接调用即可。

**示例参考[gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin/blob/main/server/initialize/plugin.go)的插件安装**

```go
package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/email"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/plugin"
	"github.com/gin-gonic/gin"

    chatroom "github.com/GPorter-t/gin-plugin-chatroom"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
	}
}

func InstallPlugin(Router *gin.Engine) {
	PublicGroup := Router.Group("")
	fmt.Println("无鉴权插件安装==》", PublicGroup)
	PrivateGroup := Router.Group("")
	fmt.Println("鉴权插件安装==》", PrivateGroup)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	//  添加跟角色挂钩权限的插件 示例 本地示例模式于在线仓库模式注意上方的import 可以自行切换 效果相同
	PluginInit(PrivateGroup, email.CreateEmailPlug(
		global.GVA_CONFIG.Email.To,
		global.GVA_CONFIG.Email.From,
		global.GVA_CONFIG.Email.Host,
		global.GVA_CONFIG.Email.Secret,
		global.GVA_CONFIG.Email.Nickname,
		global.GVA_CONFIG.Email.Port,
		global.GVA_CONFIG.Email.IsSSL,
	))

    /* 初始化聊天室插件 */
    PluginInit(PrivateGroup, chatroom.CreateChatRoomPlugin(
		global.GVA_CONFIG.Redis.Addr,
		global.GVA_CONFIG.Redis.Password,
		global.GVA_CONFIG.Redis.DB,
	))
}
```

### 调用
#### 创建聊天室
```python
import requests

data = {
    "sender": "a",  # 发送方，聊天室的成员
    "recipients": ["b", "d"], # 接收者，聊天室的成员
    "message": "hello world", # 消息主体
    "room_name": "test", # 聊天室名称
}


def test_create_room():
    r = requests.post("http://{% your_server_address %}/chatroom/create", json=data)
    print(r.text)

if __name__ == "__main__":
    test_create_room()
    # success: {"code":0,"data":"{% chatroom_id %}","msg":"查询成功"}
```

#### 加入聊天室
```python
import requests

data = {
    "sender": "a",
    "chat_id": "ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)"
}

def test_join_room():
    r = requests.post("http://{% your_server_address %}/chatroom/join", json=data)
    print(r.text)

if __name__ == "__main__":
    test_join_room()
    # success: {"code":0,"data":"ok","msg":"查询成功"}
```

#### 退出聊天室
```python
import requests

data = {
    "sender": "a",
    "chat_id": "ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)"
}

def test_leave_room():
    r = requests.post("http://{% your_server_address %}/chatroom/leave", json=data)
    print(r.text)



if __name__ == "__main__":
    test_leave_room()
    # success: {"code":0,"data":"ok","msg":"查询成功"}
```

#### 解散聊天室
```python
import requests

data = {
    "sender": "a",
    "chat_id": "ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)"
}

def test_leave_room():
    r = requests.post("http://{% your_server_address %}/chatroom/delete", json=data)
    print(r.text)


if __name__ == "__main__":
    test_leave_room()
    # success: {"code":0,"data":"ok","msg":"查询成功"}
```

#### 发送消息
```python
import requests

data = {
    "sender": "a",  # 发送方，聊天室的成员
    "message": "hello world", # 消息主体
    "chat_id": "ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)"
}

def test_leave_room():
    r = requests.post("http://{% your_server_address %}/chatroom/leave", json=data)
    print(r.text)


if __name__ == "__main__":
    test_leave_room()
    # success: {"code":0,"data":"ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)","msg":"查询成功"}
```

#### 获取历史数据
```python
import requests

data = {
    "sender": "a",  # 操作方，聊天室的成员，为请求者
    "chat_id": "ecb1b4f6-2d28-42af-a62c-ff996a3454a2(chatroom_id)"
}

def test_leave_room():
    r = requests.post("http://{% your_server_address %}/chatroom/leave", json=data)
    print(r.text)


if __name__ == "__main__":
    test_leave_room()
    # success: {"code":0,"data":[{"ts":"2022-08-31T13:50:57.5146853+08:00","sender":"t","message":"hello world"},{"ts":"2022-08-31T13:49:30.4761773+08:00","sender":"b","message":"hello world"}],"msg":"查询成功"}
```

## 设计思路
- 在调用过程中使用用户的uuid进行唯一标识用户，可避免用户的其他数据泄露；
- 从时间的一维性考虑，使得数据存储的结构为list，且一个uuid4 进行维护一个聊天室
- 仅聊天室的创作者有权解散所创建的聊天室，参考微信的建群
- 一个用户可以加入多个聊天室，数据结构：set
- 用户可自主退出聊天室，当聊天室在线人数为0时，删除该聊天室，防止无效缓存
- 可选用并指定聊天室的名称
- 取消未读标记，而是默认获取最近的历史数据，可提高其拓展性，参考微信聊天方式
- 理论上可使用超文本语言进行编写、发送数据

## 适用范围
- 弹幕，弹幕可抽象为临时聊天室，当用户进入时可看见历史数据，当弹幕为0且人数为0时，为假死聊天室
- 会话，不管是1v1的聊天室，XvX的聊天室，在获取数据时均使用窗口模式
- 评论，评论区可抽象为按照时间顺序发展的聊天室，此插件亦能适用，不足的地方是暂未添加删除数据功能

## 可优化项
- 信息的收发关乎数据的安全，可考虑对数据进行加解密处理，提高数据的安全性
- 消息的撤回 与 删除
- 消息的审核，规避不良发言、不良图片
