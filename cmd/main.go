package main

import (
	"fmt"
	"log"

	"github.com/bartosz-skejcik/go-terminal-chat/internal/chat"
	"github.com/bartosz-skejcik/go-terminal-chat/internal/helper"
	twitchClient "github.com/bartosz-skejcik/go-terminal-chat/internal/twitch"
)

func main() {
	// Load configuration
	config, err := helper.NewConfigManager()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create channel manager
	channelManager := helper.NewChannelManager(config)

	// Get the Twitch channel
	channel, err := channelManager.GetChannel()
	if err != nil {
		log.Fatalf("Failed to get channel: %v", err)
	}

	// Clear the terminal screen
	helper.ClearScreen()

	// Greeting
	fmt.Printf("Welcome to %s's chat!\n", channel)
	fmt.Printf("Press Ctrl+C to exit.\n\n")

	// Create chat instance with configuration settings
	chatInstance := chat.NewChat(config.Runner.ShowTimestamps, config.Runner.LogMessages)

	// Create Twitch client
	client := twitchClient.NewClient(channel, chatInstance)

	// Start the client
	if err := client.Connect(); err != nil {
		log.Fatalf("Error connecting to Twitch chat: %v", err)
	}
}
