package bot

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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
	u = tb.modifyAndAddTweetsToUser(u, tweets)
	c <- u
}

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
