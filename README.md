# Discord DCA player

A simple discordgo library for streaming audio files to voice channel. 

# Usage

```go
player.Play("/path/to/audiofile", vc)
```

The audiofile can be of any file type. For example a mp3-file will be auto translated to dca and persisted in storage.

# Requirements
You will need to have ffmpeg installed and the dca binary needs to be placed into the directory of the binary that runs the discord bot
in order to ensure safe translation of the audiofiles into dca format.