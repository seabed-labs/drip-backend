package integrationtest

import (
	"context"
	"reflect"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/test-go/testify/assert"
)

func TestHandler_UpsertProtoConfigByAddress(t *testing.T) {
	t.Run("should upsert vault proto config", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.APIRecorderOptions{
				Path: "./fixtures/upsert-protoconfig-by-address",
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				protoConfigAddress := "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr"
				// protoConfig: https://explorer.solana.com/address/Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr?cluster=devnet
				assert.NoError(t, processor.UpsertProtoConfigByAddress(context.Background(), protoConfigAddress))

				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), protoConfigAddress)
				assert.NoError(t, err)
				assert.NotNil(t, protoConfig)

				assert.Equal(t, "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr", protoConfig.Pubkey)
				assert.Equal(t, "3CTkqdcjzn1ptnNYUcFq4Sk1Smy91kz5p9JhgJBHGe3e", protoConfig.Admin)
				assert.Equal(t, uint64(60), protoConfig.Granularity)
				assert.Equal(t, uint16(50), protoConfig.TokenADripTriggerSpread)
				assert.Equal(t, uint16(25), protoConfig.TokenBWithdrawalSpread)
				assert.Equal(t, uint16(25), protoConfig.TokenBReferralSpread)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*protoConfig).NumField(), 6)
			})
	})
}
