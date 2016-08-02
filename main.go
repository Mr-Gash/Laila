package main

import (
	"flag"
	"fmt"
)

var (
	// TODO: Make Whitelist/Admins be in a database and not hardcoded.
	ChatWhitelist = map[string]bool{}
	BotAdmins = map[string]bool{}

	Laila    *Bot
	BotToken string
)

func IsAdmin(id string) bool {
	IsAdmin, ok := BotAdmins[id]
	return ok && IsAdmin
}

func IsWhitelisted(channel string) bool {
	IsListed, ok := ChatWhitelist[channel]
	return ok && IsListed
}

func init() {
	flag.StringVar(&BotToken, "t", BotToken, "Bot Account Token")
	flag.Parse()
}

func main() {
	fmt.Println("Token: ", BotToken)

	var err error
	Laila, err := NewBot(BotToken)
	if err != nil {
		fmt.Println("Failed to run bot: ", err.Error())
		return
	}

	// Chat hooks added here.
	chat_toys(Laila)
	chat_yt(Laila)

	// Help hook should always be added last.
	chat_help(Laila)

	fmt.Println("Initialized Laila. Press Ctrl+C to end process.")

	<-make(chan struct{})
	return
}
