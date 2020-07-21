package twitchext

import (
	"os"
	"testing"

	"github.com/jackmcguire1/go-twitch-ext/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	twitchPkg *Twitch

	channelID          string
	publishedExtension string
)

type (
	UtilTests               struct{ Test *testing.T }
	JWTTests                struct{ Test *testing.T }
	ConfigurationTests      struct{ Test *testing.T }
	MessagingAndPubSubTests struct{ Test *testing.T }
	ChannelInfoTests        struct{ Test *testing.T }
	SecretTests             struct{ Test *testing.T }

	TestObj struct{ Name string }
)

func init() {
	twitchPkg = NewClient(
		os.Getenv("OWNER_ID"),
		os.Getenv("CLIENT_ID"),
		os.Getenv("EXT_SECRET"),
		os.Getenv("EXT_VERSION"),
		os.Getenv("EXT_CONFIG_VER"),
	)

	channelID = os.Getenv("CHANNEL_ID")
	publishedExtension = os.Getenv("PUBLISHED_EXT")
}

func TestRunner(t *testing.T) {

	t.Run("A=util", func(t *testing.T) {
		test := UtilTests{Test: t}
		test.TestToJSON()
		test.TestToJSONRawMessage()
		test.TestToJSONArray()
	})

	//TODO test without Twitch production API
	//t.Run("A=jwt", func(t *testing.T) {
	//	test := JWTTests{Test: t}
	//	test.TestCreateClaims()
	//	test.TestJWTSign()
	//	test.TestJWTVerify()
	//})
	//
	//t.Run("A=messaging", func(t *testing.T) {
	//	test := MessagingAndPubSubTests{Test: t}
	//	test.TestFormPubSubPermissions()
	//	test.TestSendTwitchChatMessageWithMessageLimitExceeded()
	//	test.TestSendTwitchChatMessage()
	//	test.TestPublishChannelNotification()
	//	test.TestPublishWhisperNotification()
	//	test.TestPublishGlobalNotification()
	//})
	//
	//t.Run("A=config", func(t *testing.T) {
	//	test := ConfigurationTests{Test: t}
	//	test.TestSetGlobalSegment()
	//	test.TestSetBroadcasterSegment()
	//	test.TestSetDeveloperSegment()
	//	test.TestSetExtensionRequired()
	//	test.TestGetBroadcasterSegment()
	//	test.TestGetDeveloperSegment()
	//	test.TestGetAllChannelConfigurations()
	//	test.TestGetGlobalSegment()
	//})
	//
	//t.Run("A=channel", func(t *testing.T) {
	//	test := ChannelInfoTests{Test: t}
	//	test.TestGetLiveChannelsWithExtensionEnabled()
	//})
	//
	//t.Run("A=secrets", func(t *testing.T) {
	//	test := SecretTests{Test: t}
	//	test.TestCreateSecret()
	//	test.TestGetExtensionSecrets()
	//	test.TestRevokeAllSecrets()
	//})

}

func (t *UtilTests) TestToJSON() {
	assert := assert.New(t.Test)

	data := struct {
		Message string
	}{
		Message: "test",
	}

	testData := utils.ToJSON(data)
	assert.EqualValues(`{"Message":"test"}`, testData)
}

func (t *UtilTests) TestToJSONRawMessage() {
	assert := assert.New(t.Test)

	data := struct {
		Message string
	}{
		Message: "test",
	}

	testData := utils.ToRawMessage(data)
	assert.NotEmpty(testData)
}

func (t *UtilTests) TestToJSONArray() {
	assert := assert.New(t.Test)

	data := struct {
		Message string
	}{
		Message: "test",
	}

	testData := utils.ToJSONArray(data)
	assert.NotEmpty(testData)
}

func (t *JWTTests) TestCreateClaims() {
	assert := assert.New(t.Test)

	claims := twitchPkg.CreateClaims("", "", FormBroadcastSendPubSubPermissions())
	assert.EqualValues(toAllChannels, claims.ChannelID)
	assert.EqualValues(ExternalRole, claims.Role)
	assert.EqualValues(claims.UserID, twitchPkg.OwnerID)

	claims = twitchPkg.CreateClaims(channelID, BroadcasterRole, FormBroadcastSendPubSubPermissions())
	assert.EqualValues(channelID, claims.ChannelID)
	assert.EqualValues(BroadcasterRole, claims.Role)
	assert.EqualValues(claims.UserID, twitchPkg.OwnerID)
}

func (t *JWTTests) TestJWTSign() {
	assert := assert.New(t.Test)

	claims := twitchPkg.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())
	signedToken, err := twitchPkg.JWTSign(claims)
	assert.NoError(err)
	assert.NotEmpty(signedToken)
}

func (t *JWTTests) TestJWTVerify() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.JWTVerify("")
	assert.NotEmpty(err)

	claims := twitchPkg.CreateClaims(channelID, ExternalRole, FormBroadcastSendPubSubPermissions())
	signedToken, err := twitchPkg.JWTSign(claims)
	assert.NoError(err)

	claims, err = twitchPkg.JWTVerify(signedToken)
	assert.NoError(err)
	assert.NotEmpty(claims)

	assert.EqualValues(claims.ChannelID, channelID)
	assert.EqualValues(claims.Role, ExternalRole)
}

func (t *ConfigurationTests) TestSetGlobalSegment() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.SetGlobalSegment(TestObj{Name: "test"})
	assert.NoError(err)
}

func (t *ConfigurationTests) TestGetGlobalSegment() {
	assert := assert.New(t.Test)

	data, err := twitchPkg.GetGlobalSegment()
	assert.NoError(err)
	assert.EqualValues(data.Configuration.Segment.Segment, GlobalSegment)
	assert.EqualValues(data.Configuration.Record.Version, twitchPkg.ConfigVersion)
	assert.EqualValues(data.Configuration.Record.Content, utils.ToJSON(TestObj{Name: "test"}))
}

func (t *ConfigurationTests) TestSetBroadcasterSegment() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.SetBroadcasterSegment(TestObj{Name: "test"}, channelID)
	assert.NoError(err)
}

func (t *ConfigurationTests) TestGetBroadcasterSegment() {
	assert := assert.New(t.Test)

	data, err := twitchPkg.GetBroadcasterSegment(channelID)
	assert.NoError(err)
	assert.EqualValues(data.Configuration.Segment.Segment, BroadcasterSegment)
	assert.EqualValues(data.Configuration.Segment.ChannelID, channelID)
	assert.EqualValues(data.Configuration.Record.Version, twitchPkg.ConfigVersion)
	assert.EqualValues(data.Configuration.Record.Content, utils.ToJSON(TestObj{Name: "test"}))
}

func (t *ConfigurationTests) TestSetDeveloperSegment() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.SetDeveloperSegment(TestObj{Name: "test"}, channelID)
	assert.NoError(err)
}

func (t *ConfigurationTests) TestGetDeveloperSegment() {
	assert := assert.New(t.Test)

	data, err := twitchPkg.GetDeveloperSegment(channelID)
	assert.NoError(err)
	assert.EqualValues(data.Configuration.Segment.Segment, DeveloperSegment)
	assert.EqualValues(data.Configuration.Segment.ChannelID, channelID)
	assert.EqualValues(data.Configuration.Record.Version, twitchPkg.ConfigVersion)
	assert.EqualValues(data.Configuration.Record.Content, utils.ToJSON(TestObj{Name: "test"}))
}

func (t *ConfigurationTests) TestGetAllChannelConfigurations() {
	assert := assert.New(t.Test)

	configs, err := twitchPkg.GetAllChannelConfigurations(channelID)
	assert.NoError(err)
	assert.NotEmpty(configs)
	assert.EqualValues(2, len(configs.Configurations))
}

func (t *ConfigurationTests) TestSetExtensionRequired() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.SetExtensionRequired(channelID)
	assert.NoError(err)
}

func (t *MessagingAndPubSubTests) TestFormPubSubPermissions() {
	assert := assert.New(t.Test)

	permissions := FormBroadcastSendPubSubPermissions()
	assert.NotEmpty(permissions)
	assert.EqualValues(1, len(permissions.Send))
	assert.EqualValues(BroadcastPublish, permissions.Send[0])

	permissions = FormGlobalSendPubSubPermissions()
	assert.NotEmpty(permissions)
	assert.EqualValues(1, len(permissions.Send))
	assert.EqualValues(GlobalPublish, permissions.Send[0])

	permissions = FormWhisperSendPubSubPermissions("UnywsWXUjrEcUMVzt_qhB")
	assert.NotEmpty(permissions)
	assert.EqualValues(1, len(permissions.Send))
	assert.EqualValues("whisper-UnywsWXUjrEcUMVzt_qhB", permissions.Send[0])

	permissions = FormGenericPubSubPermissions()
	assert.NotEmpty(permissions)
	assert.EqualValues(1, len(permissions.Send))
	assert.EqualValues(GenericPublish, permissions.Send[0])
}

func (t *MessagingAndPubSubTests) TestSendTwitchChatMessageWithMessageLimitExceeded() {
	assert := assert.New(t.Test)

	message := `this message is longer than 280 characters long
	this message is longer than 280 characters long
	this message is longer than 280 characters long
	this message is longer than 280 characters long
	this message is longer than 280 characters long
	this message is longer than 280 characters`

	_, err := twitchPkg.SendTwitchChatMessage(channelID, message)
	assert.NotEmpty(err)
	assert.Errorf(err, `this message `+message+` exceeds 280 character limit`)
}

func (t *MessagingAndPubSubTests) TestSendTwitchChatMessage() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.SendTwitchChatMessage(channelID, "hello")
	assert.NoError(err)
}

func (t *MessagingAndPubSubTests) TestPublishChannelNotification() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.PublishChannelNotification(channelID, "Get Users")
	assert.NoError(err)
}

func (t *MessagingAndPubSubTests) TestPublishWhisperNotification() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.PublishWhisperNotification(channelID, "UnywsWXUjrEcUMVzt_qhB", "Get Users")
	assert.NoError(err)
}

func (t *MessagingAndPubSubTests) TestPublishGlobalNotification() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.PublishGlobalNotification("Get Users")
	assert.NoError(err)
}

// Currently this endpoint doesn't seem to be working
// most likely because i am not testing with my channel
// live broadcasting.
func (t *ChannelInfoTests) TestGetLiveChannelsWithExtensionEnabled() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.GetLiveChannelsWithExtensionEnabled(publishedExtension, "")
	assert.NoError(err)
}

func (t *SecretTests) TestCreateSecret() {
	assert := assert.New(t.Test)

	resp, err := twitchPkg.CreateExtensionSecret(500)
	assert.NoError(err)
	assert.NotEmpty(resp.Secrets)
}

func (t *SecretTests) TestGetExtensionSecrets() {
	assert := assert.New(t.Test)

	resp, err := twitchPkg.GetExtensionSecrets()
	assert.NoError(err)
	assert.NotEmpty(resp.Secrets)
}

func (t *SecretTests) TestRevokeAllSecrets() {
	assert := assert.New(t.Test)

	_, err := twitchPkg.RevokeExtensionSecrets()
	assert.NoError(err)
}
