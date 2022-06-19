package main

import (
	"path/filepath"
	"testing"
	"time"
)

type TestReminder struct {
	t              *testing.T
	countReminders int
}

func (r *TestReminder) send(formattedName string, birthday time.Time) error {
	r.t.Log(formattedName)
	r.countReminders++
	return nil
}

// TestEvaluateVCards checks if vcard evaluation is working
func TestEvaluateVCards(t *testing.T) {
	reminder = &TestReminder{t, 0}
	simDate := "0507"
	simulateDate = &simDate

	// walk all files in directory
	err := filepath.Walk("./test_csvs/test0507.vcf", evaluateVCards)
	if err != nil {
		t.Error(err)
	}

	if reminder.(*TestReminder).countReminders != 3 {
		t.Log("Found", reminder.(*TestReminder).countReminders, "reminders")
		t.Error("Found != 3 reminders!")
	}
}
