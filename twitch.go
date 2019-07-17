package twitchext

import (
	"net/http"
)

// Twitch package struct
type Twitch struct {
	client        *http.Client
	OwnerID       string
	Secret        string
	ClientID      string
	Version       string
	ConfigVersion string
}

// Options optional parameters for the
// twitch ext client for configuration
type Options struct {
	Client *http.Client
}

// NewClient create reference to twitch-ext package
// https://dev.twitch.tv/docs/extensions/reference/#endpoints
func NewClient(
	ownerID string,
	clientID string,
	secret string,
	extVersion string,
	configVersion string,
	opts ...*Options,
) (
	twitch *Twitch,
) {
	twitch = &Twitch{
		client:        &http.Client{},
		OwnerID:       ownerID,
		Secret:        secret,
		ClientID:      clientID,
		Version:       extVersion,
		ConfigVersion: configVersion,
	}

	if len(opts) > 0 {
		if opts[0].Client != nil {
			twitch.client = opts[0].Client
		}
	}

	return
}
