package sendCode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	tool2 "zhcy/tool/selfTool"
	"zhcy/util/dbSQL"
)

// 发送邮箱验证码
func SendEmailcode1(c *gin.Context) {
	userss := make(map[string]string)
	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	fmt.Println(userss["email"])
	tool2.Toolkit.Connect(1)
	emailcode := tool2.GenerateTokenCode()
	fmt.Println(emailcode)
	content := tool2.Toolkit.Template(1, userss["username"], "邮箱验证", emailcode, "系统管理员")
	//fmt.Println(content)
	Passage := make(chan bool)
	fmt.Println(userss["email"])
	go tool2.Toolkit.Send("创泓度网络", content, userss["email"], Passage)

	ctx1, _ := dbSQL.RedisConnect()
	dbSQL.ConnRedis.Set(ctx1, userss["email"], emailcode, time.Duration(240)*time.Minute)
	c.JSON(200, gin.H{"message": "Verification code sent"})

}

// 发送邮箱验证码1
func SendEmailcode(c *gin.Context) {
	db := tool2.GetMysqlConnect("select * from users where email like ?;")
	userss := make(map[string]string)
	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(userss["email"])
	data1 := tool2.Toolkit.JsonSlice(data)
	data2, err := tool2.Toolkit.JsonList(data1["list"])
	fmt.Println(data2[0]["name"])
	fmt.Println(data1["code"])
	if err != nil {
		fmt.Println("error:", err.Error())
		c.JSON(200, gin.H{"code": 0, "message": tool2.Dbsql.Wrong})
		return
	}

	//判断
	if data1["code"].(float64) == 1 {

		tool2.Toolkit.Connect(1)
		emailcode := tool2.GenerateTokenCode()
		fmt.Println(emailcode)
		content := tool2.Toolkit.Template(1, data2[0]["name"].(string), "邮箱验证", emailcode, "系统管理员")
		//fmt.Println(content)
		Passage := make(chan bool)
		fmt.Println(userss["email"])
		go tool2.Toolkit.Send("创泓度网络", content, userss["email"], Passage)

		ctx1, _ := dbSQL.RedisConnect()
		dbSQL.ConnRedis.Set(ctx1, userss["email"], emailcode, time.Duration(240)*time.Minute)
		c.JSON(200, gin.H{"message": "邮箱验证码已发送"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "邮箱不存在，请先注册"})
}
