package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ramonpizana/NVDIA/internal/vcard"
	"gopkg.in/mail.v2"
)

// mails over information
func MailInfo(pass string, information []*vcard.Vcard) error {
	s := make([]string, 0)

	if pass == "" {
		return errors.New("could not find password, received " + pass)
	}
	for i := 0; i < len(information); i++ {
		s = append(s, fmt.Sprintf("%+v \n %+v \n", strconv.Itoa(information[i].Price), information[i].Link))
	}
	m := mail.NewMessage()

	m.SetHeader("From", "ramonpizana@hotmail.com")
	m.SetHeader("To", "ramonpizana123gmail.com")
	m.SetHeader("Subject", "Found cards!")

	m.SetBody("text/html", strings.Join(s[:], "\n"))
	//verify stmp env var via gpass
	d := mail.NewDialer("smtp.gmail.com", 587, "ramonpizana123gmail.com", pass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
