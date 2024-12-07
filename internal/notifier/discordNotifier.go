package notifier

import (
	"github.com/bwmarrin/discordgo"
)

type discordNotifier struct {
	session *discordgo.Session
}

func NewDiscordNotifier(s *discordgo.Session) *discordNotifier {
	return &discordNotifier{
		session: s,
	}
}

func (n *discordNotifier) SendResponse(i *discordgo.InteractionCreate, ir *discordgo.InteractionResponse) error {
	err := n.session.InteractionRespond(i.Interaction, ir)
	if err != nil {
		return err
	}
	return nil
}

func (n *discordNotifier) SendComplexMessage(userID string, data *discordgo.MessageSend) error {
	dmChannel, err := n.session.UserChannelCreate(userID)
	if err != nil {
		return err
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
