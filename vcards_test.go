package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type TestReminder struct {
	t              *testing.T
	countReminders int
}

func (r *TestReminder) send(formattedName string, birthday time.Time, c Config) error {
	r.t.Log(formattedName)
	r.countReminders++
	return nil
}

// TestEvaluateVCards checks if vcard evaluation is working
func TestEvaluateVCards(t *testing.T) {
	var c Config

	c.reminder = &TestReminder{t, 0}
	simDate := "0507"
	c.simulateDate = &simDate

	// test test0507.vcf
	err := filepath.Walk("./test_csvs/test0507.vcf", c.evaluateVCards)
	if err != nil {
		t.Error(err)
	}

	if c.reminder.(*TestReminder).countReminders != 3 {
		t.Log("Found", c.reminder.(*TestReminder).countReminders, "reminders")
		t.Error("Found != 3 reminders!")
	}

	// test wrong path
	err = filepath.Walk("./test_csvs/notexistant.vcf", c.evaluateVCards)
	//if errors.Is(err, o) {
	if !strings.HasPrefix(err.Error(), "lstat") {
		t.Log(err)
		t.Error("err attribute is not checked correctly:", err)
	}

}
