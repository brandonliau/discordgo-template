package command

import (
	"fmt"
	"time"

	"discordgo-template/internal/application/usecase"
	"discordgo-template/internal/interfaces/discord/interaction"
	"discordgo-template/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type statusCommand struct {
	systemService *usecase.SystemService
}

func StatusCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "status",
		Description: "View the status of the bot.",
	}
}

func StatusCommandHandler(systemService *usecase.SystemService) interaction.HandleFunc {
	c := &statusCommand{
		systemService: systemService,
	}
	return c.execute
}

func (c *statusCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	req := usecase.GetStatusRequest{}

	res, err := c.systemService.GetStatus(req)
	if err != nil {
		return nil, err
	}

	rsp := interaction.InteractionInitialResponse(
		interaction.WithEmbeds(c.statusEmbed(s, res)),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}

func (c *statusCommand) statusEmbed(session *discordgo.Session, res *usecase.GetStatusResult) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Status",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: ":stopwatch: Latency",
				Value: fmt.Sprintf(
					"> **Discord gateway:** `%d ms`\n",
					session.HeartbeatLatency().Milliseconds(),
				),
				Inline: false,
			},
			{
				Name: ":alarm_clock: Uptime",
				Value: fmt.Sprintf(
					"> **Last restart:** <t:%d:R>\n"+
						"> **Uptime:** %d days, %d hours, %d min, %d sec",
					res.StartTime,
					res.UptimeDays,
					res.UptimeHours,
					res.UptimeMinutes,
					res.UptimeSeconds,
				),
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
