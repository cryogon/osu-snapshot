# Osu Snapshot
take a snapshot of osu!

copies
- Replays
- *.db files
- and beatmap ids to reduce size which can be downloaded later

## My motivation to make this
so i can uninstall the game and take a break indefinately (or potentially quit w)

## Prerequistes
- go

## Steps to use
- update osuPath variable in `backup.go` if needed
- run `go run backup.go` to generate a backup
- upload the backup folder to any cloud storage :)


## WIP
### backup.go
- update it so it generates a zip file at the end instead of directory
- add a proper progress bar
- show logs only if debug parameter gets passed
- make osuPath more easily configarable, potential in some other config file

### restore.go (not implemented yet)
- need to create a script to restore the backup
    - downloading all the beatmaps from different mirrors, to not get ip banned
    - copying all the data to osu! directory
    - which includes proper merging of data, so all data remains
