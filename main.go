package main

import (
	"context"
	"database/sql"
	"io"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/log2webhook"
)

var (
	db *sql.DB
)

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
	client, err := disgo.New(os.Getenv("BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildMembers,
				gateway.IntentDirectMessages,
			),
		),
		bot.WithEventListenerFunc(ready),
		bot.WithEventListenerFunc(applicationCommandInteractionCreate),
		bot.WithEventListenerFunc(autocompleteInteractionCreate),
		bot.WithEventListenerFunc(componentInteractionCreate),
		bot.WithEventListenerFunc(modalSubmitInteractionCreate),
		bot.WithEventListenerFunc(messageCreate),
		bot.WithEventListenerFunc(messageDelete),
		bot.WithEventListenerFunc(guildMemberJoin),
	)
	if err != nil {
		logrus.Fatal("error creating Discord session,", err)
		return
	} else {
		logrus.Info("Discord session created")
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		logrus.Error("error opening connection,", err)
		return
	}
	app, err := client.Rest().GetCurrentApplication()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("Bot is now running as '%s'!", app.Bot.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logrus.Info("Shutting down...")
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
