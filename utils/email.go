package utils

import (
	"bytes"
	"crypto/tls"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dangtran47/go_crud/initializers"
	"github.com/dangtran47/go_crud/models"
	"github.com/k3a/html2text"
	"github.com/thanhpk/randstr"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendMail(user *models.User, data *EmailData) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	from := config.EmailFrom
	to := user.Email
	smtpPassword := config.SMTPPassword
	smtpUser := config.SMTPUser
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	if template.ExecuteTemplate(&body, "verificationCode.html", &data) != nil {
		log.Fatal("Could not execute template", err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email", err)
	}
}

func SendVerificationEmail(user *models.User) {
	config, _ := initializers.LoadConfig(".")

	code := randstr.String(6)
	verificationCode := Encode(code)

	firstName := user.Name

	if strings.Contains(user.Name, " ") {
		firstName = strings.Split(user.Name, " ")[0]
	}

	emailData := EmailData{
		URL:       config.ClientOrigin + "/verify/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	user.VerificationCode = verificationCode
	initializers.DB.Save(user)
	SendMail(user, &emailData)
}
