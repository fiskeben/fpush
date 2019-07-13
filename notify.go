package main

import (
	"fmt"
	"github.com/gregdel/pushover"
	"log"
	"os"
)

func sendPushNotification(filename string) error {
	p := pushover.New(os.Getenv("PUSHOVER_KEY"))
	r := pushover.NewRecipient(os.Getenv("PUSHOVER_RECIPIENT_KEY"))
	m := pushover.NewMessageWithTitle("Movement!", "Movement has been detected.")

	if err := addAttachment(m, filename); err != nil {
		log.Println(err.Error())
	}

	res, err := p.SendMessage(m, r)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %v", err)
	}
	verboseLog(res.String())
	return nil
}

func addAttachment(m *pushover.Message, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %v", filename, err)
	}
	if err := m.AddAttachment(file); err != nil {
		return fmt.Errorf("failed to add attachment '%s', skipping: %v", filename, err)
	}
	return nil
}
