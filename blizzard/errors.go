package blizzard

import "github.com/whitesmith/unplugg-of-warcraft"

const (
	errInvalidURL      warcraft.Error = "invalid url"
	errFailedRequest   warcraft.Error = "failed http request"
	errInvalidResponse warcraft.Error = "invalid http response"
	errGetDump         warcraft.Error = "failed to get new api dump"
	errFilterDump      warcraft.Error = "failed to validate auctions"
	errSaveDump        warcraft.Error = "failed to save new api dump"
)
