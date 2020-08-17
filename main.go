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
	c := make(chan User)
	go findUserTweets(client, "FluffyHookers", c)
	go findUserTweets(client, "elpidophoros", c)
	fh, el := <-c, <-c
	fmt.Println("==============================================================")
	fmt.Println(fh)
	fmt.Println("==============================================================")
	fmt.Println(el)
}
