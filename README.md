# lastfm-cli

A fast, minimal **command-line interface for Last.fm** written in Go.

The goal of this project is to provide a clean terminal interface for interacting with the Last.fm API: viewing scrobbles, checking currently playing tracks, and exploring listening statistics.

---

# Project Goals

* Simple CLI interface
* Fast execution (single binary)
* Clean terminal output
* No unnecessary dependencies
* Works well on Linux systems

---

# Core Features (MVP)

Implement these first.

### 1. Configuration

Store configuration locally.

Config file:

```
~/.config/lastfm-cli/config.yaml
```

Example:

```
api_key: YOUR_LASTFM_API_KEY
username: YOUR_LASTFM_USERNAME
```

Tasks:

* load config file
* fallback to environment variables
* allow command to initialize config

Command:

```
lastfm init
```

---

### 2. Recent Tracks

Command:

```
lastfm recent
```

Example output:

```
1. Radiohead — Paranoid Android
2. Massive Attack — Teardrop
3. Daft Punk — Voyager
```

API endpoint:

```
user.getRecentTracks
```

Tasks:

* fetch recent tracks
* parse JSON
* print clean terminal output

---

### 3. Currently Playing

Command:

```
lastfm now
```

Output example:

```
Now playing:
Daft Punk — Digital Love
Album: Discovery
```

Implementation:

* call `user.getRecentTracks`
* check `@attr.nowplaying`

---

### 4. Top Artists

Command:

```
lastfm top artists
```

Example output:

```
1. Radiohead
2. Daft Punk
3. Massive Attack
```

API endpoint:

```
user.getTopArtists
```

---

# CLI Structure

Recommended structure:

```
lastfm-cli/
├─ cmd/
│  ├─ root.go
│  ├─ now.go
│  ├─ recent.go
│  └─ top.go
│
├─ internal/
│  ├─ api/
│  │  └─ client.go
│  ├─ config/
│  │  └─ config.go
│  └─ models/
│     └─ track.go
│
├─ main.go
└─ go.mod
```

---

# API Client

Create a reusable client:

```
type Client struct {
    APIKey string
}
```

Base URL:

```
https://ws.audioscrobbler.com/2.0/
```

Requests should include:

```
method
api_key
format=json
```

Example:

```
user.getRecentTracks
```

---

# Libraries

Recommended Go libraries:

CLI framework:

* Cobra

Configuration:

* Viper

HTTP:

* net/http (standard library)

---

# Development Steps (Do Tomorrow)

1. Initialize Go module
2. Create CLI structure
3. Implement config loader
4. Build API client
5. Implement `recent` command
6. Implement `now` command
7. Implement `top artists`
8. Improve output formatting

---

# Future Features

Ideas to add later.

### Watch Mode

```
lastfm watch
```

Live updates of currently playing track.

---

### Export Scrobbles

```
lastfm export --format csv
```

Download listening history.

---

### Terminal Dashboard

Add a TUI mode:

```
lastfm tui
```

Possible library:

* Bubble Tea

Dashboard could show:

* current track
* recent scrobbles
* top artists
* listening stats

---

# Installation

Build locally:

```
go build -o lastfm
```

Install:

```
install -Dm755 lastfm ~/.local/bin/lastfm
```

---

# Long-Term Ideas

* caching API responses
* colored terminal output
* album artwork support
* support for authentication actions
* scrobbling from CLI

---

# Notes

* The Last.fm API key is safe to include in local configs
* Authenticated actions require a session key
* Respect API rate limits
* Avoid unnecessary requests

---

# End Goal

Create a **fast, minimal, Unix-friendly Last.fm CLI** that behaves similarly to tools like `gh` or `yt-dlp`, but focused on music listening data.
