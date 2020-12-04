package main

// https://discord.com/api/oauth2/authorize?client_id=784151923225657364&permissions=388160&scope=bot

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token string
	LastURL string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	if Token == "" {
		flag.PrintDefaults()
		return
	}

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Could not start Discord session", err)
		return
	}

	discord.AddHandler(watchChat)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()

	if err != nil {
		fmt.Println("Could not open connection: ", err)
		return
	}

	fmt.Println("Bot is currently running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<- sc

	fmt.Println("Terminating bot..")
	discord.Close()
}

func generateArchiveLink() string {
	u, err := url.Parse(LastURL)
	if err != nil {
		fmt.Println("Invalid URL: ", LastURL)
	}

	cleaned := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

	return fmt.Sprintf("https://archive.is/%s*", cleaned)
}


func watchChat(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Embeds) > 0 {
		for _, embed := range m.Embeds {
			LastURL = embed.URL
		}
	}

	if len(m.Mentions) == 0 {
		return
	}

	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			if LastURL != "" {
				s.ChannelMessageSend(m.ChannelID, generateArchiveLink())
			}
		}
	}
}