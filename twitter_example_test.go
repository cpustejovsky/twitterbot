package twitterbot_test

import (
	"fmt"
	"os"

	bot "github.com/cpustejovsky/twitterbot"
	"github.com/mailgun/mailgun-go/v4"
)

func ExampleEmailUnreadTweets() {
	n := []string{"rob_pike", "golang"}
	var mgInstance *mailgun.MailgunImpl
	var EmailAddress = "your_email_address@example.com"
	creds := bot.TwitterCredentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	tc, err := bot.NewClient(creds)
	if err != nil {
		fmt.Println(err)
	}
	if err := bot.EmailUnreadTweets(tc, mgInstance, n, 5, EmailAddress); err != nil {
		fmt.Printf("No email was sent.\n%v", err)
	} else {
		fmt.Printf("Email is being sent")
	}
}
