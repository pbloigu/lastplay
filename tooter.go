package lastplay

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type emoji = rune

const (
	BLACK_HEART     emoji = rune(0x1F5A4)
	HORNS           emoji = rune(0x1F918)
	GUITAR          emoji = rune(0x1F3B8)
	SKULL           emoji = rune(0x1F480)
	SKULL_AND_BONES emoji = rune(0x2620)
	DEVIL           emoji = rune(0x1F608)
)

var emojis = []emoji{BLACK_HEART, HORNS, GUITAR, SKULL, SKULL_AND_BONES, DEVIL}

type tooter struct {
	mUrl   string
	mToken string
	tf     TootFormatter
}

func (t tooter) toot(artist, track string) error {
	msg := t.tf(artist, track)
	r, err := t.createRequest(msg)
	if err != nil {
		return err
	}
	if err = t.doToot(r); err != nil {
		return err
	}
	return nil
}

// func (t t) createToot(artist, track string) string {
// 	fn := func() emoji {
// 		return emojis[rand.Intn(len(emojis))]
// 	}
// 	return fmt.Sprintf("%c %s ~~ %s %c", fn(), artist, track, fn())
// }

func (t tooter) createRequest(toot string) (http.Request, error) {

	data := url.Values{}
	data.Set("status", toot)

	u, err := url.ParseRequestURI(t.mUrl)
	if err != nil {
		return http.Request{}, err
	}
	u.Path = "/api/v1/statuses"

	r, _ := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	r.Header.Add("Authorization", "Bearer "+t.mToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Idempotency-Key", uuid.NewString())

	return *r, nil
}

func (t tooter) doToot(r http.Request) error {
	client := &http.Client{}
	resp, err := client.Do(&r)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Toot failed.")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error().Int("status", resp.StatusCode).Msg("Mastodon responded with a failude code.")
		return fmt.Errorf("Mastodon response code: %d", resp.StatusCode)
	}
	return nil
}
