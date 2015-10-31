package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	log.Println("Ready to serve.")
	http.HandleFunc("/script.js", jsHandler)
	http.HandleFunc("/", twitchtrackHandler)
	http.HandleFunc("/refresh", refreshHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "script.js")
}

func twitchtrackHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "twitchtrack.html")
}

type follows struct {
	Follows []struct {
		Channel struct {
			Game   string
			Name   string
			Status string
			Url    string
		}
	}
}

type streams struct {
	Stream struct {
		Viewers int
	}
}

type data struct {
	Channels []*channel `json:"channels"`
}

type channel struct {
	Channel string `json:"channel"`
	Game    string `json:"game"`
	Viewers int    `json:"viewers"`
	Stream  string `json:"stream"`
	Url     string `json:"url"`
	Online  bool   `json:"online"`
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	f, err := getFollows()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	enc := json.NewEncoder(w)
	var tc []*channel
	for _, e := range f.Follows {
		tc = append(tc, &channel{
			Channel: e.Channel.Name,
			Game:    e.Channel.Game,
			Viewers: 0,
			Stream:  e.Channel.Status,
			Url:     e.Channel.Url,
			Online:  false,
		})
	}
	var wg sync.WaitGroup
	views := make([]int, len(tc))
	for i, e := range tc {
		wg.Add(1)
		go func(i int, c string) {
			defer wg.Done()
			v, err := getViewers(c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err)
				return
			}
			views[i] = v
		}(i, e.Channel)
	}
	wg.Wait()
	for i, e := range views {
		tc[i].Viewers = e
		if e > 0 {
			tc[i].Online = true
		}
	}
	err = enc.Encode(data{tc})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func getFollows() (follows, error) {
	var f follows
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/users/ElTheodus/follows/channels", nil)
	if err != nil {
		return f, err
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return f, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return f, err
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

func getViewers(channel string) (viewers int, err error) {
	var s streams
	viewers = 0
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/", "https://api.twitch.tv/kraken/streams/", channel), nil)
	if err != nil {
		return viewers, err
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return viewers, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return viewers, err
	}
	if err = json.Unmarshal(b, &s); err != nil {
		return viewers, err
	}
	viewers = s.Stream.Viewers
	return viewers, nil
}
