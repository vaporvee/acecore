package main

import "github.com/bwmarrin/discordgo"

var cmd_ask Command = Command{
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
		respondEmbed(s, i.Interaction, discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeImage,
			Color: hexToDecimal(color["primary"]),
			Image: &discordgo.MessageEmbedImage{
				URL: simpleGetFromAPI("image", "https://yesno.wtf/api").(string),
			}}, false)
	},
}
