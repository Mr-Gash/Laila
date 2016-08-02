package main

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"strings"
)

const (
	CHAT_PUBLIC  = iota // For all users on all channels.
	CHAT_PRIVATE = iota // For all users on whitelisted channels, or for admins on all.
	CHAT_ADMIN   = iota // For admins on all channels.
)

func TernaryString(b bool, t string, f string) string {
	if b {
		return t
	}

	return f
}

type ChatFunc func(*Bot, string, *discord.Message) bool

type ChatCommand struct {
	Func    ChatFunc
	Cmd     string
	Aliases []string
	HasArgs bool
	ArgHelp string // "!cmd <a> <b>"
	Help    string // "!cmd lets you do <x>, while <a> defines <y> and <b> does <z>."
	Access  int
}

func getCommandBody(content, cmd string) (bool, string) {
	low := strings.ToLower(content)
	if len(low) >= len(cmd)+1 &&
		low[0:len(cmd)+1] == cmd+" " {

		return true, content[len(cmd)+1 : len(content)]
	}

	return false, ""
}

func (c *ChatCommand) IsCommand(str string) (bool, string) {
	if !c.HasArgs {
		low := strings.ToLower(str)
		if low == c.Cmd {
			return true, ""
		}

		for x := range c.Aliases {
			if c.Aliases[x] == low {
				return true, ""
			}
		}

		return false, ""
	} else {
		good, body := getCommandBody(str, c.Cmd)
		if good {
			return true, body
		}

		for x := range c.Aliases {
			good, body = getCommandBody(str, c.Aliases[x])
			if good {
				return true, body
			}
		}

		return false, ""
	}
}

type Bot struct {
	Session   *discord.Session
	ID        string
	Self      *discord.User
	ChatHooks map[string]*ChatCommand
}

func NewBot(Token string) (*Bot, error) {
	session, err := discord.New(Token)
	if err != nil {
		return nil, err
	}

	err = session.Open()
	if err != nil {
		return nil, err
	}

	user, err := session.User("@me")
	if err != nil {
		return nil, err
	}

	fmt.Println("Username: ", user.Username)
	fmt.Println("UID: ", user.ID)

	bot := &Bot{
		Session:   session,
		ID:        user.ID,
		Self:      user,
		ChatHooks: make(map[string]*ChatCommand, 8),
	}

	session.AddHandler(bot.messageCreate)

	return bot, nil
}

func (b *Bot) Send(channel string, msg string) (string, error) {
	m, err := b.Session.ChannelMessageSend(channel, msg)
	if err != nil {
		return "", err
	}

	return m.ID, nil
}

func (b *Bot) messageCreate(s *discord.Session, m *discord.MessageCreate) {
	if m.Author.ID == b.ID {
		return
	}

	fmt.Println(m.Author.Username + "> " + m.Content)

	Admin := IsAdmin(m.Author.ID)
	Whitelisted := IsWhitelisted(m.ChannelID)

	var data *ChatCommand
	for i := range b.ChatHooks {
		data = b.ChatHooks[i]

		if data.Access != CHAT_PUBLIC {
			if data.Access == CHAT_ADMIN && !Admin {
				continue
			}

			if data.Access == CHAT_PRIVATE && !Whitelisted && !Admin {
				continue
			}
		}

		t, body := data.IsCommand(m.Content)
		if t {
			stop := data.Func(b, body, m.Message)
			if stop {
				return
			}
		}
	}
}
