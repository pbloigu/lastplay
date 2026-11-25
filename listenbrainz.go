package lastplay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

const ListenBrainzApi string = "https://api.listenbrainz.org/"

type lb struct {
	userName  string
	userAgent string
}
type lbResponse struct {
	Payload playingNow `json:"payload"`
}

type playingNow struct {
	UserId     string    `json:"user_id"`
	PlayingNow bool      `json:"playing_now"`
	Listens    []listens `json:"listens"`
}

type listens struct {
	TrackMetadata struct {
		ArtistName string `json:"artist_name"`
		TrackName  string `json:"track_name"`
	} `json:"track_metadata"`
}

func (lb lb) getPlayingNow() (*playingNow, error) {
	url := fmt.Sprintf("%s/1/user/%s/playing-now", ListenBrainzApi, lb.userName)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to initialize request.")
		return nil, err
	}
	req.Header["User-Agent"] = []string{lb.userAgent}

	resp, err := client.Do(req)

	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to access ListenBrainz.")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to read ListenBrainz response.")
		return nil, err
	}

	target := lbResponse{}
	err = json.Unmarshal(body, &target)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Failed to parse ListenBrainz response.")
		return nil, err
	}

	return &target.Payload, nil

}
