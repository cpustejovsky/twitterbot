package main

import (
	"fmt"
	"net/http"

	bot "github.com/cpustejovsky/go_twitter_bot/bot"
	mailer "github.com/cpustejovsky/go_twitter_bot/mailer"
)

func handleSendEmail(w http.ResponseWriter, r *http.Request) {
	client, err := bot.GetClient()
	if err != nil {
		fmt.Printf("Error getting Twitter Client:\n%v\n", err)
		return
	}

	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan bot.User)
	var u []bot.User

	for _, name := range n {
		go bot.FindUserTweets(client, name, c)
		u = append(u, <-c)
	}

	if err := mailer.SendEmail(u); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}
