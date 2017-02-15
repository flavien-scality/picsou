package main

import (
	"fmt"
	"github.com/scality/picsou/pkg/stats"
	"github.com/scality/picsou/pkg/report"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/smtp"
)

var auth smtp.Auth

func main() {
	auth := smtp.PlainAuth("", "maxime.vaude@gmail.com", "toto", "smtp.gmail.com")
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s := stats.New(sess)
	templateData := &report.TemplateData{
		Name: "Mr Freeze",
		Stats: s,
	}
	r := report.NewRequest(&auth, []string{"maxime.vaude@scality.com"}, "AWS Daily Report", "Hello, world!", "./assets/reports/daily.html", templateData).ParseTemplate().SendEmail()
	if r != nil {
		fmt.Println("SendEmail Failure: ", err)
	}
}
