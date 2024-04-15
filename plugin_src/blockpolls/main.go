package main

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/shared"
)

var db *sql.DB

var dbCreateQuery string = `
CREATE TABLE IF NOT EXISTS blockpolls (
	guild_id TEXT NOT NULL,
	channel_id TEXT,
	global BOOLEAN,
	allowed_role TEXT,
	PRIMARY KEY (guild_id)
)
`

var Plugin = &shared.Plugin{
	Name: "Block Polls",
	Init: func(d *sql.DB) error {
		db = d
		_, err := d.Exec(dbCreateQuery)
		if err != nil {
			return err
		}
		shared.BotConfigs = append(shared.BotConfigs, bot.WithEventListenerFunc(messageCreate))
		return nil
	},
	Commands: []shared.Command{
		{
			Definition: discord.SlashCommandCreate{
				Name:                     "block-polls",
				Description:              "Block polls from beeing posted in this channel.",
				DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageChannels),
				Contexts: []discord.InteractionContextType{
					discord.InteractionContextTypeGuild,
					discord.InteractionContextTypePrivateChannel},
				IntegrationTypes: []discord.ApplicationIntegrationType{
					discord.ApplicationIntegrationTypeGuildInstall},
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionSubCommand{
						Name:        "toggle",
						Description: "Toggle blocking polls from beeing posted in this channel.",
						Options: []discord.ApplicationCommandOption{
							&discord.ApplicationCommandOptionBool{
								Name:        "global",
								Description: "If polls are blocked server wide or only in the current channel.",
							},
							&discord.ApplicationCommandOptionRole{
								Name:        "allowed-role",
								Description: "The role that bypasses this block role.",
							},
						},
					},
					/*&discord.ApplicationCommandOptionSubCommand{
						Name:        "list",
						Description: "List the current block polls rules for this server.",
					},*/
				},
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				switch *e.SlashCommandInteractionData().SubCommandName {
				case "toggle":
					isGlobal := isGlobalBlockPolls(e.GuildID().String())
					if isGlobal && !e.SlashCommandInteractionData().Bool("global") {
						e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Polls are currently globally blocked. Disable global blocking to enable channel specific blocking.").SetEphemeral(true).Build())
					} else {
						exists, isGlobal := toggleBlockPolls(e.GuildID().String(), e.Channel().ID().String(), e.SlashCommandInteractionData().Bool("global"), e.SlashCommandInteractionData().Role("allowed-role").ID.String())
						if exists {
							if e.SlashCommandInteractionData().Bool("global") {
								err := e.CreateMessage(discord.NewMessageCreateBuilder().
									SetContent("Polls are now globally unblocked.").SetEphemeral(true).
									Build())
								if err != nil {
									logrus.Error(err)
								}
							} else {
								err := e.CreateMessage(discord.NewMessageCreateBuilder().
									SetContent("Polls are now unblocked in " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
									Build())
								if err != nil {
									logrus.Error(err)
								}
							}
						} else {
							if isGlobal {
								err := e.CreateMessage(discord.NewMessageCreateBuilder().
									SetContent("Polls are now globally blocked.").SetEphemeral(true).
									Build())
								if err != nil {
									logrus.Error(err)
								}
							} else {
								err := e.CreateMessage(discord.NewMessageCreateBuilder().
									SetContent("Polls are now blocked in " + discord.ChannelMention(e.Channel().ID())).SetEphemeral(true).
									Build())
								if err != nil {
									logrus.Error(err)
								}
							}
						}
					}
					/*case "list":
					list := listBlockPolls(e.GuildID().String())*/
				}
			},
		},
	},
}

func messageIsPoll(channelID string, messageID string, client bot.Client) bool {
	url := rest.DefaultConfig().URL + "/channels/" + channelID + "/messages/" + messageID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err)
		return false
	}

	auth := "Bot " + client.Token()
	req.Header.Set("Authorization", auth)

	resp, err := client.Rest().HTTPClient().Do(req)
	if err != nil {
		logrus.Error(err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return false
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Error(err)
		return false
	}

	_, ok := data["poll"]
	return ok
}
