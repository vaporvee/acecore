package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
)

type UserExtend struct {
	GlobalName           string           `json:"global_name"`
	AvatarDecorationData AvatarDecoration `json:"avatar_decoration_data"`
}
type AvatarDecoration struct {
	Asset string `json:"asset"`
	SkuID string `json:"sku_id"`
	URL   string
}

var userFlagsString map[discord.UserFlags]string = map[discord.UserFlags]string{
	discord.UserFlagDiscordEmployee:           "<:Discord_Employee:1224708831419043942>[`Discord Employee`](https://discord.com/company)",
	discord.UserFlagPartneredServerOwner:      "<:Discord_Partner:1224708689190060092>[`Discord Partner`](https://discord.com/partners)",
	discord.UserFlagHypeSquadEvents:           "<:Hypesquad_Events:1224708685494747237>[`HypeSquad Events`](https://discord.com/hypesquad)",
	discord.UserFlagBugHunterLevel1:           "<:Bug_Hunter_Level_1:1224708828415918231>[`Bug Hunter Level 1`](https://support.discord.com/hc/en-us/articles/360046057772-Discord-Bugs)",
	discord.UserFlagHouseBravery:              "<:Hypesquad_Bravery:1224708678905630801>[`HypeSquad Bravery`](https://discord.com/settings/hypesquad-online)",
	discord.UserFlagHouseBrilliance:           "<:Hypesquad_Brilliance:1224708677584424961>[`HypeSquad Brilliance`](https://discord.com/settings/hypesquad-online)",
	discord.UserFlagHouseBalance:              "<:Hypequad_Balance:1224708826901516309>[`HypeSquad Balance`](https://discord.com/settings/hypesquad-online)",
	discord.UserFlagEarlySupporter:            "<:Early_Supporter:1224708674065272873>[`Early Supporter`](https://discord.com/settings/premium)",
	discord.UserFlagTeamUser:                  "`TeamUser`",
	discord.UserFlagBugHunterLevel2:           "<:Bug_Hunter_Level_2:1224708682378383461>[`Bug Hunter Level 2`](https://support.discord.com/hc/en-us/articles/360046057772-Discord-Bugs)",
	discord.UserFlagVerifiedBot:               "",
	discord.UserFlagEarlyVerifiedBotDeveloper: "<:Early_Verified_Bot_Developer:1224708675294203934>`Early Verified Bot Developer`",
	discord.UserFlagDiscordCertifiedModerator: "<:Discord_Certified_Moderator:1224708830223532124>[`Discord Certified Moderator`](https://discord.com/safety)",
	discord.UserFlagBotHTTPInteractions:       "`BotHTTPInteractions`",
	discord.UserFlagActiveDeveloper:           "<:Active_Developer:1224708676611215380>[`Active Developer`](https://support-dev.discord.com/hc/en-us/articles/10113997751447?ref=badge)",
}

var cmd_userinfo Command = Command{
	Definition: discord.SlashCommandCreate{
		Name:        "info",
		Description: "Gives you information about a user or this bot.",
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "user",
				Description: "Gives you information about a user and its profile images.",
				Options: []discord.ApplicationCommandOption{
					&discord.ApplicationCommandOptionUser{
						Name:        "user",
						Description: "The user you need information about.",
						Required:    true,
					},
				},
			},
			&discord.ApplicationCommandOptionSubCommand{
				Name:        "bot-service",
				Description: "Gives you information about this bot's server service.",
			},
		},
	},
	Interact: func(e *events.ApplicationCommandInteractionCreate) {
		switch *e.SlashCommandInteractionData().SubCommandName {
		case "user":
			var user discord.User = e.SlashCommandInteractionData().User("user")
			var userHasFlags string = fetchFlagStrings(user)
			var userType string = "User"
			if user.Bot {
				userType = "Unverified Bot"
				if user.PublicFlags&discord.UserFlagVerifiedBot != 0 {
					userType = "Verified Bot"
				}
			} else if user.System {
				userType = "System"
			}
			embedBuilder := discord.NewEmbedBuilder()
			embedBuilder.SetThumbnail(checkDefaultPb(user))
			embedBuilder.AddField("ID", user.ID.String(), false)
			embedBuilder.AddField("Type", userType, true)
			if user.GlobalName != nil {
				embedBuilder.AddField("Global name", *user.GlobalName, true)
			}
			embedBuilder.AddField("Username", user.Username, true)
			if userHasFlags != "" {
				embedBuilder.AddField("Badges", userHasFlags, false)
			}
			if user.Discriminator != "0" {
				embedBuilder.AddField("Discriminator", user.Discriminator, false)
			}
			if user.AccentColor != nil {
				embedBuilder.AddField("Accent color", strconv.Itoa(*user.AccentColor), true)
			}
			if user.AvatarDecorationURL() != nil {
				value := fmt.Sprintf("[PNG (animated)](%s)\n[PNG](%s)", *user.AvatarDecorationURL(), *user.AvatarDecorationURL()+"?passthrough=false")
				embedBuilder.AddField("Avatar decoration", value, true)
			}
			creation := "<:discord_member:1224717530078253166> " + discord.TimestampStyleLongDateTime.FormatTime(user.CreatedAt()) + "-" + discord.TimestampStyleRelative.FormatTime(user.CreatedAt())
			embedBuilder.AddField("Created at", creation, false)

			if user.BannerURL() != nil {
				value := fmt.Sprint(*user.BannerURL())
				embedBuilder.SetImage(value)
			}
			embedBuilder.SetTitle("User info")
			embedBuilder.SetDescription(user.Mention())
			embedBuilder.SetColor(hexToDecimal(color["primary"]))
			embedBuilder.SetFooterText("Currently a bit broken because of Discord's constant user API changes")
			err := e.CreateMessage(discord.NewMessageCreateBuilder().
				SetEmbeds(embedBuilder.Build()).
				Build())
			if err != nil {
				logrus.Error(err)
			}

		case "bot-service":

		}

	},
	AllowDM: true,
}

func checkDefaultPb(user discord.User) string {
	if user.AvatarURL() == nil {
		return "https://discord.com/assets/ac6f8cf36394c66e7651.png"
	}
	return *user.AvatarURL()
}

func fetchFlagStrings(user discord.User) string {
	var userHasFlagsString string
	for flag, flagName := range userFlagsString {
		if flag&user.PublicFlags != 0 {
			userHasFlagsString += flagName + ", "
		}
	}
	if user.AvatarDecorationData != nil && user.AvatarDecorationData.Asset == "a_5e1210779d99ece1c0b4f438a5bc6e72" {
		userHasFlagsString += "<:Limited_Lootbox_Clown:1224714172705804300>[`Lootbox Clown`](https://discord.com/settings/Lootboxes)"
	}
	/*
		if user.PremiumType > 0 {
			userHasFlagsString += "<:Nitro:1224708672492666943>[`Nitro`](https://discord.com/settings/premium), "
		}
		if user.Bot {
			appuser := bot.State.Application
			if appuser.Flags&1<<23 != 0 {
				userHasFlagsString += "<:Supports_Commands:1224848976201646100>[`Supports Commands`](https://discord.com/blog/welcome-to-the-new-era-of-discord-apps?ref=badge)"
			}
			if appuser.Flags&1<<6 != 0 {
				userHasFlagsString += "<:Uses_Automod:1224862880982106202>`Uses Automod`"
			}
		}
	*/
	returnString := strings.TrimSuffix(userHasFlagsString, ", ")
	return returnString
}
