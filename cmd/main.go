package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ramonpizana/NVDIA/internal"
	"github.com/ramonpizana/NVDIA/search"
)

func main() {
	pass := os.Getenv("SMTP_PASSWORD")
	if pass == "" {
		// Compatibility with the previous setup. Prefer SMTP_PASSWORD going forward.
		pass = os.Getenv("PSWRD")
	}

	var err error
	if sendTestEmail() {
		fmt.Println("enviando correo de prueba")
		err = internal.MailTest(pass)
	} else {
		err = search.SearchNewEgg(pass)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func sendTestEmail() bool {
	enabled, err := strconv.ParseBool(os.Getenv("SEND_TEST_EMAIL"))
	return err == nil && enabled
}
