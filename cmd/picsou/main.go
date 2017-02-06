package main

import (
	"github.com/scality/picsou/pkg/stats"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/smtp"
	"log"
	"bytes"
	"strconv"
	"html/template"
	"strings"
	// "fmt"
)


// TODO: fix email in specific pkg
type EmailUser struct {
    Username    string
    Password    string
    EmailServer string
    Port        int
}

type SmtpTemplateData struct {
    From    string
    To      string
    Subject string
    Body    string
}

var email = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"> 
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>Email title or subject</title>

<style type="text/css">
	
  /* reset */
  #outlook a {padding:0;} /* Force Outlook to provide a "view in browser" menu link. */ 
  .ExternalClass {width:100%;} /* Force Hotmail to display emails at full width */  
  .ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div {line-height: 100%;} /* Forces Hotmail to display normal line spacing.  More on that: http://www.emailonacid.com/forum/viewthread/43/ */ 
  p {margin: 0; padding: 0; font-size: 0px; line-height: 0px;} /* squash Exact Target injected paragraphs */
  table td {border-collapse: collapse;} /* Outlook 07, 10 padding issue fix */
  table {border-collapse: collapse; mso-table-lspace:0pt; mso-table-rspace:0pt; } /* remove spacing around Outlook 07, 10 tables */
  
  /* bring inline */
  img {display: block; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic;}
  a img {border: none;} 
  a {text-decoration: none; color: #000001;} /* text link */
  a.phone {text-decoration: none; color: #000001 !important; pointer-events: auto; cursor: default;} /* phone link, use as wrapper on phone numbers */
  span {font-size: 13px; line-height: 17px; font-family: monospace; color: #000001;}

</style>
<!--[if gte mso 9]>
  <style>
  /* Target Outlook 2007 and 2010 */
  </style>
<![endif]-->
</head>
<body style="width:100%; margin:0; padding:0; -webkit-text-size-adjust:100%; -ms-text-size-adjust:100%;">

<!-- body wrapper -->
<table cellpadding="0" cellspacing="0" border="0" style="margin:0; padding:0; width:100%; line-height: 100% !important;">
  <tr>
    <td valign="top">
      <!-- edge wrapper -->
      <table cellpadding="0" cellspacing="0" border="0" align="center" width="600" style="background: #efefef;">
        <tr>
          <td valign="top">
            <!-- content wrapper -->
            <table cellpadding="0" cellspacing="0" border="0" align="center" width="560" style="background: #cfcfcf;">
              <tr>
                <td valign="top" style="vertical-align: top;">
<!-- ///////////////////////////////////////////////////// -->

<table cellpadding="0" cellspacing="0" border="0" align="center">
  <tr>
    <td valign="top" style="vertical-align: top;">
      <span style="">text</span>
    </td>
  </tr>
</table>
<table cellpadding="0" cellspacing="0" border="0" align="center">
  <tr>
    <td valign="top" style="vertical-align: top;">
      <img src="full path to image" alt="alt text" title="title text" width="50" height="50" style="width: 50px; height: 50px;"/>
    </td>
  </tr>
</table>
<table cellpadding="0" cellspacing="0" border="0" align="center">
  <tr height="30">
    <td valign="top" style="vertical-align: top; background: #efefef;" width="600" >
    </td>
  </tr>
</table>

<!-- //////////// -->
                </td>
              </tr>
            </table>
            <!-- / content wrapper -->
          </td>
        </tr>
      </table>
      <!-- / edge wrapper -->
    </td>
  </tr>
</table>  
<!-- / page wrapper -->
</body>
</html>
`

func main() {
    var err error
    sess, err := session.NewSession()
    if err != nil {
        panic(err)
    }
    _ = stats.New(sess, stats.Regions)
    emailUser := &EmailUser{"", "", "smtp.gmail.com", 587}
    auth := smtp.PlainAuth("",
        emailUser.Username,
        emailUser.Password,
        emailUser.EmailServer,
    )
    emailTemplate := []byte("From: {{ .From }}\r\n" +
        "To: {{ .To }}\r\n" +
        "Subject: {{ .Subject }}\r\n" +
        "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
	"\n" +
        "{{ .Body }}" +
        "Sincerely,\n" +
        "{{ .From }}\r\n")

    var doc bytes.Buffer

    context := &SmtpTemplateData{
        "maxime.vaude@gmail.com",
        "maxime.vaude@scality.com",
        "AWS Cost Report",
        //data.Res.String(),
	email,
    }
    t := template.New("emailTemplate")
    t, err = t.Parse(string(emailTemplate))
    if err != nil {
        log.Print("error trying to parse mail template")
    }
    err = t.Execute(&doc, context)
    if err != nil {
        log.Print("error trying to execute mail template")
    }
    err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
        auth,
        emailUser.Username,
        []string{"maxime.vaude@scality.com"},
        doc.Bytes())
    if err != nil {
        log.Print("ERROR: attempting to send a mail ", strings.Replace(err.Error(), "\n5.7.14 ", "", -1))
    }
}
