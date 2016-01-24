package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const RefreshTime = 120 // time in seconds

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/script.js", scriptHandler)
	http.HandleFunc("/longpoll", longpollHandler)
	log.Println("Serving.")
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "twitchtrack.html")
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "script.js")
}

var lastData data

func longpollHandler(w http.ResponseWriter, r *http.Request) {
	first, err := strconv.ParseBool(r.URL.Query().Get("first"))
	if err != nil || first == false {
		first = false
	}
	log.Println("First Request:", first)
	enc := json.NewEncoder(w)
	if first {
		d, err := refresh(w)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		lastData = d
		log.Println("First: ", d)
		if err := enc.Encode(d); err != nil {
			log.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	for _ = range time.Tick(time.Duration(RefreshTime) * time.Second) {
		d, err := refresh(w)
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if len(lastData.Channels) != len(d.Channels) {
			lastData = d
			if err := enc.Encode(d); err != nil {
				log.Println("Error: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Println("Refresh: ", d)
			return
		}
		for i, tc := range d.Channels {
			if tc.Stream != lastData.Channels[i].Stream {
				lastData = d
				if err := enc.Encode(d); err != nil {
					log.Println("Error: ", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				log.Println("Refresh: ", d)
				return
			}
		}
	}
}

func refresh(w http.ResponseWriter) (d data, err error) {
	f, err := getFollows()
	if err != nil {
		return d, err
	}
	var tc []channel
	for _, e := range f.Follows {
		tc = append(tc, channel{
			Channel: e.Channel.Name,
			Game:    e.Channel.Game,
			Stream:  e.Channel.Status,
			Url:     e.Channel.Url,
		})
	}
	var wg sync.WaitGroup
	for i, c := range tc {
		wg.Add(1)
		go func(i int, c channel) {
			defer wg.Done()
			v, err := getViewers(c.Channel)
			if err != nil {
				log.Println("Error: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			tc[i].Viewers = v
		}(i, c)
	}
	wg.Wait()
	return data{tc}, nil
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

type data struct {
	Channels []channel `json:"channels"`
}

type channel struct {
	Channel string `json:"channel"`
	Game    string `json:"game"`
	Stream  string `json:"stream"`
	Url     string `json:"url"`
	Viewers int    `json:"viewers"`
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

type streams struct {
	Stream struct {
		Viewers int
	}
}

func getViewers(channel string) (viewers int, err error) {
	var s streams
	viewers = 0
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/streams/%s/", channel), nil)
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
