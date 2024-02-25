package main

import (
	"bytes"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
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
				Name:        "custom",
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
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Adds existing forms to this channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionChannel,
						Name:         "result_channel",
						Description:  "Where the form results should appear",
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "type",
						Description:  "Which type of form you want to add",
						Autocomplete: true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "title",
						Description: "The title the form should have",
					},
					{
						Type:         discordgo.ApplicationCommandOptionChannel,
						Name:         "accept_channel_id",
						Description:  "Channel for results that need to be accepted by a moderator before sending it to the result channel",
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "mods_can_comment",
						Description: "Moderators can open a new channel on the form result, which then pings the user who submitted it",
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
		case "custom":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Feature not available yet use `/form add` instead",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		case "add":
			var title, formID, overwriteTitle, acceptChannelID string
			var modsCanComment bool
			options := i.ApplicationCommandData().Options[0]

			if len(options.Options) <= 1 {
				formID = "template_general"
			} else {
				formID = options.Options[1].StringValue()
			}
			switch formID {
			case "template_feedback":
				title = "Submit Feedback"
			case "template_ticket":
				title = "Make a new ticket"
			case "template_url":
				title = "Add your URL"
			case "template_general":
				title = "Form"
			}

			var exists bool = true
			formManageID := uuid.New()
			for exists {
				formManageID = uuid.New()
				exists = getFormManageIdExists(i.GuildID, formManageID)
			}

			if !(len(options.Options) <= 2) {
				overwriteTitle = options.Options[2].StringValue()
				title = overwriteTitle
			}
			if !(len(options.Options) <= 3) {
				acceptChannelID = options.Options[3].ChannelValue(s).ID
			}
			if !(len(options.Options) <= 4) {
				modsCanComment = options.Options[4].BoolValue()
			}

			message, _ := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Color:       hexToDecimal(color["primary"]),
					Title:       title,
					Description: "Press the bottom button to open a form popup.",
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								CustomID: "form:" + formManageID.String(),
								Style:    discordgo.SuccessButton,
								Label:    "Submit",
								Emoji: discordgo.ComponentEmoji{
									Name: "📥",
								},
							},
						},
					},
				},
			})
			addFormButton(i.GuildID, i.ChannelID, message.ID, formManageID.String(), formID, options.Options[0].ChannelValue(s).ID, overwriteTitle, acceptChannelID, modsCanComment)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Successfully added form button!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	},
	ComponentIDs: getFormButtonIDs(),
	ComponentInteract: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Interaction.MessageComponentData().CustomID == "form_demo" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "form_demo" + i.Interaction.Member.User.ID,
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
		}
	},
	ModalIDs: getFormIDs(),
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The form data would be send to a specified channel. 🤲",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		choices := []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Feedback",
				Value: "template_feedback",
			},
			{
				Name:  "Support Ticket",
				Value: "template_ticket",
			},
			{
				Name:  "Submit URL",
				Value: "template_url",
			},
			{
				Name:  "General",
				Value: "template_general",
			},
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	},
}

func getFormIDs() []string {
	//needs custom IDs from databank
	return []string{"form_demo", "template_feedback", "template_ticket", "template_url", "template_general"}
}

func getFormButtonIDs() []string {
	var IDs []string = []string{"form_demo"}
	var formButtonIDs []string = getFormManageIDs()
	for _, buttonID := range formButtonIDs {
		IDs = append(IDs, "form:"+buttonID)
	}
	return IDs
}
