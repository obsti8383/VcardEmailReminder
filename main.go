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

type Reminder interface {
	send(formattedName string, birthday time.Time) error
}

var (
	emailRecipient *string
	emailSender    *string
	smtpServer     *string
	smtpUsername   *string
	smtpPassword   *string
	simulateDate   *string
	debugLog       bool = false
	reminder       Reminder
)

func main() {
	// use email reminder
	reminder = &Email{}
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
