package alert

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/rest"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/disgoorg/disgo/discord"
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
	appConfig config.AppConfig,
) (Service, error) {
	logrus.WithField("discordWebhookID", appConfig.GetDiscordWebhookID()).Info("initiating alert service")
	service := serviceImpl{}
	if appConfig.GetDiscordWebhookID() != "" && appConfig.GetDiscordWebhookAccessToken() != "" {
		service = serviceImpl{
			network: appConfig.GetNetwork(),
			enabled: true,
		}
		webhookID, err := strconv.ParseInt(appConfig.GetDiscordWebhookID(), 10, 64)
		if err != nil {
			return nil, err
		}
		client := webhook.New(snowflake.ID(webhookID), appConfig.GetDiscordWebhookAccessToken(),
			webhook.WithLogger(logrus.New()),
			webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
				RepliedUser: false,
			}),
		)
		service.client = client
	}
	if err := service.SendInfo(context.Background(), "Info", "Initialized alert service"); err != nil {
		return nil, err
	}
	return service, nil
}

type serviceImpl struct {
	network config.Network
	enabled bool
	client  webhook.Client
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
	ScaledTokenADepositAmount float64
	ScaledDripAmount          float64
	Granularity               uint64
	NumberOfSwaps             uint64
	Owner                     string
	Position                  string
	USDValue                  *float64
}

func (a serviceImpl) SendNewPositionAlert(
	ctx context.Context,
	alertParams NewPositionAlert,
) error {
	log := logrus.WithField("msg", "new position").WithField("position", alertParams.Position)
	if !a.enabled {
		log.Info("alert service disabled, skipping info alert")
		return nil
	}
	log.Info("attempting to send notification")

	granularityStr := getGranularityString(alertParams.Granularity)
	tokenA := utils.GetWithDefault(alertParams.TokenASymbol, alertParams.TokenAMint)
	tokenB := utils.GetWithDefault(alertParams.TokenBSymbol, alertParams.TokenBMint)
	usdValue := utils.GetWithDefault(alertParams.USDValue, 0)

	inLineTrue := utils.GetBoolPtr(true)
	embed := discord.NewEmbedBuilder().
		SetTitle("New Position!").
		SetColor(int(SuccessColor)).
		SetURL(a.getExplorerURL(alertParams.Position)).
		SetFields(
			discord.EmbedField{Name: "Token A", Value: tokenA, Inline: inLineTrue},
			discord.EmbedField{Name: "Token B", Value: tokenB, Inline: inLineTrue},
			discord.EmbedField{Name: "Token A Deposit", Value: strconv.FormatFloat(alertParams.ScaledTokenADepositAmount, 'f', -1, 32), Inline: inLineTrue},
			discord.EmbedField{Name: "Granularity", Value: granularityStr, Inline: inLineTrue},
			discord.EmbedField{Name: "Drip Amount", Value: strconv.FormatFloat(alertParams.ScaledDripAmount, 'f', -1, 32), Inline: inLineTrue},
			discord.EmbedField{Name: "Number of swaps", Value: strconv.FormatUint(alertParams.NumberOfSwaps, 10), Inline: inLineTrue},
			discord.EmbedField{Name: "USD Value", Value: strconv.FormatFloat(usdValue, 'f', -1, 32)},
			discord.EmbedField{Name: "Owner", Value: alertParams.Owner},
		).
		Build()
	embeds := []discord.Embed{embed}

	tokenAEmbed := discord.NewEmbedBuilder().
		SetURL(a.getExplorerURL(alertParams.Position)).
		SetImage(utils.GetWithDefault(alertParams.TokenAIconURL, "https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg")).
		Build()
	embeds = append(embeds, tokenAEmbed)

	tokenBEmbed := discord.NewEmbedBuilder().
		SetTitle("TokenB").
		SetColor(int(SuccessColor)).
		SetURL(a.getExplorerURL(alertParams.Position)).
		SetImage(utils.GetWithDefault(alertParams.TokenBIconURL, "https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg")).
		Build()
	embeds = append(embeds, tokenBEmbed)

	return a.send(ctx, embeds...)
}
func getGranularityString(granularity uint64) (granularityStr string) {
	minuteInS := uint64(time.Minute.Seconds())
	hourInS := uint64(time.Hour.Seconds())
	dayInS := uint64((time.Hour * 24).Seconds())
	if uint64(granularity/minuteInS) <= 1 {
		granularityStr = "Every Minute"
	} else if granularity > minuteInS && granularity < hourInS {
		granularityStr = fmt.Sprintf("Every %d Minutes", uint64(granularity/minuteInS))
	} else if uint64(granularity/hourInS) <= 1 {
		granularityStr = "Every Hour"
	} else if granularity > hourInS && granularity < dayInS {
		granularityStr = fmt.Sprintf("Every %d Hours", uint64(granularity/hourInS))
	} else if uint64(granularity/dayInS) <= 1 {
		granularityStr = "Every Day"
	} else {
		granularityStr = fmt.Sprintf("Every %d Days", uint64(granularity/dayInS))
	}
	return granularityStr
}

func (a serviceImpl) send(ctx context.Context, embeds ...discord.Embed) error {
	_, err := a.client.CreateMessage(
		discord.NewWebhookMessageCreateBuilder().
			SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
			SetEmbeds(embeds...).
			Build(),
		// delay each request by 2 seconds
		rest.WithDelay(2*time.Second),
		rest.WithCtx(ctx),
	)
	return err
}

func (a serviceImpl) getExplorerURL(account string) string {
	switch a.network {
	case config.MainnetNetwork:
		return fmt.Sprintf("https://explorer.solana.com/address/%s", account)
	case config.DevnetNetwork:
		return fmt.Sprintf("https://explorer.solana.com/address/%s?cluster=devnet", account)
	default:
		return fmt.Sprintf(
			"https://explorer.solana.com/address/%s?cluster=%s",
			account,
			url.QueryEscape("custom&customUrl=http://localhost:8899"))
	}
}
