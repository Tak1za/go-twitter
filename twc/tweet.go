package twc

//Retweet Response
type Retweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

//GetRetweets returns a slice of Retweet information
func (api TWApi) GetRetweets(id string) ([]Retweet, error) {
	var ret []Retweet
	route := "statuses/retweets/"
	url := baseUrl + route + id + ".json"
	err := api.getRequest(url, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
