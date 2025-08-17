package lastplay

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequest(t *testing.T) {
	tooter := tooter{
		mUrl:   "http://foo.bar",
		mToken: "api token",
	}

	req, err := tooter.createRequest("Hello world")
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "status=Hello+world", string(body))
	assert.Equal(t, tooter.mUrl+"/api/v1/statuses", req.URL.String())
	assert.Equal(t, "Bearer "+tooter.mToken, req.Header.Get("Authorization"))
	assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
	assert.NotEmpty(t, req.Header.Get("Idempotency-Key"))
}

func TestDoToot(t *testing.T) {
	tooter := tooter{

		mToken: "api token",
		tf: func(artist, track string) string {
			return artist + ":" + track
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, "/api/v1/statuses", req.URL.String())
		assert.Equal(t, "Bearer "+tooter.mToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		assert.NotEmpty(t, req.Header.Get("Idempotency-Key"))

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "status=Hello+world", string(body))
		rw.WriteHeader(200)
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	tooter.mUrl = server.URL
	r, err := tooter.createRequest("Hello world")
	if err != nil {
		t.Error(r)
	}
	if err = tooter.doToot(r); err != nil {
		t.Error(err)
	}
}
