package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type Email struct{}

func (e *Email) send(formattedName string, birthday time.Time) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		*smtpUsername,
		*smtpPassword,
		strings.Split(*smtpServer, ":")[0],
	)

	from := mail.Address{Name: "Birthday Reminder", Address: *emailSender}
	to := mail.Address{Name: "", Address: *emailRecipient}
	title := formattedName + " birthday is on " + birthday.Format("Jan 2")

	body := formattedName + " birthday is on " + birthday.Format("Jan 2") + "!"
	if birthday.Year() != 1 {
		var age int
		age = time.Now().Year() - birthday.Year()
		log.Printf("Age: %d", age)
		body = body + fmt.Sprintf("\r\nHe/she gets %d years old.\r\n", age)
	}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		*smtpServer,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	encoded := mail.Address{Name: String, Address: ""}
	return strings.Trim(encoded.String(), "\" <@>")
}
