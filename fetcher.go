package lastplay

import (
	"github.com/rs/zerolog/log"
)

type f struct {
	lbUser      string
	lbUserAgent string
}

func (f f) getLatestTrack() (artist, track *string, err error) {
	lb := lb{
		userName:  f.lbUser,
		userAgent: f.lbUserAgent,
	}

	rt, err := lb.getPlayingNow()

	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to access ListenBrainz.")
		return nil, nil, err
	}

	if rt.PlayingNow && len(rt.Listens) > 0 {
		return &rt.Listens[0].TrackMetadata.ArtistName, &rt.Listens[0].TrackMetadata.TrackName, nil
	} else {
		return nil, nil, nil
	}
}
