package main

import (
	"github.com/vaporvee/acecore/cmd"
)

var Plugin = &cmd.Plugin{
	Name:     "Simple Fun",
	Commands: []cmd.Command{cmd_ask, cmd_cat, cmd_dadjoke},
}
