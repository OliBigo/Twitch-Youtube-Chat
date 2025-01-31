package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

var mu sync.Mutex

func handleChannelYoutube(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var data struct {
		ChannelName string `json:"channelName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	connectedYoutube = data.ChannelName

	fmt.Printf("Received YouTube channel name: %s\n", data.ChannelName)
	YOUTUBE_API_KEY := os.Getenv("YOUTUBE_API_KEY")

	channelID := getYoutubeChannelID(data.ChannelName, YOUTUBE_API_KEY)
	liveStreamID := getYoutubeLiveStreamID(channelID, YOUTUBE_API_KEY)
	fmt.Printf("Live Stream ID: %s\n", liveStreamID)
	liveChatID := getYoutubeLiveChatID(liveStreamID, YOUTUBE_API_KEY)
	fmt.Printf("Live Chat ID: %s\n", liveChatID)

	go pollYoutubeChatMessages(liveChatID, YOUTUBE_API_KEY)

	w.WriteHeader(http.StatusOK)
}

func pollYoutubeChatMessages(liveChatID string, YOUTUBE_API_KEY string) {
    var nextPageToken string
    firstRequest := true

    for {
        fmt.Println("Polling YouTube chat messages...")

        // Ensure the first request does not include a nextPageToken
        var ytsr YouTubeLiveChatResponse
        var err error
        if firstRequest {
            ytsr, err = getYoutubeChatMessages(liveChatID, "", YOUTUBE_API_KEY, true)
            firstRequest = false
        } else {
            ytsr, err = getYoutubeChatMessages(liveChatID, nextPageToken, YOUTUBE_API_KEY, false)
        }

        if err != nil {
            log.Printf("Error polling YouTube chat messages: %v", err)
            time.Sleep(5 * time.Second) // Retry delay on error
            continue
        }

        // Process received messages
        if len(ytsr.Items) == 0 {
            log.Println("No new messages received")
        }

        for _, item := range ytsr.Items {
            msg := Message{
                Platform: "YouTube",
                User:     item.AuthorDetails.DisplayName,
                Content:  item.Snippet.DisplayMessage,
            }
            fmt.Println("Received message:", msg)
            broadcast <- msg
        }

        // Update nextPageToken for the next request
        nextPageToken = ytsr.NextPageToken

        // Respect YouTube's polling interval or use a default delay
        interval := time.Duration(ytsr.PollingIntervalMillis) * time.Millisecond
        if interval == 0 {
            interval = 5 * time.Second // Default polling interval
        }
        time.Sleep(interval)
    }
}

func handleChannelTwitch(w http.ResponseWriter, r *http.Request) {
	TWITCH_OAUTH_TOKEN := os.Getenv("TWITCH_OAUTH_TOKEN")
	TWITCH_BOT_NAME := os.Getenv("TWITCH_BOT_NAME")

	if twitchClient == nil {
		twitchClient = twitch.NewClient(TWITCH_BOT_NAME, TWITCH_OAUTH_TOKEN)
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var data struct {
		ChannelName string `json:"channelName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	connectedTwitch = data.ChannelName

	fmt.Printf("Received Twitch channel name: %s\n", data.ChannelName)

	twitchClient.Join(data.ChannelName)

	twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("[Twitch] %s: %s\n", message.User.Name, message.Message)
		broadcast <- Message{Platform: "Twitch", User: message.User.Name, Content: message.Message}
	})

	go func() {
		if err := twitchClient.Connect(); err != nil {
			log.Fatal(err)
		}
	}()
}