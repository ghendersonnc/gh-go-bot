package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofor-little/env"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func initBot(s *discordgo.Session) {
	s.AddHandler(messageHandler)
	s.AddHandler(commandHandler)
	s.Identify.Intents = discordgo.IntentsGuildMessages
	s.Identify.Presence.Status = string(discordgo.StatusOnline)

	err := s.Open()
	if err != nil {
		panic(err)
	}
	log.Println("Bot now running!")
}

func main() {
	if err := env.Load(".env"); err != nil {
		panic(err)
	}

	botToken := env.Get("BOT_TOKEN", "NO_TOKEN")
	if botToken == "NO_TOKEN" {
		panic("NO TOKEN AVAILABLE")
	}

	discord, err := discordgo.New("Bot " + botToken)

	if err != nil {
		panic(err)
	}

	initBot(discord)

	log.Println("Registering commands")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "833851489301692436", v)
		if err != nil {
			log.Panicf("Fuck: %v", err)
		}
		registeredCommands[i] = cmd
	}

	log.Println("Commands registered :)")

	defer func(discord *discordgo.Session) {
		err := discord.Close()
		if err != nil {
			log.Panicf("Error closing bot :: %v", err)
		}
	}(discord)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err != nil {
		panic("Bot could not close. Forcing...")
	}

	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, "833851489301692436", v.ID)
		if err != nil {
			log.Panicf("Cannot delete %v :: %v", v.Name, err)
		}
	}

	log.Println("Bot shutdown UwU")
}
