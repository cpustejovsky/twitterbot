package bot

import (
	"fmt"
	"sync"

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

	var wg sync.WaitGroup

	for _, name := range userNames {
		wg.Add(1)
		go tb.findUserTweets(&wg, name, count)
	}

	wg.Wait()

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
func (tb *TwitterBot) findUserTweets(wg *sync.WaitGroup, userName string, count int) {
	defer wg.Done()

	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      count,
		TweetMode:  "extended",
	}
	tweets, resp, err := tb.client.Timelines.UserTimeline(params)
	u := User{
		name: userName,
	}
	if err != nil && resp.StatusCode == 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		return
	}
	u = tb.modifyAndAddTweetsToUser(u, tweets)
	tb.users = append(tb.users, u)
}

func (tb *TwitterBot) modifyAndAddTweetsToUser(u User, tweets []twitter.Tweet) User {
	for _, tweet := range tweets {
		if tweet.Favorited == false {
			tb.likeTweet(tweet)
			ut := userTweet{
				text: tweet.FullText,
				id:   tweet.IDStr,
				link: fmt.Sprintf("https://twitter.com/%v/status/%v", u.name, tweet.IDStr),
			}
			u.tweets = append(u.tweets, ut)
		}
	}
	u.name = tweets[0].User.Name
	return u
}

func (tb *TwitterBot) likeTweet(tweet twitter.Tweet) {
	var p twitter.FavoriteCreateParams
	p.ID = tweet.ID
	tb.client.Favorites.Create(&p)
}
