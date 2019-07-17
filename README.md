# go-twitch-ext

[git]:      https://git-scm.com/
[golang]:   https://golang.org/
[releases]: https://github.com/jackmcguire1/go-twitch-ext/releases/
[modules]:  https://github.com/golang/go/wiki/Modules
[discord]: https://discord.gg/8NXaEyV 
[MIT]: https://opensource.org/licenses/MIT

[![GoDoc](https://godoc.org/github.com/jackmcguire1/go-twitch-ext?status.svg)](https://godoc.org/github.com/jackmcguire1/go-twitch-ext)
[![Build Status](https://travis-ci.com/jackmcguire1/go-twitch-ext.svg?branch=master)](hhttps://travis-ci.org/jackmcguire1/go-twitch-ext)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/go-twitch-ext)](https://goreportcard.com/report/github.com/jackmcguire1/go-twitch-ext)
[![codecov](https://codecov.io/gh/jackmcguire1/go-twitch-ext/branch/master/graph/badge.svg)](https://codecov.io/gh/jackmcguire1/go-twitch-ext)
[![GitHub Release](https://img.shields.io/github/release-pre/jackmcguire1/go-twitch-ext.svg)](releases)

> A library to help with the development of an EBS for a [Twitch Extension](https://dev.twitch.tv/docs/extensions "twitch Extension")

> For any help please consult FAQ section
## Supported Endpoints & Features

**Features:**

> Twitch JWT
- [x] Twitch Claims structure supported
- [x] Sign Twitch claims into JWT Tokens
- [x] Verify Client/EBS Created Twitch JWT tokens into claims obj

**API Endpoint:**
>This package supports the following [Twitch Extension API endpoints](https://dev.twitch.tv/docs/extensions/reference/)


- [x] Get Live Channels with Extension Activated
- [x] Create Extension Secret
- [x] Get Extension Secret
- [x] Revoke Extension Secrets
- [x] Set Extension Required Configuration
- [x] Set Extension Configuration Segment
- [x] Get Extension Channel Configuration
- [x] Get Extension Configuration Segment
- [x] Send Extension PubSub Message
- [x] Send Extension Chat Message

## Installing
`go get github.com/jackmcguire1/go-twitch-ext`

## Example
```Go
package main

import (
	"log"
	"os"

	twitch "github.com/jackmcguire1/go-twitch-ext"
)

var twitchPkg *twitch.Twitch

func init() {
	twitchPkg = twitch.NewClient(
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
		twitch.BroadcasterRole,
		twitch.FormBroadcastSendPubSubPermissions(),
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


```

## Development

To develop `go-twitch-ext` or interact with its source code in any meaningful way, be
sure you have the following installed:

### Prerequisites

- [Git][git]
- [Go 1.12][golang]+

You will need to activate [Modules][modules] for your version of Go, generally
by invoking `go` with the support `GO111MODULE=on` environment variable set.


### Twitch Extension Configuration
From your [Twitch Extension Dashboard](https://dev.twitch.tv/dashboard/extensions) you can get the following:
- Client ID
- Base64 Secret
- Extension Version
- Extension (**Broadcaster**/**Developer**) Config Version - *OPTIONAL*

### Owner ID
To get the owner ID, you will need to execute a simple CURL command against the Twitch `/users` endpoint. You'll need your extension <b>client ID</b> as part of the query (this will be made consistent with the Developer Rig shortly, by using _owner name_).

```bash
curl -H "Client-ID: <client id>" -X GET "https://api.twitch.tv/helix/users?login=<owner name>"
```
## Create your own extension!
Get started and create your extension [today!](https://dev.twitch.tv/extensions).

## FAQ & SUPPORT
For any questions or suggestions please join the 'go-twitch-ext' channel on [Discord][discord]!

## License
The source code for go-twitch-ext is released under the [MIT License][MIT].

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)