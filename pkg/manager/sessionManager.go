package manager

import (
	"log"

	"DiscordTemplate/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type sessionManager struct {
	session  *discordgo.Session
	commands map[string]command.Command
}

func NewSessionManager(s *discordgo.Session) *sessionManager {
	return &sessionManager{
		session:  s,
		commands: make(map[string]command.Command),
	}
}

func (m *sessionManager) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &command.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := m.commands[i.ApplicationCommandData().Name]; ok {
			rd, err := command.Execute(cmdArgs)
			if err != nil {
				log.Printf("[ERROR] Failed to execute %s: %v", command.GetCommand().Name, err)
			}
			err = m.SendResponse(i, rd)
			if err != nil {
				log.Printf("[ERROR] Failed to message user %s: %v", userID, err)
			}
		}
		log.Printf("[INFO] %s executed %s", cmdArgs.UserID, i.ApplicationCommandData().Name)
	}
}

func (m *sessionManager) RegisterCommand(c command.Command) {
	cname := c.GetCommand().Name
	if _, ok := m.commands[cname]; ok {
		log.Printf("[ERROR] Application command %s already registered", cname)
	}
	_, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.GetCommand())
	if err != nil {
		log.Printf("[ERROR] Failed to add application command %s : %v", cname, err)
	}
	m.commands[cname] = c
}

func (m *sessionManager) SendResponse(i *discordgo.InteractionCreate, rd *discordgo.InteractionResponseData) error {
	err := m.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: rd,
	})
	if err != nil {
		return err
	}
	return nil
}
