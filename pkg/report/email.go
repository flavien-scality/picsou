package report

import (
	"bytes"
	"fmt"
	"github.com/scality/picsou/pkg/settings"
	"github.com/scality/picsou/pkg/stats"
	"html/template"
	"net/smtp"
	"os"
	"time"
)

// TODO: 3 different example report levels: total, regions, person (SSH KEY)

// TemplateData struct
type TemplateData struct {
	Settings *settings.Settings
	Data *stats.Stats
}

// Template struct
type Template struct {
	filename string
	data     *TemplateData
}

// Request struct
type Request struct {
	auth     *smtp.Auth
	from     string
	to       []string
	subject  string
	body     string
	template *Template
}

// NewRequest create a new Request struct with the parameters given
func NewRequest(auth *smtp.Auth, to []string, subject, body, templatePath string, templateData *TemplateData) *Request {
	return &Request{
		auth:    auth,
		to:      to,
		subject: subject,
		body:    body,
		template: &Template{
			filename: templatePath,
			data:     templateData,
		},
	}
}

// SendEmail handles sending an email with the formated template
func (r *Request) SendEmail() *Request {
	fmt.Println("r.to: ", r.to)
	date := time.Now()
	year, month, day := date.Date()
	hour, min, sec := date.Clock()
	now := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", day, month, year, hour, min, sec)
	from := fmt.Sprintf("From: %s\r\n", os.Getenv("PICSOU_USER"))
	to := "To: " + r.to[0] + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	subject := fmt.Sprintf("Subject: %s %s\r\n\r\n", r.subject, now)
	msg := []byte(from + to + mime + subject + "\n" + r.body + "\r\n")
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, *r.auth, fmt.Sprintf("%s", os.Getenv("PICSOU_USER")), r.to, msg); err != nil {
		fmt.Println("pb during SendMail: ", err)
		return nil
	}
	return r
}

// ParseTemplate format the template using the metadata giver by the struct data
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
