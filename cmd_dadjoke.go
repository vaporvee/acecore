package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var dadjoke_command Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "dadjoke",
		Description: "Gives you a random joke that is as bad as your dad would tell them",
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return
		}

		type Joke struct {
			Joke string `json:"joke"`
		}
		var joke Joke
		err = json.Unmarshal(body, &joke)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: joke.Joke,
			},
		})
	},
}
