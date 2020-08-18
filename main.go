package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handleSendEmail(w http.ResponseWriter, r *http.Request) {
	client, err := getClient()
	if err != nil {
		fmt.Printf("Error getting Twitter Client:\n%v\n", err)
		return
	}

	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan User)
	var u []User

	for _, name := range n {
		go findUserTweets(client, name, c)
		u = append(u, <-c)
	}

	if err := sendEmail(u); err != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err)
	} else {
		fmt.Fprintf(w, "Email is being sent")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my Go Twitter Bot!")
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			fmt.Println(err)
		}
	}
	http.HandleFunc("/run-twitter-bot", handleSendEmail)
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
