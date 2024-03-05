package main

import (
	"github.com/bwmarrin/discordgo"
)

var cmd_dadjoke Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "dadjoke",
		Description: "Gives you a random joke that is as bad as your dad would tell them",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		respond(s, i.Interaction, simpleGetFromAPI("joke", "https://icanhazdadjoke.com/").(string), false)
	},
}
