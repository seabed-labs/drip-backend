package integrationtest

import (
	"context"
	"testing"
	"time"

	repository2 "github.com/dcaf-labs/drip/pkg/service/repository/analytics"
	repository3 "github.com/dcaf-labs/drip/pkg/service/repository/queue"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/integrationtest"
	solanaClient "github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/assert"
)

func Test_ProcessTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)

	t.Run("should upsert v1 deposit (depositWithMetadata) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test4",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "5otEnwWY8z5MJhrxaXRgGeN9jyM3nAcC8gGJQY1asE3wfGB3QbjPQ1DUFXzEtPuLcdWpdBuw3XZuL177jRiwMBiv"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetDepositMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.DepositMetric{
					Signature:           sig,
					IxIndex:             int32(4),
					IxName:              "DepositWithMetadata",
					IxVersion:           1,
					Slot:                int32(185407760),
					Time:                time.Unix(1680117259, 0),
					Vault:               "BmJs2b1PnepHwBaiWqzuX8LywBLkbBymqj7Cpjz5WjuY",
					Referrer:            pointer.ToString("9d1wBhpKd24y1XwqawL3WWUjw14H7JqB9jcqkrhB3eHW"),
					TokenADepositAmount: uint64(200000000000),
					TokenAUsdPriceDay:   nil,
				})
			})
	})

	t.Run("should upsert v1 drip (dripOrcaWhirlpool) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test5",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "5uECxzjML1a5sXuMPqwcpP8BMAKCXCpUVwFydPHWBnaW9V72ockWQKevfSQBaqiSxtGKtmstkVaxmo5Jcfnp9Lcb"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetDripMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.DripMetric{
					Signature:                  sig,
					IxIndex:                    1,
					IxName:                     "DripOrcaWhirlpool",
					IxVersion:                  int32(1),
					Slot:                       int32(175080460),
					Time:                       time.Unix(1675037066, 0),
					Vault:                      "6PnzoW2Bcs6WGqYvecfSxN9C2EeDmQCjUCeFA7JDDfZG",
					VaultTokenASwappedAmount:   uint64(1035325),
					VaultTokenBReceivedAmount:  uint64(39648400),
					KeeperTokenAReceivedAmount: uint64(3636),
					TokenAUsdPriceDay:          nil,
					TokenBUsdPriceDay:          nil,
				})
			})
	})

	t.Run("should upsert v1 withdraw (withdrawB) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test6",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "3QCyiPWTdbeEi9VKMzRrYL2Rw8vdAHXcqFjKcbYx11yv55mttRmLwPFNepZEZ9Qd1UHjfnUCSmffnsdSTL3axXmh"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetWithdrawalMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.WithdrawalMetric{
					Signature:                    sig,
					IxIndex:                      int32(0),
					IxName:                       "WithdrawB",
					IxVersion:                    int32(1),
					Slot:                         int32(185936032),
					Time:                         time.Unix(1680358672, 0),
					Vault:                        "BmJs2b1PnepHwBaiWqzuX8LywBLkbBymqj7Cpjz5WjuY",
					UserTokenAWithdrawAmount:     uint64(0),
					UserTokenBWithdrawAmount:     uint64(23055242594808),
					TreasuryTokenBReceivedAmount: uint64(23089877410),
					ReferralTokenBReceivedAmount: uint64(11544938705),
					TokenAUsdPriceDay:            nil,
					TokenBUsdPriceDay:            nil,
				})
			})
	})

	t.Run("should upsert v1 withdraw (closePosition) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test7",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "shVfUMeDU38tQxdWvrYwFrBCjDve5889JdEAhgmJB3bXejhgx8VpRwtXaYhG4FhZXKHrT7pz2SZk3UvpkaHmfE6"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetWithdrawalMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.WithdrawalMetric{
					Signature:                    sig,
					IxIndex:                      int32(2),
					IxName:                       "ClosePosition",
					IxVersion:                    int32(1),
					Slot:                         int32(188060646),
					Time:                         time.Unix(1681330759, 0),
					Vault:                        "BmJs2b1PnepHwBaiWqzuX8LywBLkbBymqj7Cpjz5WjuY",
					UserTokenAWithdrawAmount:     uint64(110040160587),
					UserTokenBWithdrawAmount:     uint64(80279369044978),
					TreasuryTokenBReceivedAmount: uint64(80399968998),
					ReferralTokenBReceivedAmount: uint64(40199984499),
					TokenAUsdPriceDay:            nil,
					TokenBUsdPriceDay:            nil,
				})
			})
	})

	t.Run("should upsert v0 deposit (depositWithMetadata) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test8",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "2tf1Mmg2bir7BHd3zbxt1R3C44xWezxqmho9XHVQxPPPvt62WMksFpm4b6WjzNYJLaSwjLs7vcPNFMwE71z2TQ1B"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetDepositMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.DepositMetric{
					Signature:           sig,
					IxIndex:             int32(1),
					IxName:              "DepositWithMetadata",
					IxVersion:           0,
					Slot:                int32(146790548),
					Time:                time.Unix(1660975474, 0),
					Vault:               "2dqTZ6Q3UDQTez6HviDjXFg9BNNvCpE9mEcV92peRacj",
					Referrer:            nil,
					TokenADepositAmount: uint64(100000000),
					TokenAUsdPriceDay:   nil,
				})
			})
	})

	t.Run("should upsert v0 drip (dripOrcaWhirlpool) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test9",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "3ySFiaC3VPnPUL8ywCFgMKGG7XMWCLJaz3WFdKGukkBEMcouAv9jw3Xi7eNssLYsYxg4rW3indsVyjupjUWSKYQr"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetDripMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.DripMetric{
					Signature:                  sig,
					IxIndex:                    1,
					IxName:                     "DripOrcaWhirlpool",
					IxVersion:                  int32(0),
					Slot:                       int32(146908472),
					Time:                       time.Unix(1661040016, 0),
					Vault:                      "2dqTZ6Q3UDQTez6HviDjXFg9BNNvCpE9mEcV92peRacj",
					VaultTokenASwappedAmount:   uint64(25740000),
					VaultTokenBReceivedAmount:  uint64(729375291),
					KeeperTokenAReceivedAmount: uint64(260000),
					TokenAUsdPriceDay:          nil,
					TokenBUsdPriceDay:          nil,
				})
			})
	})

	t.Run("should upsert v0 withdraw (withdrawB) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test10",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "4TjArAvd9ESBubQV3Z1bDXahX2uu1J8a4wYTFhomB1i8Xxskzer8KpsXaJJ238edDTyKp1z8LabtxTfHeiWE2vyH"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetWithdrawalMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.WithdrawalMetric{
					Signature:                    sig,
					IxIndex:                      int32(1),
					IxName:                       "WithdrawB",
					IxVersion:                    int32(0),
					Slot:                         int32(146634888),
					Time:                         time.Unix(1660884165, 0),
					Vault:                        "BrkNC3vpj17h8hwDoqyYrCvEF1BqV6wNQLxP4DfhBiLb",
					UserTokenAWithdrawAmount:     uint64(0),
					UserTokenBWithdrawAmount:     uint64(6748651),
					TreasuryTokenBReceivedAmount: uint64(68168),
					ReferralTokenBReceivedAmount: uint64(0),
					TokenAUsdPriceDay:            nil,
					TokenBUsdPriceDay:            nil,
				})
			})
	})

	t.Run("should upsert v0 withdraw (closePosition) metric", func(t *testing.T) {
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test11",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				client solanaClient.Solana,
				repo repository2.AnalyticsRepository,
				txQueue repository3.TransactionUpdateQueue,
			) {
				sig := "5AA9RecspxtdsXVPPcj43ZnkuGQCwFBF3iDz76hKPfNvFXy5sJPGtFrFhco5f2oo4bb5WAVnutuk8JcnNB3hdNHG"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				assert.NoError(t, processor.ProcessTransaction(ctx, *tx))
				metric, err := repo.GetWithdrawalMetricBySignature(ctx, sig)
				assert.NoError(t, err)
				assert.NotNil(t, metric)
				assert.Equal(t, *metric, model.WithdrawalMetric{
					Signature:                    sig,
					IxIndex:                      int32(2),
					IxName:                       "ClosePosition",
					IxVersion:                    int32(0),
					Slot:                         int32(146791839),
					Time:                         time.Unix(1660976178, 0),
					Vault:                        "2dqTZ6Q3UDQTez6HviDjXFg9BNNvCpE9mEcV92peRacj",
					UserTokenAWithdrawAmount:     uint64(2000000),
					UserTokenBWithdrawAmount:     uint64(0),
					TreasuryTokenBReceivedAmount: uint64(0),
					ReferralTokenBReceivedAmount: uint64(0),
					TokenAUsdPriceDay:            nil,
					TokenBUsdPriceDay:            nil,
				})
			})
	})
}
