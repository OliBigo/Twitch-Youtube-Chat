package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
    conn *websocket.Conn
    id   string
}

type Message struct {
    Platform string `json:"platform"`
    User     string `json:"user"`
    Content  string `json:"content"`
}

type YouTubeSearchChannelResponse struct {
    Kind       string       `json:"kind"`
    Etag       string       `json:"etag"`
    RegionCode string       `json:"regionCode"`
    PageInfo   PageInfo     `json:"pageInfo"`
    Items      []ChannelItem `json:"items"`
}

type YouTubeSearchVideoResponse struct {
    Kind       string      `json:"kind"`
    Etag       string      `json:"etag"`
    RegionCode string      `json:"regionCode"`
    PageInfo   PageInfo    `json:"pageInfo"`
    Items      []VideoItem `json:"items"`
}

type YoutubeLiveStreamResponse struct {
    Kind  string          `json:"kind"`
    Etag  string          `json:"etag"`
    Items []LiveStreamItem `json:"items"`
    Page  PageInfo        `json:"pageInfo"`
}

type YouTubeLiveChatResponse struct {
    Kind                  string      `json:"kind"`
    Etag                  string      `json:"etag"`
    PollingIntervalMillis uint        `json:"pollingIntervalMillis"`
    PageInfo              PageInfo    `json:"pageInfo"`
    NextPageToken         string      `json:"nextPageToken"`
    Items                 []LiveItem  `json:"items"`
}

type PageInfo struct {
    TotalResults   int `json:"totalResults"`
    ResultsPerPage int `json:"resultsPerPage"`
}

type ChannelItem struct {
    Kind    string   `json:"kind"`
    Etag    string   `json:"etag"`
    ID      ChannelID `json:"id"`
    Snippet Snippet  `json:"snippet"`
}

type VideoItem struct {
    Kind                string              `json:"kind"`
    Etag                string              `json:"etag"`
    ID                  VideoID             `json:"id"`
    Snippet             Snippet             `json:"snippet"`
    LiveStreamingDetails LiveStreamingDetails `json:"liveStreamingDetails"`
}

type LiveStreamItem struct {
    Kind                string              `json:"kind"`
    Etag                string              `json:"etag"`
    ID                  string              `json:"id"`
    Snippet             SnippetLive         `json:"snippet"`
    LiveStreamingDetails LiveStreamingDetails `json:"liveStreamingDetails"`
}

type LiveItem struct {
    Kind          string        `json:"kind"`
    Etag          string        `json:"etag"`
    ID            string        `json:"id"`
    Snippet       SnippetLive   `json:"snippet"`
    AuthorDetails AuthorDetails `json:"authorDetails"`
}

type ChannelID struct {
    Kind      string `json:"kind"`
    ChannelID string `json:"channelId"`
}

type VideoID struct {
    Kind    string `json:"kind"`
    VideoID string `json:"videoId"`
}

type Snippet struct {
    PublishedAt         string     `json:"publishedAt"`
    ChannelID           string     `json:"channelId"`
    Title               string     `json:"title"`
    Description         string     `json:"description"`
    Thumbnails          Thumbnails `json:"thumbnails"`
    ChannelTitle        string     `json:"channelTitle"`
    LiveBroadcastContent string     `json:"liveBroadcastContent"`
    PublishTime         string     `json:"publishTime"`
}

type SnippetLive struct {
    Type              string            `json:"type"`
    LiveChatID        string            `json:"liveChatId"`
    AuthorChannelID   string            `json:"authorChannelId"`
    PublishedAt       string            `json:"publishedAt"`
    HasDisplayContent bool              `json:"hasDisplayContent"`
    DisplayMessage    string            `json:"displayMessage"`
    TextMessageDetails TextMessageDetails `json:"textMessageDetails"`
}

type LiveStreamingDetails struct {
    ActualStartTime    string `json:"actualStartTime"`
    ScheduledStartTime string `json:"scheduledStartTime"`
    ConcurrentViewers  string `json:"concurrentViewers"`
    ActiveLiveChatID   string `json:"activeLiveChatId"`
}

type TextMessageDetails struct {
    MessageText string `json:"messageText"`
}

type AuthorDetails struct {
    ChannelID        string `json:"channelId"`
    ChannelUrl       string `json:"channelUrl"`
    DisplayName      string `json:"displayName"`
    ProfileImageUrl  string `json:"profileImageUrl"`
    IsVerified       bool   `json:"isVerified"`
    IsChatOwner      bool   `json:"isChatOwner"`
    IsChatSponsor    bool   `json:"isChatSponsor"`
    IsChatModerator  bool   `json:"isChatModerator"`
}

type Thumbnails struct {
    Default Thumbnail `json:"default"`
    Medium  Thumbnail `json:"medium"`
    High    Thumbnail `json:"high"`
}

type Thumbnail struct {
    URL    string `json:"url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}