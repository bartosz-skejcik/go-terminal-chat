package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ChannelManager handles channel-related configuration and user interactions
type ChannelManager struct {
	config *Config
}

// NewChannelManager creates a new ChannelManager
func NewChannelManager(config *Config) *ChannelManager {
	return &ChannelManager{
		config: config,
	}
}

// GetChannelFromArgs extracts channel name from command line arguments
func (cm *ChannelManager) GetChannelFromArgs(args []string) (string, bool) {
	if len(args) > 1 {
		return args[1], true
	}
	return "", false
}

// PromptForChannel interactively helps user configure Twitch channel
func (cm *ChannelManager) PromptForChannel() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	// If a channel is already configured, offer to change it
	if cm.config.Twitch.Channel != "" {
		fmt.Printf("Current Twitch channel is set to: %s\n", cm.config.Twitch.Channel)
		fmt.Println("Would you like to change the channel? (yes/no)")

		response, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %v", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "yes" || response == "y" {
			return cm.askAndUpdateChannel(reader)
		}

		return cm.config.Twitch.Channel, nil
	}

	// No channel configured, prompt for a new one
	fmt.Println("No Twitch channel is currently configured.")
	return cm.askAndUpdateChannel(reader)
}

// askAndUpdateChannel handles the process of asking for and updating a channel name
func (cm *ChannelManager) askAndUpdateChannel(reader *bufio.Reader) (string, error) {
	fmt.Println("Please enter the Twitch channel name:")

	channelInput, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	channel := strings.TrimSpace(channelInput)

	// Validate channel name
	if channel == "" {
		fmt.Println("Invalid channel name. No channel was set.")
		return "", fmt.Errorf("empty channel name")
	}

	// Update the channel in the configuration
	if err := cm.config.Update("twitch.channel", channel); err != nil {
		return "", fmt.Errorf("failed to update channel: %v", err)
	}

	fmt.Printf("Twitch channel set to: %s\n", channel)
	return channel, nil
}

// GetChannel determines the Twitch channel to use
func (cm *ChannelManager) GetChannel() (string, error) {
	// First, check command-line arguments
	args := os.Args
	if channel, found := cm.GetChannelFromArgs(args); found {
		return channel, nil
	}

	// If no channel in args, use configured or prompt
	return cm.PromptForChannel()
}

// InitializeChannelManager is a convenience function to create a new ChannelManager
func InitializeChannelManager() (*ChannelManager, error) {
	// Load configuration
	config, err := NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return NewChannelManager(config), nil
}
