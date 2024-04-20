# discordbot
A discord bot on golang 1.22.1 which supports registered commands */servers*, */roll*, */updatelog*.

More commands will be added soon, keep in touch with changes.

# Information
I decided to make a bot which tracks online statistics on Majestic RP servers, updates info about them every 10 seconds in a json file and in an embed message via */servers* command. */roll* will throw a random number from 1 to 6 just like in a chatbox on gamesense.pub

# Usage
Just replace your bot token and application id on line 28, 29. You can find them on [Discord Dev Portal](https://discord.com/developers/applications)
```go
var botToken = "Replace with your bot token here"
var applicationID = "Replace with application id here"
```

Then run your application and have fun!
