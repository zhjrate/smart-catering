package selfTool

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
	tool2 "zhcy/tool"
	"zhcy/util/dbSQL"
)

type ToolK struct {
	tool2.Toolkit
	tool2.MailBox
}

var Toolkit = ToolK{}

var Dbsql dbSQL.Database

func GetMysqlConnect(sql string) dbSQL.Database {

	dbSQL.Connect()
	Dbsql.Sql = sql
	Dbsql.Right = "加载商品信息成功"
	Dbsql.Wrong = "加载商品信息失败"
	return Dbsql
}

func GetCommidity(data []map[string]interface{}) []map[string]string {
	var strArray []map[string]string //装商品信息的二维数组
	//// 初始化二维切片
	//for i := 0; i < len(data); i++ {
	//
	//	// 将内部切片添加到外部切片
	//	strArray = append(strArray, strArray[i])
	//}
	//strArray = make([]map[string]string, 100)
	fmt.Println(len(data))
	fmt.Println(len(strArray))
	for i := 0; i < len(data); i++ {
		strArray = append(strArray, map[string]string{})
		for j := 0; j < 7; j++ {
			switch j {
			case 0:

				strArray[i]["name"] = data[i]["name"].(string)

			case 1:
				strArray[i]["description"] = data[i]["description"].(string)

			case 2:
				strArray[i]["price"] = fmt.Sprintf("%v", data[i]["price"])

			case 3:
				strArray[i]["num"] = fmt.Sprintf("%v", data[i]["num"])

			case 4:
				if data[i]["startTime"] == nil {
					strArray[i]["startTime"] = ""
				} else {
					strArray[i]["startTime"] = data[i]["startTime"].(string)
				}
			case 5:
				if data[i]["updateTime"] == nil {
					strArray[i]["updateTime"] = ""
				} else {
					strArray[i]["updateTime"] = data[i]["updateTime"].(string)
				}
			case 6:
				strArray[i]["id"] = fmt.Sprintf("%v", data[i]["id"])

			}
		}

	}
	return strArray
}

// 进行获取存入redis缓存的数据
func GetRedis(c *gin.Context) string {
	userss := make(map[string]string)
	token := GenerateToken()
	fmt.Println("登录时产生的token", token)
	ctx1, redisC := dbSQL.RedisConnect()
	b, s := dbSQL.RedisAdd(userss["username"], token, time.Duration(240)*time.Minute)
	if b != true {
		fmt.Println("查询到token不存在,登录时产生的token1", s)
		return s
	}
	// 从 Redis 缓存中获取数据

	a := userss["username"]
	val, err := redisC.Get(ctx1, a).Result() //获取成功
	//fmt.Println(a)
	//
	if err == redis.Nil {
		fmt.Println("Key does not exist")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("查询到token存在,获取原来的token:", val)
	}
	return val
}

// 定期轮询数据库，实时更新
func WatchDatabase(dataChannel chan<- string) {
	db := GetMysqlConnect("SELECT * FROM commodity WHERE some_condition LIMIT 1")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 查询数据库中的数据并获取需要实时监测的列的值
			var updatedData string
			a := db.SelW(updatedData)
			if a != nil {
				log.Println("查询数据库出错:", a)
				continue
			}

			// 将更新的数据发送到通道中
			dataChannel <- updatedData
		}
	}
}

var way = dbSQL.Database{}

// 设置定时任务
func SetDs(status string, id string) bool {
	// 设置定时任务
	go func() {
		ticker := time.NewTicker(10 * time.Second) // 每隔5秒执行一次任务
		for {
			select {
			case <-ticker.C:
				way.Sql = "SELECT status from commodity where id = ?;"
				data := way.SelS(id)
				data1 := Toolkit.JsonSlice(data)
				data2, _ := Toolkit.JsonList(data1["list"])

				if data1["code"] == 1 {
					way.Sql = "update commodity set status=? where id = ?;"
					if data2[0]["status"].(string) == "1" {
						way.Right = "商品下架成功"
						way.Wrong = "商品下架失败"
						way.UPS(status, id)

					}
					way.Right = "商品上架成功"
					way.Wrong = "商品上架失败"
					fmt.Println("定时任务执行:", time.Now())

				}
			}
		}
	}()

	return true
}
