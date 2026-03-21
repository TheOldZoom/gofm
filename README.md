# gofm

A small terminal CLI for checking a Last.fm user's recent tracks.

`gofm` is a minimal Go command-line app that talks to the Last.fm API, stores your local config, and prints recent listening history in a clean terminal format.

## What you get

- Automatic first-run setup for your Last.fm username and API key
- Local config stored in your XDG config directory
- Recent-track lookup for your saved user or any username you pass explicitly
- Simple single-binary workflow
- No auth flow or session handling required for read-only lookups

---

## Install

### Option 1: Download from Releases

Open [Releases](https://github.com/theOldZoom/gofm/releases) and download the file for your system:

| I use...              | Download this file  |
| --------------------- | ------------------- |
| Linux (x86_64)        | `gofm-linux-amd64`  |
| Linux (ARM64)         | `gofm-linux-arm64`  |
| macOS (Intel)         | `gofm-darwin-amd64` |
| macOS (Apple Silicon) | `gofm-darwin-arm64` |
| Windows (x86_64)      | `gofm-windows-amd64.exe` |
| Windows (ARM64)       | `gofm-windows-arm64.exe` |

Then make it executable and run it:

```bash
chmod +x gofm-linux-amd64
./gofm-linux-amd64 recent
```

### Option 2: Install to `~/.local/bin`

This installs the latest matching release as `~/.local/bin/gofm`:

```bash
curl -fsSL https://raw.githubusercontent.com/theoldzoom/gofm/master/install.sh | sh
```

If `~/.local/bin` is already in your `PATH`, you can then run:

```bash
gofm recent
```

### Option 3: Build from source

You need [Go](https://go.dev/) installed.

```bash
git clone https://github.com/theOldZoom/gofm.git
cd gofm
make build
./build/gofm recent
```

## Usage

If no config exists yet, `gofm` launches interactive setup and saves your details to:

```text
~/.config/gofm/config.yaml
```

Example config:

```yaml
username: your_lastfm_username
api_key: your_lastfm_api_key
```

### API key

`gofm` needs a Last.fm API key for requests to the Last.fm API.

Create one here:

[Last.fm API account page](https://www.last.fm/api/account/create)

or view your accounts here:

[https://www.last.fm/api/accounts](https://www.last.fm/api/accounts)



## Notes

- First-run setup validates both your API key and username before saving config
- Use `--config` to point to a custom config file if needed
- `gofm` also reads environment variables through Viper, so `USERNAME` and `API_KEY` can be used as fallbacks
- The project currently focuses on read-only lookups

---

## Tech stack

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra) for the CLI
- [Viper](https://github.com/spf13/viper) for configuration
- `net/http` from the Go standard library for API requests

---

Here’s a more structured, actionable, and hierarchical **TODO for your Last.fm CLI**, reflecting commands, subcommands, flags, and implementation priorities. This is closer to a “developer roadmap” than a flat list.

---

## GoFM TODO 

### 1. Core Commands (Must-Have)

* `now` – Show currently playing track

  * Flags: `--user`, `--json`, `--verbose`
* `recent` – Show recently played tracks

  * Flags: `--user`, `--limit`, `--json`

### 2. Top Lists

* `top artists` – Show user’s top artists
* `top tracks` – Show user’s top tracks
* Shared flags: `--user`, `--limit`, `--period`, `--json`

### 3. Info Commands

* `info artist <name>` – Show bio, tags, stats
* `info track <artist> <track>` – Show track info, playcount, listeners
* `info album <artist> <album>` – Show album details and tags
* Flags: `--json`, `--verbose`

### 4. User & Stats

* `user info <username>` – Show general user stats
* `stats weekly` – Show weekly listening stats
* `stats monthly` – Show monthly listening stats
* `stats favorite` – Show most played track/artist/album
* Flags: `--user`, `--json`, `--limit`

### 5. Interaction Commands

* `love <track>` / `unlove <track>` – Mark tracks as loved/unloved
* `scrobble <track> <artist>` – Manually scrobble a track
* `tag track <track>` / `album <album>` / `artist <artist>` – Add/remove tags

### 6. Advanced Features

* `compare <user1> <user2>` – Compare top tracks/artists
* `trending` – Show trending tracks/artists globally
* `recommend` – Suggest tracks based on top artists
* `playlist export <format>` – Export top tracks to M3U/CSV


### 7. Output & UX Improvements

* Support JSON output globally (`--json`)
* Add terminal formatting: tables, colors, alignment
* Pagination for long lists

### 8. Developer & Maintenance

* Add tests for config loading and API calls
* Automatic rate-limit handling or retries


---

## Contributing

Contributions are welcome.

```bash
git clone https://github.com/theOldZoom/gofm.git
cd gofm
go build ./...
```

Before opening a PR:

- Keep changes focused
- Run `go test ./...`
- Run `go build ./...`
- Follow the existing project style

## License

MIT. See `[LICENSE](LICENSE)`.
