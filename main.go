package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handleSendEmail(w http.ResponseWriter, r *http.Request) {
	client, err1 := getClient()
	if err1 != nil {
		fmt.Printf("Error getting Twitter Client:\n%v\n", err1)
		return
	}

	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan User)
	var u []User

	for _, name := range n {
		go findUserTweets(client, name, c)
		u = append(u, <-c)
	}

	err2 := sendEmail(u)
	if err2 != nil {
		fmt.Fprintf(w, "No email was sent.\n%v", err2)
	}
	fmt.Fprintf(w, "Email is being sent")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	enverr := godotenv.Load()
	if enverr != nil {
		fmt.Println(enverr)
	}
	http.HandleFunc("/run-twitter-bot", handleSendEmail)
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
