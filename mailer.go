package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

func sendEmail(u []User) {
	// Create an instance of the Mailgun Client
	err := godotenv.Load()
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		fmt.Println(err)
	}
	sender := "twitter-updates@estuaryapp.com"
	subject := "Twitter Updates"
	body := "" //use template
	recipient := "charles.pustejovsky@gmail.com"

	// Create a new template

	// The message object allows you to add attachments and Bcc recipients
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	m := mg.NewMessage(sender, subject, body, recipient)
	m.SetTemplate("my-template")

	// Add the variables to be used by the template
	m.AddVariable("title", "Testing How Users Look")
	m.AddVariable("body", fmt.Sprintf("<ul><li>%v</li><li>%v</li></ul>", u[0].name, u[1].name))

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
