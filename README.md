# ribbot
:warning: this is code is very weekend-warrior and hacked together. It barely works. But it does work. Anyway - use with discretion :warning:

## What does it do?

Ribbot will respond to commands that include reddit links and proceed to download, and upload the video associated with the link, if any.

example command:

```user41039: r/https://reddit.com/r/TikTokCringe/comments/31jfk3/Some-Post-Title/```

## Hosting

- Build the binary for the target machine. As an example, this is the command I ran to build a binary for the Ubuntu, intel powered VM I use:
`env GOOS=linux GOARCH=386 go build -o ribbot`
- Run the binary. If you SSH into your VM like I do, you'll want to use tmux so the running binary doesn't die when you exit ssh. To do this simply:
```
$tmux
$./ribbot
```
- press ctrl+b then press d to exit tmux.

### Acknowledgements
This code is powered by [disgord](https://github.com/andersfylling/disgord/) - a discord module for golang.
