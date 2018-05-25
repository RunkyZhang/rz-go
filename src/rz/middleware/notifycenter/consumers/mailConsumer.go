//package consumers
//
//import (
//	"gopkg.in/gomail.v2"
//	"net/smtp"
//	"rz/middleware/notifycenter/global"
//)
//
//var (
//	MailConsumer *mailConsumer
//)
//
//func init() {
//	MailConsumer = &mailConsumer{
//		Host: global.Config.Mail.Host,
//		Port: global.Config.Mail.Port,
//		UserName: global.Config.Mail.UserName,
//		password: global.Config.Mail.Password,
//		From: global.Config.Mail.From,
//		ContentType: global.Config.Mail.ContentType,
//	}
//
//	MailConsumer.dialer = gomail.NewDialer(MailConsumer.Host, MailConsumer.Port, MailConsumer.UserName, MailConsumer.password)
//	MailConsumer.dialer. = unencryptedAuth{
//		smtp.PlainAuth(
//			"",
//			"monitor@gangtiequn.com",
//			"399GCG",
//			"mail.gangtiequn.com",
//		),
//	}
//}
//
//type mailConsumer struct {
//	Host        string
//	Port        int
//	UserName    string
//	password    string
//	From        string
//	ContentType string
//
//	dialer *Dialer
//}
//
//func SendMail(subject string, content string, mails []string) error {
//	m := gomail.NewMessage()
//	m.SetHeader("From", "notifycenter <monitor@gangtiequn.com>")
//	m.SetHeader("To", mails...)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", content)
//	//	m.Attach("/home/Alex/lolcat.jpg")
//	if err := d.DialAndSend(m); err != nil {
//		return err
//	}
//	return nil
//}
//
//type unencryptedAuth struct {
//	smtp.Auth
//}
//
//func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
//	s := *server
//	s.TLS = true
//	return a.Auth.Start(&s)
//}
//
