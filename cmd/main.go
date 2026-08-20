package main

import (
	"fmt"
	"os"

	"github.com/ramonpizana/NVDIA/search"
	// "github.com/gocolly/colly"
	// "net/smtp"
)

func main() {

	pass := os.Getenv("SMTP_PASSWORD")
	if pass == "" {
		// Compatibility with the previous setup. Prefer SMTP_PASSWORD going forward.
		pass = os.Getenv("PSWRD")
	}

	err := search.SearchNewEgg(pass)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
