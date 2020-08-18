package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

func sendEmail(u []User) {
	empties := 0
	for _, user := range u {
		fmt.Println(len(user.tweets))
		if len(user.tweets) == 0 {
			empties++
		}
	}
	fmt.Println(len(u))
	if len(u) == empties {
		log.Fatal("no users to send email to.")
	}

	// Create an instance of the Mailgun Client
	enverr := godotenv.Load()
	if enverr != nil {
		fmt.Println(enverr)
	}
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		fmt.Println(err)
	}
	sender := "twitter-updates@estuaryapp.com"
	subject := "Twitter Updates"
	html := ""
	recipient := "charles.pustejovsky@gmail.com"

	m := mg.NewMessage(sender, subject, html, recipient)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: make use of html/templates for templating
	var tweets bytes.Buffer
	tweets.WriteString("<h1>Daily Tweet Update</h1>")
	for _, user := range u {
		tweets.WriteString("<h3>Tweets from" + user.name + "</h3><ul>")
		if len(user.tweets) == 0 {
			tweets.WriteString("<li>No new tweets.</li>")
		}
		for _, tweet := range user.tweets {
			tweets.WriteString("<li>" + tweet.text + "\n<a target='_blank' rel='noopener noreferrer' href=" + tweet.link + ">(link)</a></li>")
		}
		tweets.WriteString("</ul>")
	}

	m.SetHtml(tweets.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MailGun API:\nID: %s\nResp: %s\n", id, resp)
}
