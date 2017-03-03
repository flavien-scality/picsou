package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/scality/picsou/pkg/report"
	"github.com/scality/picsou/pkg/settings"
	"github.com/scality/picsou/pkg/stats"
	"github.com/urfave/cli"
	"net/smtp"
	"os"
	"sort"
)

var auth smtp.Auth

func getReport() {
	auth := smtp.PlainAuth("", os.Getenv("PICSOU_USER"), os.Getenv("PICSOU_PSD"), "smtp.gmail.com")
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	meta := settings.New("./assets/settings.yml")
	fmt.Println("parsing settings done: ", meta)
	s := stats.New(sess)
	templateData := &report.TemplateData{
		Settings: meta,
		Data: s,
	}
	r := report.NewRequest(&auth, meta.GetUsersEmail(), "AWS Daily Report", "Hello, world!", "./assets/reports/daily.html", templateData).ParseTemplate().SendEmail()
	if r != nil {
		fmt.Println("SendEmail Failure: ", err)
	}
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "upload",
			Aliases: []string{"u"},
			Usage: "upload secret file to S3",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "download",
			Aliases: []string{"d"},
			Usage: "download secret file from S3",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "report",
			Aliases: []string{"r"},
			Usage: "send AWS costs report",
			Action: func(c *cli.Context) error {
				getReport()
				return nil
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Run(os.Args)
}
