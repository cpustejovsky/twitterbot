package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

func sendEmail() {
	// Create an instance of the Mailgun Client
	err := godotenv.Load()
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		fmt.Println(err)
	}
	sender := "sender@estuaryapp.com"
	subject := "Fancy subject!"
	body := "Hello from Mailgun Go!"
	recipient := "charles.pustejovsky@gmail.com"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
