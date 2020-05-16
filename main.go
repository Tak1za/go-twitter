package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Tak1za/go-twitter/twc"
)

type keys struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

func main() {
	keys := getKeysFromJson()

	accessKeys := twc.AccessKeys{
		Key:    keys.Key,
		Secret: keys.Secret,
	}

	tClient := twc.GetTwitterClient(accessKeys)
	res, err := tClient.GetRetweets("1261544718133071874")
	if err != nil {
		panic(err)
	}

	var usernames []string
	for _, retweet := range res {
		usernames = append(usernames, retweet.User.ScreenName)
	}

	fmt.Println(usernames)
}

func getKeysFromJson() keys {
	var keys keys
	f, err := os.Open("keys.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&keys)

	return keys
}
