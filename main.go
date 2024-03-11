package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//TODO: add more error handlings

var db *sql.DB
var bot *discordgo.Session

func main() {
	godotenv.Load()

	var err error
	connStr := "postgresql://" + os.Getenv("DB_USER") + ":" + url.QueryEscape(os.Getenv("DB_PASSWORD")) + "@" + os.Getenv("DB_SERVER") + ":" + string(os.Getenv("DB_PORT")) + "/" + os.Getenv("DB_NAME") + "?sslmode=disable&application_name=Discord Bot"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	initTables()
	bot, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	} else {
		fmt.Println("Discord session created")
	}
	bot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentMessageContent | discordgo.IntentGuildMembers
	bot.AddHandler(ready)
	bot.AddHandler(interactionCreate)
	bot.AddHandler(messageCreate)
	bot.AddHandler(messageDelete)
	bot.AddHandler(guildMemberJoin)
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Printf("\nBot is now running as \"%s\"!", bot.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("\nShutting down...")
	bot.Close()
}
