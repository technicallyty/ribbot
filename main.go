package main

import (
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/technicallyty/vidbot/redditbot"
	"os"
	"strings"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

const MaxThreads = 10

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}
var noCtx = context.Background()
var threads int32 = 0

// handleMsg is a basic command handler
func handleMsg(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message
	go func() {
		// TODO: rate limiting the bot to 10 go routines. will update later
		if val := atomic.LoadInt32(&threads); val >= MaxThreads {
			data.Message.Reply(noCtx, s, "maximum requests reached. please try again later.")
			return
		}
		atomic.AddInt32(&threads, 1)
		s.Channel(msg.ChannelID).TriggerTypingIndicator()
		stripped := strings.ReplaceAll(msg.Content, " ", "")
		vb := redditbot.NewVidBot(stripped)
		log.Infof("received link: %v\n", stripped)
		if vb.SetResourceURL(stripped) {
			path, resource, err := vb.Download()
			if err != nil {
				log.Errorf("download error: %v\n", err)
				msg.Reply(noCtx, s, err)
				atomic.AddInt32(&threads, -1)
				return
			}

			log.Infof("opening file: %v\n", path)
			videoFile, err := os.Open(path)
			if err != nil {
				msg.Reply(noCtx, s, err)
				atomic.AddInt32(&threads, -1)
				return
			}
			defer videoFile.Close()
			log.Info("uploading video....")
			_, err = msg.Reply(noCtx, s, &disgord.CreateMessageParams{
				Content:    "",
				Nonce:      "",
				Tts:        false,
				Embed:      nil,
				Components: nil,
				Files: []disgord.CreateMessageFileParams{
					{
						Reader:     videoFile,
						FileName:   resource + ".mp4",
						SpoilerTag: false,
					},
				},
				SpoilerTagContent:        false,
				SpoilerTagAllAttachments: false,
				AllowedMentions:          nil,
				MessageReference:         nil,
			})
			if err != nil {
				log.Errorf("reply error: %v\n", err)
				msg.Reply(noCtx, s, err)
				atomic.AddInt32(&threads, -1)
				return
			} else {
				log.Info("video uploaded")
			}
			split := strings.Split(path, "/")
			split = split[:len(split)-1]
			err = os.RemoveAll(strings.Join(split, "/"))
			if err != nil {
				log.Errorf("err removing dir: %v\n", err)
			}
		} else {
			log.Infof("invalid reddit link: %v\n", msg.Content)
			msg.Reply(noCtx, s, fmt.Sprintf("%v is not a valid reddit link", msg.Content))
		}
		atomic.AddInt32(&threads, -1)
	}()
}

const prefix = "r/"

func main() {
	client := disgord.New(disgord.Config{
		ProjectName: "ribbot",
		BotToken:    "ODY1MDYwMzgyOTk3NjEwNTA3.YO-gQw.mZSJmD_AMMlmjQzeEX9VXZ8ihac",
		Logger:      log,
		RejectEvents: []string{
			// rarely used, and causes unnecessary spam
			disgord.EvtTypingStart,

			// these require special privilege
			// https://discord.com/developers/docs/topics/gateway#privileged-intents
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		// ! Non-functional due to a current bug, will be fixed.
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "write " + prefix + "[reddit link]",
			},
		},
		DMIntents: disgord.IntentDirectMessages | disgord.IntentDirectMessageReactions | disgord.IntentDirectMessageTyping,
		// comment out DMIntents if you do not want the bot to handle direct messages
	})

	defer client.Gateway().StayConnectedUntilInterrupted()

	logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	filter.SetPrefix(prefix)

	// create a handler and bind it to new message events
	// thing about the middlewares are whitelists or passthrough functions.
	client.Gateway().WithMiddleware(
		filter.NotByBot,    // ignore bot messages
		filter.HasPrefix,   // message must have the given prefix
		logFilter.LogMsg,   // log command message
		filter.StripPrefix, // remove the command prefix from the message
	).MessageCreate(handleMsg)

	// create a handler and bind it to the bot init
	// dummy log print
	client.Gateway().BotReady(func() {
		log.Info("Bot is ready!")
	})
}
