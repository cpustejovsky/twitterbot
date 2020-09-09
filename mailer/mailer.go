package twitterBot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	bot "github.com/cpustejovsky/go_twitter_bot/bot"
	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

type EmptyError struct{}

func (e *EmptyError) Error() string {
	return "no users to send email to."
}

func CheckUsers(u []bot.User) error {
	empties := 0
	for _, user := range u {
		if len(user.Tweets) == 0 {
			empties++
		}
	}
	if len(u) == empties {
		return &EmptyError{}
	}
	return nil
}

func SetUpMailGun() *mailgun.MailgunImpl {
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	return mg
}

func formatHtml(u []bot.User, m *mailgun.Message) {
	//TODO: make use of html/templates for templating
	var tweets bytes.Buffer
	tweets.WriteString("<h1>Daily Tweet Update</h1>")
	for _, user := range u {
		tweets.WriteString("<h3>Tweets from " + user.Name + "</h3><ul>")
		if len(user.Tweets) == 0 {
			tweets.WriteString("<li>No new tweets.</li>")
		}
		for _, tweet := range user.Tweets {
			tweets.WriteString("<li>" + tweet.Text + " <a target='_blank' rel='noopener noreferrer' href=" + tweet.Link + ">(link)</a></li>")
		}
		tweets.WriteString("</ul>")
	}

	m.SetHtml(tweets.String())
}

func SendEmail(u []bot.User) error {
	err := CheckUsers(u)
	if err != nil {
		return err
	}
	mg := SetUpMailGun()

	sender := "twitter-updates@estuaryapp.com"
	subject := "Twitter Updates"
	html := ""
	recipient := "charles.pustejovsky@gmail.com"

	m := mg.NewMessage(sender, subject, html, recipient)

	formatHtml(u, m)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MailGun API:\nID: %s\nResp: %s\n", id, resp)
	return nil
}
