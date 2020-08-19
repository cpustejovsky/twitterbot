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

## To Dos
* ~~Set up as a web app on a Heroku dyno that [dyno-waker](https://github.com/cpustejovsky/dyno-waker) can hit daily.~~
* Add tests
* Add information about unread notifications
* Add information about unread messages
* Use `html/templates` for templating email body instead of string manipulation
* Pass in twitter usernames as parameters
* Look into if it's possible for a twitter bot to like tweets on behalf of a user.
  * If Possible:
    * Allow users to connect their Twitter account to estuaryapp.com
    * refactor the web app to accept a POST request from estuaryapp.com that contains the usernames to use along with the email address.