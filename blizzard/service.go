package blizzard

import (
	"github.com/whitesmith/unplugg-of-warcraft"
)

// Service represents a service for interacting with the blizzard API.
type Service struct {
	client *Client
}

// GetAPIDump retrieves an HA API dump.
func (s *Service) GetAPIDump(realm, locale, key string, last int64) (*warcraft.APIDump, error) {
	// get request url.
	r, err := s.client.GetDumpURL(realm, locale, key)
	if err != nil {
		return nil, warcraft.ErrGetDump
	} else if r.Modified <= last {
		return nil, warcraft.ErrDumpExists
	}

	// get api dump.
	d, err := s.client.GetDump(r.URL)
	if err != nil {
		return nil, warcraft.ErrGetDump
	}

	// set dump timestamp.
	d.Timestamp = r.Modified

	return d, nil
}

// ValidateAuctions filters invalid/unnecessary auctions.
func (s *Service) ValidateAuctions(dump *warcraft.APIDump) ([]warcraft.Auction, error) {
	as := make([]warcraft.Auction, 0)
	for _, a := range dump.Auctions {
		if a.Timeleft != "SHORT" && a.Buyout != 0 {
			a.Timestamp = dump.Timestamp
			as = append(as, a)
		}
	}
	return as, nil
}
