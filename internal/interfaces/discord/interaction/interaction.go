package interaction

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type HandleFunc func(*discordgo.Session, *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)

func EncodeCustomID(routingKey string, params url.Values) (string, error) {
	if strings.TrimSpace(routingKey) == "" || strings.ContainsAny(routingKey, "?&=") {
		return "", fmt.Errorf("invalid routing key %q", routingKey)
	}
	customID := routingKey
	if len(params) > 0 {
		customID += "?" + params.Encode()
	}
	if len(customID) > 100 {
		return "", fmt.Errorf("custom ID exceeds Discord's 100-character limit")
	}
	return customID, nil
}

func DecodeCustomID(customID string) (string, url.Values, error) {
	routingKey, query, hasQuery := strings.Cut(customID, "?")
	if strings.TrimSpace(routingKey) == "" {
		return "", nil, fmt.Errorf("custom ID has no routing key")
	}
	if !hasQuery {
		return routingKey, make(url.Values), nil
	}
	params, err := url.ParseQuery(query)
	if err != nil {
		return "", nil, fmt.Errorf("decode custom ID: %w", err)
	}
	return routingKey, params, nil
}

func RoutingKey(customID string) (string, error) {
	key, _, err := DecodeCustomID(customID)
	return key, err
}
