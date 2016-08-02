package main

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"google.golang.org/api/youtube/v3"
	"net/http"
)

var YouTube *youtube.Service

// Search code ripped from https://github.com/nlordell/utub

type Yeolcha struct {
	Transport http.RoundTripper
}

func (t *Yeolcha) RoundTrip(req *http.Request) (*http.Response, error) {
	// API key from https://developers.google.com/youtube/v3/docs/search/list
	// Note that we need to trick the YouTube search service into thinking the
	// page actually made the request by setting the propery referer and origin
	// headers

	newReq := *req

	args := newReq.URL.Query()
	args.Set("key", "AIzaSyD-a9IF8KKYgoC3cpgS-Al7hLQDbugrDcw")
	newReq.URL.RawQuery = args.Encode()

	newReq.Header = make(http.Header)
	for k, v := range req.Header {
		newReq.Header[k] = v
	}
	newReq.Header.Add("referer", "https://content.googleapis.com/static/proxy.html?jsh=m%3B%2F_%2Fscs%2Fapps-static%2F_%2Fjs%2Fk%3Doz.gapi.en_GB.Pc_OA3os_Rw.O%2Fm%3D__features__%2Fam%3DAQ%2Frt%3Dj%2Fd%3D1%2Frs%3DAGLTcCPi4tWKbCZjJQ1Tpnq94gY9Shvgag")
	newReq.Header.Add("x-origin", "https://developers.google.com")
	newReq.Header.Add("x-referer", "https://developers.google.com")

	transport := t.Transport
	if t.Transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(&newReq)
}

func chat_yt(b *Bot) {
	b.ChatHooks["!yt"] = &ChatCommand{
		Func: func(x *Bot, search string, m *discord.Message) bool {
			if YouTube == nil {
				service, err := youtube.New(&http.Client{
					Transport: &Yeolcha{},
				})

				if err != nil {
					fmt.Println("YouTube error: ", err.Error())
					return true
				}

				YouTube = service
			}

			response, err := YouTube.Search.List("id,snippet").
				Q(search).
				MaxResults(1).
				Type("video").
				Do()

			if err != nil {
				x.Send(m.ChannelID, "Failed to perform search! ["+err.Error()+"]")
				return true
			}

			for _, item := range response.Items {
				// item.Id.VideoId, item.Snippet.Title
				x.Send(m.ChannelID, "https://www.youtube.com/watch?v="+item.Id.VideoId)
				return true
			}

			x.Send(m.ChannelID, "No results! :frowning:")
			return true
		},
		Cmd:     "!yt",
		Access:  CHAT_PRIVATE,
		HasArgs: true,
		Help:    "Lists the first video result of a YouTube search.",
		ArgHelp: "!yt <title>",
	}
}
