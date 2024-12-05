package notifier

import (
	"DiscordTemplate/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type discordNotifier struct {
	session *discordgo.Session
	logger  logger.Logger
}

func NewDiscordNotifier(s *discordgo.Session) *discordNotifier {
	return &discordNotifier{
		session: s,
	}
}

func (n *discordNotifier) SendResponse(i *discordgo.InteractionCreate, rd *discordgo.InteractionResponseData) error {
	err := n.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: rd,
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *discordNotifier) SendComplexMessage(userID string, data *discordgo.MessageSend) error {
	dmChannel, err := n.session.UserChannelCreate(userID)
	if err != nil {
		n.logger.Error("Failed to create a private channel with user %s: %v", userID, err)
	}
	_, err = n.session.ChannelMessageSendComplex(
		dmChannel.ID,
		data,
	)
	if err != nil {
		return err
	}
	return nil
}
