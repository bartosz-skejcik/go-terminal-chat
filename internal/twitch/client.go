package twitch

import (
	"github.com/bartosz-skejcik/go-terminal-chat/internal/chat"
	"github.com/gempir/go-twitch-irc/v4"
)

func NewClient(channel string, chat *chat.Chat) *twitch.Client {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		chat.ShowMessage(message)
	})

	client.Join(channel)
	return client
}

