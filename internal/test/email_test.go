package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T){
	e := email.NewEmail()
	e.From = "Jordan Wright <cx_13535071701@163.com>"
	e.To = []string{"chefffy1223@gmail.com"}
	//e.Bcc = []string{"test_bcc@example.com"}
	//e.Cc = []string{"test_cc@example.com"}
	e.Subject = "TestSendEmail"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	//err := e.Send("smtp.163.com:465", smtp.PlainAuth("", "cx_13535071701@163.com", "HAMTYVKJZFAFPDGR", "smtp.163.com"))

	//返回EOF时，关闭SSL重试
	err := e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "cx_13535071701@163.com", "HAMTYVKJZFAFPDGR", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true,ServerName: "smtp.163.com"})

	if err !=nil {
		t.Fatal(err)
	}

}