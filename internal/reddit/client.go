package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Listing struct {
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Post struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Subreddit   string `json:"subreddit"`
	Author      string `json:"author"`
	Score       int    `json:"score"`
	NumComments int    `json:"num_comments"`
	Thumbnail   string `json:"thumbnail"`
	Permalink   string `json:"permalink"`
}

func GetHome() (*Listing, error) {
	cookie := os.Getenv("REDDIT_SESSION_COOKIE")
	cookie = strings.TrimSpace(cookie)
	if cookie == "" {
		return nil, errors.New("Needed Reddit Session Cookie")
	}

	username := os.Getenv("REDDIT_USERNAME")
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("Needed Reddit User Name")
	}

	req, err := http.NewRequest("GET", "https://old.reddit.com/.json?limit=25", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "reddit_session="+cookie)

	// reddit blocks generic/missing UA
	req.Header.Set("User-Agent", "reddit-pwa-client/0.1 by "+username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Reddit returned %d", resp.StatusCode)
	}

	var listing Listing
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return nil, err
	}

	return &listing, nil
}
