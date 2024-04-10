package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var cmd_form Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "form",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageChannels),
		Description:              "Create custom forms right inside Discord",
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "help",
				Description: "Gives you an example file and demo for creating custom forms",
			},
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "custom",
				Description: "Create a new custom form right inside Discord",
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionAttachment{
						Name:        "json",
						Description: "Your edited form file",
						Required:    true,
					},
				},
			},
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "add",
				Description: "Adds existing forms to this channel",
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionChannel{
						Name:         "result_channel",
						Description:  "Where the form results should appear",
						ChannelTypes: []discord.ChannelType{discord.ChannelTypeGuildText},
					},
					&discord.ApplicationCommandOptionMentionable{
						Name:        "moderator",
						Description: "Who can interact with moderating buttons.",
					},
					&discord.ApplicationCommandOptionString{
						Name:         "type",
						Description:  "Which type of form you want to add",
						Autocomplete: true,
					},
					&discord.ApplicationCommandOptionString{
						Name:        "title",
						Description: "The title the form should have",
					},
					&discord.ApplicationCommandOptionChannel{
						Name:         "approve_channel",
						Description:  "Channel for results that need to be accepted by a moderator before sending it to the result channel",
						ChannelTypes: []discord.ChannelType{discord.ChannelTypeGuildText},
					},
					&discord.ApplicationCommandOptionBool{
						Name:        "mods_can_answer",
						Description: "Moderators can open a new channel on the form result, which then pings the user who submitted it",
					},
				},
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		switch *e.SlashCommandInteractionData().SubCommandName {
		case "help":
			fileData, err := formTemplates.ReadFile("form_templates/form_demo.json")
			if err != nil {
				logrus.Error(err)
				return
			}
			fileReader := bytes.NewReader(fileData)
			err = e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("NOT SUPPORTED YET!(use `/form add` instead)\n\nGet the example file edit it (make sure to have a unique \"form_type\") and submit it via `/form create`.\nOr use the demo button to get an idea of how the example would look like.").
				SetFiles(discord.NewFile("example.json", "json", fileReader)).
				SetContainerComponents(discord.ActionRowComponent{discord.NewPrimaryButton("Demo", "form_demo").WithEmoji(discord.ComponentEmoji{Name: "📑"})}).SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		case "custom":
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Feature not available yet use `/form add` instead").SetEphemeral(true).
				Build())
			if err != nil {
				logrus.Error(err)
			}
		case "add":
			var title, formID, overwriteTitle, acceptChannelID string
			var modsCanAnswer bool
			var resultChannelID string
			data := e.SlashCommandInteractionData()
			if data.Channel("result_channel").ID.String() != "0" {
				resultChannelID = data.Channel("result_channel").ID.String()
			}
			moderator := data.Role("moderator").ID.String()
			if moderator == "0" {
				moderator = e.User().ID.String()
			}
			formID = data.String("type")
			overwriteTitle = data.String("title")
			if overwriteTitle != "" {
				title = overwriteTitle
			}
			if data.Channel("approve_channel").ID.String() != "0" {
				acceptChannelID = data.Channel("approve_channel").ID.String()
			}
			modsCanAnswer = data.Bool("mods_can_answer")

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
			messagebuild := discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().
				SetTitle(title).SetDescription("Press the bottom button to open a form popup.").SetColor(hexToDecimal(color["primary"])).
				Build()).SetContainerComponents(discord.ActionRowComponent{
				discord.NewSuccessButton("Submit", "form:"+formManageID.String()).WithEmoji(discord.ComponentEmoji{
					Name:     "anim_rocket",
					ID:       snowflake.MustParse("1215740398706757743"),
					Animated: true,
				})}).
				Build()
			message, err := e.Client().Rest().CreateMessage(e.Channel().ID(), messagebuild)
			if err != nil {
				logrus.Error(err)
			}
			var category string
			if modsCanAnswer {
				c, err := e.Client().Rest().CreateGuildChannel(*e.GuildID(), discord.GuildCategoryChannelCreate{Name: title + " mod answers"})
				if err != nil {
					logrus.Error(err)
				}
				category = c.ID().String()
			}

			addFormButton(e.GuildID().String(), e.Channel().ID().String(), message.ID.String(), formManageID.String(), formID, resultChannelID, overwriteTitle, acceptChannelID, category, moderator)
			err = e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Successfully added form button!").SetEphemeral(true).Build())
			if err != nil {
				logrus.Error(err)
			}
		}
	},
	DynamicComponentIDs: func() []string { return getFormButtonIDs() },
	DynamicModalIDs:     func() []string { return getFormButtonIDs() },
	ComponentInteract: func(e *events.ComponentInteractionCreate) {
		if e.Data.Type() == discord.ComponentTypeButton {
			if strings.ContainsAny(e.ButtonInteractionData().CustomID(), ";") {
				var form_manage_id string = strings.TrimPrefix(strings.Split(e.ButtonInteractionData().CustomID(), ";")[0], "form:")
				switch strings.Split(e.ButtonInteractionData().CustomID(), ";")[1] {
				case "decline":
					err := e.Client().Rest().DeleteMessage(e.Channel().ID(), e.Message.ID)
					if err != nil {
						logrus.Error(err)
					}
					e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Submission declined!").SetEphemeral(true).Build())
				case "approve":
					embed := e.Message.Embeds[0]
					embed.Description = fmt.Sprintf("This submission was approved by <@%s>.", e.User().ID)
					_, err := e.Client().Rest().CreateMessage(snowflake.MustParse(getFormResultValues(form_manage_id).ResultChannelID), discord.NewMessageCreateBuilder().
						SetEmbeds(embed).
						Build())
					if err != nil {
						logrus.Error(err)
					}
					e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Submission accepted!").SetEphemeral(true).Build())
					err = e.Client().Rest().DeleteMessage(e.Channel().ID(), e.Message.ID)
					if err != nil {
						logrus.Error(err)
					}
				case "comment":
					author := strings.TrimSuffix(strings.Split(e.Message.Embeds[0].Fields[len(e.Message.Embeds[0].Fields)-1].Value, "<@")[1], ">")
					embed := e.Message.Embeds[0]
					moderator := e.User().ID
					channel := createFormComment(form_manage_id, snowflake.MustParse(author), moderator, "answer", embed, *e.GuildID(), e.Client())
					e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Created channel " + discord.ChannelMention(channel.ID())).SetEphemeral(true).Build())
				}
			} else {
				if strings.HasPrefix(e.ButtonInteractionData().CustomID(), "form:") {
					var formManageID string = strings.TrimPrefix(e.ButtonInteractionData().CustomID(), "form:")
					e.Modal(jsonStringBuildModal(e.User().ID.String(), formManageID, getFormType(formManageID), getFormOverwriteTitle(formManageID)))
				} else if e.ButtonInteractionData().CustomID() == "form_demo" {
					e.Modal(jsonStringBuildModal(e.User().ID.String(), "form_demo", "form_demo"))
				}
			}
		}
	},
	ModalSubmit: func(e *events.ModalSubmitInteractionCreate) {
		if !strings.HasPrefix(e.Data.CustomID, "form_demo") {
			var form_manage_id string = strings.Split(e.Data.CustomID, ":")[1]
			var result FormResult = getFormResultValues(form_manage_id)
			var fields []discord.EmbedField
			var modal ModalJson = getModalByFormID(getFormType(form_manage_id))
			var overwrite_title string = getFormOverwriteTitle(form_manage_id)
			if overwrite_title != "" {
				modal.Title = overwrite_title
			}
			var inline bool
			var index int = 0
			for _, component := range e.Data.Components {
				var input discord.TextInputComponent = component.(discord.TextInputComponent)
				inline = input.Style == discord.TextInputStyleShort
				fields = append(fields, discord.EmbedField{
					Name:   modal.Form[index].Label,
					Value:  input.Value,
					Inline: &inline,
				})
				index++
			}

			fields = append(fields, discord.EmbedField{
				Value: "From <#" + e.Channel().ID().String() + "> by " + e.User().Mention(),
			})
			if result.ResultChannelID == "" {
				if result.CommentCategoryID != "" {
					channel := createFormComment(form_manage_id, e.User().ID, snowflake.MustParse(result.ModeratorID), "answer", discord.NewEmbedBuilder().
						SetAuthorName(*e.User().GlobalName).SetAuthorIcon(*e.User().AvatarURL()).SetTitle("\""+modal.Title+"\"").SetDescription("This is the submitted result").
						SetColor(hexToDecimal(color["primary"])).SetFields(fields...).
						Build(), *e.GuildID(), e.Client())
					err := e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Created channel " + discord.ChannelMention(channel.ID())).SetEphemeral(true).Build())
					if err != nil {
						logrus.Error(err)
					}
				} else {
					e.CreateMessage(discord.NewMessageCreateBuilder().
						SetContent("You need to provide either a `result_channel` or enable `mods_can_answer` to create a valid form.").SetEphemeral(true).
						Build())
				}
			} else {
				logrus.Debug(result.AcceptChannelID)
				if result.AcceptChannelID == "" {
					_, err := e.Client().Rest().CreateMessage(snowflake.MustParse(result.ResultChannelID), discord.NewMessageCreateBuilder().
						SetEmbeds(discord.NewEmbedBuilder().
							SetAuthorName(*e.User().GlobalName).SetAuthorIcon(*e.User().AvatarURL()).SetTitle("\""+modal.Title+"\"").SetDescription("This is the submitted result").
							SetColor(hexToDecimal(color["primary"])).SetFields(fields...).
							Build()).
						SetContainerComponents(discord.NewActionRow(discord.
							NewButton(discord.ButtonStylePrimary, "Comment", "form:"+form_manage_id+";comment", "").
							WithEmoji(discord.ComponentEmoji{Name: "👥"}))).
						Build())
					if err != nil {
						logrus.Error(err)
					} else {
						err = e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Submitted!").SetEphemeral(true).Build())
						if err != nil {
							logrus.Error(err)
						}
					}
				} else {
					logrus.Debug("HEERE")
					var buttons []discord.InteractiveComponent
					if result.CommentCategoryID != "" {
						buttons = []discord.InteractiveComponent{discord.
							NewButton(discord.ButtonStylePrimary, "Comment", "form:"+form_manage_id+";comment", "").
							WithEmoji(discord.ComponentEmoji{Name: "👥"})}
					}
					buttons = append(buttons, discord.
						NewButton(discord.ButtonStyleDanger, "Decline", "form:"+form_manage_id+";decline", "").
						WithEmoji(discord.ComponentEmoji{Name: "🛑"}),
						discord.
							NewButton(discord.ButtonStyleSuccess, "Approve", "form:"+form_manage_id+";approve", "").
							WithEmoji(discord.ComponentEmoji{Name: "🎉"}))
					_, err := e.Client().Rest().CreateMessage(snowflake.MustParse(result.AcceptChannelID), discord.NewMessageCreateBuilder().
						SetEmbeds(discord.NewEmbedBuilder().
							SetAuthorName(*e.User().GlobalName).SetAuthorIcon(*e.User().AvatarURL()).SetTitle("\""+modal.Title+"\"").SetDescription("**This submission needs approval.**").
							SetColor(hexToDecimal(color["primary"])).SetFields(fields...).
							Build()).
						SetContainerComponents(discord.NewActionRow(buttons...)).
						Build())

					if err != nil {
						logrus.Error(err)
					} else {
						err = e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Submitted!").SetEphemeral(true).Build())
						if err != nil {
							logrus.Error(err)
						}
					}
				}
			}
		} else {
			err := e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("The results would be submited...").SetEphemeral(true).Build())
			if err != nil {
				logrus.Error(err)
			}
		}

	},
	Autocomplete: func(e *events.AutocompleteInteractionCreate) {
		err := e.AutocompleteResult([]discord.AutocompleteChoice{
			&discord.AutocompleteChoiceString{
				Name:  "Support Ticket",
				Value: "template_ticket",
			},
			&discord.AutocompleteChoiceString{
				Name:  "Submit URL",
				Value: "template_url",
			},
			&discord.AutocompleteChoiceString{
				Name:  "General",
				Value: "template_general",
			},
		})
		if err != nil {
			logrus.Error(err)
		}
	},
}

var cmd_ticket_form Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:                     "ticket",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageChannels),
		Description:              "A quick command to create Ticketpanels. (/form for more)",
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:        "title",
				Description: "The title the ticket should have",
			},
			&discord.ApplicationCommandOptionMentionable{
				Name:        "moderator",
				Description: "Who can interact with moderating buttons.",
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		var title string = "Ticket"
		var moderator string
		data := e.SlashCommandInteractionData()
		if data.String("title") != "" {
			title = data.String("title")
		}
		moderator = data.Role("moderator").ID.String()
		if moderator == "" {
			moderator = data.User("moderator").ID.String()
		}
		var exists bool = true
		var formManageID uuid.UUID = uuid.New()
		for exists {
			formManageID = uuid.New()
			exists = getFormManageIdExists(formManageID)
		}
		messagebuild := discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().
			SetTitle(title).SetDescription("Press the bottom button to open a form popup.").SetColor(hexToDecimal(color["primary"])).
			Build()).SetContainerComponents(discord.ActionRowComponent{
			discord.NewSuccessButton("Submit", "form:"+formManageID.String()).WithEmoji(discord.ComponentEmoji{
				Name:     "anim_rocket",
				ID:       snowflake.MustParse("1215740398706757743"),
				Animated: true,
			})}).
			Build()
		message, err := e.Client().Rest().CreateMessage(e.Channel().ID(), messagebuild)
		if err != nil {
			logrus.Error(err)
			return
		}
		if title == "" {
			title = "Ticket"
		}
		var category string
		c, err := e.Client().Rest().CreateGuildChannel(*e.GuildID(), discord.GuildCategoryChannelCreate{Name: title + " mod answers"})
		if err != nil {
			logrus.Error(err)
		}
		category = c.ID().String()
		if title == "Ticket" {
			title = ""
		}

		addFormButton(e.GuildID().String(), e.Channel().ID().String(), message.ID.String(), formManageID.String(), "template_ticket", "", title, "", category, moderator)
		err = e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Successfully added ticket panel!\n(`/form` for more options or custom ticket forms.)").SetEphemeral(true).Build())
		if err != nil {
			logrus.Error(err)
		}
	},
}

// moderator can be userID as well as  roleID
func createFormComment(form_manage_id string, author snowflake.ID, moderator snowflake.ID, commentName string, embed discord.Embed, guildID snowflake.ID, client bot.Client) discord.Channel {
	var category snowflake.ID = snowflake.MustParse(getFormResultValues(form_manage_id).CommentCategoryID)
	_, err := client.Rest().GetChannel(category)
	if err != nil {
		c, err := client.Rest().CreateGuildChannel(guildID, discord.GuildCategoryChannelCreate{Name: strings.Trim(embed.Title, "\"") + " mod " + commentName + "s"})
		if err != nil {
			logrus.Error(err)
		}
		category = c.ID()
		updateFormCommentCategory(form_manage_id, category.String())
	}
	ch, err := client.Rest().CreateGuildChannel(guildID, discord.GuildTextChannelCreate{
		ParentID: category,
		Name:     strings.ToLower(embed.Author.Name) + "-" + commentName,
	})
	if err != nil {
		logrus.Error(err)
	}
	var permissionOverwrites []discord.PermissionOverwrite = []discord.PermissionOverwrite{
		discord.RolePermissionOverwrite{
			RoleID: guildID,
			Deny:   discord.PermissionViewChannel,
		}}

	if isIDRole(client, guildID, moderator) {
		permissionOverwrites = append(permissionOverwrites, discord.RolePermissionOverwrite{
			RoleID: moderator,
			Allow:  discord.PermissionViewChannel,
		})
	} else {
		permissionOverwrites = append(permissionOverwrites, discord.MemberPermissionOverwrite{
			UserID: moderator,
			Allow:  discord.PermissionViewChannel,
		})
	}
	permissionOverwrites = append(permissionOverwrites, discord.RolePermissionOverwrite{
		RoleID: author,
		Allow:  discord.PermissionViewChannel,
	})
	_, err = client.Rest().UpdateChannel(ch.ID(), discord.GuildTextChannelUpdate{PermissionOverwrites: &permissionOverwrites})
	if err != nil {
		logrus.Error(err)
	}
	modTypeChar := ""
	if isIDRole(client, guildID, moderator) {
		modTypeChar = "&"
	}
	embed.Description = "This was submitted"
	_, err = client.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().
		SetContent("<@"+modTypeChar+moderator.String()+"> <@"+author.String()+">").SetEmbeds(embed).
		Build())
	if err != nil {
		logrus.Error(err)
	}
	return ch
}

func getFormButtonIDs() []string {
	var IDs []string = []string{"form_demo"}
	var formButtonIDs []string = getFormManageIDs()
	for _, buttonID := range formButtonIDs {
		IDs = append(IDs, "form:"+buttonID)
	}
	return IDs
}
