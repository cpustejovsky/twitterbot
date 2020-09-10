package main

import (
	"fmt"
	"net/http"

	t "github.com/cpustejovsky/go_twitter_bot"
)

func handleSendEmail(w http.ResponseWriter, r *http.Request) {
	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan t.User)
	tb, _ := t.NewBot()

	for _, name := range n {
		go tb.FindUserTweets(name, c)
		tb.AddUsers(c)
	}

	if err := tb.SendEmail(); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}
