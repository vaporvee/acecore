package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

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
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "json",
						Description: "Your edited form file",
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
					},
					{
						Type:        discordgo.ApplicationCommandOptionMentionable,
						Name:        "moderator",
						Description: "Who can interact with moderating buttons.",
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
						Name:         "approve_channel",
						Description:  "Channel for results that need to be accepted by a moderator before sending it to the result channel",
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "mods_can_answer",
						Description: "Moderators can open a new channel on the form result, which then pings the user who submitted it",
					},
				},
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "help":
			fileData, err := formTemplates.ReadFile("form_templates/form_demo.json")
			if err != nil {
				logrus.Error(err)
				return
			}
			fileReader := bytes.NewReader(fileData)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "NOT SUPPORTED YET!(use `/form add` instead)\n\nGet the example file edit it (make sure to have a unique \"form_type\") and submit it via `/form create`.\nOr use the demo button to get an idea of how the example would look like.",
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
			if err != nil {
				logrus.Error(err)
			}
		case "custom":
			err := respond(i.Interaction, "Feature not available yet use `/form add` instead", true)
			if err != nil {
				logrus.Error(err)
			}
		case "add":
			var title, formID, overwriteTitle, acceptChannelID string
			var modsCanAnswer bool
			var resultChannelID string
			moderator := i.Member.User.ID
			if i.ApplicationCommandData().Options != nil {
				options := i.ApplicationCommandData().Options[0]
				for _, opt := range options.Options {
					switch opt.Name {
					case "result_channel":
						resultChannelID = opt.ChannelValue(s).ID
					case "type":
						formID = opt.StringValue()
					case "title":
						overwriteTitle = opt.StringValue()
						title = overwriteTitle
					case "approve_channel":
						acceptChannelID = opt.ChannelValue(s).ID
					case "mods_can_answer":
						modsCanAnswer = opt.BoolValue()
					case "moderator":
						moderator = opt.RoleValue(s, i.GuildID).ID
						if moderator == "" {
							moderator = opt.UserValue(s).ID
						}
					}
				}
			}
			if formID == "" {
				formID = "template_general"
			}
			if title == "" {
				formTitles := map[string]string{
					"template_ticket":  "Make a new ticket",
					"template_url":     "Add your URL",
					"template_general": "Form",
				}
				if val, ok := formTitles[formID]; ok {
					title = val
				}
			}
			var exists bool = true
			var formManageID uuid.UUID = uuid.New()
			for exists {
				formManageID = uuid.New()
				exists = getFormManageIdExists(formManageID)
			}

			message, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
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
									Name:     "anim_rocket",
									ID:       "1215740398706757743",
									Animated: true,
								},
							},
						},
					},
				},
			})
			if err != nil {
				logrus.Error(err)
				return
			}
			var category string
			if modsCanAnswer {
				c, err := s.GuildChannelCreate(i.GuildID, title+" mod answers", discordgo.ChannelTypeGuildCategory)
				if err != nil {
					logrus.Error(err)
				}
				category = c.ID
			}
			addFormButton(i.GuildID, i.ChannelID, message.ID, formManageID.String(), formID, resultChannelID, overwriteTitle, acceptChannelID, category, moderator)
			err = respond(i.Interaction, "Successfully added form button!", true)
			if err != nil {
				logrus.Error(err)
			}
		}
	},
	DynamicComponentIDs: func() []string { return getFormButtonIDs() },
	DynamicModalIDs:     func() []string { return getFormButtonIDs() },
	ComponentInteract: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if strings.ContainsAny(i.MessageComponentData().CustomID, ";") {
			var form_manage_id string = strings.TrimPrefix(strings.Split(i.MessageComponentData().CustomID, ";")[0], "form:")
			switch strings.Split(i.MessageComponentData().CustomID, ";")[1] {
			case "decline":
				err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
				if err != nil {
					logrus.Error(err)
				}
				respond(i.Interaction, "Submission declined!", true)
			case "approve":
				embed := i.Message.Embeds[0]
				embed.Description = fmt.Sprintf("This submission was approved by <@%s>.", i.Member.User.ID)
				_, err := s.ChannelMessageSendComplex(getFormResultValues(form_manage_id).ResultChannelID, &discordgo.MessageSend{
					Embed: embed,
				})
				if err != nil {
					logrus.Error(err)
				}
				respond(i.Interaction, "Submission accepted!", true)
				err = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
				if err != nil {
					logrus.Error(err)
				}

			case "comment":
				author := strings.TrimSuffix(strings.Split(i.Message.Embeds[0].Fields[len(i.Message.Embeds[0].Fields)-1].Value, "<@")[1], ">")
				embed := i.Message.Embeds[0]
				moderator := i.Member.User.ID
				createFormComment(form_manage_id, author, moderator, "answer", embed, i)
			}
		} else {
			if strings.HasPrefix(i.Interaction.MessageComponentData().CustomID, "form:") {
				var formManageID string = strings.TrimPrefix(i.Interaction.MessageComponentData().CustomID, "form:")
				jsonStringShowModal(i.Interaction, i.Interaction.MessageComponentData().CustomID, getFormType(formManageID), getFormOverwriteTitle(formManageID))
			} else if i.Interaction.MessageComponentData().CustomID == "form_demo" {
				jsonStringShowModal(i.Interaction, "form_demo", "form_demo")
			}
		}
	},
	ModalSubmit: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if !strings.HasPrefix(i.ModalSubmitData().CustomID, "form_demo") {
			var form_manage_id string = strings.Split(i.ModalSubmitData().CustomID, ":")[1]
			var result FormResult = getFormResultValues(form_manage_id)
			var fields []*discordgo.MessageEmbedField
			var modal ModalJson = getModalByFormID(getFormType(form_manage_id))
			var overwrite_title string = getFormOverwriteTitle(form_manage_id)
			if overwrite_title != "" {
				modal.Title = overwrite_title
			}
			for index, component := range i.ModalSubmitData().Components {
				var input *discordgo.TextInput = component.(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   modal.Form[index].Label,
					Value:  input.Value,
					Inline: input.Style == discordgo.TextInputShort,
				})
			}

			channel, _ := s.Channel(i.ChannelID)
			fields = append(fields, &discordgo.MessageEmbedField{
				Value: "From <#" + channel.ID + "> by <@" + i.Member.User.ID + ">",
			})
			if result.ResultChannelID == "" {
				if result.CommentCategoryID != "" {
					createFormComment(form_manage_id, i.Member.User.ID, result.ModeratorID, "answer", &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{
							Name:    i.Member.User.Username,
							IconURL: i.Member.AvatarURL("256"),
						},
						Title:       "\"" + modal.Title + "\"",
						Color:       hexToDecimal(color["primary"]),
						Description: "This is the submitted result",
						Fields:      fields,
					}, i)
				} else {
					respond(i.Interaction, "You need to provide either a `result_channel` or enable `mods_can_answer` to create a valid form.", true)
				}
			} else {
				if result.AcceptChannelID == "" {
					var buttons []discordgo.MessageComponent
					if result.CommentCategoryID != "" {
						buttons = []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Style: discordgo.PrimaryButton,
										Emoji: discordgo.ComponentEmoji{
											Name: "👥",
										},
										Label:    "Comment",
										CustomID: "form:" + form_manage_id + ";comment",
									},
								},
							},
						}
					}
					_, err := s.ChannelMessageSendComplex(result.ResultChannelID, &discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{
								Name:    i.Member.User.Username,
								IconURL: i.Member.AvatarURL("256"),
							},
							Title:       "\"" + modal.Title + "\"",
							Color:       hexToDecimal(color["primary"]),
							Description: "This is the submitted result",
							Fields:      fields,
						},
						Components: buttons,
					})
					if err != nil {
						logrus.Error(err)
					} else {
						err = respond(i.Interaction, "Submited!", true)
						if err != nil {
							logrus.Error(err)
						}
					}
				} else {
					var buttons []discordgo.MessageComponent
					if result.CommentCategoryID != "" {
						buttons = []discordgo.MessageComponent{
							discordgo.Button{
								Style: discordgo.PrimaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "👥",
								},
								Label:    "Comment",
								CustomID: "form:" + form_manage_id + ";comment",
							},
						}
					}
					buttons = append(buttons,
						discordgo.Button{
							Style: discordgo.DangerButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🛑",
							},
							Label:    "Decline",
							CustomID: "form:" + form_manage_id + ";decline",
						},
						discordgo.Button{
							Style: discordgo.SuccessButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🎉",
							},
							Label:    "Approve",
							CustomID: "form:" + form_manage_id + ";approve",
						})
					_, err := s.ChannelMessageSendComplex(result.AcceptChannelID, &discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{
								Name:    i.Member.User.Username,
								IconURL: i.Member.AvatarURL("256"),
							},
							Title:       "\"" + modal.Title + "\"",
							Color:       hexToDecimal(color["primary"]),
							Description: "**This submission needs approval.**",
							Fields:      fields,
						},
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: buttons,
							},
						}},
					)
					if err != nil {
						logrus.Error(err)
					} else {
						err = respond(i.Interaction, "Submited!", true)
						if err != nil {
							logrus.Error(err)
						}
					}
				}
			}
		} else {
			err := respond(i.Interaction, "The results would be submited...", true)
			if err != nil {
				logrus.Error(err)
			}
		}

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
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			logrus.Error(err)
		}
	},
}

var cmd_ticket_form Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:                     "ticket",
		DefaultMemberPermissions: int64Ptr(discordgo.PermissionManageChannels),
		Description:              "A quick command to create Ticketpanels. (/form for more)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "The title the ticket should have",
			},
			{
				Type:        discordgo.ApplicationCommandOptionMentionable,
				Name:        "moderator",
				Description: "Who can interact with moderating buttons.",
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var title string = "Ticket"
		var moderator string
		if i.ApplicationCommandData().Options != nil {
			for _, opt := range i.ApplicationCommandData().Options {
				switch opt.Name {
				case "title":
					title = opt.StringValue()
				case "moderator":
					moderator = opt.RoleValue(s, i.GuildID).ID
					if moderator == "" {
						moderator = opt.UserValue(s).ID
					}
				}
			}
		}
		if moderator == "" {
			moderator = i.Member.User.ID
		}
		var exists bool = true
		var formManageID uuid.UUID = uuid.New()
		for exists {
			formManageID = uuid.New()
			exists = getFormManageIdExists(formManageID)
		}
		message, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
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
								Name:     "anim_rocket",
								ID:       "1215740398706757743",
								Animated: true,
							},
						},
					},
				},
			},
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		if title == "" {
			title = "Ticket"
		}
		var category string
		c, err := s.GuildChannelCreate(i.GuildID, title+" mod answers", discordgo.ChannelTypeGuildCategory)
		if err != nil {
			logrus.Error(err)
		}
		category = c.ID
		if title == "Ticket" {
			title = ""
		}

		addFormButton(i.GuildID, i.ChannelID, message.ID, formManageID.String(), "template_ticket", "", title, "", category, moderator)
		err = respond(i.Interaction, "Successfully added ticket panel!\n(`/form` for more options or custom ticket forms.)", true)
		if err != nil {
			logrus.Error(err)
		}
	},
}

// moderator can be userID as well as  roleID
func createFormComment(form_manage_id string, author string, moderator string, commentName string, embed *discordgo.MessageEmbed, i *discordgo.InteractionCreate) {
	category := getFormResultValues(form_manage_id).CommentCategoryID
	_, err := bot.Channel(category)
	if err != nil {
		c, err := bot.GuildChannelCreate(i.GuildID, strings.Trim(embed.Title, "\"")+" mod "+commentName+"s", discordgo.ChannelTypeGuildCategory)
		if err != nil {
			logrus.Error(err)
		}
		category = c.ID
		updateFormCommentCategory(form_manage_id, c.ID)
	}
	ch, err := bot.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		ParentID: category,
		Name:     strings.ToLower(embed.Author.Name) + "-" + commentName,
	})
	if err != nil {
		logrus.Error(err)
	}
	err = bot.ChannelPermissionSet(ch.ID, i.GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionViewChannel)
	if err != nil {
		logrus.Error(err)
	}
	modType := discordgo.PermissionOverwriteTypeMember
	if isIDRole(i.GuildID, moderator) {
		modType = discordgo.PermissionOverwriteTypeRole
	}
	err = bot.ChannelPermissionSet(ch.ID, moderator, modType, discordgo.PermissionViewChannel, 0)
	if err != nil {
		logrus.Error(err)
	}
	err = bot.ChannelPermissionSet(ch.ID, author, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionViewChannel, 0)
	if err != nil {
		logrus.Error(err)
	}
	modTypeChar := "&"
	if modType == discordgo.PermissionOverwriteTypeMember {
		modTypeChar = ""
	}
	_, err = bot.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
		Content: "<@" + modTypeChar + moderator + "> <@" + author + ">",
		Embed:   embed,
	})
	if err != nil {
		logrus.Error(err)
	}
	respond(i.Interaction, "Created channel <#"+ch.ID+">", true)
}

func getFormButtonIDs() []string {
	var IDs []string = []string{"form_demo"}
	var formButtonIDs []string = getFormManageIDs()
	for _, buttonID := range formButtonIDs {
		IDs = append(IDs, "form:"+buttonID)
	}
	return IDs
}
