package main

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"time"
)

func chat_toys(b *Bot) {
	rand.Seed(time.Now().UnixNano())

	b.ChatHooks["!roll"] = &ChatCommand{
		Func: func(x *Bot, _ string, m *discord.Message) bool {
			x.Send(m.ChannelID, strconv.Itoa(rand.Intn(100)+1))
			return true
		},
		Cmd:    "!roll",
		Access: CHAT_PRIVATE,
		Help:   "Generates a random number of the range of 1 to 100.",
	}

	b.ChatHooks["!flip"] = &ChatCommand{
		Func: func(x *Bot, _ string, m *discord.Message) bool {
			x.Send(m.ChannelID, TernaryString(rand.Intn(2) == 0, "Heads", "Tails"))
			return true
		},
		Cmd:    "!flip",
		Access: CHAT_PRIVATE,
		Help:   "Flips a coin. Gives heads or tails.",
	}

	b.ChatHooks["!name"] = &ChatCommand{
		Func: func(x *Bot, name string, m *discord.Message) bool {
			channel, err := x.Session.Channel(m.ChannelID)
			if err != nil {
				fmt.Println("Error changing name: ", err.Error())
				x.Send(m.ChannelID, ":frowning: :x:")
				return true
			}

			err = x.Session.GuildMemberNickname(channel.GuildID, "@me/nick", name)
			if err != nil {
				x.Send(m.ChannelID, "Error: "+err.Error())
				return true
			}

			x.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			x.Send(m.ChannelID, ":ok_hand:")
			return true
		},
		Cmd:     "!name",
		Access:  CHAT_ADMIN,
		HasArgs: true,
		Help:    "Changes my name to something else.",
		ArgHelp: "!name <new name>",
	}

	b.ChatHooks[":shrug:"] = &ChatCommand{
		Func: func(x *Bot, _ string, m *discord.Message) bool {
			channel, err := x.Session.Channel(m.ChannelID)
			if err != nil {
				fmt.Println("Error shrugging: ", err.Error())
				x.Send(m.ChannelID, ":frowning: :x:")
				return true
			}

			original := x.Self.Username

			err = x.Session.GuildMemberNickname(channel.GuildID, "@me/nick", m.Author.Username)
			if err != nil {
				x.Send(m.ChannelID, "Error: "+err.Error())
				return true
			}

			err = x.Session.ChannelMessageDelete(m.ChannelID, m.ID)
			if err == nil {
				x.Send(m.ChannelID, "¯\\_(ツ)_/¯")
				x.Session.GuildMemberNickname(channel.GuildID, "@me/nick", original)
			}

			return true
		},
		Cmd: ":shrug:",
		Access: CHAT_PRIVATE,
		Help: "¯\\_(ツ)_/¯",
	}
}
