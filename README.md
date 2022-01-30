# Introduction

![Go Vet and Lint Status](https://github.com/obsti8383/VcardEmailReminder/actions/workflows/go_lint_vet_and_testBuild.yml/badge.svg)

VcardEmailReminder parses one or multiple addressbook VCARD files (.vcf) and 
checks the Birthday field ("BDAY") to verify if its the persons birthday today.
If yes a email is sent as an reminder.

Best to be used with vdirsyncer to download the VCARD files from a CardDav server
and scan this local files with VcardEmailReminder and put both programms in a CronJob.

References:
- vdirsyncer: https://github.com/pimutils/vdirsyncer 

# Command line parameters
The following parameters are available:

  --password : password for smtp server (required)
		
  --path : path where the vcf files reside (or vcf file directly) (required)
		
  --recipient : recipients email address (required)
		
  --sender : senders email address (required)
		
  --simulateDate : simulate date string, e.g. "0716" for the 16th of July (optional)
		
  --smtp : smtp server adress, e.g. "smtp.variomedia.de:25" (required)
		
  --username : username for smtp server (required)



# Example command line
./VcardEmailReminder --password "123456" --path ~/.contacts/ --recipient recipient@test.test --sender bdreminder@test.test --smtp smtp.testprovider.com:25 --username username@test.test --simulateDate "1218"

# Dependencies
go get -u github.com/mapaiva/vcard-go

