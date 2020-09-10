package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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
