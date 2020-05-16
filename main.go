package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Tak1za/go-twitter/twc"
)

type keys struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

func main() {
	var (
		keyFile    string
		usersFile  string
		tweetID    string
		numWinners int
	)
	flag.StringVar(&keyFile, "key", "keys.json", "The file where we store consumer key and secret for twitter API")
	flag.StringVar(&usersFile, "users", "users.csv", "The file where users who have retweeted the tweet are stored. This will be created if it does not exist")
	flag.StringVar(&tweetID, "tweet", "", "The ID of the tweet you wish to find the retweeters of")
	flag.IntVar(&numWinners, "winners", 0, "The number of winners to pick for the contest")
	flag.Parse()

	keys, err := getKeysFromJson(keyFile)
	if err != nil {
		panic(err)
	}

	accessKeys := twc.AccessKeys{
		Key:    keys.Key,
		Secret: keys.Secret,
	}

	tClient, err := twc.GetTwitterClient(accessKeys)
	if err != nil {
		panic(err)
	}
	if tweetID != "" {
		retweets, err := tClient.GetRetweets(tweetID)
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

		if numWinners == 0 {
			return
		}

		fmt.Println("The winners are: ", pickWinners(usersFile, numWinners))
	} else {
		panic(errors.New("Please enter a tweet ID using the -tweet flag"))
	}
}

func pickWinners(usersFile string, num int) []string {
	winningUsers := make([]string, 0, num)
	users := existing(usersFile)

	seed := rand.NewSource(time.Now().Unix())
	random := rand.New(seed).Perm(len(users))[:num]
	for _, j := range random {
		winningUsers = append(winningUsers, users[j])
	}
	return winningUsers
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
