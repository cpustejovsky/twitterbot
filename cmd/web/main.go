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
	tb, _ := t.NewBot()

	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan t.User)
	var u []t.User

	for _, name := range n {
		go tb.FindUserTweets(name, c)
		u = append(u, <-c)
	}

	if err := t.SendEmail(u); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}

func main() {
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/run-twitter-bot", handleSendEmail)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
