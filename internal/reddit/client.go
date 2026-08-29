package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PostListing struct {
	Data struct {
		After    string `json:"after"`
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Post struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Subreddit   string `json:"subreddit"`
	Author      string `json:"author"`
	Score       int    `json:"score"`
	NumComments int    `json:"num_comments"`
	Thumbnail   string `json:"thumbnail"`
	Permalink   string `json:"permalink"`
}

func GetHome() (*PostListing, error) {
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

	var listing PostListing
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return nil, err
	}

	return &listing, nil
}

type CommentListing struct {
	Data struct {
		Children []struct {
			Kind string  `json:"kind"`
			Data Comment `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Comment struct {
	Author  string    `json:"author"`
	Body    string    `json:"body"`
	Score   int       `json:"score"`
	Name    string    `json:"name"`
	Replies []Comment `json:"-"`
}

func (c *Comment) UnmarshalJSON(data []byte) error {
	type Alias Comment // avoid infinite recursion into UnmarshalJSON
	aux := &struct {
		Replies json.RawMessage `json:"replies"`
		*Alias
	}{Alias: (*Alias)(c)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if len(aux.Replies) == 0 || string(aux.Replies) == `""` {
		return nil // no replies
	}

	var nested CommentListing
	if err := json.Unmarshal(aux.Replies, &nested); err != nil {
		return err
	}
	for _, child := range nested.Data.Children {
		c.Replies = append(c.Replies, child.Data)
	}
	return nil
}

type PostDetail struct {
	Post     Post
	Comments []Comment
}

func GetPost(postId string) (*PostDetail, error) {
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

	req, err := http.NewRequest("GET", "https://old.reddit.com/comments/"+postId+"/.json", nil)
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

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw) != 2 {
		return nil, errors.New("unexpected reddit response shape")
	}

	var postListing PostListing
	if err := json.Unmarshal(raw[0], &postListing); err != nil {
		return nil, err
	}

	var commentListing CommentListing
	if err := json.Unmarshal(raw[1], &commentListing); err != nil {
		return nil, err
	}

	detail := PostDetail{}
	if len(postListing.Data.Children) > 0 {
		detail.Post = postListing.Data.Children[0].Data
	}
	for _, child := range commentListing.Data.Children {
		if child.Kind != "t1" {
			continue // skip "more" stubs
		}
		detail.Comments = append(detail.Comments, child.Data)
	}

	return &detail, nil
}
