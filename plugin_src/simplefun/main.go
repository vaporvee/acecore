package main

import (
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/cmd"
)

var Plugin = &cmd.Plugin{
	Name: "Simple Fun",
	Register: func(e *events.Ready) error {
		app, err := e.Client().Rest().GetCurrentApplication()
		if err != nil {
			return err
		}
		logrus.Infof("%s has a working plugin called \"testplugin\"", app.Bot.Username)
		return nil
	},
	Commands: []cmd.Command{cmd_ask, cmd_cat, cmd_dadjoke},
}
