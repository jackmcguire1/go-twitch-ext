package twitchext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ExtensionEnabledChannels Response type of the getLiveChannelsWithExtensionEnabled
type ExtensionEnabledChannels struct {
	ResponseCommon
	Channels []*Channel `json:"channels"`
	Bookmark string     `json:"cursor"`
}

// Channel A struct representative of an individual
// channel returned by getLiveChannelsWithExtensionEnabled
type Channel struct {
	Game     string `json:"game"`
	ID       string `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Viewers  string `json:"view_count"`
}

// GetLiveChannelsWithExtensionEnabled Retrieve all live
// twitch channels which have the extension enabled.
// https://dev.twitch.tv/docs/extensions/reference/#get-live-channels-with-extension-activated
func (t *Twitch) GetLiveChannelsWithExtensionEnabled(
	extensionId string,
	bookmark string,
) (
	channels *ExtensionEnabledChannels,
	err error,
) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/live_activated_channels",
		extensionId,
	)

	q := url.Values{}

	if bookmark != "" {
		q.Set("cursor", bookmark)
	}

	resp, headers, err := t.do(http.MethodGet, addr, nil, nil, q)
	if err != nil {
		return
	}

	channels = &ExtensionEnabledChannels{}
	err = json.Unmarshal(resp, &channels)
	if err != nil {
		return
	}
	channels.Headers = headers

	return
}
