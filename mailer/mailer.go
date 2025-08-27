package mailer

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
)

type Mailer struct {
	from     string
	password string
	host     string
	port     string
	tmpl     *template.Template
}

func absoluteTemplatesDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "templates")
	return dir
}

// Load the template and return a new ready-to-use mailer
func New() *Mailer {
	tmpl := template.Must(template.ParseGlob(filepath.Join(absoluteTemplatesDir(), "*.html")))

	return &Mailer{
		from:     os.Getenv("SENDER_EMAIL"),
		password: os.Getenv("SMTP_PASS"),
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		tmpl:     tmpl,
	}
}

func (m *Mailer) Send(to, subject string, data any) error {
	var body bytes.Buffer

	if err := m.tmpl.ExecuteTemplate(&body, "otp_mailer.html", data); err != nil {
		return err
	}

	msg := "From: " + m.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body.String()

	// addr := m.host + ":" + m.port
	// auth := smtp.PlainAuth("", m.from, m.password, m.host)

	conn, err := tls.Dial("tcp", m.host+":465", &tls.Config{ServerName: m.host})
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}
	defer c.Close()

	// 1. authenticate
	auth := smtp.PlainAuth("", m.from, m.password, m.host)
	if err := c.Auth(auth); err != nil {
		return err
	}

	// 2. set sender / recipient
	if err := c.Mail(m.from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	// 3. send body
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err1 := w.Close(); err1 != nil && err == nil {
		err = err1
	}

	return err
}
