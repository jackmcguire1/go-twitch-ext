package main

import (
	"log"
	"os"

	"github.com/jackmcguire1/go-twitch-ext"
)

var twitchPkg *twitchext.Twitch

func init() {
	twitchPkg = twitchext.NewClient(
		os.Getenv("OWNER_ID"),
		os.Getenv("CLIENT_ID"),
		os.Getenv("EXT_SECRET"),
		os.Getenv("EXT_VERSION"),
		os.Getenv("EXT_CONFIG_VER"),
	)
}

func main() {
	claims := twitchPkg.CreateClaims(
		"35851594",
		twitchext.BroadcasterRole,
		twitchext.FormBroadcastSendPubSubPermissions(),
	)
	token, err := twitchPkg.JWTSign(claims)
	if err != nil {
		log.Fatal(err)
	}

	claims, err = twitchPkg.JWTVerify(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(claims)
}
