package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

type keys struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

func main() {
	keys := getKeys()

	token := getToken(keys)

	tClient := getTwitterClient(token)
	res2, err := tClient.Get(fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", "1162771373426540545"))
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()

	io.Copy(os.Stdout, res2.Body)
}

func getKeys() keys {
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

func getToken(keys keys) oauth2.Token {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(keys.Key, keys.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	dec.Decode(&token)
	if err != nil {
		panic(err)
	}

	return token
}

func getTwitterClient(token oauth2.Token) *http.Client {
	var conf oauth2.Config
	tClient := conf.Client(context.Background(), &token)
	return tClient
}
