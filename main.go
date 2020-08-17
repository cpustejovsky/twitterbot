package main

import (
	"fmt"
)

func main() {
	creds, err := loadCreds()
	client, err := getClient(&creds)
	if err != nil {
		fmt.Printf("Error getting Twitter Client:\n%v\n", err)
		return
	}
	
	n := []string{"yechielk", "FluffyHookers", "elpidophoros"}
	c := make(chan User)
	var u []User

	for _, name := range n {
		go findUserTweets(client, name, c)
		u = append(u, <-c)
	}
	sendEmail(u)
}
