This is a CLI tool to upload your match history to ballchasing.com after each Rocket League session. This is needed because the bakkesmod replay uploader will not work anymore when Rocket League implements anti-cheat.

It is meant to be added to the steam launch options so it runs every time you close the game, but you can also run it manually or as a cron job.

### Advantages of this tool over [mark-codes-stuff/ballchasing_replay_uploader](https://github.com/mark-codes-stuff/ballchasing_replay_uploader)

- You don't need to save every replay manually
- Works with games played on any platform (Epic Games, Steam, PlayStation, Xbox, Nintentdo Switch)

### Disadvantage

- Up to 20 replays per session https://www.rocketleague.com/en/news/introducing-rocket-league-match-history-and-player-profiles, so you'll have to run the script manually or close Rocket League every 20 games (you should also go touch grass).

In theory you could just cron this script to run every hour and you'd never have that issue but the game disconnects from epic services when the script runs (because only one auth session at a time is allowed), so that might be annoying if it triggers while playing. I have not tried it but at the very least you will be disconnected from your party.

## Install

- Download a binary from releases
- Run it once eg in the terminal to set the auth tokens for EGS and ballchasing.com
- Add the following launch argument to Rocket League on Steam

```sh
%command%; ./<the path to the binary you downloaded>
```

- Play a game, close the game, and see it appear on ballchasing.com

### For golang enjoyers

If you have `go` installed you can run

```sh
go install github.com/Amzd/upload-match-history-to-ballchasing@latest
```
and use it like `%command%; ~/go/bin/upload-match-history-to-ballchasing`
