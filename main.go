// VcardEmailReminder
// Copyright (C) 2018 Florian Probst
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mapaiva/vcard-go"
)

var (
	emailRecipient *string
	emailSender    *string
	smtpServer     *string
	smtpUsername   *string
	smtpPassword   *string
	simulateDate   *string
	debugLog       bool = false
)

func main() {
	// parse command line parameters/flags
	Path := flag.String("path", "", "path where the vcf files reside (or vcf file directly) (required)")
	emailRecipient = flag.String("recipient", "", "recipients email address (required)")
	emailSender = flag.String("sender", "", "senders email address (required)")
	smtpServer = flag.String("smtp", "", "smtp server adress, e.g. \"smtp.variomedia.de:25\" (required)")
	smtpUsername = flag.String("username", "", "username for smtp server (required)")
	smtpPassword = flag.String("password", "", "password for smtp server (required)")
	simulateDate = flag.String("simulateDate", "", "simulate date string, e.g. \"0716\" for the 16th of July (optional)")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name != "simulateDate" && f.Value.String() == "" {
			fmt.Println("Required parameter", f.Name, "is missing. Aborting.\nThe following parameters are available:")
			flag.PrintDefaults()
			os.Exit(1)
		} else {
			if debugLog {
				log.Printf("%s = \"%s\"", f.Name, f.Value.String())
			}
		}
	})

	// walk all files in directory
	err := filepath.Walk(*Path, evaluateVCards)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// function for filepath.Walk() which does all the work of parsing VCF files
// and checking if a birthday date matches
func evaluateVCards(path string, info os.FileInfo, err error) error {
	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	// check if we should use current date or if the user wants us to simulate a date
	var now time.Time
	if *simulateDate != "" {
		now, _ = time.Parse("0102", *simulateDate)
	} else {
		now = time.Now().Local()
	}

	if !info.IsDir() {
		// parse file for VCARDs
		cards, err := vcard.GetVCards(path)

		if err != nil {
			log.Println(err)
			return err
		}

		// verify if any card was in the file
		if len(cards) == 0 {
			log.Println("No vcard found in file ", path)
			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				log.Print(err)
			}
			// print the fileContent for debugging purposes
			str := string(fileContent) // convert content to a 'string'
			log.Println(str)           // print the content as a 'string'
		}

		// iterate over all found cards and check if birthday == now
		for _, card := range cards {
			if card == (vcard.VCard{}) {
				log.Println("VCard seems to be empty")
			} else {
				bd := card.BirthDay
				if bd != "" {
					if debugLog {
						log.Println(card.FormattedName, "BirthDay: ", bd)
					}

					// check the different date formats which are used in VCARDs
					bdTime, err := time.Parse("20060102", bd)
					if err != nil {
						bdTime, err = time.Parse("2006-01-02", bd)
						if err != nil {
							if strings.HasPrefix(bd, "--") {
								// year of birth unknown
								bd = strings.TrimPrefix(bd, "--")
								bd = "0001" + bd
								bdTime, err = time.Parse("20060102", bd)
								if err != nil {
									log.Println("Could not parse birthday date with suffix -- correctly: ", bd)
								}
							} else {
								log.Println(card.FormattedName, ": BirthDay has unknown format: ", bd)
							}
						}
					}

					// if we have found a birthday date, then check if birthday == now
					if err != nil || !bdTime.IsZero() {
						if bdTime.Month() == now.Month() && bdTime.Day() == now.Day() {
							log.Println("Today", card.FormattedName, "has his/her birthday")
							err = sendEmailReminder(card.FormattedName, bdTime)
							if err != nil {
								log.Fatal("Error sending reminder Email!: ", err.Error())
							}
						} else {
							if debugLog {
								log.Println(card.FormattedName, "hasn't birthday today (but on", bdTime.Month(), bdTime.Day(), ")")
							}
						}
					} else {
						log.Println(card.FormattedName, ": Could not evaluate birthday")
					}
				}

			}
		}

	}

	return nil
}

func sendEmailReminder(formattedName string, birthday time.Time) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		*smtpUsername,
		*smtpPassword,
		strings.Split(*smtpServer, ":")[0],
	)

	from := mail.Address{"Birthday Reminder", *emailSender}
	to := mail.Address{"", *emailRecipient}
	title := formattedName + " has birthday today"

	body := formattedName + " has birthday on " + birthday.Format("Jan 2") + "!"
	if birthday.Year() != 1 {
		var age int
		age = time.Now().Year() - birthday.Year()
		log.Printf("Age: %d", age)
		body = body + fmt.Sprintf("\r\nHe/she is then %d years old.\r\n", age)
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
	encoded := mail.Address{String, ""}
	return strings.Trim(encoded.String(), " <@>")
}
