package dbSQL

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mssola/user_agent"
	"github.com/redis/go-redis/v9"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	tool2 "zhcy/tool"
)

type LimitS struct {
	Limit string
	Max   int
}

type Code struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

type CodeToken struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Cid     string `json:"cid"`
	Title   string `json:"title"`
}

type Code_list struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	List    interface{} `json:"list""`
	Key     string      `json:"key"`
}

type Database struct {
	Sql   string
	Right string
	Wrong string
}

type DataBool struct {
	Sql       string
	Right     string
	Wrong     string
	User      string
	PassWorld string
}

type SQL_Array struct {
	Sql string
}

var DB *sql.DB

var key bool = true

var TookRsa = &tool2.ToolRsa{}
var Toolkit = tool2.Toolkit{}

func init() {
	dbName := "zhihui:1122334455@tcp(49.233.1.189:3030)/zhihui"

	for {
		_, ery := net.LookupHost("www.baidu.com")
		if ery == nil {
			db, err := sql.Open("mysql", dbName)
			if err != nil {
				return
			}

			DB = db
			RasConnect()
			break
		}
		time.Sleep(1 * time.Second)
	}

}
func Connect() error {
	dbName := "zhihui:1122334455@tcp(49.233.1.189:3030)/zhihui"
	var (
		ess error
	)
	for {
		_, ery := net.LookupHost("www.baidu.com")
		if ery == nil {
			db, err := sql.Open("mysql", dbName)
			if err != nil {
				return ess
			}
			ess = err
			DB = db
			RasConnect()
			break
		}
		time.Sleep(1 * time.Second)
	}

	return ess

}

func RasConnect() {
	sql := Database{}
	sql.Sql = "SELECT privatekey,publickey,whether FROM certificate where id = 1;"
	sql.Right = ""
	sql.Wrong = ""
	retu := sql.SelA()
	suessec, err := Toolkit.JsonList(Toolkit.JsonSlice(retu)["list"])
	if err != nil {
		fmt.Println(err)
	} else {
		if suessec[0]["privatekey"].(string) != "" {
			private := TookRsa.RecoverRsaPrivate([]byte(suessec[0]["privatekey"].(string)))
			if private == nil || private.N == nil || private.E == 0 {
				fmt.Println("Rsa PrivateKey:", "error")
			} else {
				fmt.Println("Rsa PrivateKey:", "Acquisition and binding conversion succeeded")
				TookRsa.FileCreateKey([]byte(suessec[0]["privatekey"].(string)), "file/rsa/privatekey.pem")
			}

		}
		if suessec[0]["publickey"].(string) != "" {
			public := TookRsa.RecoverRsaPublic([]byte(suessec[0]["publickey"].(string)))
			if public == nil || public.N == nil || public.E == 0 {
				fmt.Println("Rsa PublicKey:", "error")
			} else {
				fmt.Println("Rsa PublicKey:", "Acquisition and binding conversion succeeded")
				TookRsa.FileCreateKey([]byte(suessec[0]["publickey"].(string)), "file/rsa/public.pem")
			}
		}

		if suessec[0]["whether"] != nil {
			if suessec[0]["whether"].(string) == "0" {
				TookRsa.Code = false
				fmt.Println("Symmetric Data Signature Status:", "close")
			} else if suessec[0]["whether"].(string) == "1" {
				TookRsa.Code = true
				fmt.Println("Symmetric Data Signature Status:", "enable")
			}
		} else {
			fmt.Println("Symmetric Data Signature Status:", "Uninitialized!!!")
		}
	}
}

var CTX context.Context
var ConnRedis *redis.Client

//var ConnRedisCtx context.Context

func RedisConnect() (context.Context, *redis.Client) {
	CTX = context.Background()
	ConnRedis = redis.NewClient(
		&redis.Options{
			//38.55.188.153
			Addr:     "49.233.1.189:6379",
			Password: "zhao132909",
			DB:       0,
		})
	i := ConnRedis.Ping(CTX)
	fmt.Println("Redis: ", i)
	return CTX, ConnRedis

}

func CookieToken(c *gin.Context) (bool, map[string]string) {

	list := make(map[string]string)
	header := c.Request.Header
	name := header.Get("username")
	token := header.Get("token")
	fmt.Println("username:", name)
	//fmt.Println("token:", token)

	//username, err := c.Request.Cookie(name)
	if name != "" {
		return false, nil
	}
	user, _ := url.QueryUnescape(name)
	fmt.Println("user:", user)

	//token, err := c.Request.Cookie(header["token"][0])
	if token != "" {
		return false, nil
	}
	tokens, _ := url.QueryUnescape(token)
	fmt.Println("tokens:", tokens)
	list = map[string]string{
		"UserName": user,
		"Token":    tokens,
	}
	Title, _ := ConnRedis.Get(CTX, user).Result()
	if Title == "" {
		//c.JSON(200, gin.H{
		//	"code":    3,
		//	"message": "胆子这么大？",
		//})
		return false, list
	} else {
		//TTL = 查询过期时间 以秒为单位
		if tokens == Title {
			val := ConnRedis.TTL(CTX, user).Val()
			time1 := float64(val) / float64(time.Minute)
			var time2 float64 = 50
			if time1 < time2 {
				_, err := ConnRedis.Expire(CTX, user, 6*time.Hour).Result()
				if err != nil {
					return false, list
				} else {
					return true, list
				}
			} else {
				return true, list
			}
		} else {
			//c.JSON(200, gin.H{
			//	"code":    3,
			//	"message": "令箭超时 请重新登陆",
			//})
			return false, list
		}

	}
	return false, nil

}

func RedisHeaders(c *gin.Context, name string) bool {
	header := c.Request.Header
	if len(header["Token"]) == 0 {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "Authentication restriction",
		})
		return false
	} else {
		Title, _ := ConnRedis.Get(CTX, name).Result()
		if Title == "" {
			c.JSON(200, gin.H{
				"code":    0,
				"message": "Authentication restriction",
			})
			return false
		} else {
			//TTL = 查询过期时间 以秒为单位
			if header["Token"][0] == Title {
				val := ConnRedis.TTL(CTX, name).Val()
				time1 := time.Now().Add(1*time.Second).Unix() + int64(val)
				var time2 int64 = time.Now().Add(1 * time.Second).Unix()
				if time1 < time2 { //说明token已经过期
					_, err := ConnRedis.Expire(CTX, name, 240*time.Minute).Result()
					if err != nil {
						c.JSON(200, gin.H{
							"code":    0,
							"message": "重新定义token出错:" + err.Error(),
						})
						return false
					} else {
						c.JSON(200, gin.H{
							"code":    1,
							"message": "token已经重置",
						})
						return true
					}
				} else {
					c.JSON(200, gin.H{
						"code":    1,
						"message": "token未过期",
					})
					return true
				}
			} else {
				c.JSON(200, gin.H{
					"code":    0,
					"message": "前端的Token与储存的不一致",
				})
				return false
			}

		}
	}
}

func RedisAdd(key string, value string, time time.Duration) (bool, string) {
	title, _ := ConnRedis.Get(CTX, key).Result()
	if title != "" {
		return true, title
	} else {
		err := ConnRedis.Set(CTX, key, value, time).Err()
		if err != nil {
			return false, ""
		} else {
			return true, value
		}
	}
}

func RequestHeader(header http.Header) {

}

func (sql Database) SelBool() bool {
	rows, err := DB.Query(sql.Sql)
	defer rows.Close()
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	} else {
		return false
	}

}

func (sql Database) SelBoolTrue(args ...any) bool {
	stmp, err := DB.Prepare(sql.Sql)
	if err != nil {
		//code := Code{
		//	Code:    0,
		//	Message: err.Error(),
		//}
		fmt.Println(err.Error())
		return false
	}
	defer stmp.Close()

	rows, err := stmp.Query(args...)
	if err != nil {
		//code := Code{
		//	Code:    0,
		//	Message: sql.Wrong,
		//}
		return false
	}
	defer rows.Close()
	if rows.Next() {
		return true
	} else {
		return false
	}

}

func (sql Database) SelBoolFalse(args ...any) (bool, error) {
	stmp, err := DB.Prepare(sql.Sql)
	if err != nil {
		//code := Code{
		//	Code:    0,
		//	Message: err.Error(),
		//}
		//fmt.Println(err.Error())
		return false, err
	}
	defer stmp.Close()

	rows, err := stmp.Query(args...)
	if err != nil {
		//code := Code{
		//	Code:    0,
		//	Message: sql.Wrong,
		//}
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return false, nil
	} else {
		return true, err
	}

}

func (sql Database) SelSRow(args ...any) (int, []string, error) {
	var retu []string
	stmp, err := DB.Prepare(sql.Sql)
	if err != nil {
		return 0, retu, err
	}
	defer stmp.Close()
	rows, err := stmp.Query(args...)
	if err != nil {
		return 0, retu, err
	}
	defer rows.Close()
	for rows.Next() {
		var number string
		rows.Scan(&number)
		retu = append(retu, number)
	}
	return len(retu), retu, nil
}

func (sql Database) UPS(args ...interface{}) interface{} {
	stmp, err := DB.Prepare(sql.Sql)
	if err != nil {
		code := Code{
			Code:    0,
			Message: err.Error(),
		}
		return code
	}

	defer stmp.Close()

	_, err = stmp.Exec(args...)
	if err != nil {
		code := Code{
			Code:    0,
			Message: err.Error(),
		}
		return code
	}
	code := Code{
		Code:    1,
		Message: sql.Right,
	}
	if TookRsa.Code {
		jsonData, _ := json.Marshal(code)
		key, err := TookRsa.SignRSAdata(jsonData)
		if err != nil {
			fmt.Println(err.Error())
		}
		code.Key = key
	}
	return code
}

func (sql Database) SelS(args ...interface{}) interface{} {
	stmt, err := DB.Prepare(sql.Sql)
	if err != nil {
		code := Code{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		code := Code{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	list := make([]map[string]string, 0)

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			// handle error
			continue
		}

		record := make(map[string]string)
		for i, v := range values {
			if v != nil {
				switch val := v.(type) {
				case []byte:
					record[columns[i]] = string(val)
				default:
					record[columns[i]] = fmt.Sprintf("%v", val)
				}
			}
		}

		list = append(list, record)
	}

	if len(list) == 0 {
		code := Code{
			Code:    0,
			Message: "暂无更多数据",
		}

		return code
	} else {
		code := Code_list{
			Code:    1,
			Message: sql.Right,
			List:    list,
		}
		if TookRsa.Code {
			jsonData, _ := json.Marshal(code)
			key, err := TookRsa.SignRSAdata(jsonData)
			if err != nil {
				fmt.Println(err.Error())
			}
			code.Key = key
		}
		return code
	}

}

func (sql Database) SelW(updatedData string) interface{} {

	a := DB.QueryRow("SELECT * FROM commodity WHERE some_condition LIMIT 1").Scan(&updatedData)
	return a
}
func (sql Database) Sel() interface{} {
	rows, err := DB.Query(sql.Sql)

	defer rows.Close()
	if err != nil {
		code := Code_list{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	} else {
		columns, _ := rows.Columns()                  // 获取列
		scanArgs := make([]interface{}, len(columns)) // 填充数据
		values := make([]interface{}, len(columns))   // 存放数据值
		// 将values的指针存放到scanArgs里
		for i := range values {
			scanArgs[i] = &values[i]
		}
		// 用来装整个查询结果的map
		list := make([]any, 0)
		for rows.Next() {
			// Scan方法赋值
			rows.Scan(scanArgs...)
			// 因为指针绑定，所以values里有值了，遍历赋值给record
			record := make(map[string]string)
			for i, v := range values {
				if v != nil {
					record[columns[i]] = string(v.([]byte)) // 需要进行断言，然后string方法转化为字符串
				}
			}
			list = append(list, record)
		}
		if len(list) == 0 {
			code := Code_list{
				Code:    1,
				Message: "暂无更多数据",
			}
			return code
		} else {
			code := Code_list{
				Code:    1,
				Message: sql.Right,
				List:    list,
			}
			if TookRsa.Code {
				jsonData, _ := json.Marshal(code)
				key, err := TookRsa.SignRSAdata(jsonData)
				if err != nil {
					fmt.Println(err.Error())
				}
				code.Key = key
			}
			return code
		}

	}

}

func (sql Database) SelA() interface{} {
	rows, err := DB.Query(sql.Sql)
	defer rows.Close()
	if err != nil {
		code := Code_list{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	} else {
		columns, _ := rows.Columns()                  // 获取列
		scanArgs := make([]interface{}, len(columns)) // 填充数据
		values := make([]interface{}, len(columns))   // 存放数据值
		// 将values的指针存放到scanArgs里
		for i := range values {
			scanArgs[i] = &values[i]
		}
		// 用来装整个查询结果的map
		list := make([]any, 0)
		for rows.Next() {
			// Scan方法赋值
			rows.Scan(scanArgs...)
			// 因为指针绑定，所以values里有值了，遍历赋值给record
			record := make(map[string]string)
			for i, v := range values {
				if v != nil {
					record[columns[i]] = string(v.([]byte)) // 需要进行断言，然后string方法转化为字符串
				}
			}
			list = append(list, record)
		}
		if len(list) == 0 {
			code := Code_list{
				Code:    2,
				Message: "暂无更多数据",
			}
			return code
		} else {
			code := Code_list{
				Code:    1,
				Message: sql.Right,
				List:    list,
			}
			return code
		}

	}

}

func (sql Database) Limit(s int, count int) LimitS {
	res, err := DB.Query(sql.Sql)
	retu := LimitS{
		Limit: "",
		Max:   0,
	}
	if err != nil {
		return retu
	}
	defer res.Close()
	var ins string
	if res.Next() {
		res.Scan(&ins)
	}
	iis, _ := strconv.Atoi(ins)
	ress := iis / s
	y := iis % s
	if y > 0 {
		ress++
	}
	unt := (count - 1) * s
	limit := "LIMIT " + strconv.Itoa(unt) + "," + strconv.Itoa(s)
	retu.Limit = limit
	retu.Max = ress
	return retu
}
func (sql Database) LoseUsername() interface{} {
	code := Code{
		Code:    0,
		Message: "此用户不存在",
	}
	return code
}
func (sql Database) LosePassword() interface{} {
	code := Code{
		Code:    0,
		Message: "密码错误",
	}
	return code
}
func (sql Database) LoseToken() interface{} {
	code := Code{
		Code:    0,
		Message: "token已过期",
	}
	return code
}
func (sql Database) OkUser() interface{} {
	code := Code{
		Code:    1,
		Message: "用户已存在",
	}
	return code
}
func (sql Database) Lose() interface{} {
	code := Code{
		Code:    0,
		Message: "失败",
	}
	return code
}
func (sql Database) Ok() interface{} {
	code := Code{
		Code:    1,
		Message: "成功",
	}
	return code
}
func (sql Database) LoseSet(message string) interface{} {
	code := Code{
		Code:    0,
		Message: message,
	}
	return code
}

func (sql Database) InsS(args ...any) interface{} {
	stmt, err := DB.Prepare(sql.Sql)
	if err != nil {
		fmt.Println(err.Error())
		code := Code{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println(err.Error())
		code := Code{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	}
	defer rows.Close()
	code := Code{
		Code:    1,
		Message: sql.Right,
	}
	if TookRsa.Code {
		jsonData, _ := json.Marshal(code)
		key, err := TookRsa.SignRSAdata(jsonData)
		if err != nil {
			fmt.Println(err.Error())
		}
		code.Key = key
	}
	return code
}

func (sql Database) Ins() interface{} {
	_, err := DB.Query(sql.Sql)
	if err != nil {
		code := Code{
			Code:    0,
			Message: sql.Wrong,
		}
		return code
	} else {
		code := Code{
			Code:    1,
			Message: sql.Right,
		}
		if TookRsa.Code {
			jsonData, _ := json.Marshal(code)
			key, err := TookRsa.SignRSAdata(jsonData)
			if err != nil {
				fmt.Println(err.Error())
			}
			code.Key = key
		}
		return code
	}
}

func GetClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// The format of X-Forwarded-For can be: client, proxy1, proxy2
		ips := strings.Split(forwardedFor, ", ")
		return ips[0]
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func (sql DataBool) QueryBool(c *gin.Context, names string) interface{} {
	root, err := DB.Query(sql.Sql, sql.User, sql.PassWorld)
	defer root.Close()
	if err != nil {
		code := Code{
			Code:    0,
			Message: "系统错误",
		}
		return code
	} else {

		if root.Next() {
			var name string
			var cid string
			var title string
			root.Scan(&name, &cid, &title)
			code := CodeToken{
				Code:    1,
				Message: sql.Right,
				Token:   "123",
				Cid:     cid,
			}

			//uaparsers := user_agent.New(c.Request.UserAgent())
			//oss := uaparsers.OS()
			//browser, Version := uaparsers.Browser()
			//DataD := Database{
			//	Sql:   "INSERT INTO db.fqc_login (name,ip,system,browser,version) VALUES (?,?,?,?,?);",
			//	Right: "",
			//	Wrong: "",
			//}
			//Request := make(map[string]string)
			//c.BindJSON(&Request)
			//var username string = Request["username"]
			//parts := strings.Split(c.Request.RemoteAddr, ":")
			//
			//DataD.UPS(names, parts[0], oss, browser, Version)
			//
			//return code

			booll, Token := RedisAdd(name, tool2.RandAuto(10, 150), 8*time.Hour)
			if booll {
				code = CodeToken{
					Code:    1,
					Message: sql.Right,
					Token:   Token,
					Cid:     cid,
					Title:   title,
				}

				uaparsers := user_agent.New(c.Request.UserAgent())
				oss := uaparsers.OS()
				browser, Version := uaparsers.Browser()
				DataD := Database{
					Sql:   "INSERT INTO db.fqc_login (name,ip,system,browser,version) VALUES (?,?,?,?,?);",
					Right: "",
					Wrong: "",
				}
				//Request := make(map[string]string)
				//c.BindJSON(&Request)
				//var username string = Request["username"]
				//parts := strings.Split(c.Request.RemoteAddr, ":")

				//DataD.UPS(names, parts[0], oss, browser, Version)
				DataD.UPS(names, GetClientIP(c.Request), oss, browser, Version)
				return code
			} else {
				code := Code{
					Code:    2,
					Message: "Redis Error",
				}
				return code
			}
		} else {
			code := Code{
				Code:    2,
				Message: sql.Wrong,
			}
			return code
		}
	}
}

func (sql SQL_Array) Sel_List() (int, interface{}) {
	root, err := DB.Query(sql.Sql)

	defer root.Close()
	if err != nil {
		return 0, 1
	} else {
		columns, _ := root.Columns()                  // 获取列
		scanArgs := make([]interface{}, len(columns)) // 填充数据
		values := make([]interface{}, len(columns))   // 存放数据值
		// 将values的指针存放到scanArgs里
		for i := range values {
			scanArgs[i] = &values[i]
		}
		// 用来装整个查询结果的map
		list := make([]any, 0)
		for root.Next() {
			// Scan方法赋值
			root.Scan(scanArgs...)
			// 因为指针绑定，所以values里有值了，遍历赋值给record
			record := make(map[string]string)
			for i, v := range values {
				if v != nil {
					record[columns[i]] = string(v.([]byte)) // 需要进行断言，然后string方法转化为字符串
				}
			}
			list = append(list, record)
		}
		return 1, list
	}
}
