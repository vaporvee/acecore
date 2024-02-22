package main

import (
	"bytes"
	"os"

	"github.com/bwmarrin/discordgo"
)

var fileData []byte

var form_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "form",
		Description: "Create custom forms right inside Discord",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "help",
				Description: "Gives you a example file and demo for creating custom forms",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "create",
				Description: "Create a new custom form right inside Discord",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "title",
						Description: "The title inside the form window",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "json",
						Description: "Your edited form file",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "results_channel",
						Description: "The channel where the form results should be posted",
						Required:    true,
					},
				},
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "help":
			fileData, _ = os.ReadFile("./attachments/example_modal.json")
			fileReader := bytes.NewReader(fileData)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Get the example file edit it and submit it via `/form create`.\nOr use the demo button to get an idea of how the example would look like.",
					Flags:   discordgo.MessageFlagsEphemeral,
					Files: []*discordgo.File{
						{
							Name:        "example.json",
							ContentType: "json",
							Reader:      fileReader,
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "📑",
									},
									CustomID: "form_demo",
									Label:    "Demo",
									Style:    discordgo.PrimaryButton,
								},
							},
						},
					},
				},
			})
		case "create":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Placeholder",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
	ComponentIDs: []string{"form_demo"},
	ComponentInteract: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "form_demo_modal" + i.Interaction.Member.User.ID,
				Title:    "Demo form",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "demo_short",
								Label:       "This is a simple textline",
								Style:       discordgo.TextInputShort,
								Placeholder: "...and it is required!",
								Value:       "",
								Required:    true,
								MaxLength:   20,
								MinLength:   0,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "demo_paragraph",
								Label:       "This is a paragraph",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "...and it is not required!",
								Value:       "We already have some input here",
								Required:    false,
								MaxLength:   2000,
								MinLength:   0,
							},
						},
					},
				},
			},
		})
	},
	ModalID: "form_demo_modal",
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The form data would be send to a specified channel. 🤲",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	},
}
