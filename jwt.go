package twitchext

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// RoleType The user role type
type RoleType string

// Types of user roles used within the JWT Claims
const (
	BroadcasterRole RoleType = "broadcaster"
	ExternalRole    RoleType = "external"
	ModeratorRole   RoleType = "moderator"
	ViewerRole      RoleType = "viewer"

	toAllChannels = "all"
)

// TwitchJWTClaims contains information
// containing twitch specific JWT information.
type TwitchJWTClaims struct {
	OpaqueUserID string             `json:"opaque_user_id,omitempty"`
	UserID       string             `json:"user_id"`
	ChannelID    string             `json:"channel_id,omitempty"`
	Role         RoleType           `json:"role"`
	Unlinked     bool               `json:"is_unlinked,omitempty"`
	Permissions  *PubSubPermissions `json:"pubsub_perms"`
	jwt.StandardClaims
}

// CreateClaims will construct a claims suitable for generating a JWT token,
// containing necessary information required by the Twitch API.
// @param channelID if this value is empty it will default to 'all'
// @param role if this value is empty it will default to 'external'
func (t *Twitch) CreateClaims(
	channelID string,
	role RoleType,
	permissions *PubSubPermissions,
) (
	claims *TwitchJWTClaims,
) {
	expiration := time.Now().Add(time.Minute*3).UnixNano() / int64(time.Millisecond)
	//expiration := time.Now().AddDate(1000, 0 ,0 ).UnixNano() / int64(time.Millisecond)
	if role == "" {
		role = ExternalRole
	}

	if channelID == "" {
		channelID = toAllChannels
	}

	claims = &TwitchJWTClaims{
		UserID:      t.OwnerID,
		ChannelID:   channelID,
		Role:        role,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	return
}

// JWTSign Sign the a JWT Claim to produce a base64 token.
func (t *Twitch) JWTSign(claims *TwitchJWTClaims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := base64.StdEncoding.DecodeString(t.Secret)
	if err != nil {
		return
	}

	tokenString, err = token.SignedString(key)
	if err != nil {
		return
	}

	return
}

// JWTVerify validates a extension client side twitch base64 token and converts it
// into a twitch claim type, containing relevant information.
func (t *Twitch) JWTVerify(token string) (claims *TwitchJWTClaims, err error) {
	if token == "" {
		err = fmt.Errorf("JWT token string missing")
		return
	}

	parsedToken, err := jwt.ParseWithClaims(token, &TwitchJWTClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", tkn.Header["alg"])
		}

		key, err := base64.StdEncoding.DecodeString(t.Secret)

		if err != nil {
			return nil, err
		}
		return key, nil
	})

	if err != nil {
		return
	}

	claims, ok := parsedToken.Claims.(*TwitchJWTClaims)
	if !ok || !parsedToken.Valid {
		err = fmt.Errorf("Could not parse JWT")
		return
	}

	return
}
