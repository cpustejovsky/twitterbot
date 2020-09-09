package twitterBot

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
	Text  string
	Id    string
	Link  string
	Liked bool
}

type User struct {
	Name   string
	Tweets []UserTweet
}

func greek(tweet string) bool {
	for _, char := range tweet {
		if char >= 945 && char <= 1023 {
			return true
		}
	}
	return false
}

func loadCreds() Credentials {
	if os.Getenv("PORT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Loading Credentials...")
	creds := Credentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	return creds
}

func GetClient() (*twitter.Client, error) {
	creds := loadCreds()
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func FindUserTweets(client *twitter.Client, userName string, c chan User) {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      5,
		TweetMode:  "extended",
	}
	tweets, resp, err := client.Timelines.UserTimeline(params)
	u := User{
		Name: userName,
	}
	if err != nil && resp.StatusCode == 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		c <- u
	}
	u.Name = tweets[0].User.Name
	for _, tweet := range tweets {
		greek := greek(tweet.FullText)
		if greek == false {
			if tweet.Favorited == false {
				var p twitter.FavoriteCreateParams
				p.ID = tweet.ID
				client.Favorites.Create(&p)
				ut := UserTweet{
					Text: tweet.FullText,
					Id:   tweet.IDStr,
					Link: fmt.Sprintf("https://twitter.com/%v/status/%v", userName, tweet.IDStr),
				}
				u.Tweets = append(u.Tweets, ut)
			}
		}
	}
	fmt.Printf("found %v unliked tweets from %v\n", len(u.Tweets), u.Name)
	c <- u
}
