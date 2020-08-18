package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

type EmptyError struct{}

func (e *EmptyError) Error() string {
	return "no users to send email to."
}

func checkUsers(u []User) error {
	empties := 0
	for _, user := range u {
		if len(user.tweets) == 0 {
			empties++
		}
	}
	if len(u) == empties {
		return &EmptyError{}
	}
	return nil
}

func setUpMailGun() *mailgun.MailgunImpl {
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	return mg
}

func formatHtml(u []User, m *mailgun.Message) {
	//TODO: make use of html/templates for templating
	var tweets bytes.Buffer
	tweets.WriteString("<h1>Daily Tweet Update</h1>")
	for _, user := range u {
		tweets.WriteString("<h3>Tweets from " + user.name + "</h3><ul>")
		if len(user.tweets) == 0 {
			tweets.WriteString("<li>No new tweets.</li>")
		}
		for _, tweet := range user.tweets {
			tweets.WriteString("<li>" + tweet.text + " <a target='_blank' rel='noopener noreferrer' href=" + tweet.link + ">(link)</a></li>")
		}
		tweets.WriteString("</ul>")
	}

	m.SetHtml(tweets.String())
}

func sendEmail(u []User) error {
	err := checkUsers(u)
	if err != nil {
		return err
	}
	mg := setUpMailGun()

	sender := "twitter-updates@estuaryapp.com"
	subject := "Twitter Updates"
	html := ""
	recipient := "charles.pustejovsky@gmail.com"

	m := mg.NewMessage(sender, subject, html, recipient)

	formatHtml(u, m)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MailGun API:\nID: %s\nResp: %s\n", id, resp)
	return nil
}
