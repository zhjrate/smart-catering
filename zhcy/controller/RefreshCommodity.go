package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	_select "zhcy/controller/select"
	"zhcy/tool/selfTool"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// 实时监控商品
func RefreshCommodity(c *gin.Context) {
	// 建立WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"code": 0, "message": err})
		return
	}

	// 将新连接添加到clients映射中
	clientsMu.Lock()
	clients[ws] = true
	clientsMu.Unlock()

	// 接收 WebSocket 消息的 goroutine
	go func() {
		defer func() {
			// 连接断开时从clients映射中删除
			clientsMu.Lock()
			delete(clients, ws)
			clientsMu.Unlock()

			// 关闭 WebSocket 连接
			ws.Close()
		}()

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	// 运行数据定期刷新函数
	UpdateDataPeriodically()
}

func UpdateDataPeriodically() {

	for {
		// 实际应用中，你需要根据你的数据库和数据更新的逻辑来获取最新的数据
		data := FetchDataFromDatabase()

		// 将数据发送给所有连接的客户端
		clientsMu.Lock()
		for client := range clients {
			err := client.WriteJSON(data)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()

		// 模拟每秒钟更新一次数据
		sleepTime := 1 // seconds
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func FetchDataFromDatabase() []map[string]string {
	// 实际应用中，你需要根据你的数据库和数据更新的逻辑来获取最新的数据
	// 这里只是简单地返回一个示例数据

	db := selfTool.GetMysqlConnect("select * from commodity;")
	data := db.Sel()
	data1 := _select.Toolkit.JsonSlice(data)
	data2, _ := _select.Toolkit.JsonList(data1["list"])
	if data1["code"].(float64) == 1 {

		strArray := selfTool.GetCommidity(data2)
		return strArray
	}
	return nil
}
