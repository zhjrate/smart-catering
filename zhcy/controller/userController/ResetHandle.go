package userController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"zhcy/tool/selfTool"
	"zhcy/util/dbSQL"
)

// 修改用户密码通过手机验证码
func ResetPassword(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from users where phone like ?;")
	userss := make(map[string]string)

	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(userss["phone"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	data2, _ := selfTool.Toolkit.JsonList(data1["list"])
	var id string
	if len(data2) > 0 {
		id = data2[0]["id"].(string)
	} else {
		id = ""
	}

	fmt.Println("从数据库查到的手机号对应的id为:", id)
	ctx1, redisC := dbSQL.RedisConnect()
	fmt.Println(userss["phone"])
	fmt.Println(data1["code"].(float64))
	// 从 Redis 缓存中获取数据
	a := userss["phone"]
	phoneCode, err := redisC.Get(ctx1, a).Result() //获取emailcode
	//fmt.Println(a)
	//
	if err == redis.Nil {

		fmt.Println("手机号key does not exist")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("查询到phoneCode,获取phoneCode:", phoneCode)
	}

	selfTool.Dbsql.Sql = "update users set password=? where id=?;"
	selfTool.Dbsql.Right = "修改用户密码成功"
	selfTool.Dbsql.Wrong = "修改用户密码失败"
	//判断redis中的emailcode和前端传来的emailcode
	if phoneCode == userss["code"] {
		if data1["code"].(float64) == 1 {

			a := selfTool.Dbsql.UPS(userss["password"], id)
			fmt.Println(a)
			c.JSON(200, selfTool.Dbsql.Right)
			return

		}
		c.JSON(200, "手机号不正确")
	}
	c.JSON(200, "手机号验证码不正确")

}
