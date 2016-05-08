package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/elm.js", elmHandler)
	http.HandleFunc("/data", dataHandler)
	var port string
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		port = "80"
	}
	log.Println("Serving on port", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func elmHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "elm.js")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	d, err := refresh()
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	enc := json.NewEncoder(w)
	if err = enc.Encode(d); err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Response: ", d)
	return
}

func refresh() (d data, err error) {
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
			URL:     e.Channel.URL,
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
				return
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
			URL    string
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
	URL     string `json:"url"`
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
