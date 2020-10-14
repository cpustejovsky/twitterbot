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

//NewBot creates a twitter bot based on Twitter API credentials
func NewBot(creds TwitterCredentials) (TwitterBot, error) {
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

// creds := t.TwitterCredentials{
// 	AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
// 	AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
// 	ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
// 	ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
// }
// mg, err := mailgun.NewMailgunFromEnv()
// if err != nil {
// 	log.Fatal(err)
// }

//EmailUnreadTweets takes Twitter API credentials, a MailGun implementation, a slice of Twitter usernames, and a count of how many tweets to check and sends emails of unread tweets to the recipient's email address
func EmailUnreadTweets(creds TwitterCredentials, mg *mailgun.MailgunImpl, userNames []string, count int, recipient string) error {
	n := []string{"FluffyHookers", "elpidophoros"}
	c := make(chan User)
	tb, err := NewBot(creds)
	if err != nil {
		return err
	}
	for _, name := range n {
		go tb.FindUserTweets(name, c, count)
		tb.AddUsers(c)
	}

	if err := tb.SendEmail(mg, recipient); err != nil {
		return err
	} else {
		return nil
	}
}

//FindUserTweets takes finds count tweets for userName and passes a User struct to channel
func (tb *TwitterBot) FindUserTweets(userName string, c chan User, count int) {
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
		c <- u
	}
	u = tb.modifyAndAddTweetsToUser(u, tweets)
	c <- u
}

//AddUsers takes a user channel and appends it to the twitter bots users slice
func (tb *TwitterBot) AddUsers(c chan User) {
	tb.users = append(tb.users, <-c)
}

func greek(tweet string) bool {
	for _, char := range tweet {
		if char >= 945 && char <= 1023 {
			return true
		}
	}
	return false
}

func (tb *TwitterBot) modifyAndAddTweetsToUser(u User, tweets []twitter.Tweet) User {
	for _, tweet := range tweets {
		prevLiked := true
		greek := greek(tweet.FullText)
		if greek == false {
			if tweet.Favorited == false {
				tb.likeTweet(tweet)
				prevLiked = true
			}
			ut := userTweet{
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

func (tb *TwitterBot) likeTweet(tweet twitter.Tweet) {
	var p twitter.FavoriteCreateParams
	p.ID = tweet.ID
	tb.client.Favorites.Create(&p)
}
