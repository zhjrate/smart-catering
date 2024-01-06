package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Toolkit struct {
}

// 换页取值 中间5位
func (receiver Toolkit) Page(sm int, max int) []int {
	Number := make([]int, min(max, 5, max-(sm-1)))
	for k := 0; k < len(Number); k++ {
		Number[k] = sm + k
	}
	fmt.Println(Number)
	return Number
}

func min(a int, b int, c int) int {
	if a > b {
		if c < b {
			return c
		} else {
			return b
		}
	} else if a < b {
		return c
	}
	return b
}

func (e Toolkit) MaxInt(number []int) (int, error) {
	if len(number) == 0 {
		return 0, errors.New("切片长度为0")
	}
	max := number[0]
	for _, num := range number {
		if num > max {
			max = num
		}
	}
	fmt.Println(number)
	return max, nil
}

func (receiver Toolkit) MinInt(number []int) (int, error) {
	if len(number) == 0 {
		return 0, errors.New("切片长度为0")
	}
	min := number[0]
	for _, num := range number {
		if num < min {
			min = num
		}
	}
	return min, nil
}

func (receiver Toolkit) MaxFloat(number []float64) (float64, error) {
	if len(number) == 0 {
		return 0, errors.New("切片长度为0")
	}
	max := number[0]
	for _, num := range number {
		if num > max {
			max = num
		}
	}
	return max, nil
}

func (receiver Toolkit) MinFloat(number []float64) (float64, error) {
	if len(number) == 0 {
		return 0, errors.New("切片长度为0")
	}
	min := number[0]
	for _, num := range number {
		if num < min {
			min = num
		}
	}
	return min, nil
}

func (receiver Toolkit) MaxTime(number []time.Time) (time.Time, error) {
	if len(number) == 0 {
		return time.Time{}, errors.New("切片长度为0")
	}
	max := number[0]
	for _, date := range number {
		if date.After(max) {
			max = date
		}
	}
	return max, nil
}

func (receiver Toolkit) MinTime(number []time.Time) (time.Time, error) {
	if len(number) == 0 {
		return time.Time{}, errors.New("切片长度为0")
	}
	min := number[0]
	for _, date := range number {
		if date.Before(min) {
			min = date
		}
	}
	return min, nil
}

func (receiver Toolkit) RandSeed(i int) string {

	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	result := make([]byte, i)
	for t := 0; t < i; t++ {
		result[t] = digits[rand.Intn(len(digits))]
	}
	return string(result)

}

// MD5二次封装
func MD5Bytes(s []byte) string {
	ret := md5.Sum(s)
	return hex.EncodeToString(ret[:])
}

// 生成随机数 长度可控制
func RandAuto(count int, max int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var retu string = ""
	for i := 0; i < count; i++ {
		c := r.Intn(max)
		//if i == 1 {
		//	retu = strconv.Itoa(c)
		//} else {
		//	retu = retu + "-" + strconv.Itoa(c)
		//}
		retu = retu + "" + strconv.Itoa(c)
	}
	return retu
}

//加密解密第三方

// var PwdKey = []byte("DIS**#KKKDJJSKDI")
var PwdKey = "linkbook1qaz*WSX"

// PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func PKCS7UnPadding1(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

// 实现加密Aes
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 实现解密Aes
func AesDeCrypt(cypted []byte, key []byte) (string, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding1(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}

// 加密base64
func EnPwdCode(pwdStr string) string {
	pwd := []byte(pwdStr)
	result, err := AesEcrypt(pwd, []byte(PwdKey))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(result)
}

// 解密base64
func DePwdCode(pwd string) string {
	temp, _ := hex.DecodeString(pwd)
	//执行AES解密
	res, _ := AesDeCrypt(temp, []byte(PwdKey))
	return res
}

//func main() {
//
//	//aes加密
//	destring:=`{"name":"菜鸟教程11","site":"http://www.runoob.com"}`
//	deStr := EnPwdCode(destring)
//	fmt.Println(deStr) //4f4d74c15e0ad4afb323a17927b1176ecb0c95ecbdf8e776ceb093499e3ff4c45157b007ae7dff1688ac2d2bf9fef28644922a1b3bbc6ef5881cb1ed0dff298a
//
//	//aes解密
//	decodeStr := DePwdCode("4f4d74c15e0ad4afb323a17927b1176ecb0c95ecbdf8e776ceb093499e3ff4c45157b007ae7dff1688ac2d2bf9fef28644922a1b3bbc6ef5881cb1ed0dff298a")
//	fmt.Println(decodeStr) //{"name":"菜鸟教程11","site":"http://www.runoob.com"}
//}

func CFBEncrypt(key []byte, message string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	ct := make([]byte, aes.BlockSize+len(message))
	copy(ct[:aes.BlockSize], iv)

	cfb := cipher.NewCFBEncrypter(c, iv)
	cfb.XORKeyStream(ct[aes.BlockSize:], []byte(message))

	return base64.StdEncoding.EncodeToString(ct), nil
}

func CFBDecode(key []byte, message string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ct, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	iv := ct[:aes.BlockSize]
	ct = ct[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(c, iv)
	cfb.XORKeyStream(ct, ct)

	return string(ct), nil
}

func (receiver Toolkit) JsonSlice(data any) map[string]interface{} {
	var rute map[string]interface{}
	jsonDate, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshaling error:", err)
		return rute
	}
	if err := json.Unmarshal(jsonDate, &rute); err != nil {
		fmt.Println("JSON unmarshaling error:", err)
		return rute
	}
	return rute
}

func (receiver Toolkit) JsonList(data interface{}) ([]map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON encoding error:", err)
		return nil, err
	}

	var parsedData []map[string]interface{}
	err = json.Unmarshal(jsonData, &parsedData)
	if err != nil {
		fmt.Println("JSON decoding error:", err)
		return nil, err
	}

	return parsedData, nil
}
