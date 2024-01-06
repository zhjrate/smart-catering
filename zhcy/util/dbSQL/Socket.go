package dbSQL

import (
	json2 "encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var Socket string

func SocketServer() {
	Listen, err := net.Listen("tcp", ":82")
	if err != nil {
		fmt.Println(err)
		Socket = ":82通讯建立失败"
		os.Exit(1)
	}
	defer Listen.Close()
	Socket = ":82通讯已建立"
	for {
		conn, erc := Listen.Accept()
		if erc != nil {
			fmt.Println(erc)
			continue
		}
		go HandleSocket(conn)

	}
}

//fmt.Println("等待 2 秒钟...")
//<-time.After(2 * time.Second)
//fmt.Println("2 秒钟已经过去了")

//func HandleSocket(conn net.Conn) {
//	defer conn.Close()
//	fmt.Println("设备上线：", conn.RemoteAddr())
//	for {
//		buf := make([]byte, 512)
//		_, err := conn.Read(buf)
//		if err != nil {
//			fmt.Println("设备下线: ", conn.RemoteAddr())
//			go ConnDel(&ConnList, &UserList, conn)
//			//for i := 0; i < 3; i++ {
//			//	<-time.After(2 * time.Second)
//			//	newConn, err := net.Dial("tcp", conn.RemoteAddr().String())
//			//	if err != nil {
//			//		fmt.Println("重新连接客户机失败：", err)
//			//	} else {
//			//		fmt.Println("重新连接客户机成功：", newConn.RemoteAddr())
//			//		conn = newConn
//			//		break
//			//	}
//			//}
//			return
//		}
//		DataSocket(buf, conn)
//	}
//}

func HandleSocket(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client Connect Open：", conn.RemoteAddr())
	buf := make([]byte, 512)
	for {
		//err := conn.SetReadDeadline(time.Now().Add(6 * time.Second))
		//if err != nil {
		//	fmt.Println("设置缓冲区超时时间出错:", err)
		//	return
		//}
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//fmt.Println("Client Connect Close: ", conn.RemoteAddr())
				//go ConnDel(&ConnList, &UserList, conn)
				//for i := 0; i < 3; i++ {
				//	<-time.After(2 * time.Second)
				//	newConn, err := net.Dial("tcp", conn.RemoteAddr().String())
				//	if err != nil {
				//		fmt.Println("重新连接客户机失败：", err)
				//	} else {
				//		fmt.Println("重新连接客户机成功：", newConn.RemoteAddr())
				//		conn = newConn
				//		break
				//	}
				//}
				return
			}
			//fmt.Println("该连接已超时:", err)
			fmt.Println("Client Connect Close: ", conn.RemoteAddr())
			go ConnDel(&ConnList, &UserList, conn)
			return
		}
		//buf[:n] = 读取缓冲区内字节 buf=读取整个缓冲区 两个有区别 n=缓冲区内查出来的字节数
		DataSocket(buf[:n], conn)
	}
}

type Condition struct {
	Conn net.Conn
	Mac  string
	Type string
}

type Usertion struct {
	Conn net.Conn
	User string
	Key  string
}

type DataOrder struct {
	Order string
	Mac   string
}

func ControlSocket(conn net.Conn, str string) {
	var data DataOrder
	err := json2.Unmarshal([]byte(str), &data)
	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("命令格式错误 操作受限"))
		return
	}

	cone := ConnSelMac(&ConnList, data.Mac)
	switch cone.Type {
	case "鸣驹":
		switch data.Order {
		case "a":

			break
		case "b":
			cone.Conn.Write([]byte("a2"))
			<-time.After(100 * time.Millisecond)
			cone.Conn.Write([]byte("b2"))
			conn.Write([]byte("流程执行成功"))
			break
		case "c":
			cone.Conn.Write([]byte("a3"))
			<-time.After(100 * time.Millisecond)
			cone.Conn.Write([]byte("b3"))
			conn.Write([]byte("流程执行成功"))
			break
		}
		break
	case "":

	case "创度":

		break
	}
}

type Mac struct {
	Type string
	Mac  string
}
type User struct {
	Action string
	User   string
}

var ConnList = make([]Condition, 0, 0)
var UserList = make([]Usertion, 0, 0)

func DataSocket(data []byte, conn net.Conn) { //处理消息
	var mac Mac
	jsonStr := string(data)
	jsonStr = strings.ReplaceAll(string(data), "\x00", "")
	err := json2.Unmarshal([]byte(jsonStr), &mac)

	if ConnSel(&UserList, conn) {
		go ControlSocket(conn, jsonStr)
		return
	}

	if jsonStr == " " {
		return
	}

	if err != nil {
		fmt.Printf("IP:%v Byte：%v Cause:%v", conn.RemoteAddr().String(), string(data), "尝试与当前服务端建立通讯 通讯被拒绝")
		conn.Write([]byte("通讯被拒绝"))
		return
	}

	//校验 注册包 当前为用户检测 控制端用户检查 并创建切片数据
	if mac.Type == "" {
		var user User
		es := json2.Unmarshal([]byte(jsonStr), &user)
		if es != nil {
			fmt.Printf("IP:%v Byte：%v Cause:%v", conn.RemoteAddr().String(), string(data), "尝试与当前服务端建立通讯 通讯被拒绝")
			conn.Write([]byte("非法操作"))
			return
		}

		//校验 注册包 当前为用户检测 控制端用户检查 并创建切片数据
		if user.Action == "enrol" {
			row, err := DB.Query("SELECT ins FROM db.socket_user WHERE USER LIKE ? ;", user.User)
			if err != nil {
				fmt.Printf("IP:%v Byte：%v Cause:%v", conn.RemoteAddr().String(), string(data), "数据库执行查询异常")
				conn.Write([]byte("服务端数据库执行查询异常"))
				return
			}
			defer row.Close()
			if row.Next() {
				condition := Usertion{
					Conn: conn,
					User: user.User,
					Key:  "123",
				}
				go UserAdd(&UserList, conn, condition)
				conn.Write([]byte("用户验证"))
				return
			}
			return
		} else {
			fmt.Printf("IP:%v Byte：%v Cause:%v", conn.RemoteAddr().String(), string(data), "尝试与当前服务端建立通讯 通讯被拒绝")
			conn.Write([]byte("非法操作"))
			return
		}

	}

	//校验 注册包 当前为硬件检测 附带厂家及mac地址
	row, err := DB.Query("select sm.mac from db.socket_type st join db.socket_mac sm on sm.cid like st.ins where st.type  like ? and sm.mac like ? ;", mac.Type, mac.Mac)
	defer row.Close()
	if err != nil {
		return
	}
	type title struct {
		Mac string
	}
	if row.Next() {
		condition := Condition{
			Conn: conn,
			Mac:  mac.Mac,
			Type: mac.Type,
		}
		go ConnAdd(&ConnList, conn, condition)
		conn.Write([]byte("校检通过"))
	} else {
		fmt.Printf("IP:%v Byte：%v Cause:%v", conn.RemoteAddr().String(), string(data), "不属于本公司设备")
		conn.Write([]byte("校检失败"))
	}
}

func UserAdd(UserList *[]Usertion, Conn net.Conn, condition Usertion) {
	found := false
	for i := 0; i < len(*UserList); i++ {
		if (*UserList)[i].Conn.RemoteAddr().String() == Conn.RemoteAddr().String() {
			found = true
			break
		}
	}
	if !found {
		*UserList = append(*UserList, condition)
	}
	fmt.Println(*UserList)
}

func ConnAdd(ConnList *[]Condition, Conn net.Conn, condition Condition) {
	found := false
	for i := 0; i < len(*ConnList); i++ {
		if (*ConnList)[i].Conn.RemoteAddr().String() == Conn.RemoteAddr().String() {
			found = true
			break
		}
	}
	if !found {
		*ConnList = append(*ConnList, condition)
	}
	//fmt.Println(*ConnList)
}

func ConnDel(ConnList *[]Condition, UserList *[]Usertion, Conn net.Conn) {
	//切片删除合并 当前Conn切片内的
	for i := 0; i < len(*ConnList); i++ {
		if (*ConnList)[i].Conn.RemoteAddr().String() == Conn.RemoteAddr().String() {
			*ConnList = append((*ConnList)[:i], (*ConnList)[i+1:]...)
			break
		}
	}
	//当前为User切片内
	for i := 0; i < len(*UserList); i++ {
		if (*UserList)[i].Conn.RemoteAddr().String() == Conn.RemoteAddr().String() {
			*UserList = append((*UserList)[:i], (*UserList)[i+1:]...)
			break
		}
	}
}

func ConnSel(UserList *[]Usertion, Conn net.Conn) bool {
	var retu bool = false
	for i := 0; i < len(*UserList); i++ {
		if (*UserList)[i].Conn.RemoteAddr().String() == Conn.RemoteAddr().String() {
			retu = true
			break
		}
	}
	return retu
}

func ConnSelMac(ConnList *[]Condition, mac string) Condition {
	var conn Condition
	for i := 0; i < len(*ConnList); i++ {
		if (*ConnList)[i].Mac == mac {
			conn.Conn = (*ConnList)[i].Conn
			conn.Mac = (*ConnList)[i].Mac
			conn.Type = (*ConnList)[i].Type
			break
		}
	}
	return conn
}
