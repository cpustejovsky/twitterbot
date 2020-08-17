package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type UserTweet struct {
	text string
	id   string
	link string
}

type User struct {
	name   string
	tweets []UserTweet
}

func notGreek(tweet string) bool {
	notGreek := true
	for _, char := range tweet {
		if char >= 945 && char <= 1023 {
			notGreek = false
		}
	}
	return notGreek
}

func loadCreds() (Credentials, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loading Credentials...")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}
	return creds, err
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("\nGot %+v's Twitter Account\n", user.ScreenName)
	return client, nil
}

func findUserTweets(client *twitter.Client, userName string, c chan User) {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      5,
		TweetMode:  "extended",
	}
	tweets, resp, err := client.Timelines.UserTimeline(params)
	u := User{
		name: userName,
	}
	if err != nil && resp.StatusCode == 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		c <- u
	}
	u.name = tweets[0].User.Name
	for _, tweet := range tweets {
		notGreek := notGreek(tweet.FullText)
		if notGreek == true {
			ut := UserTweet{
				text: tweet.FullText,
				id:   tweet.IDStr,
				link: fmt.Sprintf("https://twitter.com/%v/status/%v", userName, tweet.IDStr),
			}
			u.tweets = append(u.tweets, ut)
		}
	}
	c <- u
}