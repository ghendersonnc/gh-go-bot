package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gofor-little/env"
	"io"
	"log"
	"net/http"
	"time"
)

func gelbooruRequest(tags string) GelbooruPosts {
	gelbooruBaseUrl := fmt.Sprintf(
		"https://gelbooru.com/index.php?page=dapi&s=post&q=index&api_key=%s&user_id=%s&tags=%s%%20rating:general%%20sort:random&json=1&limit=1",
		env.Get("GELBOORU_API_KEY", "NO_KEY"),
		env.Get("GELBOORU_USER_ID", "NO_USER_ID"),
		tags)
	response, err := http.Get(gelbooruBaseUrl)
	if err != nil {
		log.Fatalf("Could not GET :: %v\n", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Could not close request body :: %v\n", err)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Could not read response body :: %v\n", err)
	}

	var result GelbooruPosts
	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("Could not decode JSON response :: %v\n", err)
	}

	return result
}

func genericResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string, embedUrl string, embedDescription string) {
	embeds := make([]*discordgo.MessageEmbed, 0)
	if embedUrl != "" {
		imageEmbed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x346beb,
			Description: embedDescription,
			Image: &discordgo.MessageEmbedImage{
				URL: embedUrl,
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     "Image",
		}

		embeds = append(embeds, imageEmbed)
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Embeds:  embeds,
		},
	})

	if err != nil {
		log.Printf("Could not respond :: %v\n", err)
	}
}
