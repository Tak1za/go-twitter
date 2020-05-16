package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"os"

	"github.com/Tak1za/go-twitter/twc"
)

type keys struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

func main() {
	var (
		keyFile   string
		usersFile string
		tweetID   string
	)
	flag.StringVar(&keyFile, "key", "keys.json", "the file where we store consumer key and secret for twitter API")
	flag.StringVar(&usersFile, "users", "users.csv", "the file where users who have retweeted the tweet are stored. This will be created if it does not exist")
	flag.StringVar(&tweetID, "tweet", "", "The ID of the tweet you wish to find the retweeters of")
	flag.Parse()

	keys, err := getKeysFromJson(keyFile)
	if err != nil {
		panic(err)
	}

	accessKeys := twc.AccessKeys{
		Key:    keys.Key,
		Secret: keys.Secret,
	}

	tClient := twc.GetTwitterClient(accessKeys)
	retweets, err := tClient.GetRetweets("1261710171979513857")
	if err != nil {
		panic(err)
	}

	usernames := make([]string, 0, len(retweets))
	for _, retweet := range retweets {
		usernames = append(usernames, retweet.User.ScreenName)
	}

	existingUsernames := existing(usersFile)

	uniqueUsers := merge(usernames, existingUsernames)

	err = writeUsers(usersFile, uniqueUsers)
	if err != nil {
		panic(err)
	}
}

func writeUsers(usersFile string, users []string) error {
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)
	for _, user := range users {
		if err := w.Write([]string{user}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func existing(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}

	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, user := range a {
		uniq[user] = struct{}{}
	}

	for _, user := range b {
		uniq[user] = struct{}{}
	}

	ret := make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}

	return ret
}

func getKeysFromJson(keyFile string) (keys, error) {
	var keys keys
	f, err := os.Open(keyFile)
	if err != nil {
		return keys, err
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&keys)

	return keys, nil
}
