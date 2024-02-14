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
			Name:        "test",
			Description: "A test command.",
		},
		{
			Name:        "secondtest",
			Description: "A second test command.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The query to search for.",
					Required:    true,
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

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "test" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You tested me!",
				},
			})
		}
		if i.ApplicationCommandData().Name == "secondtest" {
			// Check if the command has options
			if len(i.ApplicationCommandData().Options) > 0 {
				// Loop through the options and handle them
				for _, option := range i.ApplicationCommandData().Options {
					switch option.Name {
					case "query":
						value := option.Value.(string)
						response := fmt.Sprintf("You provided the query: %s", value)
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
