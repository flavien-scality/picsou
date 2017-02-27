package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/apex/go-apex"
	"github.com/scality/picsou/pkg/stats"
	"github.com/scality/picsou/pkg/report"
	"github.com/aws/aws-sdk-go/aws/session"
)

var auth smtp.Auth

type message struct {
	Value string `json:"value"`
}

func main() {
	auth := smtp.PlainAuth("", os.Getenv("PICSOU_USER"), os.Getenv("PICSOU_PSD"), "smtp.gmail.com")
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s := stats.New(sess)
	templateData := &report.TemplateData{
		Name: "Mr Freeze",
		Stats: s,
	}
	r := report.NewRequest(&auth, []string{"maxime.vaude@scality.com", "thibault.riviere@scality.com", "mathieu.cassagne@scality.com"}, "AWS Daily Report", "Hello, world!", "./assets/reports/daily.html", templateData).ParseTemplate().SendEmail()
	if r != nil {
		fmt.Println("SendEmail Failure: ", err)
	}
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}
		m.Value = strings.ToUpper(m.Value)
		return m, nil
	})
}
