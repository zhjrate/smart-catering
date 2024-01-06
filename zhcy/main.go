package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"zhcy/controller"

	classification2 "zhcy/controller/classification"
	"zhcy/controller/loginAndRegister"
	_select "zhcy/controller/select"
	"zhcy/controller/sendCode"
	"zhcy/controller/userController"
	tool2 "zhcy/tool/selfTool"
	"zhcy/util/dbSQL"
)

func main() {

	r := gin.Default()

	// 自定义跨域中间件配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"} // 允许的来源
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	// 添加自定义的登录token检查中间件
	//r.Use(middleware)

	toolGroup := r.Group("/toolGroup")
	//验证码发送
	toolGroup.POST("/sendEmailCode", sendCode.SendEmailcode) //发送邮箱验证码
	toolGroup.POST("/sendPhoneCode", sendCode.SendPhoneCode) //发送手机号验证码

	//登录注册路由组
	loginRegister := r.Group("/loginRegister")
	//用户名登录注册
	loginRegister.POST("/nameLogin", loginAndRegister.LoginHandler)       //用户名登录
	loginRegister.POST("/nameRegister", loginAndRegister.RegisterHandler) //用户名注册
	//邮箱登录注册
	loginRegister.POST("/emailLogin", loginAndRegister.EmailLoginHandler)       //邮箱登录
	loginRegister.POST("/emailRegister", loginAndRegister.EmailRegisterHandler) //邮箱注册
	//手机号登录
	loginRegister.POST("/phoneLogin", loginAndRegister.PhoneLogin) //手机号登录
	//手机号注册
	loginRegister.POST("/phoneRegister", loginAndRegister.PhoneRegister) //手机号注册

	//用户信息路由组
	updateGroup := r.Group("/updateGroup")
	//用户密码重设
	updateGroup.PUT("/resetPassword", userController.ResetPassword) //修改用户密码

	//商品信息路由组
	selectGroup := r.Group("/selectGroup")
	//查看所有的商品信息
	selectGroup.GET("/commdityAll", _select.CommdityAll) //获取所有商品信息
	//查看某个的商品信息
	selectGroup.POST("/commdityOne", _select.CommdityOne) //获取某个商品信息
	//查看某个的商品信息
	selectGroup.POST("/commdityAdd", _select.CommdityAdd) //添加某个商品信息
	//查看某个的商品信息
	selectGroup.PUT("/commdityReset", _select.CommdityReset) //修改某个商品信息
	//删除某个的商品信息
	selectGroup.DELETE("/commdityDel", _select.CommdityDel) //删除某个商品信息
	//上架下架某个的商品信息
	selectGroup.PUT("/commdityUpDown", _select.CommdityUpDown)
	//上架下架某个的商品信息
	selectGroup.PUT("/setComUpDown", _select.CommdityUpDown)
	//上架下架某个的商品信息(定时)
	selectGroup.PUT("/setDsComUpDown", _select.SetDsCommdityUpDown)
	//商品信息的图片管理
	selectGroup.POST("/sendPicture", _select.SendPicture)

	//商品分类管理路由组
	classification := r.Group("/classification")
	//商品分类信息添加
	classification.POST("/classificaAdd", classification2.ClassificaAdd)
	classification.PUT("/classificaAReset", classification2.ClassificaAReset)
	classification.DELETE("/classificaDel", classification2.ClassificaDel)

	//实时刷新所有的商品信息
	selectGroup.GET("/refreshCommodity", controller.RefreshCommodity) //获取所有商品信息

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	go dbSQL.Connect()
	go controller.UpdateDataPeriodically() // 启动一个goroutine用于定期更新数据

	// 设置定时任务
	go func(status string, id string) {
		ticker := time.NewTicker(10 * time.Second)
		fmt.Println(_select.S.Status, _select.S.Id)
		// 每隔10秒执行一次任务
		for {
			select {
			case <-ticker.C:
				tool2.Dbsql.Sql = "SELECT status from commodity where id = ?;"
				data := tool2.Dbsql.SelS(id)
				data1 := tool2.Toolkit.JsonSlice(data)
				data2, _ := tool2.Toolkit.JsonList(data1["list"])
				fmt.Println(data1["code"].(float64))
				fmt.Println(data2[0]["status"].(string))
				if data2[0]["status"].(string) == status {
					fmt.Println("商品已经被定时上架")
					return
				}
				if data1["code"].(float64) == 1 {
					tool2.Dbsql.Sql = "update commodity set status=? where id = ?;"
					final := tool2.Dbsql.UPS(status, id)
					final1 := tool2.Toolkit.JsonSlice(final)
					if final1["code"].(float64) == 1 {

						if data2[0]["status"].(string) == "1" {

							tool2.Dbsql.Right = "商品定时下架成功"
							tool2.Dbsql.Wrong = "商品定时下架失败"
							if final1["code"].(float64) == 1 {
								fmt.Println("商品定时下架成功")
								return
							}
							fmt.Println("商品定时下架失败")
							return
						}
						tool2.Dbsql.Right = "商品定时上架成功"
						tool2.Dbsql.Wrong = "商品定时上架失败"
						if final1["code"].(float64) == 1 {
							fmt.Println("商品定时上架成功")
							return
						}
						fmt.Println("商品定时上架失败")
						return
					}
				}
				//tool2.Dbsql.Sql = "SELECT status from commodity where id = ?;"
				//data := tool2.Dbsql.SelS(id)
				//data1 := tool2.Toolkit.JsonSlice(data)
				//data2, _ := tool2.Toolkit.JsonList(data1["list"])
				//fmt.Println(data1["code"].(float64))
				//fmt.Println(data2[0]["status"].(string))
				//if data2[0]["status"].(string) == status {
				//	//c.JSON(200, gin.H{"code": 0, "message": "商品已经被定时上架"})
				//	fmt.Println("商品已经被定时上架")
				//}
				//if data1["code"].(float64) == 1 {
				//	tool2.Dbsql.Sql = "update commodity set status=? where id = ?;"
				//	final := tool2.Dbsql.UPS(status, id)
				//	final1 := tool2.Toolkit.JsonSlice(final)
				//	if final1["code"].(float64) == 1 {
				//
				//		if data2[0]["status"].(string) == "1" {
				//
				//			tool2.Dbsql.Right = "商品定时下架成功"
				//			tool2.Dbsql.Wrong = "商品定时下架失败"
				//			if final1["code"].(float64) == 1 {
				//
				//				//c.JSON(200, gin.H{"code": 1, "message": tool2.Dbsql.Right})
				//				fmt.Println("商品定时下架成功")
				//			}
				//			//c.JSON(200, gin.H{"code": 0, "message": tool2.Dbsql.Wrong})
				//
				//		}
				//		tool2.Dbsql.Right = "商品定时上架成功"
				//		tool2.Dbsql.Wrong = "商品定时上架失败"
				//		if final1["code"].(float64) == 1 {
				//			fmt.Println("商品定时上架成功")
				//			//c.JSON(200, gin.H{"code": 1, "message": tool2.Dbsql.Right})
				//
				//		}
				//		fmt.Println("商品定时上架失败")
				//		//c.JSON(200, gin.H{"code": 0, "message": tool2.Dbsql.Wrong})
				//		//return true
				//	}
				//}
			}
		}
	}(_select.S.Status, _select.S.Id)

	r.Run(":8080")

}

func middleware(c *gin.Context) {
	req := c.Request

	username := req.Header.Get("username")
	exist := dbSQL.RedisHeaders(c, username)
	fmt.Println(username)
	//ctx1, redisC := dbSQL.RedisConnect()
	//token, err := redisC.Get(ctx1, username).Result()
	if exist {

		c.JSON(200, tool2.Dbsql.Ok())
		c.Next()
	}
	c.JSON(200, gin.H{"code": 1, "message": "Unauthorized - Login token expired"})
	c.Abort()
	return

}
