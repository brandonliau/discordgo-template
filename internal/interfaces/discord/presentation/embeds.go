package presentation

import (
	"fmt"
	"strings"
	"time"

	"discordgo-skeleton/internal/application/view"

	"github.com/bwmarrin/discordgo"
)

func NoticeEmbed(title string, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func WeatherEmbed(wv view.WeatherView) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s", wv.Location()),
		Color:       Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🌀 Condition",
				Value:  fmt.Sprintf("> `%s`", wv.Condition),
				Inline: false,
			},
			{
				Name:   "🌡️ Temperature",
				Value:  fmt.Sprintf("> Actual: `%.0f°F` \n > Feels like: `%.0f°F`", wv.TempF, wv.FeelsLikeF),
				Inline: false,
			},
			{
				Name:   "⬆️ High / ⬇️ Low",
				Value:  fmt.Sprintf("> `%.0f°F / %.0f°F`", wv.HighF, wv.LowF),
				Inline: false,
			},
			{
				Name:   "💧 Humidity",
				Value:  fmt.Sprintf("> `%d%%`", wv.Humidity),
				Inline: false,
			},
			{
				Name:   "💨 Wind",
				Value:  fmt.Sprintf("> `%.0f mph %s`", wv.WindMph, wv.WindDir),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func PinsEmbed(views []view.WeatherView) *discordgo.MessageEmbed {
	var builder strings.Builder
	for _, v := range views {
		fmt.Fprintf(&builder, "**%s** — %.0f°F, %s\n", v.Location(), v.TempF, v.Condition)
	}
	return &discordgo.MessageEmbed{
		Title:       "Pinned locations",
		Description: builder.String(),
		Color:       Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
