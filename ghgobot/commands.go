package ghgobot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var (
	Commands = []*discordgo.ApplicationCommand{
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

func respondPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	GenericResponse(s, i, "Pong!", "", "")
}

func gelbooru(s *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(i.ApplicationCommandData().Options))
	for _, opt := range i.ApplicationCommandData().Options {
		optionMap[opt.Name] = opt
	}

	var tags string
	var post GelbooruPosts
	var formattedTags string
	if len(optionMap) == 0 {
		tags = ""
		post = GelbooruRequest("")
	} else if len(optionMap) > 0 {
		if val, ok := optionMap["tags"]; ok {
			formattedTags = strings.Replace(val.StringValue(), " ", "%20", -1)
			tags = val.StringValue()
		}

		if val, ok := optionMap["rating"]; ok {
			formattedTags = formattedTags + "%20rating:" + val.StringValue()
		}

		post = GelbooruRequest(formattedTags)
	}

	if len(post.Post) == 0 {
		GenericResponse(s, i, "No images found!", "", "")
		return
	}
	originalUrl := fmt.Sprintf("https://gelbooru.com/index.php?page=post&s=view&id=%d&tags=%s", post.Post[0].ID, formattedTags)
	GenericResponse(s, i, "Post from Gelbooru", post.Post[0].FileURL, fmt.Sprintf("%s\n%s", originalUrl, tags))
}
