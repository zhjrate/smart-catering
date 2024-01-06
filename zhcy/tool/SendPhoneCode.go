package tool

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PhoneCode struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

func (receiver *PhoneCode) Init() *PhoneCode {
	receiver.AccessKeyId = "LTAI5tJmmAEr4qTGAskfLkG1"
	receiver.AccessKeySecret = "4k98VwrHCYxJqYMrcMzoinGJGElJK9"
	receiver.SignName = "创泓度网络"
	receiver.TemplateCode = "SMS_464445625"
	return receiver
}

func (receiver *PhoneCode) SendCode(Phone string, PhoneCode string) error {
	// 构建请求参数
	TemplateParam := fmt.Sprintf(`{"code":"%s"}`, PhoneCode)
	paras := url.Values{}
	paras.Add("AccessKeyId", receiver.AccessKeyId)
	paras.Add("Format", "JSON")
	paras.Add("SignatureMethod", "HMAC-SHA1")
	paras.Add("SignatureVersion", "1.0")
	paras.Add("SignatureNonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	paras.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	paras.Add("Action", "SendSms")
	paras.Add("Version", "2017-05-25")
	paras.Add("RegionId", "cn-hangzhou")
	paras.Add("PhoneNumbers", Phone)
	paras.Add("SignName", receiver.SignName)
	paras.Add("TemplateCode", receiver.TemplateCode)
	paras.Add("TemplateParam", TemplateParam)

	// 构建用于签名的字符串
	sortedQueryString := paras.Encode()
	stringToSign := "POST&" + url.QueryEscape("/") + "&" + url.QueryEscape(sortedQueryString)

	// 签名
	signature := sign(receiver.AccessKeySecret+"&", stringToSign)
	paras.Add("Signature", signature)
	resp, err := http.Post("http://dysmsapi.aliyuncs.com/", "application/x-www-form-urlencoded", strings.NewReader(paras.Encode()))
	if err != nil {
		fmt.Println("发送请求失败:", err)

		return errors.New("发送请求失败:" + err.Error())
	}
	defer resp.Body.Close()
	return nil
}

func sign(accessKeySecret, stringToSign string) string {
	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
