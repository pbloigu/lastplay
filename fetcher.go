package lastplay

import (
	"github.com/rs/zerolog/log"

	"github.com/pazuzu156/lastfm-go"
)

type f struct {
	lfUser    string
	lfmKey    string
	lfmSecret string
}

func (f f) getLatestTrack() (artist, track string, err error) {
	api := lastfm.New(f.lfmKey, f.lfmSecret)

	rt, err := api.User.GetRecentTracks(lastfm.P{
		"limit":   1,
		"user":    f.lfUser,
		"api_key": f.lfmKey,
	})
	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to access LastFm.")
		return "", "", err
	}

	return rt.Tracks[0].Artist.Name, rt.Tracks[0].Name, nil
}
