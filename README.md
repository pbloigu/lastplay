# lastplay

## Summary
This is a Mastodon bot backend to post (that is, to toot) the latest played track as seen by LastFm.

## Details
Ever since having found out Deezer effectively has no API (they have, they just don't allow new registrations), I've been looking for a way to post my "Now playing" status to Mastodon (using a bot account separate of my main account).

This little piece of software tries to fulfill that need. It does it by performing 3 steps in succession:

1. Fetch the last played track from Last.fm
2. See if it differs from the one fetched on the previous run
3. If it does, post a toot and make note of the track for the next run.

This can then be set to run perodically e.g. with cron or systemd timers on *nix systems.

It's obvious the updates are not real time because first of all it dependes on the interval of consecutive executions and secondly because LastFm status does not update in real time. Therefore instead of "now playing" this is more like "last played", hence the name of the app. This however is as close to "now playing" as it gets with the integration options currently available.

## Usage
Configure the following environment variables:

|Variable|Value|
|-|-|
|LP_LFM_USER|Your Last.fm username|
|LP_LFM_KEY|Your Last.fm API key|
|LP_LFM_SECRET|Your Last.fm API secret|
|LP_MASTODON_URL|Url of your Mastodon instance|
|LP_MASTODON_TOKEN|Your Mastodon API token|

### Where / how to acquire these details?
To obtain Last.fm API key & secret you need to register an application by going here: https://www.last.fm/api/account/create Obviously you need to be logged in to Last.fm. Please use the API responsibly and make sure to read any TOS associated with using the API. It's your problem if you get banned.

As for the Mastodon API token: you need to log in to your instance, go to Prefereces (in the default Mastodon web UI), select Development and then New application. Give your application a name and access to scope "write:statuses" in order for this bot backend to be able to post toots.

Once the application is added you need to click on it once more to see the authentication details.  From the page that opens, grab the token access token.

### Running
Then just set the app to run periodically, e.g. once a minute. Your bot should then toot something like the following:

<img width="606" height="185" alt="image" src="https://github.com/user-attachments/assets/d0b21688-2dd8-4ce6-8dc6-d2d403a2e3a5" />


## Caveats
It should be noted that the status file which stores the last track fetched from Last.fm API is stored in a file in the system's TEMP directory (which usually on *nix systems is /tmp). Depending on the OS configuration the TEMP directory might be wiped upon boot, resulting in possible duplicate toots if the system is restarted. This might be improved in the future versions.

## As a library
In addition to running the published binary, there's the option of writing your own bot backend. More details on this to follow, for now check [main.go](./cmd/main.go)
