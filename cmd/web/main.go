package main

import (
	"log"
	"net/http"
	"os"
	"time"

	t "github.com/cpustejovsky/go_twitter_bot"
	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	creds      t.TwitterCredentials
	mgInstance *mailgun.MailgunImpl
}

func main() {
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	port := ":" + os.Getenv("PORT")

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)

	creds := t.TwitterCredentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	mgInstance := mg

	app := &application{
		creds:      creds,
		errorLog:   errorLog,
		infoLog:    infoLog,
		mgInstance: mgInstance,
	}

	srv := &http.Server{
		Addr:         port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
