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
		option := i.ApplicationCommandData().Options[0].Name
		role := i.ApplicationCommandData().Options[0].Options[0].RoleValue(s, i.GuildID).ID
		var content string
		if setAutoJoinRole(i.GuildID, option, role) {
			content = "Setup auto join role for " + option + "s as <@&" + role + ">"
		} else {
			content = "Updated auto join role for " + option + "s as <@&" + role + ">"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

	},
}
