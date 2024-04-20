// Token: Replace with your bot token here
// App ID: Replace with application id here
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Permission  int    `json:"permission"`
}

var botToken = "Replace with your bot token here"
var applicationID = "Replace with application id here"
var enableLogs = true

var emojiMap = map[string]string{
	"New York":      "<:new_york:1216351445884993657>",
	"Detroit":       "<:detroit:1216351353505579078>",
	"Chicago":       "<:chicago:1216351335398768741>",
	"San Francisco": "<:san_francisco:1216351588579409990>",
	"Atlanta":       "<:atlanta:1216351298597683280>",
	"San Diego":     "<:san_diego:1216351558967623701>",
	"Los Angeles":   "<:los_angeles:1216351410015440997>",
	"Miami":         "<:miami:1216351428025516162>",
	"Las Vegas":     "<:las_vegas:1216351381355495424>",
	"Washington":    "<:washington:1216351607915286528>",
}

func generateServersEmbed() (discordgo.MessageEmbed, error) {

	data, err := os.ReadFile("servers.json")
	if err != nil {
		return discordgo.MessageEmbed{}, err
	}

	var servers []ServerData
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return discordgo.MessageEmbed{}, err
	}

	var response string
	for _, info := range servers {
		emoji1, ok := emojiMap[info.Name]
		if !ok {
			emoji1 = "❌" // Default emoji if server name is not found in the map
		}

		emoji := "❌"
		if info.Available {
			emoji = "✅"
		}

		response += fmt.Sprintf("%s%s: %s\n- `Online: %d`\n\n", emoji1, info.Name, emoji, info.PlayersCount)
	}

	return discordgo.MessageEmbed{
		Title:       "Majestic server monitoring",
		Description: response,
		Color:       0x81D8D0,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info is current as of",
		},
	}, nil
}

func main() {
	if !enableLogs {
		log.SetOutput(io.Discard)
	}
	go parser()
	registerRollCommand()
	registerServersCommand()
	registerUpdateLogCommand()

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreate(s, i)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		fmt.Println("Новый пользователь присоединился к серверу")
	})
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord connection: ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// CTRL+C for exit.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {

		case "roll":
			roll := rand.Intn(6) + 1

			resultMessage := fmt.Sprintf("%s rolled a %d", i.Member.Mention(), roll)
			if roll == 1 {
				resultMessage += "... snake eyes!"
			}

			embed := &discordgo.MessageEmbed{
				Description: resultMessage,
				Color:       0x81D8D0,
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})

		case "servers":
			embed, err := generateServersEmbed()
			if err != nil {
				fmt.Println("Error generatind embed: ", err)
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{&embed},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔄",
									},
									Label: "Refresh",
									Style: discordgo.SuccessButton,
									// URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
									CustomID: "refresh_servers",
								},
							},
						},
					},
					CustomID: "2",
				},
			})

		case "updatelog":
			updateLogMessage :=
				"- Added /updatelog. Info about latest update will be posted here\n" +
					"- Reworked /servers. Now information about majestic servers is updating every 10 seconds\n" +
					" - Note: type /servers again if you want to get latest statistics. The messege isn't updating automatically"

			embed := &discordgo.MessageEmbed{
				Title:       "Build: 1.0.1b [March 10, 2024]",
				Description: updateLogMessage,
				Color:       0x81D8D0,
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
		}

	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "refresh_servers":
			embed, err := generateServersEmbed()
			if err != nil {
				fmt.Println("Error generatind embed: ", err)
				return
			}
			s.ChannelMessageEditEmbed(i.Interaction.Message.ChannelID, i.Interaction.Message.ID, &embed)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Successfully refreshed statistics",
				},
			})
		}
	}
}

func registerRollCommand() {
	command := Command{
		Name:        "roll",
		Description: "Rolls a number from 1 to 6.",
		Type:        1,
		Permission:  0, // Permission integer (0 for everyone)
	}

	registerCommand(botToken, applicationID, command)
}

func registerServersCommand() {
	command := Command{
		Name:        "servers",
		Description: "Shows majestic online statistics.",
		Type:        1,
		Permission:  0, // Permission integer (0 for everyone)
	}

	registerCommand(botToken, applicationID, command)
}

func registerUpdateLogCommand() {
	command := Command{
		Name:        "updatelog",
		Description: "Shows info about latest update.",
		Type:        1,
		Permission:  0, // Permission integer (0 for everyone)
	}

	registerCommand(botToken, applicationID, command)
}

func registerCommand(botToken, applicationID string, command Command) {
	commandJSON, err := json.Marshal(command)
	if err != nil {
		fmt.Println("Error marshaling command JSON:", err)
		return
	}

	url := fmt.Sprintf("https://discord.com/api/v9/applications/%s/commands", applicationID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(commandJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to register %s command. Status code: %d\n", command.Name, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Println("Response Body:", string(body))

		return
	}

	fmt.Printf("%s command was registered successfully.\n", command.Name)
}
