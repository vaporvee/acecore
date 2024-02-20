package main

import "github.com/bwmarrin/discordgo"

var ask_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "ask",
		Description: "Ask anything and get a gif as response!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "The question you want to ask",
				Required:    true,
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:  discordgo.EmbedTypeImage,
						Color: hexToDecimal(color["primary"]),
						Image: &discordgo.MessageEmbedImage{
							URL: simpleGetFromAPI("image", "https://yesno.wtf/api").(string),
						},
					},
				},
			},
		})
	},
}
