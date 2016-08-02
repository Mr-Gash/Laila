package main

import (
	discord "github.com/bwmarrin/discordgo"
)

func chat_help(b *Bot) {
	// No args chat command, just lists everything.
	b.ChatHooks["!help"] = &ChatCommand{
		Func: func(x *Bot, _ string, m *discord.Message) bool {
			str := "Type \"**!help !cmd**\" to learn more about \"**!cmd**\"\n\n"

			public := ""
			private := ""
			admin := ""

			for k := range b.ChatHooks {
				data := b.ChatHooks[k]
				if data.Help == "" && data.ArgHelp == "" {
					continue
				}

				switch data.Access{
				case CHAT_PUBLIC:
					public += data.Cmd + "\n"
				case CHAT_PRIVATE:
					private += data.Cmd + "\n"
				case CHAT_ADMIN:
					admin += data.Cmd + "\n"
				}
			}

			str += "```\n"
			if public != "" {
				str += "- PUBLIC -\n" + public + "\n"
			}

			if private != "" {
				str += "- WHITELISTED CHANNELS -\n" + private + "\n"
			}

			if admin != "" {
				str += "- BOT ADMINS ONLY -\n" + admin
			}
			str += "```"

			if public != "" || private != "" || admin != "" {
				b.Send(m.ChannelID, str)
			} else {
				b.Send(m.ChannelID, "No chat commands to display :frowning:")
			}

			return true
		},
		Cmd:     "!help",
		Help:    "Lists all chat commands or lists specific help for a chat command.",
		ArgHelp: "!help\n!help <cmd>",
	}
}
