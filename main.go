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

//TODO: make the codebase less garbage

var db *sql.DB

func main() {
	godotenv.Load()

	var err error
	connStr := "postgresql://" + os.Getenv("DB_USER") + ":" + url.QueryEscape(os.Getenv("DB_PASSWORD")) + "@" + os.Getenv("DB_SERVER") + ":" + string(os.Getenv("DB_PORT")) + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	initTable()

	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	} else {
		fmt.Println("Discord session created")
	}
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds
	defer removeCommandFromAllGuilds(discord)
	discord.AddHandler(ready)
	discord.AddHandler(interactionCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Printf("Bot is now running as \"%s\"!", discord.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("\nShutting down...")
	defer removeCommandFromAllGuilds(discord)
	discord.Close()
}
