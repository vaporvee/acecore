package main

import (
	"database/sql"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/shared"
)

var db *sql.DB

var dbCreateQuery string = ` 
CREATE TABLE IF NOT EXISTS autojoinroles (
	guild_id TEXT NOT NULL,
	bot_role TEXT,
	user_role TEXT,
	PRIMARY KEY (guild_id)
);
`

var Plugin = &shared.Plugin{
	Name: "Auto Join Roles",
	Init: func(d *sql.DB) error {
		db = d
		_, err := d.Exec(dbCreateQuery)
		if err != nil {
			return err
		}
		shared.BotConfigs = append(shared.BotConfigs, bot.WithEventListenerFunc(guildMemberJoin))
		return nil
	},
	Commands: []shared.Command{
		{
			Definition: discord.SlashCommandCreate{
				Name:                     "autojoinroles",
				Description:              "Give users a role when they join",
				DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionManageRoles),
				Contexts: []discord.InteractionContextType{
					discord.InteractionContextTypeGuild,
					discord.InteractionContextTypePrivateChannel},
				IntegrationTypes: []discord.ApplicationIntegrationType{
					discord.ApplicationIntegrationTypeGuildInstall},
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionSubCommand{
						Name:        "app",
						Description: "Give apps a role when they join (Leave empty to remove current)",
						Options: []discord.ApplicationCommandOption{
							&discord.ApplicationCommandOptionRole{
								Name:        "role",
								Description: "The role apps should get when they join the server",
							},
						},
					},
					&discord.ApplicationCommandOptionSubCommand{
						Name:        "user",
						Description: "Give users a role when they join (Leave empty to remove current)",
						Options: []discord.ApplicationCommandOption{
							&discord.ApplicationCommandOptionRole{
								Name:        "role",
								Description: "The role users should get when they join the server",
							}},
					},
				},
			},
			Interact: func(e *events.ApplicationCommandInteractionCreate) {
				var role string
				option := *e.SlashCommandInteractionData().SubCommandName
				var content string
				if len(e.SlashCommandInteractionData().Options) == 1 {
					var givenRole discord.Role = e.SlashCommandInteractionData().Role("role")
					role = givenRole.ID.String()
					botrole, err := getHighestRole(e.GuildID().String(), e.Client())
					if err != nil {
						logrus.Error(err)
					}
					if givenRole.Position >= botrole.Position {
						content = "<@&" + role + "> is not below the Bot's current highest role(<@&" + botrole.ID.String() + ">). That makes it unable to manage it."
					} else {
						if setAutoJoinRole(e.GuildID().String(), option, role) {
							content = "Updated auto join role for " + option + "s as <@&" + role + ">"
						} else {
							content = "Setup auto join role for " + option + "s as <@&" + role + ">"
						}
					}
				} else if setAutoJoinRole(e.GuildID().String(), option, role) {
					content = "Deleted auto join role for " + option + "s"
				}
				if content == "" {
					content = "No auto join role set for " + option + "s to delete."
				}
				err := e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(content).SetEphemeral(true).Build())
				if err != nil {
					logrus.Error(err)
				}
				purgeUnusedAutoJoinRoles(e.GuildID().String())
			},
		},
	},
}

func getHighestRole(guildID string, c bot.Client) (*discord.Role, error) {
	botmember, err := c.Rest().GetMember(snowflake.MustParse(guildID), c.ApplicationID())
	if err != nil {
		return nil, err
	}
	roles, err := c.Rest().GetRoles(snowflake.MustParse(guildID))
	if err != nil {
		return nil, err
	}
	var highestRole *discord.Role
	for _, roleID := range botmember.RoleIDs {
		for _, role := range roles {
			if role.ID == roleID {
				if highestRole == nil || role.Position > highestRole.Position {
					highestRole = &role
				}
				break
			}
		}
	}
	return highestRole, nil
}
