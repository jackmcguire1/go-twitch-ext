package twitchext

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
)

// ResponseCommon ...
type ResponseCommon struct {
	Headers http.Header
}

func (rc *ResponseCommon) convertHeaderToInt(header string) (v int) {
	v, _ = strconv.Atoi(rc.Headers.Get(header))
	return
}

// GetPubSubChannelRateLimit returns the "Ratelimit-Ratelimitermessagesbychannel-Limit" header as an int
// for the pubsub messages to channel
func (rc *ResponseCommon) GetPubSubChannelRateLimit() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimitermessagesbychannel-Limit")
}

// GetPubSubChannelRateLimitRemaining returns the "Ratelimit-Ratelimitermessagesbychannel-Remaining" header as an int
func (rc *ResponseCommon) GetPubSubChannelRateLimitRemaining() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimitermessagesbychannel-Remaining")
}

// GetExtSetConfigurationRateLimit returns the "Ratelimit-Ratelimiterextensionsetconfiguration-Limit" header as an int
func (rc *ResponseCommon) GetExtSetConfigurationRateLimit() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimiterextensionsetconfiguration-Limit")
}

// GetExtSetConfigurationRateLimitRemaining returns the "Ratelimit-Ratelimiterextensionsetconfiguration-Remaining" header as an int.
func (rc *ResponseCommon) GetExtSetConfigurationRateLimitRemaining() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimiterextensionsetconfiguration-Remaining")
}

// GetExtSendChatMessageRateLimit returns the "Ratelimit-Ratelimit-Ratelimiterextensionchatmessages-Limit" header as an int.
func (rc *ResponseCommon) GetExtSendChatMessageRateLimit() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimiterextensionchatmessages-Limit")
}

// GetExtSendChatMessageRateLimitRemaining returns the "Ratelimit-Ratelimit-Ratelimiterextensionchatmessages-Remaining" header as an int.
func (rc *ResponseCommon) GetExtSendChatMessageRateLimitRemaining() int {
	return rc.convertHeaderToInt("Ratelimit-Ratelimiterextensionchatmessages-Remaining")
}

func (t *Twitch) do(
	method string,
	url string,
	claims *TwitchJWTClaims,
	b []byte,
	q url.Values,
) (
	data []byte,
	headers http.Header,
	err error,
) {
	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		err = fmt.Errorf("failed to construct request err:%s", err)
		return
	}
	err = t.setExtensionRequestHeaders(req, claims)
	if err != nil {
		return
	}
	req.URL.RawQuery = q.Encode()

	resp, err := t.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	headers = resp.Header

	switch resp.StatusCode {
	case http.StatusOK:
		data, err = ioutil.ReadAll(resp.Body)
	case http.StatusNoContent:
	case http.StatusTooManyRequests:
		err = fmt.Errorf(
			"rate limit exceeded, response headers:%q",
			utils.ToJSON(headers),
		)
		return
	default:
		err = fmt.Errorf(
			"unsupported response httpCode:%d status:%s",
			resp.StatusCode,
			resp.Status,
		)
		return
	}

	return
}

func (t *Twitch) setExtensionRequestHeaders(
	req *http.Request,
	claims *TwitchJWTClaims,
) (
	err error,
) {
	if claims != nil {
		var token string
		token, err = t.JWTSign(claims)
		if err != nil {
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Client-ID", t.ClientID)
	req.Header.Set("Content-Type", "application/json")

	return
}
