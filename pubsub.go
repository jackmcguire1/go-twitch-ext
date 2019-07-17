package twitchext

import (
	"fmt"
	"net/http"

	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
)

// PublishType The Pub/Sub broadcast type
type PublishType string

// Types of Pub/Sub Permissions or targets
const (
	GenericPublish   PublishType = "*"
	BroadcastPublish PublishType = "broadcast"
	GlobalPublish    PublishType = "global"
)

type pubSubNotification struct {
	Message     string        `json:"message"`
	Targets     []PublishType `json:"targets"`
	ContentType string        `json:"content_type"`
}

// PubSubPermissions publish permissions used within
// JWT claims
type PubSubPermissions struct {
	Send   []PublishType `json:"send,omitempty"`
	Listen []PublishType `json:"listen,omitempty"`
}

func createWhisper(opaqueId string) PublishType {
	return PublishType("whisper-" + opaqueId)
}

// FormWhisperSendPubSubPermissions create the pubsub permissions
// for publishing a whisper message type
func FormWhisperSendPubSubPermissions(opaqueId string) *PubSubPermissions {
	return &PubSubPermissions{
		Send: []PublishType{createWhisper(opaqueId)},
	}
}

// FormBroadcastSendPubSubPermissions create the pubsub permissions
// for publishing a broadcast message type
func FormBroadcastSendPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []PublishType{BroadcastPublish},
	}
}

// FormGlobalSendPubSubPermissions create the pubsub permissions
// for publishing a global targeted message
func FormGlobalSendPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []PublishType{GlobalPublish},
	}
}

// FormGenericPubSubPermissions create the pubsub permissions
// for publishing to message for any target type
func FormGenericPubSubPermissions() *PubSubPermissions {
	return &PubSubPermissions{
		Send: []PublishType{GenericPublish},
	}
}

// PublishChannelNotification publish a notification to
// a specific channel with the twitch extension enabled.
func (t *Twitch) PublishChannelNotification(channelID string, i interface{}) (res *ResponseCommon, err error) {
	url := fmt.Sprintf("https://api.twitch.tv/extensions/message/%s", channelID)
	claims := t.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())

	data := &pubSubNotification{
		Message:     utils.ToJSON(i),
		Targets:     []PublishType{BroadcastPublish},
		ContentType: "application/json",
	}
	_, headers, err := t.do(http.MethodPost, url, claims, utils.ToRawMessage(data), nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{Headers: headers}

	return
}

// PublishWhisperNotification publish a notification to
// a specific user viewing the twitch extension.
func (t *Twitch) PublishWhisperNotification(channelID string, opaqueId string, i interface{}) (res *ResponseCommon, err error) {
	url := fmt.Sprintf("https://api.twitch.tv/extensions/message/%s", channelID)
	claims := t.CreateClaims(channelID, ExternalRole, FormWhisperSendPubSubPermissions(opaqueId))

	data := &pubSubNotification{
		Message:     utils.ToJSON(i),
		Targets:     []PublishType{createWhisper(opaqueId)},
		ContentType: "application/json",
	}
	_, headers, err := t.do(http.MethodPost, url, claims, utils.ToRawMessage(data), nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{Headers: headers}

	return
}

// PublishGlobalNotification publish a notification to
// all channels with the twitch extension enabled.
// https://dev.twitch.tv/docs/extensions/reference/#send-extension-pubsub-message
func (t *Twitch) PublishGlobalNotification(i interface{}) (res *ResponseCommon, err error) {
	url := fmt.Sprintf("https://api.twitch.tv/extensions/message/all")
	claims := t.CreateClaims("", ExternalRole, FormGlobalSendPubSubPermissions())

	data := &pubSubNotification{
		Message:     utils.ToJSON(i),
		Targets:     []PublishType{GlobalPublish},
		ContentType: "application/json",
	}

	_, headers, err := t.do(http.MethodPost, url, claims, utils.ToRawMessage(data), nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{Headers: headers}

	return
}
