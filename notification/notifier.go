package notification

import (
	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/notification/email"
	"gitag.ir/armogroup/armo/services/reality/notification/notif"
	"gitag.ir/armogroup/armo/services/reality/notification/sms"
)

type Notifier interface {
	SendSMS(user models.User, message string, template string) error
	SendEmail(user models.User, message string, subject string, template string) error
	SendPhoneVerifyCode(phone string, message string) (err error)
	SendEmailVerifyCode(to string, message string) (err error)
}

type notifier struct {
	sms   notif.Driver
	email notif.Driver
}

func MakeNotifier() Notifier {
	tokenKey := config.AppConfig.KavenegarApiKey

	kavenegar := sms.NewKavenegar(tokenKey)

	host := config.AppConfig.SMTPHost
	port := config.AppConfig.SMTPPort
	username := config.AppConfig.SMTPUsername
	password := config.AppConfig.SMTPPassword
	from := config.AppConfig.SMTPFrom

	if config.AppConfig.SendEmail {
		mailer := email.NewMailer(host, port, username, password, from)
		return notifier{kavenegar, mailer}
	}

	return notifier{kavenegar, nil}
}

func (n notifier) SendSMS(user models.User, message string, template string) (err error) {
	err = n.sms.SendWithTemplate(user.Phone, message, "", template)
	return
}

func (n notifier) SendEmail(user models.User, message string, subject string, template string) (err error) {
	err = n.email.Send(user.Email.String, message, subject)
	return
}

func (n notifier) SendPhoneVerifyCode(phone string, message string) (err error) {
	// FIXME: ctpverify should be a config
	err = n.sms.SendWithTemplate(phone, "", message, "ctpverify")
	return
}

func (n notifier) SendEmailVerifyCode(email string, message string) (err error) {
	err = n.email.Send(email, "Verify Code", message)
	return
}
