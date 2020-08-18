package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	}
	fmt.Fprintf(w, "Email is being sent")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/run-twitter-bot", handleSendEmail)
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
