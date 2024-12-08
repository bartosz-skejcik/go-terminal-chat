package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gempir/go-twitch-irc/v4"
)

// Badge represents a Twitch badge with style and icon.
type Badge struct {
	Name            string
	Color           string
	ForegroundColor string // New field for foreground color
	Icon            string
}

// Chat handles chat message processing and display.
type Chat struct {
	showTimestamp bool
	logMessages   bool
	badgeStyles   map[string]Badge
	userStyle     lipgloss.Style
	messageStyle  lipgloss.Style
}

// NewChat creates a new Chat instance.
func NewChat(showTimestamp bool, logMessages bool) *Chat {
	c := &Chat{
		showTimestamp: showTimestamp,
		logMessages:   logMessages,
		badgeStyles: map[string]Badge{
			"premium":           {Name: "premium", Color: "#ADD8E6", ForegroundColor: "#FFA500", Icon: " "},        // Gold
			"subscriber":        {Name: "subscriber", Color: "#32CD32", ForegroundColor: "#FFF", Icon: " "},        // Light Green
			"sub-gift-leader":   {Name: "sub-gift-leader", Color: "#FF69B4", ForegroundColor: "#FFF", Icon: " "},   // Pink
			"moderator":         {Name: "moderator", Color: "#0000FF", Icon: "󰓥 "},                                  // Blue
			"hype-train":        {Name: "hype-train", Color: "#FFA500", Icon: " "},                                 // Orange
			"subtember-2024":    {Name: "subtember-2024", Color: "#800080", ForegroundColor: "#F7820F", Icon: "󰮿 "}, //Purple
			"partner":           {Name: "partner", Color: "#D776FF", Icon: " "},                                    // Purple
			"twitch-recap-2023": {Name: "twitch-recap-2023", Color: "#9146FF", Icon: " "},                          // Red
			"glitchcon2020":     {Name: "glitchcon2020", Color: "#F0ABFC", Icon: "🦖"},                               // Pink
			"vip":               {Name: "vip", Color: "#DB2777", Icon: "󰮊 "},                                        // Green
			"broadcaster":       {Name: "broadcaster", Color: "#DC2626", Icon: " "},                                // Purple
			"cheer":             {Name: "cheer", Color: "#ffd700", Icon: " "},                                      // Purple
			// Add more badge styles as needed. Replace icons with your preferred Nerdfonts.
		},
		userStyle:    lipgloss.NewStyle().Bold(true),
		messageStyle: lipgloss.NewStyle().Padding(0, 1), //Light Pink
	}
	return c
}

// ShowMessage formats and returns a styled chat message string.
func (c *Chat) ShowMessage(message twitch.PrivateMessage) {
	userColor := message.User.Color
	if userColor == "" {
		userColor = "#FFFFFF" // Default to white if no color is specified
	}

	badges := formatBadges(message, *c)

	timestamp := formatTimestamp(*c)

	userName := c.userStyle.Foreground(lipgloss.Color(userColor)).Render(message.User.DisplayName) // Use DisplayName for better user experience

	if c.logMessages {
		c.WriteMessageToFile(message)
	}

	if message.FirstMessage {
		c.messageStyle = c.messageStyle.Foreground(lipgloss.Color("#D946EF")).Bold(true) // Pink
	} else {
		c.messageStyle = lipgloss.NewStyle().Padding(0, 1) //Light Pink
	}

	fmt.Printf(" %s%s%s:%s\n\n", timestamp, badges, userName, c.messageStyle.Render(message.Message))
}

func formatBadges(message twitch.PrivateMessage, c Chat) string {

	badges := ""

	for badgeKey, badgeLevel := range message.User.Badges {
		// Split the badge key to get the base badge name
		// baseBadgeName := strings.Split(badgeKey, "/")[0]
		splitedBadgeName := strings.Split(badgeKey, "/")
		var baseBadgeName string

		baseBadgeName = splitedBadgeName[0]

		// Only display badge if level is greater than 0
		if badgeLevel > 0 {
			b, ok := c.badgeStyles[baseBadgeName]
			if ok {
				// Use ForegroundColor if specified, otherwise default to black
				foregroundColor := "#FFF"
				if b.ForegroundColor != "" {
					foregroundColor = b.ForegroundColor
				}

				// Add level to badge if it's meaningful
				badgeText := b.Icon
				if badgeLevel > 1 {
					badgeText += strconv.Itoa(badgeLevel)
				}

				style := lipgloss.NewStyle().
					Background(lipgloss.Color(b.Color)).
					Foreground(lipgloss.Color(foregroundColor)).
					MarginRight(1)

				badges += style.Render(badgeText)
			}
		}
	}

	return badges
}

func formatTimestamp(c Chat) string {
	timestamp := ""
	if c.showTimestamp {
		timestamp = fmt.Sprintf("[%s] ", time.Now().Format("15:04:05"))
	}

	return timestamp
}

// WriteMessageToFile writes a chat message to a JSON log file.
func (c *Chat) WriteMessageToFile(message twitch.PrivateMessage) {
	messageJSON, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	file, err := os.OpenFile("chat.log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.Write(append(messageJSON, ',', '\n'))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}
}
