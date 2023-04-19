package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	solanaClient "github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/assert"
)

func Test_ProcessTransactionUpdateQueueItem(t *testing.T) {
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
				repo repository.Repository,
				client solanaClient.Solana,
				txQueue repository.TransactionUpdateQueue,
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
				repo repository.Repository,
				client solanaClient.Solana,
				txQueue repository.TransactionUpdateQueue,
			) {
				sig := "5uECxzjML1a5sXuMPqwcpP8BMAKCXCpUVwFydPHWBnaW9V72ockWQKevfSQBaqiSxtGKtmstkVaxmo5Jcfnp9Lcb"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				tx, err := client.GetTransaction(ctx, sig)
				assert.NoError(t, err)
				//assert.NotNil(t, tx)
				//assert.NotNil(t, tx.Transaction)
				//bytes, err := json.Marshal(tx)
				//assert.NoError(t, err)
				//var unamrshalTx rpc.GetTransactionResult
				//assert.NoError(t, json.Unmarshal(bytes, &unamrshalTx))
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
				repo repository.Repository,
				client solanaClient.Solana,
				txQueue repository.TransactionUpdateQueue,
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
					Signature:                sig,
					IxIndex:                  int32(0),
					IxName:                   "WithdrawB",
					IxVersion:                int32(1),
					Slot:                     int32(185936032),
					Time:                     time.Unix(1680358672, 0),
					Vault:                    "BmJs2b1PnepHwBaiWqzuX8LywBLkbBymqj7Cpjz5WjuY",
					UserTokenAWithdrawAmount: uint64(0),
					UserTokenBWithdrawAmount: uint64(23055242594808),
					// TODO: This should not be 0
					TreasuryTokenBReceivedAmount: uint64(0),
					ReferralTokenBReceivedAmount: uint64(0),
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
				repo repository.Repository,
				client solanaClient.Solana,
				txQueue repository.TransactionUpdateQueue,
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
					Signature:                sig,
					IxIndex:                  int32(2),
					IxName:                   "ClosePosition",
					IxVersion:                int32(1),
					Slot:                     int32(188060646),
					Time:                     time.Unix(1681330759, 0),
					Vault:                    "BmJs2b1PnepHwBaiWqzuX8LywBLkbBymqj7Cpjz5WjuY",
					UserTokenAWithdrawAmount: uint64(110040160587),
					UserTokenBWithdrawAmount: uint64(80279369044978),
					// TODO: This should not be 0
					TreasuryTokenBReceivedAmount: uint64(0),
					ReferralTokenBReceivedAmount: uint64(0),
					TokenAUsdPriceDay:            nil,
					TokenBUsdPriceDay:            nil,
				})
			})
	})
}
