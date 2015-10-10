package main

import (
	"html/template"
	"net/http"
    "appengine/urlfetch"
    "appengine"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "sync"
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
    ctx := appengine.NewContext(r)
    f, err := getFollows(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        ctx.Errorf(err.Error())
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
            Url:    e.Channel.Url,
            Online:  false,
        })
    }
    var wg sync.WaitGroup
    views := make([]int, len(tc))
    for i, e := range tc {
        wg.Add(1)
        go func(i int, c string, ctx appengine.Context) {
            defer wg.Done()
            v, err := getViewers(ctx, c)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                ctx.Errorf(err.Error())
                return
            }
            views[i] = v
        }(i, e.Channel, ctx)
    }
    wg.Wait()
    for i, e := range views {
        tc[i].Viewers = e
        if e>0 {
            tc[i].Online = true
        }
    }
    err = enc.Encode(data{tc})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        ctx.Errorf(err.Error())
        return
    }
}

func getFollows(ctx appengine.Context) (follows, error){
    var f follows
    client := urlfetch.Client(ctx)
    res, err := client.Get("https://api.twitch.tv/kraken/users/ElTheodus/follows/channels")
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

func getViewers(ctx appengine.Context, channel string) (viewers int, err error){
    var s streams
    viewers = 0
    client := urlfetch.Client(ctx)
    res, err := client.Get(fmt.Sprintf("%s%s/", "https://api.twitch.tv/kraken/streams/", channel))
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
