package classification

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"zhcy/tool/selfTool"
)

var classification map[string]string

// 添加某个商品分类
func ClassificaAdd(c *gin.Context) {

	if err := c.BindJSON(&classification); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}
	selfTool.Dbsql.Right = "商品分类添加成功"
	selfTool.Dbsql.Wrong = "商品分类添加错误"
	selfTool.Dbsql.Sql = "insert into classification(dishes,taste) values(?,?);"
	data := selfTool.Dbsql.InsS(classification["dishes"], classification["taste"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	if data1["code"].(float64) != 1 {
		c.JSON(200, gin.H{"code": 1, "message": "商品分类添加成功"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "商品分类添加失败"})
	return
}

// 修改某个商品
func ClassificaAReset(c *gin.Context) {

	if err := c.BindJSON(&classification); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(200, gin.H{"error": "Invalid JSON"})
		return
	}

	//修改商品分类信息
	if classification["dishes"] == "" || classification["taste"] == "" {
		c.JSON(200, selfTool.Dbsql.LoseSet("漏选了什么东西"))
		return
	}
	selfTool.Dbsql.Sql = "select dishes from classification  where id = ?;"
	selfTool.Dbsql.Right = "商品分类信息成功"
	selfTool.Dbsql.Wrong = "商品分类信息失败"
	data := selfTool.Dbsql.SelS(classification["dishes"], classification["id"])
	data1 := selfTool.Toolkit.JsonSlice(data)
	if data1["code"].(float64) == 1 {
		selfTool.Dbsql.Sql = "update classification set dishes = ? where id = ?;"
		selfTool.Dbsql.Right = "商品分类信息成功"
		selfTool.Dbsql.Wrong = "商品分类信息失败"
		selfTool.Dbsql.UPS(classification["dishes"], classification["id"])
		c.JSON(200, gin.H{"code": 1, "message": selfTool.Dbsql.Right})
		return
	}
	c.JSON(200, selfTool.Dbsql.LoseSet("商品分类添加失败"))
	return

}

// 删除某个商品分类信息
func ClassificaDel(c *gin.Context) {

	db := selfTool.GetMysqlConnect("select * from classification where id = ?;")
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

		db3 := selfTool.GetMysqlConnect("delete from classification  where id=?;")
		db3.UPS(commdity["id"])
		//db := selfTool.GetMysqlConnect("update commodity set name=?,description=?,price=?,num=? where id=?;")
		c.JSON(200, gin.H{"code": 1, "message": "校验成功,管理员已删除该商品分类"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "message": "该商品分类为空,管理员已删除该商品分类"})
	return

}
