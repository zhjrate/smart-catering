package tool

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/smtp"
)

type MailBox struct {
	Dialer *gomail.Dialer
	From   string
	Count  int
}

func (receiver *MailBox) Connect(id int) {
	switch id {
	case 0:
		receiver.Dialer = gomail.NewDialer("smtp.exmail.qq.com", 465, "xiaolin@itda.vip", "Zhao132909")
		receiver.From = "Inversion@itda.vip"
		break
	case 1:
		//ghtculliwoeubade
		//zsnztjpgozmxggff
		receiver.Dialer = gomail.NewDialer("smtp.qq.com", 587, "1121489610@qq.com", "ghtculliwoeubade")
		receiver.From = "ITDA@foxmail.com"
		break
	case 2:
		//ysimqihnaiucddih
		receiver.Dialer = gomail.NewDialer("smtp.qq.com", 587, "2844047062@qq.com", "ysimqihnaiucddih")
		receiver.From = "FQCNET@foxmail.com"
		break

	case 3:
		//meck wjzt ewgs ynsm
		receiver.Dialer = gomail.NewDialer("smtp.gmail.com", 465, "fqcnet@gmail.com", "meck wjzt ewgs ynsm")
		receiver.From = "fqcnet@gmail.com"
		break
	}
	receiver.Count = id

	if _, err := receiver.Dialer.Dial(); err != nil {
		fmt.Println("Mail server connection error.")
	} else {
		fmt.Println("Mail server connected success.")
	}

}

func (receiver *MailBox) Send(Subject string, content string, Mailbox string, Passage chan bool) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("Message-ID", "<unique-message-id@itda.vip>")
	mailer.SetHeader("From", receiver.From)
	mailer.SetHeader("To", Mailbox)
	mailer.SetHeader("Subject", Subject)
	mailer.SetBody("text/html", content)
	// 发送邮件
	err := receiver.Dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Could not send email: %v", err)
	}
	Passage <- true
	//return nil
}

func Conn() {
	// 建立与邮件服务器的连接
	client, err := smtp.Dial("smtp.exmail.qq.com:465")
	if err != nil {
		log.Fatalf("Could not connect to SMTP server: %v", err)
		return
	}
	defer client.Close()

	// 使用SMTP登录信息
	auth := smtp.PlainAuth("", "xiaolin@itda.vip", "Zhao132909", "smtp.exmail.qq.com")

	// 发送邮件
	sendEmail(client, auth, "Inversion@itda.vip", "2844047062@qq.com", "Test Subject", "Hello, this is a test email!")

	fmt.Println("Email sent successfully")
}

func sendEmail(client *smtp.Client, auth smtp.Auth, from, to, subject, body string) {
	// 设置发件人
	if err := client.Mail(from); err != nil {
		log.Printf("Could not set sender: %v", err)
		return
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		log.Printf("Could not set recipient: %v", err)
		return
	}

	// 创建邮件数据
	wc, err := client.Data()
	if err != nil {
		log.Printf("Could not create email data: %v", err)
		return
	}
	defer wc.Close()

	// 设置邮件头
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n", from, to, subject)
	if _, err := wc.Write([]byte(headers)); err != nil {
		log.Printf("Could not write email headers: %v", err)
		return
	}

	// 设置邮件正文
	if _, err := wc.Write([]byte(body)); err != nil {
		log.Printf("Could not write email body: %v", err)
		return
	}
}

func (receiver *MailBox) Template(id int, name string, action string, code string, user string) string {
	switch id {
	case 1: //验证码
		temptext := `<!DOCTYPE html>
<html lang="en">
<head>
    <base target="_blank" />
    <style type="text/css">::-webkit-scrollbar{ display: none; }</style>
    <style id="cloudAttachStyle" type="text/css">#divNeteaseBigAttach, #divNeteaseBigAttach_bak{display:none;}</style>
    <style id="blockquoteStyle" type="text/css">blockquote{display:none;}</style>
    <style type="text/css">
        body{font-size:14px;font-family:arial,verdana,sans-serif;line-height:1.666;padding:0;margin:0;overflow:auto;white-space:normal;word-wrap:break-word;min-height:100px}
        td, input, button, select, body{font-family:Helvetica, 'Microsoft Yahei', verdana}
        pre {white-space:pre-wrap;white-space:-moz-pre-wrap;white-space:-pre-wrap;white-space:-o-pre-wrap;word-wrap:break-word;width:95%}
        th,td{font-family:arial,verdana,sans-serif;line-height:1.666}
        img{ border:0}
        header,footer,section,aside,article,nav,hgroup,figure,figcaption{display:block}
        blockquote{margin-right:0px}
    </style>
</head>
<body tabindex="0" role="listitem">
<table width="700" border="0" align="center" cellspacing="0" style="width:700px;">
    <tbody>
    <tr>
        <td>
            <div style="width:700px;margin:0 auto;border-bottom:1px solid #ccc;margin-bottom:30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="700" height="39" style="font:12px Tahoma, Arial, 宋体;">
                    <tbody><tr><td width="210"></td></tr></tbody>
                </table>
            </div>
            <div style="width:680px;padding:0 10px;margin:0 auto;">
                <div style="line-height:1.5;font-size:14px;margin-bottom:25px;color:#4d4d4d;">
                    <strong style="display:block;margin-bottom:15px;">尊敬的用户：<span style="color:#f60;font-size: 20px;">%s</span>。</strong>
                    <strong style="display:block;margin-bottom:15px;">
                        您正在进行<span style="color: red">%s</span>操作，请在验证码输入框中输入：<span style="color:#f60;font-size: 24px">%s</span> ，以完成操作，有效期3分钟。
                    </strong>
                </div>
                <div style="margin-bottom:30px;">
                    <small style="display:block;margin-bottom:20px;font-size:12px;">
                        <p style="color:#747474;"><br>
                            注意：此操作可能会修改您的密码、登录邮箱或绑定手机。如非本人操作，请及时登录并修改密码以保证帐户安全。
                            <br>（工作人员不会向您索取此验证码，请勿泄漏! )
                        </p>
                    </small>
                </div>
            </div>
            <div style="width:700px;margin:0 auto;">
                <div style="padding:10px 10px 0;border-top:1px solid #ccc;color:#747474;margin-bottom:20px;line-height:1.3em;font-size:12px;text-align: right;">
                    <p>请保管好您的邮箱，避免账号被他人盗用。<br>
						此为系统邮件，支持回复反馈问题。<br>如有打扰，敬请见谅。
                    </p>
                    <p>CuBing Max %s</p>
                </div>
            </div>
        </td>
    </tr>
    </tbody>
</table>
</body>
</html>`

		Template := fmt.Sprintf(temptext, "", name, action, code, user)
		return Template

		break
	case 2: //登陆提醒
		temptext := `<!DOCTYPE html>
<html lang="en">
<head>
    <base target="_blank" />
    <style type="text/css">::-webkit-scrollbar{ display: none; }</style>
    <style id="cloudAttachStyle" type="text/css">#divNeteaseBigAttach, #divNeteaseBigAttach_bak{display:none;}</style>
    <style id="blockquoteStyle" type="text/css">blockquote{display:none;}</style>
    <style type="text/css">
        body{font-size:14px;font-family:arial,verdana,sans-serif;line-height:1.666;padding:0;margin:0;overflow:auto;white-space:normal;word-wrap:break-word;min-height:100px}
        td, input, button, select, body{font-family:Helvetica, 'Microsoft Yahei', verdana}
        pre {white-space:pre-wrap;white-space:-moz-pre-wrap;white-space:-pre-wrap;white-space:-o-pre-wrap;word-wrap:break-word;width:95%}
        th,td{font-family:arial,verdana,sans-serif;line-height:1.666}
        img{ border:0}
        header,footer,section,aside,article,nav,hgroup,figure,figcaption{display:block}
        blockquote{margin-right:0px}
    </style>
</head>
<body tabindex="0" role="listitem">
<table width="700" border="0" align="center" cellspacing="0" style="width:700px;">
    <tbody>
    <tr>
        <td>
            <div style="width:700px;margin:0 auto;border-bottom:1px solid #ccc;margin-bottom:30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="700" height="39" style="font:12px Tahoma, Arial, 宋体;">
                    <tbody><tr><td width="210"></td></tr></tbody>
                </table>
            </div>
            <div style="width:680px;padding:0 10px;margin:0 auto;">
                <div style="line-height:1.5;font-size:14px;margin-bottom:25px;color:#4d4d4d;">
                    <strong style="display:block;margin-bottom:15px;">尊敬的用户：<span style="color:#f60;font-size: 18px;">%s</span>。</strong>
                    <strong style="display:block;margin-bottom:15px;">
                        您于北京时间:<span style="color: red">%s</span>登陆系统，登陆IP为：<span style="color:#f60;font-size: 22px">%s</span> ，该消息为内测提醒功能,请忽略。谢谢合作。
                    </strong>
                </div>
                <div style="margin-bottom:30px;">
                    <small style="display:block;margin-bottom:20px;font-size:12px;">
                        <p style="color:#747474;"><br>
                            注意：此操作可能会修改您的密码、登录邮箱或绑定手机。如非本人操作，请及时登录并修改密码以保证帐户安全。
                            <br>（工作人员不会向您索取此验证码，请勿泄漏! )
                        </p>
                    </small>
                </div>
            </div>
            <div style="width:700px;margin:0 auto;">
                <div style="padding:10px 10px 0;border-top:1px solid #ccc;color:#747474;margin-bottom:20px;line-height:1.3em;font-size:12px;text-align: right;">
                    <p>请保管好您的邮箱，避免账号被他人盗用。<br>
						此为系统邮件，支持回复反馈问题。<br>如有打扰，敬请见谅。
                    </p>
                    <p>CuBing Max %s</p>
                </div>
            </div>
        </td>
    </tr>
    </tbody>
</table>
</body>
</html>`

		Template := fmt.Sprintf(temptext, "", name, action, code, user)
		return Template

		break
	}
	return ""

}
