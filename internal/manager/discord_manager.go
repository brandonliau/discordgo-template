package manager

import (
	"discord-template/internal/authenticator"
	"discord-template/internal/command"
	"discord-template/internal/component"
	"discord-template/internal/notifier"
	"discord-template/internal/repository"
	"discord-template/internal/shared"
	"discord-template/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type discordManager struct {
	session       *discordgo.Session
	repo          repository.Repository
	notifier      notifier.Notifier
	authenticator authenticator.Authenticator
	logger        logger.Logger
	commands      map[string]command.Command
	components    map[string]component.Component
}

func NewDiscordManager(s *discordgo.Session, repo repository.Repository, notifier notifier.Notifier, authenticator authenticator.Authenticator, logger logger.Logger) *discordManager {
	return &discordManager{
		session:       s,
		repo:          repo,
		notifier:      notifier,
		authenticator: authenticator,
		logger:        logger,
		commands:      make(map[string]command.Command),
		components:    make(map[string]component.Component),
	}
}

func (m *discordManager) RegisterCommand(c command.Command) {
	cname := c.Command().Name
	if _, ok := m.commands[cname]; ok {
		m.logger.Warn("Application command %s already registered", cname)
	}
	ccmd, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.Command())
	if err != nil {
		m.logger.Error("Failed to add application command %s : %v", cname, err)
	}
	m.repo.RegisterCommand(ccmd.Name, ccmd.ID)
	m.commands[cname] = c
	m.logger.Debug("Registered command %v", cname)
}

func (m *discordManager) RegisterComponent(c component.Component) {
	cname := c.CustomID()
	if _, ok := m.components[cname]; ok {
		m.logger.Warn("Application component %s already registered", cname)
	}
	m.repo.RegisterComponent(cname, c.Component())
	m.components[cname] = c
	m.logger.Debug("Registered component %v", cname)
}

func (m *discordManager) CommandInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	var ir *discordgo.InteractionResponse
	var err error
	command := m.commands[i.ApplicationCommandData().Name]
	if !m.authenticator.Authenticate(command, userID) {
		ir = shared.EphemeralContentResponse("Authorized users only!")
		m.logger.Info("Unauthorized user %s attempted to execute /%s", userID, i.ApplicationCommandData().Name)
	} else {
		ir, err = command.Execute(cmdArgs)
		if err != nil {
			m.logger.Error("Failed to execute /%s: %v", command.Command().Name, err)
		}
		m.logger.Debug("%s executed /%s", userID, i.ApplicationCommandData().Name)
	}

	err = m.notifier.SendResponse(i, ir)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}

func (m *discordManager) ComponentInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	var ir *discordgo.InteractionResponse
	var err error
	component := m.components[i.MessageComponentData().CustomID]
	ir, err = component.Execute(cmdArgs)
	if err != nil {
		m.logger.Error("Failed to execute /%s: %v", component.CustomID(), err)
		return
	}
	m.logger.Debug("%s executed /%s", userID, component.CustomID())

	err = m.notifier.SendResponse(i, ir)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}

func (m *discordManager) ModalInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}

	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	var ir *discordgo.InteractionResponse
	var err error
	component := m.components[i.ModalSubmitData().CustomID]
	ir, err = component.Execute(cmdArgs)
	if err != nil {
		m.logger.Error("Failed to execute /%s: %v", component.CustomID(), err)
		return
	}
	m.logger.Debug("%s executed /%s", userID, component.CustomID())

	err = m.notifier.SendResponse(i, ir)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}
