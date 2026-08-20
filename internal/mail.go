package internal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ramonpizana/NVDIA/internal/vcard"
	"gopkg.in/mail.v2"
)

// mails over information
func MailInfo(pass string, information []*vcard.Vcard) error {
	s := make([]string, 0)

	if pass == "" {
		return errors.New("SMTP_PASSWORD is not configured")
	}
	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO")
	username := os.Getenv("SMTP_USERNAME")
	if from == "" || to == "" || username == "" {
		return errors.New("EMAIL_FROM, EMAIL_TO and SMTP_USERNAME must be configured")
	}
	for i := 0; i < len(information); i++ {
		s = append(s, fmt.Sprintf("%s — %s — $%s MXN\n%s\n", information[i].Store, information[i].Name, strconv.Itoa(information[i].Price), information[i].Link))
	}
	m := mail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "RTX 5090 disponibles en México")

	m.SetBody("text/html", strings.Join(s[:], "\n"))
	//verify stmp env var via gpass
	d := mail.NewDialer(envOrDefault("SMTP_HOST", "smtp.gmail.com"), envIntOrDefault("SMTP_PORT", 587), username, pass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil && value > 0 {
		return value
	}
	return fallback
}
