package background

import (
	"fmt"
	"os"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nelsonmarro/bookings/config"
	"github.com/nelsonmarro/bookings/internal/email"
)

func ListenForMail(app *config.AppConfig) {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg, app)
		}
	}()
}

func sendMsg(m email.MailData, app *config.AppConfig) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println("Error connecting to SMTP server:", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("static/html/email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println("Error reading email template:", err)
		}

		mailContent := string(data)
		msgToSen := strings.Replace(mailContent, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSen)
	}

	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println("Error sending email:", err)
	} else {
		app.InfoLog.Println("Email sent successfully to", m.To)
	}
}
