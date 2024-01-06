package loginAndRegister

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
	"zhcy/tool/selfTool"

	"zhcy/util/dbSQL"
)

//type ToolK struct {
//	tool2.Toolkit
//	tool2.MailBox
//}
//
//var Dbsql dbSQL.Database
//
//var Toolkit = ToolK{}

// 用户名登录
func LoginHandler(c *gin.Context) {
	db := selfTool.GetMysqlConnect("select * from users where name like ?;")

	//用户登录
	userss := make(map[string]string)

	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(userss["username"])
	fmt.Println(userss["username"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	fmt.Println(data1["code"].(float64))
	data2, _ := selfTool.Toolkit.JsonList(data1["list"])
	//fmt.Println(data2[0]["password"])
	//判断用户是否存在以及密码是否正确

	if data1["code"].(float64) == 1 {
		if userss["password"] == data2[0]["password"] {
			token := selfTool.GetRedis(c)
			fmt.Println("获取redis token:", token)
			c.JSON(200, selfTool.Dbsql.Ok())
			return
		}
		c.JSON(200, selfTool.Dbsql.LosePassword())
		return
	}

	c.JSON(200, selfTool.Dbsql.LoseUsername())
}

// 用户名注册
func RegisterHandler(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from users where name like ?;")
	userss := make(map[string]string)

	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(userss["username"])
	data1 := selfTool.Toolkit.JsonSlice(data)

	if data1["code"].(float64) != 2 {
		if userss["username"] != "" && userss["password"] != "" {

			rand.Seed(time.Now().UnixNano())
			// 生成随机汉字，范围：0x4e00 - 0x9fff
			randomChinese := rune(rand.Intn(0x9fff-0x4e00+1) + 0x4e00)
			// 生成随机字母
			randomLetter := rune(rand.Intn(26) + 'A')
			combination := string([]rune{randomChinese, randomLetter})

			selfTool.Dbsql.Sql = "insert into  users(name,password,age,sex,phone,email) values (?,?,?,?,?,?);"
			selfTool.Dbsql.Right = "加载用户信息成功"
			selfTool.Dbsql.Wrong = "加载用户信息失败"
			selfTool.Dbsql.InsS(combination, randomLetter, "", "", "", "")

			c.JSON(200, selfTool.Dbsql.Ok())
			return
		}

		c.JSON(200, "username和password为空")
		return
	}
	c.JSON(200, "用户已存在")

}

//邮箱登录

func EmailLoginHandler(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from users where email like ?;")
	userss := make(map[string]string)

	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	data := db.SelS(userss["email"])
	fmt.Println(userss["email"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	fmt.Println(data1["code"].(float64))
	data2, err := selfTool.Toolkit.JsonList(data1["list"])
	if err != nil {
		fmt.Println("error:", err.Error())
		c.JSON(200, selfTool.Dbsql.Wrong)
		return
	}

	//判断
	if data1["code"].(float64) == 1 {
		if userss["password"] == data2[0]["password"] {
			c.JSON(200, selfTool.Dbsql.Ok())
			return
		}

		c.JSON(200, "密码错误存在")
		return
	}
	c.JSON(200, "邮箱未注册")

}

//邮箱注册

func EmailRegisterHandler(c *gin.Context) {
	db := selfTool.GetMysqlConnect("select * from users where email like ?;")
	userss := make(map[string]string)

	if err := c.BindJSON(&userss); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	data := db.SelS(userss["email"])
	data1 := selfTool.Toolkit.JsonSlice(data)

	ctx1, redisC := dbSQL.RedisConnect()
	fmt.Println(data1["code"].(float64))
	// 从 Redis 缓存中获取数据
	a := userss["email"]
	emailcode, err := redisC.Get(ctx1, a).Result() //获取emailcode
	//fmt.Println(a)
	//
	if err == redis.Nil {

		fmt.Println("邮箱key does not exist")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("查询到emailcode,获取emailcode:", emailcode)
	}
	//判断redis中的emailcode和前端传来的emailcode
	if emailcode == userss["code"] {

		if data1["code"].(float64) != 1 {

			selfTool.Dbsql.Sql = "insert into  users(name,password,age,sex,phone,email) values (?,?,?,?,?,?);"
			selfTool.Dbsql.Right = "加载用户信息成功"
			selfTool.Dbsql.Wrong = "加载用户信息失败"
			selfTool.Dbsql.InsS(userss["username"], userss["password"], userss["age"], userss["sex"], userss["phone"], userss["email"])
			c.JSON(200, selfTool.Dbsql.Ok())
			return

		}

		c.JSON(200, "邮箱存在")
		return

	}
	c.JSON(200, "email验证码不正确")

}

// 手机号登录
func PhoneLogin(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from users where phone like ?;")
	user := make(map[string]string)
	if err := c.BindJSON(&user); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(user["phone"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	fmt.Println(data1["code"].(float64))
	//data2, err := toolkit.JsonList(data1["list"])
	//if err != nil {
	//	fmt.Println("error:", err.Error())
	//	c.JSON(200, dbsql.Wrong)
	//	return
	//}
	ctx, redisC := dbSQL.RedisConnect()
	fmt.Println(data1["code"].(float64))

	// 从 Redis 缓存中获取数据
	phonecode, _ := redisC.Get(ctx, user["phone"]).Result()
	//判断
	if data1["code"].(float64) == 1 {
		if user["code"] == phonecode {
			c.JSON(200, selfTool.Dbsql.Ok())
			return
		}

		c.JSON(200, "验证码错误")
		return
	}
	c.JSON(200, "手机号未注册")

}

// 手机号注册
func PhoneRegister(c *gin.Context) {
	db := selfTool.GetMysqlConnect("select * from users where phone like ?;")
	users := make(map[string]string)
	if err := c.BindJSON(&users); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	data := db.SelS(users["phone"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	//data2, err := toolkit.JsonList(data1)
	ctx, redisC := dbSQL.RedisConnect()
	fmt.Println(data1["code"].(float64))

	// 从 Redis 缓存中获取数据
	phonecode, err := redisC.Get(ctx, users["phone"]).Result() //获取phonecode
	if err == redis.Nil {

		fmt.Println("手机号key does not exist")
	} else if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("查询到phonecode,获取phonecode:", phonecode)
	}
	//判断redis中的emailcode和前端传来的emailcode
	if phonecode == users["code"] {

		if data1["code"].(float64) != 1 {

			selfTool.Dbsql.Sql = "insert into  users(name,password,age,sex,phone,email) values (?,?,?,?,?,?);"
			selfTool.Dbsql.Right = "加载用户信息成功"
			selfTool.Dbsql.Wrong = "加载用户信息失败"
			selfTool.Dbsql.InsS(users["username"], users["password"], users["age"], users["sex"], users["phone"], users["email"])
			c.JSON(200, selfTool.Dbsql.Ok())
			return

		}

		c.JSON(200, "手机号存在")
		return

	}
	c.JSON(200, "手机号验证码不正确")

}
