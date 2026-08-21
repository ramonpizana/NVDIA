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

// MailInfo emails the available cards sorted by price.
func MailInfo(pass string, information []*vcard.Vcard) error {
	lines := make([]string, 0, len(information))
	for _, card := range information {
		lines = append(lines, fmt.Sprintf("%s — %s — $%s MXN\n%s", card.Store, card.Name, strconv.Itoa(card.Price), card.Link))
	}
	return sendMail(pass, "RTX 5090 disponibles en México", strings.Join(lines, "\n\n"))
}

// MailTest validates SMTP independently from retailer availability.
func MailTest(pass string) error {
	return sendMail(pass, "Prueba del monitor RTX 5090", "La configuración SMTP funciona correctamente. Las alertas de disponibilidad llegarán a este correo.")
}

func sendMail(pass, subject, body string) error {
	if pass == "" {
		return errors.New("SMTP_PASSWORD is not configured")
	}
	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO")
	username := os.Getenv("SMTP_USERNAME")
	if from == "" || to == "" || username == "" {
		return errors.New("EMAIL_FROM, EMAIL_TO and SMTP_USERNAME must be configured")
	}

	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := mail.NewDialer(envOrDefault("SMTP_HOST", "smtp.gmail.com"), envIntOrDefault("SMTP_PORT", 587), username, pass)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("SMTP: %w", err)
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
