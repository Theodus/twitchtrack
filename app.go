package main

import (
	"html/template"
	"net/http"
    "appengine/urlfetch"
    "appengine"
    "io/ioutil"
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

func init() {
	http.HandleFunc("/", twitchtrackHandler)
	http.HandleFunc("/refresh", refreshHandler)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    client := urlfetch.Client(ctx)
    res, err := client.Get("https://api.twitch.tv/kraken/users/ElTheodus/follows/channels")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        ctx.Errorf(err.Error())
        return
    }
    b, err := ioutil.ReadAll(res.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        ctx.Errorf(err.Error())
        return
    }
    ctx.Infof(string(b))
    /*
	channels, viewers, streams, links := refresh(ctx)
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
	*/
}

/*
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
*/
