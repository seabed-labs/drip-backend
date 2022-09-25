package alert

import (
	"context"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

type Service interface {
	SendError(ctx context.Context, err error) error
	SendInfo(ctx context.Context, title string, message string, alertType AlertColor) error
	SendAlertWithFields(ctx context.Context, title string, fields []discord.EmbedField, alertType AlertColor) error
}

func NewService(
	config *configs.AppConfig,
) (Service, error) {
	logrus.WithField("discordWebhookID", config.DiscordWebhookID).Info("initiating alert service")
	service := serviceImpl{}
	if config.DiscordWebhookID != "" && config.DiscordWebhookAccessToken != "" {
		service = serviceImpl{
			enabled:                   true,
			discordWebhookID:          config.DiscordWebhookID,
			discordWebhookAccessToken: config.DiscordWebhookAccessToken,
		}
		webhookID, err := strconv.ParseInt(service.discordWebhookID, 10, 64)
		if err != nil {
			return nil, err
		}
		client := webhook.New(snowflake.ID(webhookID), service.discordWebhookAccessToken,
			webhook.WithLogger(logrus.New()),
			webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
				RepliedUser: false,
			}),
		)
		service.client = client
		if err := service.SendInfo(context.Background(), "Info", "initialized alert service", Info); err != nil {
			return nil, err
		}
	}
	return service, nil
}

type serviceImpl struct {
	enabled                   bool
	client                    webhook.Client
	discordWebhookAccessToken string
	discordWebhookID          string
}

func (a serviceImpl) SendError(ctx context.Context, err error) error {
	if !a.enabled {
		logrus.WithError(err).Info("alert service disabled, skipping error alert")
		return nil
	}
	if _, err := a.client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
			SetEmbeds(
				discord.Embed{
					Title:       "Error",
					Description: err.Error(),
					Color:       15158332,
				},
			).
			Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
	); err != nil {
		return err
	}
	return nil
}

type AlertColor int

const (
	Success AlertColor = 1752220
	Error              = 15158332
	Warn               = 15844367
	Info               = 0
)

func (a serviceImpl) SendInfo(ctx context.Context, title string, message string, alertType AlertColor) error {
	if !a.enabled {
		logrus.WithField("msg", message).Info("alert service disabled, skipping info alert")
		return nil
	}
	if _, err := a.client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
			SetEmbeds(
				discord.Embed{
					Title:       title,
					Description: message,
					Color:       int(alertType),
				},
			).
			Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
	); err != nil {
		return err
	}
	return nil
}

func (a serviceImpl) SendAlertWithFields(ctx context.Context, title string, embedFields []discord.EmbedField, alertType AlertColor) error {
	if !a.enabled {
		logrus.WithField("msg", title).Info("alert service disabled, skipping info alert")
		return nil
	}
	if _, err := a.client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
			SetEmbeds(
				discord.Embed{
					Title:  title,
					Color:  int(alertType),
					Fields: embedFields,
				},
			).
			Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
	); err != nil {
		return err
	}
	return nil
}
