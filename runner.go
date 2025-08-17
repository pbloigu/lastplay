package lastplay

// Overall configuration parameters
type Config struct {
	// LastFm username
	LastFmUser string
	// LastFm API key
	LastFmKey string
	// LastFm API secret
	LastFmSecret string
	// The full url of your Mastodon instance
	// Trailing slash should be omitted
	MastodonUrl string
	// Access token of your Mastodon app
	// Check the Mastodon documentation for how
	// to obtain this
	MastodonAccessToken string
}

// Interface for state storage.
//
// State storage is used to store the
// artist, track combination of the last
// successfully posted Mastodon status message.
type StateStore interface {
	// SHOULD return the last artist and track that were posted
	// to Mastodon.
	//
	// This is called internally in order to decide
	// wheter a new status message should be posted.
	LastStatus() (artist, track string)

	// Called after successfully posting a new Mastodon status.
	//
	// SHOULD store the artist, track combination to the underlying
	// strorage for retrieval by LastStatus()
	Store(artist, track string) error
}

type Runner interface {
	// Runs a one cycle of:
	//
	// 1. Fetch the last played artist, track from LastFm API
	// 2. Fetch the last artist, track from the provided StateStore
	// 3. See if the two above differ
	// 4. If they do, post a Mastodon status message formatted by given TootFormatter
	//
	// In case of any error, you get to deal with it.
	//
	// That's all folks!
	Run() error
}

// Formats the Mastodon status message (Toot) content
// from the given artist and track.
//
// Status message lenght might be limited by you instance's policy.
// Other than that, go crative with this!
type TootFormatter func(artist, track string) string

func New(c Config, s StateStore, tf TootFormatter) Runner {
	return r{
		f: f{
			lfUser:    c.LastFmUser,
			lfmKey:    c.LastFmKey,
			lfmSecret: c.LastFmSecret,
		},
		ss: s,
		t: tooter{
			mUrl:   c.MastodonUrl,
			mToken: c.MastodonAccessToken,
			tf:     tf,
		},
	}
}

type r struct {
	f  f
	ss StateStore
	t  tooter
}

func (r r) Run() error {
	artist, track, err := r.f.getLatestTrack()
	if err != nil {
		return err
	}
	if r.shouldUpdate(artist, track) {
		if err := r.t.toot(artist, track); err != nil {
			return err
		}

		if err := r.ss.Store(artist, track); err != nil {
			return err
		}
	}
	return nil
}

func (r r) shouldUpdate(artist, track string) bool {
	a, t := r.ss.LastStatus()
	return a != artist || t != track
}
