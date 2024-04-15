package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

func findAndDeleteUnusedMessages(c bot.Client) {
	for _, message := range getAllSavedMessages() {
		_, err := c.Rest().GetMessage(snowflake.MustParse(message.ChannelID), snowflake.MustParse(message.ID))
		if err != nil {
			tryDeleteUnusedMessage(message.ID)
		}
	}
}
