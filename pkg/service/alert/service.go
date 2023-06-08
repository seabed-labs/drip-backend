package alert

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
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
	service := serviceImpl{
		network: appConfig.GetNetwork(),
	}
	if appConfig.GetDiscordWebhookID() != "" && appConfig.GetDiscordWebhookAccessToken() != "" {
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
		service.discordClient = &client
	}
	if appConfig.GetSlackWebhookURL() != "" {
		service.slackWebhookURL = pointer.ToString(appConfig.GetSlackWebhookURL())
	}
	if err := service.SendInfo(context.Background(), "Info", "Initialized alert service"); err != nil {
		return nil, err
	}
	return service, nil
}

type serviceImpl struct {
	network         config.Network
	discordClient   *webhook.Client
	slackWebhookURL *string
}

func (a serviceImpl) SendError(ctx context.Context, err error) (retErr error) {
	if sendErr := a.sendDiscord(ctx, discord.Embed{
		Title:       "Error",
		Description: err.Error(),
		Color:       int(InfoColor),
	}); sendErr != nil {
		retErr = multierror.Append(retErr, sendErr)
	}
	if sendErr := a.sendSlack(ctx, slack.NewTextBlockObject(slack.PlainTextType, err.Error(), false, false)); sendErr != nil {
		retErr = multierror.Append(retErr, sendErr)
	}
	return retErr
}

func (a serviceImpl) SendInfo(ctx context.Context, title string, message string) (retErr error) {
	if sendErr := a.sendDiscord(ctx, discord.Embed{
		Title:       title,
		Description: message,
		Color:       int(InfoColor),
	}); sendErr != nil {
		retErr = multierror.Append(retErr, sendErr)
	}
	if sendErr := a.sendSlack(ctx, slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, message, false, false))); sendErr != nil {
		retErr = multierror.Append(retErr, sendErr)
	}
	return retErr
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
) (err error) {
	log := logrus.WithField("msg", "new position").WithField("position", alertParams.Position)
	log.Info("attempting to send discord notification")
	if sendErr := a.sendNewDiscordPositionAlert(ctx, alertParams); sendErr != nil {
		err = multierror.Append(err, sendErr)
	}
	log.Info("attempting to send slack notification")
	if sendErr := a.sendNewSlackPositionAlert(ctx, alertParams); sendErr != nil {
		err = multierror.Append(err, sendErr)
	}
	return err
}

func (a serviceImpl) sendNewDiscordPositionAlert(
	ctx context.Context,
	alertParams NewPositionAlert,
) error {
	log := logrus.WithField("msg", "new position").WithField("position", alertParams.Position)
	log.Info("attempting to sendDiscord notification")

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

	return a.sendDiscord(ctx, embeds...)
}

func (a serviceImpl) sendNewSlackPositionAlert(
	ctx context.Context,
	alertParams NewPositionAlert,
) error {
	log := logrus.WithField("msg", "new position").WithField("position", alertParams.Position)
	log.Info("attempting to send slack notification")
	headerText := slack.NewTextBlockObject(slack.PlainTextType, "New Position! :money_mouth_face:", true, false)
	headerBlock := slack.NewHeaderBlock(headerText)

	granularityStr := getGranularityString(alertParams.Granularity)
	tokenA := utils.GetWithDefault(alertParams.TokenASymbol, alertParams.TokenAMint)
	tokenAImage := utils.GetWithDefault(alertParams.TokenAIconURL, "https://static.vecteezy.com/system/resources/previews/004/141/669/non_2x/no-photo-or-blank-image-icon-loading-images-or-missing-image-mark-image-not-available-or-image-coming-soon-sign-simple-nature-silhouette-in-frame-isolated-illustration-vector.jpg")
	tokenB := utils.GetWithDefault(alertParams.TokenBSymbol, alertParams.TokenBMint)
	tokenBImage := utils.GetWithDefault(alertParams.TokenBIconURL, "https://static.vecteezy.com/system/resources/previews/004/141/669/non_2x/no-photo-or-blank-image-icon-loading-images-or-missing-image-mark-image-not-available-or-image-coming-soon-sign-simple-nature-silhouette-in-frame-isolated-illustration-vector.jpg")
	usdValue := utils.GetWithDefault(alertParams.USDValue, 0)

	summaryBlock := slack.NewContextBlock("summary",
		slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("%f %s", alertParams.ScaledDripAmount, tokenA), false, false),
		slack.NewImageBlockElement(tokenAImage, "token a image"),
		slack.NewTextBlockObject(slack.MarkdownType, "to", false, false),
		slack.NewTextBlockObject(slack.MarkdownType, tokenB, false, false),
		slack.NewImageBlockElement(tokenBImage, "token b image"),
		slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf("%s, for a total of %d swaps, valued at *%f USD*", granularityStr, alertParams.NumberOfSwaps, usdValue), false, false),
	)
	return a.sendSlack(ctx, headerBlock, summaryBlock)
}

func getGranularityString(granularity uint64) (granularityStr string) {
	minuteInS := uint64(time.Minute.Seconds())
	hourInS := uint64(time.Hour.Seconds())
	dayInS := uint64((time.Hour * 24).Seconds())
	if uint64(granularity/minuteInS) <= 1 {
		granularityStr = "every Minute"
	} else if granularity > minuteInS && granularity < hourInS {
		granularityStr = fmt.Sprintf("every %d Minutes", uint64(granularity/minuteInS))
	} else if uint64(granularity/hourInS) <= 1 {
		granularityStr = "every Hour"
	} else if granularity > hourInS && granularity < dayInS {
		granularityStr = fmt.Sprintf("every %d Hours", uint64(granularity/hourInS))
	} else if uint64(granularity/dayInS) <= 1 {
		granularityStr = "every Day"
	} else {
		granularityStr = fmt.Sprintf("every %d Days", uint64(granularity/dayInS))
	}
	return granularityStr
}

func (a serviceImpl) sendDiscord(ctx context.Context, embeds ...discord.Embed) (err error) {
	if a.discordClient != nil {
		_, err = (*a.discordClient).CreateMessage(
			discord.NewWebhookMessageCreateBuilder().
				SetAvatarURL("https://pbs.twimg.com/profile_images/1512938686702403603/DDObiFjj_400x400.jpg").
				SetEmbeds(embeds...).
				Build(),
			// delay each request by 2 seconds
			rest.WithDelay(2*time.Second),
			rest.WithCtx(ctx),
		)
	} else {
		logrus.Info("alert service disabled, skipping info alert")
	}
	return err
}
func (a serviceImpl) sendSlack(ctx context.Context, blocks ...slack.Block) (err error) {
	if a.slackWebhookURL != nil {
		return slack.PostWebhookContext(ctx, *a.slackWebhookURL, &slack.WebhookMessage{
			//Text: "hello",
			Blocks: &slack.Blocks{BlockSet: blocks},
		})
	} else {
		logrus.Info("alert service disabled, skipping info alert")
	}
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
