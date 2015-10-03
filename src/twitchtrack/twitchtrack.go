package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/mrshankly/go-twitch/twitch"
)

func twitchtrackHandler(w http.ResponseWriter, r *http.Request) {
	p, err := template.ParseFiles("twitchtrack.html")
	if err != nil {
		panic(err)
	}
	if err := p.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", twitchtrackHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.ListenAndServe("localhost:80", nil)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	channels, viewers, streams, links := refresh()
	res := map[string]interface{}{
		"channels": channels,
		"viewers":  viewers,
		"streams":  streams,
		"links":    links,
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(res)
	if err != nil {
		fmt.Println(err)
	}
}

func refresh() ([]string, []int, []string, []string) {
	channels := []string{}
	viewers := []int{}
	streams := []string{}
	links := []string{}

	client := twitch.NewClient(&http.Client{})
	res, err := client.Users.Follows("ElTheodus", nil)
	if err != nil {
		fmt.Println(err)
	}
	for _, e := range res.Follows {
		name := e.Channel.Name
		channels = append(channels, name)
		links = append(links, e.Channel.Url)
		res, err := client.Streams.Channel(name)
		if err != nil {
			fmt.Println(err)
		}
		viewers = append(viewers, res.Stream.Viewers)
		streams = append(streams, res.Stream.Game)
	}
	return channels, viewers, streams, links
}
