package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	log.Println(err)
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds
	discord.AddHandler(messageCreate)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { //bot doesn't take its own messages
		return
	}
	if m.Message.Content == "test" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "TESTED WOOOWW",
		})
	}
}
