package user

import (
	"awesomeProject2/data/user"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

var from = "gidiboateng200@gmail.com"
var password = "ogtg kebv rrom srah"
var smtpHost = "smtp.gmail.com"
var smtpPort = "587"
var mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func EmailLogics(subject, templatePath string, emailDto *user.MailInputs, templateData interface{}) error {

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	to := []string{emailDto.Email}

	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	err = t.Execute(&body, templateData)
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func VerifyEmail(emailDto *user.MailInputs) error {
	templateData := struct {
		Code     string
		Username string
		Link     string
	}{
		Code:     emailDto.Code,
		Username: emailDto.Username,
		Link:     os.Getenv("FRONTEND_URL") + "/verify",
	}
	return EmailLogics("Verify Account", "mail/user/verify.html", emailDto, templateData)
}

func ResetPassword(emailDto *user.MailInputs) error {
	templateData := struct {
		Code     string
		Username string
		Link     string
	}{
		Code:     emailDto.Code,
		Username: emailDto.Username,
		Link:     os.Getenv("BACKEND_URL") + "api/user/reset-password?code=" + emailDto.Code,
	}
	return EmailLogics("Reset Password", "mail/user/reset.html", emailDto, templateData)
}
