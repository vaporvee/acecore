package main

import (
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"database/sql"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/log2webhook"
)

//TODO: add more error handlings

var db *sql.DB
var bot *discordgo.Session

func main() {
	logrusInitFile()
	var err error
	godotenv.Load()
	connStr := "postgresql://" + os.Getenv("DB_USER") + ":" + url.QueryEscape(os.Getenv("DB_PASSWORD")) + "@" + os.Getenv("DB_SERVER") + ":" + string(os.Getenv("DB_PORT")) + "/" + os.Getenv("DB_NAME") + "?sslmode=disable&application_name=Discord Bot"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		logrus.Fatal(err)
	}
	initTables()
	bot, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		logrus.Fatal("error creating Discord session,", err)
		return
	} else {
		logrus.Info("Discord session created")
	}
	bot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentMessageContent | discordgo.IntentGuildMembers
	bot.AddHandler(ready)
	bot.AddHandler(interactionCreate)
	bot.AddHandler(messageCreate)
	bot.AddHandler(messageDelete)
	bot.AddHandler(guildMemberJoin)
	err = bot.Open()
	if err != nil {
		logrus.Error("error opening connection,", err)
		return
	}
	logrus.Infof("Bot is now running as '%s'!", bot.State.User.Username)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logrus.Info("Shutting down...")
	bot.Close()
}

func logrusInitFile() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	timestamp := time.Now().Unix()

	var file_name string = "logs/bot." + strconv.FormatInt(timestamp, 10) + ".log"
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	log, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}

	mw := io.MultiWriter(os.Stdout, log, &log2webhook.WebhookWriter{})
	logrus.SetOutput(mw)
}
