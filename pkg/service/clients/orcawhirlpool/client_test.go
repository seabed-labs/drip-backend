package orcawhirlpool

import (
	"context"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/test-go/testify/assert"
)

func TestOrcaWhirlpoolClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(unittest.GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.StagingEnv).AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()

	client := NewOrcaWhirlpoolClient(mockConfig)

	t.Run("GetOrcaWhirlpoolQuoteEstimate should swap estimate", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		res, err := client.GetOrcaWhirlpoolQuoteEstimate(ctx,
			"GSFnjnJ7TdSsGWb6JgFhWakWrv8VGZUAghnY3EA8Xj46",
			"7ihthG4cFydyDnuA3zmJrX13ePGpLcANf3tHLmKLPN7M",
			"100000")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.AToB)
	})
}
