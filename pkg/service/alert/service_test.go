package alert

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/test-go/testify/assert"
)

func Test_SendNewPositionAlert(t *testing.T) {
	t.Skip("skipping test...")

	config.LoadEnv()

	t.Run("should send position alert", func(t *testing.T) {
		webhookID, err := strconv.ParseInt(os.Getenv("DISCORD_WEBHOOK_ID"), 10, 64)
		assert.NoError(t, err)
		accessToken := os.Getenv("DISCORD_ACCESS_TOKEN")
		client := webhook.New(snowflake.ID(webhookID), accessToken,
			webhook.WithLogger(logrus.New()),
			webhook.WithDefaultAllowedMentions(discord.AllowedMentions{
				RepliedUser: false,
			}),
		)
		newAlertService := serviceImpl{
			network: config.MainnetNetwork,
			enabled: true,
			client:  client,
		}
		ctx := context.Background()
		assert.NoError(t,
			newAlertService.SendNewPositionAlert(ctx, NewPositionAlert{
				TokenASymbol:  pointer.ToString("SAMO"),
				TokenAIconURL: pointer.ToString("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU/logo.png"),
				TokenAMint:    "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
				TokenBSymbol:  pointer.ToString("USDC"),
				TokenBIconURL: pointer.ToString("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v/logo.png	"),
				TokenBMint:                "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				ScaledTokenADepositAmount: 10000,
				ScaledDripAmount:          10,
				Granularity:               400,
				NumberOfSwaps:             45,
				Owner:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				Position:                  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				USDValue:                  pointer.ToFloat64(100),
			}))
	})

	t.Run("should return formatted granularity from getGranularityString", func(t *testing.T) {
		tests := []struct {
			input  uint64
			output string
		}{
			{
				input:  10,
				output: "Every Minute",
			},
			{
				input:  60,
				output: "Every Minute",
			},
			{
				input:  90,
				output: "Every Minute",
			},
			{
				input:  110,
				output: "Every Minute",
			},
			{
				input:  120,
				output: "Every 2 Minutes",
			},
			{
				input:  3600,
				output: "Every Hour",
			},
			{
				input:  4000,
				output: "Every Hour",
			},
			{
				input:  86400,
				output: "Every Day",
			},
			{
				input:  120400,
				output: "Every Day",
			},
			{
				input:  172800,
				output: "Every 2 Days",
			},
		}

		for _, testCase := range tests {
			assert.Equal(t, testCase.output, getGranularityString(testCase.input))
		}
	})
}
