package discord

import (
	"discordgo-skeleton/internal/interfaces/discord/interaction"

	"github.com/bwmarrin/discordgo"
)

func (g *gateway) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var customID string
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		customID = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		customID = i.MessageComponentData().CustomID
	case discordgo.InteractionModalSubmit:
		customID = i.ModalSubmitData().CustomID
	default:
		return
	}

	routingKey, _, err := interaction.DecodeCustomID(customID)
	if err != nil {
		g.logger.Error("Failed to decode custom ID %s: %v", customID, err)
		return
	}

	handleFunc, ok := g.handleFuncs[routingKey]
	if !ok {
		g.logger.Error("%s interaction handler not found", routingKey)
		return
	}

	rsp, err := handleFunc(s, i)
	if err != nil {
		g.logger.Error("Failed to execute %s interaction handler: %v", routingKey, err)
		rsp, err = interaction.InitialResponse(
			interaction.WithContent("Something went wrong."),
			interaction.WithEphemeral(),
		)
		if err != nil {
			g.logger.Error("Failed to build fallback response: %v", err)
			return
		}
	}

	if rsp == nil {
		return
	}
	if err := s.InteractionRespond(i.Interaction, rsp); err != nil {
		g.logger.Error("Failed to send interaction response: %v", err)
		return
	}

	g.logger.Debug("Successfully executed %s interaction handler", routingKey)
}
