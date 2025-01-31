package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getYoutubeChannelID(channelName, YOUTUBE_API_KEY string) string {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/search?part=snippet&q=%s&key=%s&type=channel", channelName, YOUTUBE_API_KEY)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting YouTube channel ID: %v", err)
	}
	defer res.Body.Close()

	var ytsr YouTubeSearchChannelResponse
	if err := json.NewDecoder(res.Body).Decode(&ytsr); err != nil {
		log.Fatalf("Error decoding YouTube channel ID: %v", err)
	}
	channelID := ytsr.Items[0].ID.ChannelID
	fmt.Printf("Channel ID: %s\n", channelID)
	return channelID
}

func getYoutubeLiveStreamID(channelID string, YOUTUBE_API_KEY string) string {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%s&order=date&type=video&key=%s", channelID, YOUTUBE_API_KEY)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting YouTube live stream ID: %v", err)
	}
	defer res.Body.Close()

	var ytsr YouTubeSearchVideoResponse
	if err := json.NewDecoder(res.Body).Decode(&ytsr); err != nil {
		log.Fatalf("Error decoding YouTube live stream ID: %v", err)
	}

	for _, item := range ytsr.Items {
		if item.Snippet.LiveBroadcastContent == "live" {
			return item.ID.VideoID
		}
	}
	log.Println("No live stream found with 'live' broadcast content")
	return ""
}

func getYoutubeLiveChatID(videoID string, YOUTUBE_API_KEY string) string {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,liveStreamingDetails&id=%s&order=date&type=video&key=%s", videoID, YOUTUBE_API_KEY)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting YouTube live chat ID: %v", err)
	}
	defer res.Body.Close()

	var ytsr YoutubeLiveStreamResponse
	if err := json.NewDecoder(res.Body).Decode(&ytsr); err != nil {
		log.Fatalf("Error decoding YouTube live chat ID: %v", err)
	}

	return ytsr.Items[0].LiveStreamingDetails.ActiveLiveChatID
}

func getYoutubeChatMessages(liveChatID, pageToken string, YOUTUBE_API_KEY string, isFirstRequest bool) (YouTubeLiveChatResponse, error) {
	var url string
	if isFirstRequest {
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/liveChat/messages?liveChatId=%s&part=snippet,authorDetails&maxResults=2000&key=%s", liveChatID, YOUTUBE_API_KEY)
	} else {
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/liveChat/messages?liveChatId=%s&part=snippet,authorDetails&maxResults=2000&pageToken=%s&key=%s", liveChatID, pageToken, YOUTUBE_API_KEY)
	}
	res, err := http.Get(url)
	if err != nil {
		return YouTubeLiveChatResponse{}, fmt.Errorf("error getting YouTube chat messages: %v", err)
	}
	defer res.Body.Close()

	var ytsr YouTubeLiveChatResponse
	if err := json.NewDecoder(res.Body).Decode(&ytsr); err != nil {
		return YouTubeLiveChatResponse{}, fmt.Errorf("error decoding YouTube chat messages: %v", err)
	}
	return ytsr, nil
}