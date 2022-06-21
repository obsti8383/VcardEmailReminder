package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mapaiva/vcard-go"
)

// evaluateVCards can be used as a function for filepath.Walk() which does
// all the work of parsing VCF files and checking if a birthday date matches
func (c Config) evaluateVCards(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	// check if we should use current date or if the user wants us to simulate a date
	var now time.Time
	if *c.simulateDate != "" {
		now, _ = time.Parse("0102", *c.simulateDate)
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
					if c.debugLog {
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
							err = c.reminder.send(card.FormattedName, bdTime, c)
							if err != nil {
								log.Fatal("Error sending reminder Email!: ", err.Error())
							}
						} else {
							if c.debugLog {
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
