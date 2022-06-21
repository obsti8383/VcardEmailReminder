// VcardEmailReminder
// Copyright (C) 2018-2022 Florian Probst
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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Reminder interface is used to allow sending notifications using
// different mechanisms
type Reminder interface {
	send(formattedName string, birthday time.Time, c Config) error
}

// Config is a struct that contains all configuration variables
type Config struct {
	emailRecipient *string
	emailSender    *string
	smtpServer     *string
	smtpUsername   *string
	smtpPassword   *string
	simulateDate   *string
	debugLog       bool
	reminder       Reminder
}

func main() {
	var c Config
	// default reminder, will be overwritten if test flag is not used
	c.reminder = &printlnReminder{}
	// parse command line parameters/flags
	path := flag.String("path", "", "path where the vcf files reside (or vcf file directly) (required)")
	test := flag.Bool("test", false, "test only, does not send email (optional)")
	c.emailRecipient = flag.String("recipient", "", "recipients email address (required)")
	c.emailSender = flag.String("sender", "", "senders email address (required)")
	c.smtpServer = flag.String("smtp", "", "smtp server adress, e.g. \"smtp.variomedia.de:25\" (required)")
	c.smtpUsername = flag.String("username", "", "username for smtp server (required)")
	c.smtpPassword = flag.String("password", "", "password for smtp server (required)")
	c.simulateDate = flag.String("simulateDate", "", "simulate date string, e.g. \"0716\" for the 16th of July (optional)")
	flag.Parse()
	if !*test {
		// use email reminder
		c.reminder = &EmailReminder{}

		// check for required flags for email reminder
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name != "test" && f.Name != "simulateDate" && f.Value.String() == "" {
				fmt.Println("Error: Required parameter", f.Name, "is missing.\n\nThe following parameters are available:")
				flag.PrintDefaults()
				os.Exit(1)
			} else {
				if c.debugLog {
					log.Printf("%s = \"%s\"", f.Name, f.Value.String())
				}
			}
		})
	}

	// walk all files in directory
	err := filepath.Walk(*path, c.evaluateVCards)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type printlnReminder struct{}

func (r *printlnReminder) send(formattedName string, birthday time.Time, c Config) error {
	fmt.Println(formattedName + " birthday is on " + birthday.Format("Jan 2") + "!")
	return nil
}
