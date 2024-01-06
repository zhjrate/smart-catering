package _select

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	tool2 "zhcy/tool"
	"zhcy/tool/selfTool"
	"zhcy/util/dbSQL"
)

type ToolK struct {
	tool2.Toolkit
	tool2.MailBox
}

var Toolkit = ToolK{}

type SetDs struct {
	Status string
	Id     string
}

var S = SetDs{Status: "", Id: ""}

// 所有商品展示
func CommdityAll(c *gin.Context) {
	db := selfTool.GetMysqlConnect("select * from commodity;")
	data := db.Sel()
	data1 := Toolkit.JsonSlice(data)
	data2, _ := Toolkit.JsonList(data1["list"])
	if data1["code"].(float64) == 1 {

		strArray := selfTool.GetCommidity(data2)

		c.JSON(200, gin.H{"code": 1, "message": strArray})
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "数据库不存在或者 数据库未链接"})
	return
}

var way = dbSQL.Database{}

// 某个商品展示
func CommdityOne(c *gin.Context) {
	commdity := make(map[string]string)

	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	way.Sql = "select * from commodity where id=?;"
	way.Right = "成功"
	way.Wrong = "错误"
	data := way.SelS(commdity["id"])
	c.JSON(200, data)
	return
}

// 添加某个商品
func CommdityAdd(c *gin.Context) {
	db := selfTool.GetMysqlConnect("select * from commodity where name=?;")
	commdity := make(map[string]string)
	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	dataS := db.SelS(commdity["name"])
	data1 := Toolkit.JsonSlice(dataS)
	fmt.Println(data1["code"].(float64))

	if data1["code"].(float64) != 1 {
		fmt.Println("数据库无即将添加的商品信息，可以添加")
		db := selfTool.GetMysqlConnect("insert into commodity(name,description,price,startTime,num) values(?,?,?,?,?);")
		db.InsS(commdity["name"], commdity["description"], commdity["price"], commdity["startTime"], commdity["num"])
		c.JSON(200, gin.H{"code": 1, "message": "商品添加成功"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "该条商品数据已存在,不可再次添加"})
	return
}

// 修改某个商品
func CommdityReset(c *gin.Context) {

	commdity := make(map[string]string)
	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	//修改商品信息
	if commdity["num"] == "" || commdity["price"] == "" || commdity["description"] == "" || commdity["id"] == "" {
		c.JSON(200, way.LoseSet("漏选了什么东西"))
		return
	}
	way.Sql = "update commodity set num = ? ,price = ?,cxprice=?,description = ? where id = ?;"
	way.Right = "修改成功"
	way.Wrong = "修改失败"
	way.UPS(commdity["num"], commdity["price"], commdity["description"], commdity["id"])

	c.JSON(200, way.UPS(commdity["num"], commdity["price"], commdity["cxprice"], commdity["description"], commdity["id"]))
	return

}

// 删除某个商品
func CommdityDel(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from commodity where id = ?;")
	commdity := make(map[string]string)
	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	data := db.SelS(commdity["id"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	data2, _ := selfTool.Toolkit.JsonList(data1["list"])

	if len(data2) > 0 {
		numData := data2[0]["num"].(string)
		num, err := strconv.Atoi(numData)
		if err != nil {
			fmt.Println("转换失败:", err)
			return
		}
		if num <= 0 {
			db3 := selfTool.GetMysqlConnect("delete from commodity  where id=?;")
			db3.UPS(commdity["id"])
			//db := selfTool.GetMysqlConnect("update commodity set name=?,description=?,price=?,num=? where id=?;")
			c.JSON(200, gin.H{"code": 1, "message": "该商品库存量为0,管理员已删除该商品"})
			return
		}
		c.JSON(200, gin.H{"code": 0, "message": "该商品库存量不为0,管理员不能删除该商品"})
		return
	} else {
		c.JSON(200, gin.H{"code": 0, "message": "数据库查询的信息为空"})
		return
	}

}

// 上下架某个商品
func CommdityUpDown(c *gin.Context) {

	commdity := make(map[string]string)
	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	if commdity["status"] == "" || commdity["id"] == "" {
		c.JSON(200, way.LoseSet("缺少参数"))
		return
	}

	way.Sql = "SELECT status from commodity where id = ?;"
	data := way.SelS(commdity["id"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	data2, _ := selfTool.Toolkit.JsonList(data1["list"])

	fmt.Println(data1["code"].(float64))
	fmt.Println(data2[0]["status"].(string))
	if data1["code"].(float64) == 1 {
		way.Sql = "update commodity set status=? where id = ?;"
		if data2[0]["status"].(string) == "1" {
			way.Right = "商品下架成功"
			way.Wrong = "商品下架失败"
			way.UPS(commdity["status"], commdity["id"])
			c.JSON(200, gin.H{"code": 1, "message": way.Right})
			return
		}
		way.Right = "商品上架成功"
		way.Wrong = "商品上架失败"
		way.UPS(commdity["status"], commdity["id"])
		c.JSON(200, gin.H{"code": 1, "message": way.Right})
		return

	}

	return

}

var Exist bool = false

// 定时上下架某个商品
func SetDsCommdityUpDown(c *gin.Context) {

	commdity := make(map[string]string)
	if err := c.BindJSON(&commdity); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	if commdity["status"] == "" || commdity["id"] == "" {
		c.JSON(200, way.LoseSet("缺少参数"))
		return
	}
	way.Sql = "SELECT * from commodity where id = ?;"
	data := way.SelS(commdity["id"])
	data1 := Toolkit.JsonSlice(data)
	data2, _ := Toolkit.JsonList(data1["list"])
	if data1["code"].(float64) == 1 {
		if commdity["status"] == "1" && data2[0]["status"].(string) == commdity["status"] {
			fmt.Println("商品已经被定时上架")
			c.JSON(200, gin.H{"code": 1, "message": selfTool.Dbsql.LoseSet("商品已经被定时上架")})

		}
		if commdity["status"] == "0" && data2[0]["status"].(string) == commdity["status"] {
			fmt.Println("商品已经被定时下架")
			c.JSON(200, gin.H{"code": 1, "message": selfTool.Dbsql.LoseSet("商品已经被定时下架")})
		}
		if commdity["status"] == "0" && data2[0]["status"].(string) == "1" {
			way.Right = "商品定时下架成功"
			way.Wrong = "商品定时下架失败"
			final := way.UPS(commdity["status"], commdity["id"])
			final1 := selfTool.Toolkit.JsonSlice(final)
			fmt.Println("商品定时下架成功")
			c.JSON(200, gin.H{"code": 1, "message": final1})
			return
		}
		if commdity["status"] == "1" && data2[0]["status"].(string) == "0" {
			way.Right = "商品定时上架成功"
			way.Wrong = "商品定时上架失败"

			final := way.UPS(commdity["status"], commdity["id"])
			final1 := selfTool.Toolkit.JsonSlice(final)
			fmt.Println("商品定时上架成功")
			c.JSON(200, gin.H{"code": 1, "message": final1})
			return
		}

	} else {
		fmt.Println("商品不存在,不能被上下架")
		c.JSON(200, gin.H{"code": 1, "message": selfTool.Dbsql.LoseSet("商品不存在,不能被上下架")})
		return
	}

	// 设置定时任务
	go func() {
		ticker := time.NewTicker(1 * time.Second) // 每隔1秒执行一次任务
		for {
			select {
			case <-ticker.C:

				if data1["code"].(float64) == 1 {
					way.Sql = "update commodity set status=? where id = ?;"

					way.UPS(commdity["status"], commdity["id"])
					if commdity["status"] == "1" && data2[0]["status"].(string) == commdity["status"] {
						fmt.Println("商品已经被定时上架")

					}
					if commdity["status"] == "0" && data2[0]["status"].(string) == commdity["status"] {
						fmt.Println("商品已经被定时下架")
					}
					if commdity["status"] == "0" && data2[0]["status"].(string) == "1" {
						way.Right = "商品定时下架成功"
						way.Wrong = "商品定时下架失败"
						final := way.UPS(commdity["status"], commdity["id"])
						final1 := selfTool.Toolkit.JsonSlice(final)
						fmt.Println("商品定时下架成功")
						c.JSON(200, gin.H{"code": 1, "message": final1})
						return
					}
					if commdity["status"] == "1" && data2[0]["status"].(string) == "0" {
						way.Right = "商品定时上架成功"
						way.Wrong = "商品定时上架失败"

						final := way.UPS(commdity["status"], commdity["id"])
						final1 := selfTool.Toolkit.JsonSlice(final)
						fmt.Println("商品定时上架成功")
						c.JSON(200, gin.H{"code": 1, "message": final1})
						return
					}
				} else {
					fmt.Println("商品不存在,不能被上下架")
				}

				//c.JSON(200, gin.H{"code": 1, "message": "商品不存在,不能被上下架"})
				//return
			}

		}
	}()
}

// 上传图片
func SendPicture(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存图片到服务器指定目录
	uploadPath := "uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	dstPath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})

}
