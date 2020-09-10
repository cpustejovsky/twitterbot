package main

import (
	"log"
	"net/http"
	"os"

	t "github.com/cpustejovsky/go_twitter_bot"
	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

func setUpMailGun() (*mailgun.MailgunImpl, error) {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		return mg, err
	}
	return mg, nil
}

func main() {
	//TODO: move all fatal logs to main function or at least the handlers?
	//TODO: log all errors to a /tmp/error.log file and all info to a /tmp/info.log
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	b := &Bot{}
	b.creds = t.Credentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	mg, err := setUpMailGun()
	if err != nil {
		log.Fatal(err)
	}
	b.mgInstance = mg

	http.HandleFunc("/run-twitter-bot", b.handleSendEmail)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
