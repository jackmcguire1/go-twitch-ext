package twitchext

import (
	"encoding/json"
	"fmt"
	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
	"net/http"
)

// SecretsResponse response structure received
// when generating or querying for generated secrets
type SecretsResponse struct {
	Version int       `json:"format_version"`
	Secrets []*Secret `json:"secrets"`
	ResponseCommon
}

// Secret information about a generated secret
type Secret struct {
	Active  string `json:"active"`
	Content string `json:"content"`
	Expires string `json:"expires"`
}

type secretCreation struct {
	ActivationDelay int `json:"activation_delay_secs"`
}

// CreateExtensionSecret create a new twitch secret for your twitch extension.
// The delay period, between the generation of the new secret and
// its use by Twitch, is specified by delay.
// https://dev.twitch.tv/docs/extensions/reference/#create-extension-secret
func (t *Twitch) CreateExtensionSecret(
	delay int,
) (
	resp *SecretsResponse,
	err error,
) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/auth/secret",
		t.ClientID,
	)

	claims := t.CreateClaims("", ExternalRole, FormGlobalSendPubSubPermissions())

	secret := &secretCreation{
		ActivationDelay: delay,
	}

	body, headers, err := t.do(http.MethodPost, addr, claims, utils.ToRawMessage(secret), nil)
	if err != nil {
		return
	}

	resp = &SecretsResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	resp.Headers = headers

	return
}

// GetExtensionSecrets returns a list of the extension secrets
// https://dev.twitch.tv/docs/extensions/reference/#get-extension-secret
func (t *Twitch) GetExtensionSecrets() (
	resp *SecretsResponse,
	err error,
) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/auth/secret",
		t.ClientID,
	)

	claims := t.CreateClaims("", ExternalRole, FormGlobalSendPubSubPermissions())
	body, headers, err := t.do(http.MethodGet, addr, claims, nil, nil)
	if err != nil {
		return
	}

	resp = &SecretsResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	resp.Headers = headers

	return
}

// RevokeExtensionSecrets all the secrets to the extension.
// https://dev.twitch.tv/docs/extensions/reference/#revoke-extension-secrets
func (t *Twitch) RevokeExtensionSecrets() (res *ResponseCommon, err error) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/auth/secret",
		t.ClientID,
	)

	claims := t.CreateClaims("", ExternalRole, FormGlobalSendPubSubPermissions())
	_, headers, err := t.do(http.MethodDelete, addr, claims, nil, nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{Headers: headers}

	return
}
