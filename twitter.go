package bot

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mailgun/mailgun-go/v4"
)

type userTweet struct {
	text  string
	id    string
	link  string
	liked bool
}

type User struct {
	name   string
	tweets []userTweet
}

type TwitterBot struct {
	client *twitter.Client
	users  []User
}

type TwitterCredentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//EmailUnreadTweets takes Twitter API credentials, a MailGun implementation, a slice of Twitter usernames, and a count of how many tweets to check and sends emails of unread tweets to the recipient's email address
func EmailUnreadTweets(creds TwitterCredentials, mg *mailgun.MailgunImpl, userNames []string, count int, recipient string) error {
	tb, err := newBot(creds)
	if err != nil {
		return err
	}

	c := make(chan User, len(userNames))
	for _, name := range userNames {
		go func() { c <- findUserTweets(tb.client, name, count) }()
	}
	for i := 0; i < cap(c); i++ {
		user := <-c
		tb.users = append(tb.users, user)
	}
	if err := tb.SendEmail(mg, recipient); err != nil {
		return err
	}
	return nil
}

//NewBot creates a twitter bot based on Twitter API credentials
func newBot(creds TwitterCredentials) (TwitterBot, error) {
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
		return TwitterBot{}, err
	}
	tb := &TwitterBot{}
	tb.client = client
	return *tb, nil
}

//FindUserTweets takes finds count tweets for userName and passes a User struct to channel
func findUserTweets(t *twitter.Client, userName string, count int) User {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      count,
		TweetMode:  "extended",
	}
	tweets, resp, err := t.Timelines.UserTimeline(params)
	u := User{
		name: userName,
	}
	if err != nil && resp.StatusCode == 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		return u
	}
	u = modifyAndAddTweetsToUser(t, u, tweets)
	return u
}

func modifyAndAddTweetsToUser(t *twitter.Client, u User, tweets []twitter.Tweet) User {
	for _, tweet := range tweets {
		if tweet.Favorited == false {
			ok := likeTweet(t, tweet)
			ut := userTweet{
				text:  tweet.FullText,
				id:    tweet.IDStr,
				link:  fmt.Sprintf("https://twitter.com/%v/status/%v", u.name, tweet.IDStr),
				liked: ok,
			}
			u.tweets = append(u.tweets, ut)
		}
	}
	u.name = tweets[0].User.Name
	return u
}

//likeTweet uses the Twitter API to like a tweet. If there was an error, it returns false, indicating that there was a problem liking the tweet
func likeTweet(t *twitter.Client, tweet twitter.Tweet) bool {
	var p twitter.FavoriteCreateParams
	p.ID = tweet.ID
	_, rc, err := t.Favorites.Create(&p)
	if rc.StatusCode != 200 || err != nil {
		fmt.Println("Status Code: ", rc.StatusCode)
		fmt.Println("Error:\n", err)
		return false
	}
	return true
}
