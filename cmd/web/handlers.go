package main

import (
	"fmt"
	"log"
	"net/http"

	t "github.com/cpustejovsky/go_twitter_bot"
	"github.com/mailgun/mailgun-go/v4"
)

type Bot struct {
	creds      t.Credentials
	mgInstance *mailgun.MailgunImpl
}

func (b *Bot) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan t.User)
	tb, err := t.NewBot(b.creds)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range n {
		go tb.FindUserTweets(name, c)
		tb.AddUsers(c)
	}

	if err := tb.SendEmail(b.mgInstance); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}
