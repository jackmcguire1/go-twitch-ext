package twitchext

import (
	"fmt"
	"net/http"

	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
)

const maxMessageSize = 280

// struct containing the desired
// chat message.
type chatMessage struct {
	Text string `json:"text"`
}

// SendTwitchChatMessage publish message to twitch chat of a specific channel ID.
// - Twitch extension must have this permission
// - The maximum message size is 280 character
// - There is a limit of 12 messages per minute, per channel.
// https://dev.twitch.tv/docs/extensions/reference/#send-extension-chat-message
func (t *Twitch) SendTwitchChatMessage(channelID string, message string) (res *ResponseCommon, err error) {

	if channelID == "" {
		err = fmt.Errorf("missing channelID")
		return
	}

	if len(message) > maxMessageSize {
		err = fmt.Errorf(
			"message %q exceeds %d character limit",
			message,
			maxMessageSize,
		)
		return
	}

	url := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/%s/channels/%s/chat",
		t.ClientID,
		t.Version,
		channelID,
	)
	claims := t.CreateClaims(channelID, BroadcasterRole, FormBroadcastSendPubSubPermissions())

	msg := &chatMessage{Text: message}

	_, headers, err := t.do(http.MethodPost, url, claims, utils.ToRawMessage(msg), nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{headers}

	return
}
