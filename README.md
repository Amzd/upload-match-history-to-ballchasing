This is a CLI tool to upload your match history to ballchasing.com after each Rocket League session.

It is meant to be added to the steam launch options so it runs every time you close the game.

## Install

- Download a binary from releases
- Run it once and set the auth tokens
- Add the following launch argument to rocket league

```sh
%command%; ./<the path to the binary you downloaded>
```

- Play a game and see it appear on ballchasing

### For developers

If you have `go` installed you can run

```sh
go install github.com/Amzd/upload-match-history-to-ballchasing@latest
```
and use it like `%command%; ~/go/bin/upload-match-history-to-ballchasing`
