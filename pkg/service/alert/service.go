package alert

import (
	"context"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
)

type Service interface {
	SendError(ctx context.Context, err error) error
	SendInfo(ctx context.Context, title string, message string) error
	SendNewPositionAlert(ctx context.Context, alertParams NewPositionAlert) error
}

func NewAlertService(
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
	}
	//if err := service.SendInfo(context.Background(), "Info", "initialized alert service"); err != nil {
	//	return nil, err
	//}
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
	return a.send(ctx, discord.Embed{
		Title:       "Error",
		Description: err.Error(),
		Color:       int(InfoColor),
	})
}

func (a serviceImpl) SendInfo(ctx context.Context, title string, message string) error {
	if !a.enabled {
		logrus.WithField("msg", message).Info("alert service disabled, skipping info alert")
		return nil
	}
	return a.send(ctx, discord.Embed{
		Title:       title,
		Description: message,
		Color:       int(InfoColor),
	})
}

type NewPositionAlert struct {
	TokenASymbol              *string
	TokenAIconURL             *string
	TokenAMint                string
	TokenBSymbol              *string
	TokenBIconURL             *string
	TokenBMint                string
	TokenAScaledDepositAmount float64
	Granularity               uint64
	ScaledDripAmount          float64
	NumberOfSwaps             uint64
	Owner                     string
}

func (a serviceImpl) SendNewPositionAlert(
	ctx context.Context,
	alertParams NewPositionAlert,
) error {
	granularityStr := strconv.FormatUint(alertParams.Granularity, 10)
	if alertParams.Granularity == 60 {
		granularityStr = "Minutely"
	} else if alertParams.Granularity == 3600 {
		granularityStr = "Hourly"
	}

	tokenA := alertParams.TokenAMint
	if alertParams.TokenASymbol != nil {
		tokenA = *alertParams.TokenASymbol
	}

	tokenB := alertParams.TokenBMint
	if alertParams.TokenBSymbol != nil {
		tokenB = *alertParams.TokenBSymbol
	}
	inLineTrue := utils.GetBoolPtr(true)
	embed := discord.NewEmbedBuilder().
		SetTitle("New Position!").
		SetColor(int(SuccessColor)).
		SetFields(
			discord.EmbedField{Name: "Token A", Value: tokenA, Inline: inLineTrue},
			discord.EmbedField{Name: "Token B", Value: tokenB, Inline: inLineTrue},
			discord.EmbedField{Name: "Token A Deposit", Value: strconv.FormatFloat(alertParams.ScaledDripAmount, 'f', -1, 32), Inline: inLineTrue},
			discord.EmbedField{Name: "Granularity", Value: granularityStr, Inline: inLineTrue},
			discord.EmbedField{Name: "Drip Amount", Value: strconv.FormatFloat(alertParams.ScaledDripAmount, 'f', -1, 32), Inline: inLineTrue},
			discord.EmbedField{Name: "Number of swaps", Value: strconv.FormatUint(alertParams.NumberOfSwaps, 10), Inline: inLineTrue},
			discord.EmbedField{Name: "Owner", Value: alertParams.Owner},
		).
		Build()
	embeds := []discord.Embed{embed}
	if alertParams.TokenAIconURL != nil && alertParams.TokenASymbol != nil {
		tokenAEmbed := discord.NewEmbedBuilder().
			SetTitle("TokenA").
			SetColor(int(SuccessColor)).
			SetFields(
				discord.EmbedField{Name: "Symbol", Value: *alertParams.TokenASymbol},
			).
			SetEmbedFooter(&discord.EmbedFooter{
				Text:         alertParams.TokenAMint,
				IconURL:      *alertParams.TokenAIconURL,
				ProxyIconURL: "",
			}).
			Build()

		embeds = append(embeds, tokenAEmbed)
	}
	if alertParams.TokenBIconURL != nil && alertParams.TokenBSymbol != nil {
		tokenBEmbed := discord.NewEmbedBuilder().
			SetTitle("TokenB").
			SetColor(int(SuccessColor)).
			SetFields(
				discord.EmbedField{Name: "Symbol", Value: *alertParams.TokenBSymbol},
			).
			SetEmbedFooter(&discord.EmbedFooter{
				Text:         alertParams.TokenBMint,
				IconURL:      *alertParams.TokenBIconURL,
				ProxyIconURL: "",
			}).
			Build()
		embeds = append(embeds, tokenBEmbed)
	}
	return a.send(ctx, embeds...)
}

func (a serviceImpl) send(ctx context.Context, embeds ...discord.Embed) error {
	_, err := a.client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
			SetEmbeds(embeds...).
			Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
		//rest.WithCtx(ctx),
	)
	return err
}
