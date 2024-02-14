package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	} else {
		fmt.Println("Discord session created")
	}
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds
	defer removeCommandFromAllGuilds(discord)
	discord.AddHandler(ready)
	discord.AddHandler(interactionCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Printf("Bot is now running as \"%s\"!", discord.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("\nShutting down...")
	defer removeCommandFromAllGuilds(discord)
	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "get",
			Description: "A command to get messages saved to the bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "tag",
					Description:  "Your predefined tag for the saved message",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	for _, guild := range event.Guilds {
		for _, command := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, command)
			if err != nil {
				fmt.Println("error creating command,", err)
				continue // Continue to the next guild
			}
		}
	}
}

func generateDynamicChoices(count int) []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for i := 1; i <= count; i++ {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  fmt.Sprintf("Option %d", i),
			Value: fmt.Sprintf("option_%d", i),
		})
	}
	return choices
}

var commandUseCount int

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {

		commandUseCount++

		choices := generateDynamicChoices(commandUseCount)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	}
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "get" {
			// Check if the command has options
			if len(i.ApplicationCommandData().Options) > 0 {
				// Loop through the options and handle them
				for _, option := range i.ApplicationCommandData().Options {
					switch option.Name {
					case "tag":
						value := option.Value.(string)
						response := fmt.Sprintf("You provided the tag: %s", value)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: response,
							},
						})
					}
				}
			}
		}
	}
}

func removeCommandFromAllGuilds(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		existingCommands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			fmt.Printf("error fetching existing commands for guild %s: %v\n", guild.Name, err)
			continue
		}

		for _, existingCommand := range existingCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, existingCommand.ID)
			if err != nil {
				fmt.Printf("error deleting command %s for guild %s: %v\n", existingCommand.Name, guild.Name, err)
			}
		}
	}
}
