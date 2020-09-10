package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	t "github.com/cpustejovsky/go_twitter_bot"
	"github.com/joho/godotenv"
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

func main() {
	//TODO: move all fatal logs to main function or at least the handlers?
	//TODO: log all errors to a /tmp/error.log file and all info to a /tmp/info.log
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/run-twitter-bot", handleSendEmail)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
