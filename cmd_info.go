package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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

var userFlagsString map[discordgo.UserFlags]string = map[discordgo.UserFlags]string{
	discordgo.UserFlagDiscordEmployee:           "<:Discord_Employee:1224708831419043942>[`Discord Employee`](https://discord.com/company)",
	discordgo.UserFlagDiscordPartner:            "<:Discord_Partner:1224708689190060092>[`Discord Partner`](https://discord.com/partners)",
	discordgo.UserFlagHypeSquadEvents:           "<:Hypesquad_Events:1224708685494747237>[`HypeSquad Events`](https://discord.com/hypesquad)",
	discordgo.UserFlagBugHunterLevel1:           "<:Bug_Hunter_Level_1:1224708828415918231>[`Bug Hunter Level 1`](https://support.discord.com/hc/en-us/articles/360046057772-Discord-Bugs)",
	discordgo.UserFlagHouseBravery:              "<:Hypesquad_Bravery:1224708678905630801>[`HypeSquad Bravery`](https://discord.com/settings/hypesquad-online)",
	discordgo.UserFlagHouseBrilliance:           "<:Hypesquad_Brilliance:1224708677584424961>[`HypeSquad Brilliance`](https://discord.com/settings/hypesquad-online)",
	discordgo.UserFlagHouseBalance:              "<:Hypequad_Balance:1224708826901516309>[`HypeSquad Balance`](https://discord.com/settings/hypesquad-online)",
	discordgo.UserFlagEarlySupporter:            "<:Early_Supporter:1224708674065272873>[`Early Supporter`](https://discord.com/settings/premium)",
	discordgo.UserFlagTeamUser:                  "`TeamUser`",
	discordgo.UserFlagSystem:                    "",
	discordgo.UserFlagBugHunterLevel2:           "<:Bug_Hunter_Level_2:1224708682378383461>[`Bug Hunter Level 2`](https://support.discord.com/hc/en-us/articles/360046057772-Discord-Bugs)",
	discordgo.UserFlagVerifiedBot:               "",
	discordgo.UserFlagVerifiedBotDeveloper:      "<:Early_Verified_Bot_Developer:1224708675294203934>`Early Verified Bot Developer`",
	discordgo.UserFlagDiscordCertifiedModerator: "<:Discord_Certified_Moderator:1224708830223532124>[`Discord Certified Moderator`](https://discord.com/safety)",
	1 << 19: "`BotHTTPInteractions`",
	1 << 22: "<:Active_Developer:1224708676611215380>[`Active Developer`](https://support-dev.discord.com/hc/en-us/articles/10113997751447?ref=badge)",
}

var cmd_userinfo Command = Command{
	Definition: discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Gives you information about a user or this bot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "user",
				Description: "Gives you information about a user and its profile images.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "The user you need information about.",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "bot-service",
				Description: "Gives you information about this bot's server service.",
			},
		},
	},
	Interact: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "user":
			var user *discordgo.User = i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
			var extendedUser UserExtend = extendedUserFromAPI(user.ID)
			var userHasFlags string = fetchFlagStrings(user, extendedUser.AvatarDecorationData.Asset)
			var userType string = "User"
			if user.Bot {
				userType = "Unverified Bot"
				if user.PublicFlags&discordgo.UserFlagVerifiedBot != 0 {
					userType = "Verified Bot"
				}
			} else if user.System {
				userType = "System"
			}
			createdate, err := discordgo.SnowflakeTimestamp(user.ID)
			if err != nil {
				logrus.Error(err)
			}
			err = respondEmbed(i.Interaction, discordgo.MessageEmbed{
				Title:       extendedUser.GlobalName + " user info",
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: user.AvatarURL("512")},
				Description: user.Mention(),
				Type:        discordgo.EmbedTypeArticle,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "test",
						Value: fmt.Sprint(s.State.Application.BotPublic),
					},
					{
						Name:  "ID",
						Value: user.ID,
					},
					{
						Name:   "Type",
						Value:  userType,
						Inline: true,
					},
					{
						Name:   "Global name",
						Value:  extendedUser.GlobalName,
						Inline: true,
					},
					{
						Name:   "Username",
						Value:  user.Username,
						Inline: true,
					},
					{
						Name:  "Badges",
						Value: userHasFlags,
					},
					{
						Name:   "Discriminator",
						Value:  user.Discriminator,
						Inline: true,
					},
					{
						Name:   "Accent color",
						Value:  "#" + decimalToHex(user.AccentColor),
						Inline: true,
					},
					{
						Name:   "Avatar Decoration",
						Value:  "[PNG (animated)](" + extendedUser.AvatarDecorationData.URL + ")\n[PNG](" + extendedUser.AvatarDecorationData.URL + "?size=4096&passthrough=false)\nSKU ID: `" + extendedUser.AvatarDecorationData.SkuID + "`",
						Inline: true,
					},
					{
						Name:  "Created at",
						Value: "<:discord_member:1224717530078253166> <t:" + fmt.Sprint(createdate.Unix()) + ":f> - <t:" + fmt.Sprint(createdate.Unix()) + ":R>",
					},
				},
				Color: hexToDecimal(color["primary"]),
				Image: &discordgo.MessageEmbedImage{URL: user.BannerURL("512")},
			}, false)
			if err != nil {
				logrus.Error(err)
			}
		case "bot-service":

		}

	},
	AllowDM: true,
}

func fetchFlagStrings(user *discordgo.User, decorationAsset string) string {
	var userHasFlagsString string
	for flag, flagName := range userFlagsString {
		if user.PublicFlags&flag != 0 {
			userHasFlagsString += flagName + ", "
		}
	}
	if user.PremiumType > 0 {
		userHasFlagsString += "<:Nitro:1224708672492666943>[`Nitro`](https://discord.com/settings/premium), "
	}
	if decorationAsset == "a_5e1210779d99ece1c0b4f438a5bc6e72" {
		userHasFlagsString += "<:Limited_Lootbox_Clown:1224714172705804300>[`Lootbox Clown`](https://discord.com/settings/Lootboxes)"
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

	returnString := strings.TrimSuffix(userHasFlagsString, ", ")
	return returnString
}

func extendedUserFromAPI(userID string) UserExtend {
	client := &http.Client{}
	var userExtend UserExtend
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/"+userID, nil)
	if err != nil {
		logrus.Error(err)
		return userExtend
	}
	req.Header.Add("Authorization", "Bot "+os.Getenv("BOT_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return userExtend
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfterHeader(res.Header)
		if retryAfter > 0 {
			time.Sleep(retryAfter)
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return userExtend
	}

	json.Unmarshal(body, &userExtend)
	if userExtend.AvatarDecorationData.Asset != "" {
		userExtend.AvatarDecorationData.URL = "https://cdn.discordapp.com/avatar-decoration-presets/" + userExtend.AvatarDecorationData.Asset + ".png"
	}
	return userExtend
}

func parseRetryAfterHeader(headers http.Header) time.Duration {
	retryAfterStr := headers.Get("Retry-After")
	if retryAfterStr == "" {
		return 0
	}

	retryAfter, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		return 0
	}

	return time.Duration(retryAfter) * time.Millisecond
}
