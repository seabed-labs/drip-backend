package solana

import (
	"context"
	"fmt"
	"testing"
	"time"

	api "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"

	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/test-go/testify/assert"
)

func TestSolanaClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	//mint := "31nFDfb3b4qw8JPx4FaXGgEk8omt7NuHpPkwWCSym5rC"
	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockDevnetStagingConfig(ctrl)

	client, err := NewSolanaClient(mockConfig, api.GetDefaultClientProvider())
	assert.NoError(t, err)

	t.Run("GetWalletPubKey should return public key", func(t *testing.T) {
		assert.Equal(t, client.GetWalletPubKey().String(), "J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")
	})

	t.Run("getWalletPrivKey should return private key", func(t *testing.T) {
		assert.Equal(t, client.getWalletPrivKey().String(),
			"2v28DUBnjz9eGHbgMH4fVzixzpyP8SfdBVmo19vdhgzDddqnD4HMNiFgNKtQsKErEfhnRYKFY9k4WbaGyyFKQzai")
	})

	t.Run("should getWalletPrivKey", func(t *testing.T) {
		versionResponse, err := client.GetVersion(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, versionResponse.FeatureSet)
	})

	//t.Run("ProgramSubscribe should subscribe to consumer", func(t *testing.T) {
	//	timeout := time.Second * 5
	//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//	defer cancel()
	//	err := client.ProgramSubscribe(ctx, solana.TokenProgramID.String(), func(address string, data []byte) error {
	//		assert.NotEmpty(t, address)
	//		assert.NotEmpty(t, data)
	//		return nil
	//	})
	//	assert.NoError(t, err)
	//	time.Sleep(timeout)
	//})

	t.Run("GetProgramAccounts should return all account publickeys", func(t *testing.T) {
		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		offset, err := client.GetProgramAccounts(ctx, tokenswap.ProgramID.String())
		fmt.Println(len(offset))
		assert.NoError(t, err)
		assert.NotZero(t, offset)
		assert.True(t, len(offset) > 4500)
	})

	t.Run("GetTokenMetadataAccount should return token metadata for mint with valid token metadata account", func(t *testing.T) {
		tokenMetadata, err := client.GetTokenMetadataAccount(context.Background(), "2CsU92EN1AwEcJdWaVDa9n9o6AFB5nVBozEXZWDJcDj6")
		assert.NoError(t, err)
		assert.Equal(t, tokenMetadata.Data.Symbol, "DP")
		assert.Equal(t, tokenMetadata.Data.Name, "Drip Position")
	})

	t.Run("GetTokenMetadataAccount should return err for mint without token metadata account", func(t *testing.T) {
		_, err := client.GetTokenMetadataAccount(context.Background(), "2Kp1LB7Jo5RAt5gnisbtpu9tvK24WtV9zo9tZHKqFGqL")
		assert.Error(t, err)
	})

	t.Run("GetTokenMetadataAccount should return err for non-mint address", func(t *testing.T) {
		_, err := client.GetTokenMetadataAccount(context.Background(), "F1E3YejjuQ87d73enFG1w2toeeR4igZq8eQum1UgqGeF")
		assert.Error(t, err)
	})
}
