package main

import (
	"fmt"
	"log"
	"net/http"

	t "github.com/cpustejovsky/go_twitter_bot"
)

func (app *application) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan t.User)
	tb, err := t.NewBot(app.creds)
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range n {
		go tb.FindUserTweets(name, c)
		tb.AddUsers(c)
	}

	if err := tb.SendEmail(app.mgInstance); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}
