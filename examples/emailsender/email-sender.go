package main

import (
	"fmt"
	"net/smtp"

	"github.com/abdullahnettoor/tqwp"
)

// EmailTask represents a task to send an email.
type EmailTask struct {
	tqwp.TaskModel
	Id      uint
	To      string
	Subject string
	Body    string
}

// Process sends an email using Go's net/smtp package.
func (t *EmailTask) Process() error {
	// Sender email address credentials.
	from := "your-email@example.com" // Update with your email
	password := "your-password"      // Update with your password

	// SMTP server configuration.
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	// Receiver email address.
	to := []string{t.To}

	// Message.
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", t.Subject, t.Body))

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %v", t.To, err)
	}

	fmt.Printf("Email sent to %s (Subject: %s)\n", t.To, t.Subject)
	return nil
}

func main() {
	emailList := []string{
		"john.doe@example.com",
		"jane.smith@example.com",
		"emily.johnson@example.com",
		"michael.brown@example.com",
		"sarah.wilson@example.com",
		"david.jones@example.com",
		"laura.miller@example.com",
		"chris.taylor@example.com",
		"katie.moore@example.com",
		"brian.white@example.com",
	}
	var numOfWorkers, maxRetries uint = 5, 3

	// Create and start the worker pool.
	wp := tqwp.New(&tqwp.WorkerPoolConfig{
		NumOfWorkers: numOfWorkers,
		MaxRetries:   maxRetries})
	defer wp.Summary()
	defer wp.Stop()

	wp.Start()
	for i, email := range emailList {
		t := EmailTask{
			Id:        uint(i + 1),
			To:        email,
			Subject:   "Hello!",
			Body:      "This is a test email.",
			TaskModel: tqwp.TaskModel{},
		}
		wp.EnqueueTask(&t)
	}
}
