package twitchext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
)

// SegmentType A segment configuration type
type SegmentType string

// Types of segments datastores for the configuration service
const (
	BroadcasterSegment SegmentType = "broadcaster"
	DeveloperSegment   SegmentType = "developer"
	GlobalSegment      SegmentType = "global"
)

type configurationParams struct {
	Segment   SegmentType `json:"segment"`
	ChannelID string      `json:"channel_id,omitempty"`
	Version   string      `json:"version,omitempty"`
	Content   string      `json:"content"`
}

// AllConfigurationsResponse contains all possible
// segment configurations for a broadcaster
type AllConfigurationsResponse struct {
	Configurations map[string]*Configuration `json:"configurations"`
	ResponseCommon
}

// ConfigurationResponse contains data for the queried
// segment configuration  of type global, broadcaster or developer
type ConfigurationResponse struct {
	Configuration *Configuration `json:"configurations"`
	ResponseCommon
}

// Configuration Contains SegmentType and record information
type Configuration struct {
	Segment *Segment `json:"segment"`
	Record  *Record  `json:"record"`
}

// Record contains information about the
// version and raw content stored within
// the SegmentType configuation
type Record struct {
	Version string `json:"version"`
	Content string `json:"content"`
}

// Segment contains information about the
// type and twitch channel the configuration
// is for.
type Segment struct {
	Segment   string `json:"segment_type"`
	ChannelID string `json:"channel_id,omitempty"`
}

// SetGlobalSegment sets global extension configuration
func (t *Twitch) SetGlobalSegment(data interface{}) (res *ResponseCommon, err error) {
	return t.setSegmentConfig(data, "", GlobalSegment)
}

// SetBroadcasterSegment sets channel specific broadcaster SegmentType configuration
func (t *Twitch) SetBroadcasterSegment(data interface{}, channelID string) (res *ResponseCommon, err error) {
	return t.setSegmentConfig(data, channelID, BroadcasterSegment)
}

// SetDeveloperSegment sets channel specific developer SegmentType configuration
func (t *Twitch) SetDeveloperSegment(data interface{}, channelID string) (res *ResponseCommon, err error) {
	return t.setSegmentConfig(data, channelID, DeveloperSegment)
}

// GetGlobalSegment retrieves global extension SegmentType configuration
func (t *Twitch) GetGlobalSegment() (*ConfigurationResponse, error) {
	return t.getSegmentConfig("", GlobalSegment)
}

// GetBroadcasterSegment retrieves channel specific Broadcaster SegmentType configuration
func (t *Twitch) GetBroadcasterSegment(channelID string) (*ConfigurationResponse, error) {
	return t.getSegmentConfig(channelID, BroadcasterSegment)
}

// GetDeveloperSegment retrieves channel specific developer SegmentType configuration
func (t *Twitch) GetDeveloperSegment(channelID string) (*ConfigurationResponse, error) {
	return t.getSegmentConfig(channelID, DeveloperSegment)
}

// GetAllChannelConfigurations retrieves channel specific configuration
// returns map containing both broadcaster and developer segments.
// https://dev.twitch.tv/docs/extensions/reference/#get-extension-channel-configuration
func (t *Twitch) GetAllChannelConfigurations(channelID string) (resp *AllConfigurationsResponse, err error) {
	resp = &AllConfigurationsResponse{
		Configurations: map[string]*Configuration{},
	}

	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/configurations/channels/%s",
		t.ClientID,
		channelID,
	)
	claims := t.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())

	body, headers, err := t.do(http.MethodGet, addr, claims, nil, nil)
	if err != nil {
		return
	}

	var configurations map[string]*Configuration
	err = json.Unmarshal(body, &configurations)
	if err != nil {
		return
	}

	for segmentName, segment := range configurations {
		resp.Configurations[strings.Split(segmentName, ":")[0]] = segment
	}
	resp.Headers = headers

	return
}

// https://dev.twitch.tv/docs/extensions/reference/#set-extension-configuration-segment
func (t *Twitch) setSegmentConfig(data interface{}, channelID string, segment SegmentType) (res *ResponseCommon, err error) {
	addr := fmt.Sprintf("https://api.twitch.tv/extensions/%s/configurations/", t.ClientID)

	claims := t.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())
	segmentConfig := configurationParams{
		Segment: segment,
		Content: utils.ToJSON(data),
		Version: t.ConfigVersion,
	}

	if segment != GlobalSegment {
		segmentConfig.ChannelID = channelID
	}

	_, headers, err := t.do(http.MethodPut, addr, claims, utils.ToRawMessage(segmentConfig), nil)
	if err != nil {
		return
	}
	res = &ResponseCommon{Headers: headers}

	return
}

func (t *Twitch) getSegmentConfig(channelID string, segment SegmentType) (resp *ConfigurationResponse, err error) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/configurations/segments/%s",
		t.ClientID,
		segment,
	)
	q := url.Values{}

	if segment != GlobalSegment {
		q.Set("channel_id", channelID)
	}

	claims := t.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())

	body, headers, err := t.do(http.MethodGet, addr, claims, nil, q)
	if err != nil {
		return
	}

	config := map[string]*Configuration{}
	err = json.Unmarshal(body, &config)
	if err != nil {
		return
	}

	configuration, ok := config[fmt.Sprintf("%s:%s", string(segment), channelID)]
	if !ok {
		err = fmt.Errorf(
			"Configuration missing segment:%s channelID:%s",
			segment,
			channelID,
		)
		return
	}

	resp = &ConfigurationResponse{
		Configuration: configuration,
		ResponseCommon: ResponseCommon{
			Headers: headers,
		},
	}

	return
}

type reqConfiguration struct {
	RequiredConfiguration string `json:"required_configuration"`
}

// SetExtensionRequired is used to indicate that a channel has
// updated their configration, matching that of the latest twitch
// extension configuration version configured.
// https://dev.twitch.tv/docs/extensions/reference/#set-extension-required-configuration
func (t *Twitch) SetExtensionRequired(channelID string) (res *ResponseCommon, err error) {
	addr := fmt.Sprintf(
		"https://api.twitch.tv/extensions/%s/%s/required_configuration",
		t.ClientID,
		t.Version,
	)
	q := url.Values{}
	q.Set("channel_id", channelID)

	claims := t.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())

	config := &reqConfiguration{
		RequiredConfiguration: t.ConfigVersion,
	}

	_, headers, err := t.do(http.MethodPut, addr, claims, utils.ToRawMessage(config), q)
	if err != nil {
		return
	}

	res = &ResponseCommon{Headers: headers}
	return
}
