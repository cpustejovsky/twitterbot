# Go Twitter Bot

Twitter bot written in Go to give me less reasons to spend check on Twitter.

## Set Up
To set up, create a `.env` file and fill it with the following information:
```
TWITTER_CONSUMER_KEY="your Twitter consumer key"
TWITTER_CONSUMER_SECRET="your Twitter consumer secret"
TWITTER_ACCESS_TOKEN="your Twitter access token"
TWITTER_ACCESS_TOKEN_SECRET="your Twitter access token secret"
MG_API_KEY: "your MailGun private API key"
MG_DOMAIN: "your MailGun email domain"
PORT: "the port number you want to use on local"
```

Then set up your credentials and use them to create a new client.
```go
import (
  bot "github.com/cpustejovsky/twitterbot"
  "github.com/mailgun/mailgun-go/v4"
)

creds := bot.TwitterCredentials{
	AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
	AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
	ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
}

tc, err := bot.NewClient(creds)
if err != nil {
  log.Fatal(err)
}
```
## Use

### Collecting unfavorited (unliked) tweets

Pass in the Twitter API client and a MailGun API instance along with a slice of Twitter usernames, the number of tweets you want to check.

```go
usernames := []{"foo", "bar"}

ut := CollectUserTweets(tc, usernames, 5)
```

### Email Unread Tweets

Pass in the Twitter API client and a MailGun API instance along with a slice of Twitter usernames, the number of tweets you want to check, and a recipient email address.

**Example:**
```go

mg, err := mailgun.NewMailgunFromEnv()
if err != nil {
	log.Fatal(err)
}

err := bot.EmailUnreadTweets(tc, mg, []string{"FluffyHookers", "elpidophoros"}, 5, "charles.pustejovsky@gmail.com")
if err != nil {
  log.Fatal(err)
}
```

## To Dos
* ~~Set up as a web app on a Heroku dyno that [dyno-waker](https://github.com/cpustejovsky/dyno-waker) can hit daily.~~
* ~~Add tests~~
* ~~Pass in twitter usernames as parameters~~
* ~~Look into if it's possible for a twitter bot to like tweets on behalf of a user.~~
* Add information about unread notifications
* Add information about unread messages
* Allow users to connect their Twitter account to estuaryapp.com
* Provide functionality for other users
* Use `html/templates` for templating email body instead of string manipulation
