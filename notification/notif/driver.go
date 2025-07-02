package notif

type Driver interface {
	Send(to, subject, message string) error
	SendWithTemplate(to, subject, message, template string) error
}
