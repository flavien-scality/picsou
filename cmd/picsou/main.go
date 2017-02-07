package main

import (
	"github.com/scality/picsou/pkg/stats"
	"github.com/scality/picsou/pkg/report"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/stmp"
)

var auth smtp.PlainAuth

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s = stats.New(sess, stats.Regions)
	templateData := &report.TemplateData{
		Name: "Mr Freeze",
		URL: s.data,
	}
	r := NewRequest([]string{"maxime.vaude@scality.com"}, "Hello Mr Freeze!", "Hello, world!", "./assets/reports/daily.html", templateData).ParseTemplate().SendEmail()
	if r == nil {
		fmt.Println("SendEmail Failure")
	}
