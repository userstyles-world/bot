package gif

import (
	"bot/modules/errors"

	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Simple Struct to represent a gif from the tenor API.
type Gif struct {
	Results []struct {
		URL string `json:"url"`
	} `json:"results"`
}

// GetRandomGif fetch the from the tenor API the gifs with the given query.
// And select a random gif from the result.
func GetRandomGif(query string) (string, error) {
	url := "https://api.tenor.com/v1/random?q=" + query
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var gifs Gif
	err = json.Unmarshal(body, &gifs)
	if err != nil {
		return "", err
	}
	if len(gifs.Results) == 0 {
		return "", errors.ErrNoGifsFound
	}
	rand.Seed(time.Now().UnixNano())
	gif := gifs.Results[rand.Intn(len(gifs.Results))]
	return gif.URL, nil
}
