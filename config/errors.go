package config

import "github.com/whitesmith/unplugg-of-warcraft"

const (
	errInvalidFile warcraft.Error = "failed to open config file"
	errFailedParse warcraft.Error = "failed to parse config file"
)
