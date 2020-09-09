package bot

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
	text  string
	id    string
	link  string
	liked bool
}

type User struct {
	name   string
	tweets []UserTweet
}

type TwitterBot struct {
	client *twitter.Client
}

func (tb *TwitterBot) newBot() error {
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
		log.Fatal(err)
		return err
	}
	tb.client = client
	return nil
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
	creds := Credentials{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	return creds
}

func (tb *TwitterBot) FindUserTweets(userName string, c chan User) {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      5,
		TweetMode:  "extended",
	}
	tweets, resp, err := tb.client.Timelines.UserTimeline(params)
	u := User{
		name: userName,
	}
	if err != nil && resp.StatusCode == 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		c <- u
	}
	u = tb.ModifyAndAddTweetsToUser(u, tweets)
	c <- u
}

func (tb *TwitterBot) ModifyAndAddTweetsToUser(u User, tweets []twitter.Tweet) User {
	for _, tweet := range tweets {
		prevLiked := true
		greek := greek(tweet.FullText)
		if greek == false {
			if tweet.Favorited == false {
				tb.LikeTweet(tweet)
				prevLiked = true
			}
			ut := UserTweet{
				text:  tweet.FullText,
				id:    tweet.IDStr,
				link:  fmt.Sprintf("https://twitter.com/%v/status/%v", u.name, tweet.IDStr),
				liked: prevLiked,
			}
			u.tweets = append(u.tweets, ut)
		}
	}
	u.name = tweets[0].User.Name
	return u
}

func (tb *TwitterBot) LikeTweet(tweet twitter.Tweet) {
	var p twitter.FavoriteCreateParams
	p.ID = tweet.ID
	tb.client.Favorites.Create(&p)
}
