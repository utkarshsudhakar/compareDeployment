package utils

import (
	"crypto/tls"
	"fmt"
	"../config"
	"gopkg.in/gomail.v2"
	
)

func SendMail(body string, subject string, cc string) {

	from := config.FromEmail

//	conf := ReadConfig()

	//to := conf.ToEmail
	to:=[]string{"usudhakar@informatica.com", }

	host := "mail.informatica.com"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Cc", cc)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Attach("test.jpg")

	d := gomail.NewDialer(host, 25, from, "")

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	//send mail
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Email Sent successfully!")

}
