package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gofor-little/env"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with pong",
		},
		{
			Name:        "gelbooru",
			Description: "Send a request to gelbooru using tags",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tags",
					Description: "Tags for the request. General is always included automatically",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}
)

func gelbooruRequest(tags string) GelbooruPostResponse {
	gelbooruBaseUrl := fmt.Sprintf(
		"https://gelbooru.com/index.php?page=dapi&s=post&q=index&api_key=%s&user_id=%s&tags=%s%%20rating:general%%20sort:random&json=1&limit=1",
		env.Get("GELBOORU_API_KEY", "NO_KEY"),
		env.Get("GELBOORU_USER_ID", "NO_USER_ID"),
		tags)
	response, err := http.Get(gelbooruBaseUrl)
	if err != nil {
		log.Fatalf("Could not GET :: %v\n", err)
	}
	log.Println(gelbooruBaseUrl)
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

	var result GelbooruPostResponse
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

func respondPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	genericResponse(s, i, "Pong!", "", "")
}

func gelbooru(s *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	var tags string
	var post GelbooruPostResponse

	if len(optionMap) == 0 {
		tags = ""
		post = gelbooruRequest("")
	} else if len(optionMap) > 0 {
		if val, ok := optionMap["tags"]; ok {
			tags = strings.Replace(val.StringValue(), " ", "%20", -1)
		}

		if val, ok := optionMap["rating"]; ok {
			tags = tags + "%20rating:" + val.StringValue()
		}

		post = gelbooruRequest(tags)
	}

	if len(post.Post) == 0 {
		genericResponse(s, i, "No images found!", "", "")
		return
	}

	genericResponse(s, i, "Post from Gelbooru", post.Post[0].FileURL, strings.Replace(tags, "%20", " ", -1))
}
