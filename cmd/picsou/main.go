package main

import (
	"fmt"
	"github.com/scality/picsou/pkg/stats"
	"github.com/scality/picsou/pkg/report"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/smtp"
	"strings"
)

var auth smtp.Auth

func main() {
	auth := smtp.PlainAuth("", "maxime.vaude@gmail.com", "totopassword", "smtp.gmail.com")
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s := stats.New(sess, stats.Regions)
	templateData := &report.TemplateData{
		Name: "Mr Freeze",
		URL: strings.Replace(s.Data.String(), "\n", "<br />", -1),
	}
	r := report.NewRequest(&auth, []string{"maxime.vaude@scality.com"}, "Hello Mr Freeze!", "Hello, world!", "./assets/reports/daily.html", templateData).ParseTemplate().SendEmail()
	if r != nil {
		fmt.Println("SendEmail Failure: ", err)
	}
}
