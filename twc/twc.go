package twc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

//AccessKeys datatype
type AccessKeys struct {
	Key    string
	Secret string
}

//GetTwitterClient to start making requests
func GetTwitterClient(keys AccessKeys) *http.Client {
	accessToken := getToken(keys)
	var conf oauth2.Config
	tClient := conf.Client(context.Background(), &accessToken)
	return tClient
}

func getToken(keys AccessKeys) oauth2.Token {
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
