package sendCode

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"zhcy/tool"
	"zhcy/tool/selfTool"
	"zhcy/util/dbSQL"
)

func SendPhoneCode(c *gin.Context) {

	type Data struct {
		Phone    string `json:"phone"`
		PhoneKey string `json:"phoneKey"` // 用于验证
		Key      string `json:"key"`      //随机数
	}
	var data Data
	if err := c.BindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	//err := tool.PhoneCode.SendCode()
	phoneSuiji := selfTool.GenerateTokenCode()
	key := data.Key
	//MD5加密
	hash := md5.New()
	hash.Write([]byte(data.Phone + key))
	hashInBytes := hash.Sum(nil)
	decryptedPhone := hex.EncodeToString(hashInBytes)
	fmt.Println("手机号:", data.Phone)
	fmt.Println("md5加密后手机号后的:", decryptedPhone)

	if (hex.EncodeToString(hashInBytes)) == data.PhoneKey {

		phoneCode := &tool.PhoneCode{}
		phoneCode.Init()
		phoneCode.SendCode(data.Phone, phoneSuiji)
		ctx1, _ := dbSQL.RedisConnect()
		dbSQL.ConnRedis.Set(ctx1, data.Phone, phoneSuiji, time.Duration(240)*time.Minute) //存入redis缓存
		c.JSON(200, "手机验证码发送成功,验证码为: "+phoneSuiji)
		return
		fmt.Println("手机验证码发送成功,Response:", phoneSuiji)
	}
	c.JSON(200, "md5单向加密比对不一致")
	return
}
