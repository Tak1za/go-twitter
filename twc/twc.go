package twc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const (
	baseUrl = "https://api.twitter.com/1.1/"
)

type TWApi struct {
	HttpClient *http.Client
}

//AccessKeys datatype
type AccessKeys struct {
	Key    string
	Secret string
}

//GetTwitterClient to start making requests
func GetTwitterClient(keys AccessKeys) (*TWApi, error) {
	accessToken, err := getToken(keys)
	if err != nil {
		return nil, err
	}
	var conf oauth2.Config
	tClient := conf.Client(context.Background(), &accessToken)
	return &TWApi{
		HttpClient: tClient,
	}, nil
}

func (api TWApi) getRequest(url string, data interface{}) error {
	res, err := api.HttpClient.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = decodeResponse(res.Body, data)
	if err != nil {
		return err
	}

	return nil
}

func decodeResponse(item io.ReadCloser, data interface{}) error {
	dec := json.NewDecoder(item)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func getToken(keys AccessKeys) (oauth2.Token, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return oauth2.Token{}, err
	}

	req.SetBasicAuth(keys.Key, keys.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return oauth2.Token{}, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	dec.Decode(&token)
	if err != nil {
		return oauth2.Token{}, err
	}

	return token, nil
}
