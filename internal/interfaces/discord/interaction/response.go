package interaction

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	maxEmbeds        = 10
	maxActionRows    = 5
	maxButtonsPerRow = 5
)

var (
	ErrTooManyEmbeds = errors.New("interaction: more than 10 embeds")
	ErrTooManyRows   = errors.New("interaction: components exceed 5 action rows")
)

func Response(responseType discordgo.InteractionResponseType, opts ...ResponseOption) (*discordgo.InteractionResponse, error) {
	data := &discordgo.InteractionResponseData{}
	var errs []error
	for _, opt := range opts {
		if err := opt(data); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return &discordgo.InteractionResponse{Type: responseType, Data: data}, nil
}

func InitialResponse(opts ...ResponseOption) (*discordgo.InteractionResponse, error) {
	return Response(discordgo.InteractionResponseChannelMessageWithSource, opts...)
}

func UpdateResponse(opts ...ResponseOption) (*discordgo.InteractionResponse, error) {
	return Response(discordgo.InteractionResponseUpdateMessage, opts...)
}

type ResponseOption func(*discordgo.InteractionResponseData) error

func WithContent(content string) ResponseOption {
	return func(data *discordgo.InteractionResponseData) error {
		data.Content = content
		return nil
	}
}

func WithEmbeds(embeds ...*discordgo.MessageEmbed) ResponseOption {
	return func(data *discordgo.InteractionResponseData) error {
		if len(embeds) > maxEmbeds {
			return fmt.Errorf("%w (got %d)", ErrTooManyEmbeds, len(embeds))
		}
		data.Embeds = embeds
		return nil
	}
}

func WithComponents(components ...discordgo.MessageComponent) ResponseOption {
	return func(data *discordgo.InteractionResponseData) error {
		var rows []discordgo.MessageComponent
		var buttons []discordgo.MessageComponent

		for _, c := range components {
			switch c.(type) {
			case discordgo.Button, *discordgo.Button:
				buttons = append(buttons, c)
				if len(buttons) == maxButtonsPerRow {
					rows = append(rows, discordgo.ActionsRow{Components: buttons})
					buttons = nil
				}
			case discordgo.SelectMenu, *discordgo.SelectMenu:
				if len(buttons) > 0 {
					rows = append(rows, discordgo.ActionsRow{Components: buttons})
					buttons = nil
				}
				rows = append(rows, discordgo.ActionsRow{Components: []discordgo.MessageComponent{c}})
			default:
				return fmt.Errorf("interaction: unsupported component type %T", c)
			}
		}
		if len(buttons) > 0 {
			rows = append(rows, discordgo.ActionsRow{Components: buttons})
		}

		if len(rows) > maxActionRows {
			return fmt.Errorf("%w (got %d)", ErrTooManyRows, len(rows))
		}
		data.Components = rows
		return nil
	}
}

func WithEphemeral() ResponseOption {
	return func(data *discordgo.InteractionResponseData) error {
		data.Flags |= discordgo.MessageFlagsEphemeral
		return nil
	}
}
