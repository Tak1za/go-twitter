package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	res2, err := tClient.Get(fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", "1162771373426540545"))
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()

	io.Copy(os.Stdout, res2.Body)
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
