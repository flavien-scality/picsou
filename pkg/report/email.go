package report

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
	"github.com/scality/picsou/pkg/stats"
)

// TODO: 3 different example report levels: total, regions, person (SSH KEY)

type TemplateData struct {
	Name string
	Stats *stats.Stats
}

type Template struct {
	filename string
	data *TemplateData
}

//Request struct
type Request struct {
	auth    *smtp.Auth
	from    string
	to      []string
	subject string
	body    string
	template *Template
}

func NewRequest(auth *smtp.Auth, to []string, subject, body, templatePath string, templateData *TemplateData) *Request {
	return &Request{
		auth: auth,
		to:      to,
		subject: subject,
		body:    body,
		template: &Template{
			filename: templatePath,
			data: templateData,
		},
	}
}

func (r *Request) SendEmail() *Request {
	date := time.Now()
	year, month, day := date.Date()
	hour, min, sec := date.Clock()
	now := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", day, month, year, hour, min, sec)
	from := "From: maxime.vaude@gmail.com\r\n"
	to := "To: " + r.to[0] + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	subject := fmt.Sprintf("Subject: %s %s\r\n\r\n", r.subject, now)
	msg := []byte(from + to + mime + subject + "\n" + r.body + "\r\n")
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, *r.auth, "maxime.vaude@gmail.com", r.to, msg); err != nil {
		fmt.Println("pb during SendMail: ", err)
		return nil
	}
	return r
}

func (r *Request) ParseTemplate() *Request {
	t, err := template.ParseFiles(r.template.filename)
	if err != nil {
		fmt.Println("pb during parse template: ", err)
		return nil
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, r.template.data); err != nil {
		fmt.Println("pb during execute template: ", err)
		return nil
	}
	r.body = buf.String()
	return r
}
