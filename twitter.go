package twitterbot

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mailgun/mailgun-go/v4"
)

type Liked struct {
	success bool
	msg     string
}

type userTweet struct {
	text  string
	id    string
	link  string
	liked Liked
}

type User struct {
	name   string
	tweets []userTweet
}

type TwitterCredentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//EmailUnreadTweets takes Twitter API client, a MailGun implementation, a slice of Twitter usernames, and a count of how many tweets to check and sends emails of unread tweets to the recipient's email address
func EmailUnreadTweets(tc *twitter.Client, mg *mailgun.MailgunImpl, userNames []string, count int, recipient string) error {
	users := CollectUserTweets(tc, userNames, count)

	if err := SendEmail(mg, recipient, users); err != nil {
		return err
	}
	return nil
}

//NewClient creates a Twitter client based on Twitter API credentials
func NewClient(creds TwitterCredentials) (*twitter.Client, error) {
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

/*CollectUserTweets takes a Twitter API Client, a slice of Twitter usernames, and a count;
It returns a slice of those Twitter usernames with their tweets that you account has not favorited*/
func CollectUserTweets(tc *twitter.Client, userNames []string, count int) []User {
	c := make(chan User, len(userNames))
	go func() {
		defer close(c)
		for _, name := range userNames {
			c <- findUserTweets(tc, name, count)
		}
	}()
	var users []User
	for user := range c {
		if len(user.tweets) > 0 {
			users = append(users, user)
		}
	}
	return users
}

func CollectUserTweetsV2(tc *twitter.Client, userNames []string, count int) []User {
	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)
	c := make(chan User, g)
	defer close(c)

	for i := 0; i < g; i++ {
		defer wg.Done()
		go func() {
			for user := range c {
				c <- findUserTweets(tc, user.name, count)
			}
		}()
	}

	for _, userName := range userNames {
		var u = User{
			name:   userName,
			tweets: []userTweet{},
		}
		c <- u
	}
	wg.Wait()
	var users []User
	for user := range c {
		if len(user.tweets) > 0 {
			users = append(users, user)
		}
	}
	return users
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
	if len(tweets) > 0 {
		u = modifyAndAddTweetsToUser(t, u, tweets)
	}
	return u
}

func modifyAndAddTweetsToUser(t *twitter.Client, u User, tweets []twitter.Tweet) User {
	for _, tweet := range tweets {
		if tweet.Favorited == false {
			ok, msg := likeTweet(t, tweet)
			u.tweets = append(u.tweets, userTweet{
				text: tweet.FullText,
				id:   tweet.IDStr,
				link: fmt.Sprintf("https://twitter.com/%v/status/%v", u.name, tweet.IDStr),
				liked: Liked{
					success: ok,
					msg:     msg,
				},
			})
		}
	}
	u.name = tweets[0].User.Name
	return u
}

//likeTweet uses the Twitter API to like a tweet. If there was an error, it returns false, indicating that there was a problem liking the tweet
func likeTweet(t *twitter.Client, tweet twitter.Tweet) (bool, string) {
	var p twitter.FavoriteCreateParams
	p.ID = tweet.ID
	_, rc, err := t.Favorites.Create(&p)
	if rc.StatusCode != 200 || err != nil {
		fmt.Println("Status Code: ", rc.StatusCode)
		fmt.Println("Error:\n", err)
		return false, fmt.Sprintf("Status Code: %v\n Error Message: %v\n", rc.StatusCode, err)
	}
	return true, "success"
}

//UnlikeUsersTweets takes finds count tweets for userName and passes a User struct to channel
func UnlikeUsersTweets(tc *twitter.Client, userNames []string, count int) []string {
	c := make(chan string, len(userNames))
	go func() {
		defer close(c)
		for _, name := range userNames {
			c <- unlikeUserTweets(tc, name, count)
		}
	}()
	var messages []string
	for message := range c {
		messages = append(messages, message)
	}
	return messages
}

//unlikeUserTweets takes finds count tweets for userName and passes a User struct to channel
func unlikeUserTweets(t *twitter.Client, userName string, count int) string {
	params := &twitter.UserTimelineParams{
		ScreenName: userName,
		Count:      count,
		TweetMode:  "extended",
	}
	tweets, resp, err := t.Timelines.UserTimeline(params)
	if err != nil && resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		fmt.Println(err)
		return fmt.Sprintf("Status Code: %v\n Error Message: %v\n", resp.StatusCode, err)
	}
	if len(tweets) > 0 {
		err, msg := unlikeTweets(t, tweets)
		if err != nil {
			return msg
		}
		return msg
	}
	return "no tweets found"
}

//UnlikeTweet uses the Twitter API to remove a like for a tweet. If there was an error, it returns false, indicating that there was a problem liking the tweet
func unlikeTweets(t *twitter.Client, tweets []twitter.Tweet) (error, string) {
	var p twitter.FavoriteDestroyParams
	for _, tweet := range tweets {
		p.ID = tweet.ID
		_, rc, err := t.Favorites.Destroy(&p)
		if rc.StatusCode != 200 || err != nil {
			fmt.Println("Status Code: ", rc.StatusCode)
			fmt.Println("Error:\n", err)
			return err, fmt.Sprintf("Status Code: %v\n Error Message: %v\n", rc.StatusCode, err)
		}
	}
	return nil, "success"
}
