package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"gitag.ir/armogroup/armo/services/reality/notification/notif"
)

type mailer struct {
	client *smtp.Client
	from   string
}

func NewMailer(host string, port int, username string, password string, from string) notif.Driver {
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Fatalf("Could not dial to SMTP server: %v", err)
	}

	// Upgrade the connection to TLS using STARTTLS
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false, // set to true if you don't want to verify TLS certificate
	}
	if err := conn.StartTLS(tlsConfig); err != nil {
		log.Fatalf("Could not start TLS: %v", err)
	}

	auth := smtp.PlainAuth("", username, password, host)
	if err := conn.Auth(auth); err != nil {
		log.Fatalf("Could not authenticate: %v", err)
	}

	return &mailer{
		client: conn,
		from:   from,
	}

}

func (m *mailer) Send(to, subject, message string) error {
	// Now use the established connection to send the email.
	msg := []byte("To: " + to + "\r\n" +
		"From: " + m.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	if err := m.client.Mail(m.from); err != nil {
		return fmt.Errorf("could not set sender: %v", err)
	}
	if err := m.client.Rcpt(to); err != nil {
		return fmt.Errorf("could not set recipient: %v", err)
	}
	wc, err := m.client.Data()
	if err != nil {
		return fmt.Errorf("could not get writer: %v", err)
	}
	_, err = wc.Write(msg)
	if err != nil {
		return fmt.Errorf("could not write message: %v", err)
	}
	err = wc.Close()
	if err != nil {
		return fmt.Errorf("could not close writer: %v", err)
	}
	return nil
}

func (m *mailer) SendWithTemplate(to, subject, message, template string) error {
	// Now use the established connection to send the email.
	panic("implement me. later decouple from send")
}
