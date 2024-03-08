package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var fileData []byte

var cmd_form Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "form",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageChannels),
		Description:              "Create custom forms right inside Discord",
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
						Name:         "accept_channel",
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
			fileData, _ = os.ReadFile("./form_templates/form_demo.json")
			fileReader := bytes.NewReader(fileData)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Get the example file edit it (make sure to have a unique \"form_type\") and submit it via `/form create`.\nOr use the demo button to get an idea of how the example would look like.",
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
			respond(i.Interaction, "Feature not available yet use `/form add` instead", true)
		case "add":
			var title, formID, overwriteTitle, acceptChannelID string
			var modsCanComment bool
			options := i.ApplicationCommandData().Options[0]

			for _, opt := range options.Options {
				switch opt.Name {
				case "type":
					formID = opt.StringValue()
				case "title":
					overwriteTitle = opt.StringValue()
					title = overwriteTitle
				case "accept_channel":
					acceptChannelID = opt.ChannelValue(s).ID
				case "mods_can_comment":
					modsCanComment = opt.BoolValue()
				}
			}
			if formID == "" {
				formID = "template_general"
			}

			formTitles := map[string]string{
				"template_feedback": "Submit Feedback",
				"template_ticket":   "Make a new ticket",
				"template_url":      "Add your URL",
				"template_general":  "Form",
			}
			if val, ok := formTitles[formID]; ok {
				title = val
			}

			var exists bool = true
			var formManageID uuid.UUID = uuid.New()
			for exists {
				formManageID = uuid.New()
				exists = getFormManageIdExists(formManageID)
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
			respond(i.Interaction, "Successfully added form button!", true)
		}
	},
	ComponentInteract: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if strings.HasPrefix(i.Interaction.MessageComponentData().CustomID, "form:") {
			jsonStringShowModal(i.Interaction, i.Interaction.MessageComponentData().CustomID, getFormType(strings.TrimPrefix(i.Interaction.MessageComponentData().CustomID, "form:")))
		}
		if i.Interaction.MessageComponentData().CustomID == "form_demo" {
			jsonStringShowModal(i.Interaction, "form_demo", "form_demo")
		}
	},
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		respond(i.Interaction, "The form data would be send to a specified channel. 🤲", true)
	},
	Autocomplete: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		choices := []*discordgo.ApplicationCommandOptionChoice{
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

func getFormTypes() []string {
	//needs custom IDs from databank
	return []string{"form_demo", "template_ticket", "template_url", "template_general"}
}

func getFormButtonIDs() []string {
	var IDs []string = []string{"form_demo"}
	var formButtonIDs []string = getFormManageIDs()
	for _, buttonID := range formButtonIDs {
		IDs = append(IDs, "form:"+buttonID)
	}
	return IDs
}
