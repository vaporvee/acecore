package main

import "github.com/bwmarrin/discordgo"

var autojoinroles_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "autojoinroles",
		Description: "Give users a role when they join",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "bot",
				Description: "Give bots a role when they join (Leave empty to remove current)",
				Options: []*discordgo.ApplicationCommandOption{
					{

						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "The role bots should get when they join the server",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "user",
				Description: "Give users a role when they join (Leave empty to remove current)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "The role users should get when they join the server",
					}},
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var role string
		option := i.ApplicationCommandData().Options[0].Name
		var content string
		if len(i.ApplicationCommandData().Options[0].Options) == 1 {
			role = i.ApplicationCommandData().Options[0].Options[0].RoleValue(s, i.GuildID).ID
			if setAutoJoinRole(i.GuildID, option, role) {
				content = "Updated auto join role for " + option + "s as <@&" + role + ">"
			} else {
				content = "Setup auto join role for " + option + "s as <@&" + role + ">"
			}
		} else if setAutoJoinRole(i.GuildID, option, role) {
			content = "Deleted auto join role for " + option + "s"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		purgeUnusedAutoJoinRoles(i.GuildID)
	},
}
