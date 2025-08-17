package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/pbloigu/lastplay"
)

type Emoji = rune

const (
	BLACK_HEART     Emoji = rune(0x1F5A4)
	HORNS           Emoji = rune(0x1F918)
	GUITAR          Emoji = rune(0x1F3B8)
	SKULL           Emoji = rune(0x1F480)
	SKULL_AND_BONES Emoji = rune(0x2620)
	DEVIL           Emoji = rune(0x1F608)
)

var emojis = []Emoji{BLACK_HEART, HORNS, GUITAR, SKULL, SKULL_AND_BONES, DEVIL}

type tmpFileStore struct {
}

func (ts tmpFileStore) LastStatus() (artist, track string) {
	data, err := os.ReadFile(os.TempDir() + "/lastplay.dat")
	if err != nil {
		fmt.Printf("Failed to read status file: %s\n", err)
		return "", ""
	}

	split := strings.Split(string(data), "\n")
	if len(split) != 2 {
		fmt.Printf("Status file has an incorrect number of lines: %d\n", len(split))
		return "", ""
	}
	return split[0], split[1]
}

func (ts tmpFileStore) Store(artist, track string) error {
	d1 := []byte(fmt.Sprintf("%s\n%s", artist, track))
	if err := os.WriteFile("/tmp/lastplay.dat", d1, 0644); err != nil {
		return err
	}
	return nil
}

func main() {
	lp := lastplay.New(configure(), &tmpFileStore{}, func(artist, track string) string {
		fn := func() Emoji {
			return emojis[rand.Intn(len(emojis))]
		}
		return fmt.Sprintf("%c %s ~~ %s %c", fn(), artist, track, fn())
	})
	if err := lp.Run(); err != nil {
		panic(err)
	}
}

func configure() lastplay.Config {
	return lastplay.Config{
		LastFmUser:          getEnvOrDie("LP_LFM_USER"),
		LastFmKey:           getEnvOrDie("LP_LFM_KEY"),
		LastFmSecret:        getEnvOrDie("LP_LFM_SECRET"),
		MastodonUrl:         getEnvOrDie("LP_MASTODON_URL"),
		MastodonAccessToken: getEnvOrDie("LP_MASTODON_TOKEN"),
	}
}

func getEnvOrDie(key string) string {
	env := os.Getenv(key)
	if env == "" {
		panic(fmt.Errorf("Missing environment variable: %s", key))
	}
	return env
}
