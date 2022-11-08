package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/base"
	solanaClient "github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/unittest"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_PostMint(t *testing.T) {
	mint := "31nFDfb3b4qw8JPx4FaXGgEk8omt7NuHpPkwWCSym5rC"
	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockDevnetStagingConfig(ctrl)
	e := echo.New()

	t.Run("should return an error when providing invalid amount", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "xyz",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		assert.NoError(t, h.PostMint(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, "invalid amount, must be float64", res.JSON500.Error)
	})

	t.Run("should return an error when failing to get account info", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("some error")).
			Times(1)

		assert.NoError(t, h.PostMint(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, "failed to get account info", res.JSON500.Error)
	})

	t.Run("should return an error when failing to decode borsh", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		mintDataBytes, err := rpc.DataBytesOrJSONFromBase64(base64.StdEncoding.EncodeToString([]byte{1, 0, 0, 0, 5, 234, 156}))
		assert.NoError(t, err)
		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(&rpc.GetAccountInfoResult{
				Value: &rpc.Account{
					Lamports:   0,
					Owner:      solana.PublicKey{},
					Data:       mintDataBytes,
					Executable: false,
					RentEpoch:  0,
				},
			}, nil).
			Times(1)

		assert.NoError(t, h.PostMint(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, "failed to decode mint", res.JSON500.Error)
	})

	t.Run("should return an error when api wallet is not mint authority", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		mintAuth := solana.NewWallet().PublicKey()
		mintData, err := encodeMintToBase64(token.Mint{
			MintAuthority:   &mintAuth,
			Supply:          1,
			Decimals:        6,
			IsInitialized:   true,
			FreezeAuthority: &mintAuth,
		})
		assert.NoError(t, err)
		mintDataBytes, err := rpc.DataBytesOrJSONFromBase64(mintData)
		assert.NoError(t, err)

		m.
			EXPECT().
			GetWalletPubKey().
			Return(solana.MustPublicKeyFromBase58("J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")).
			Times(2)
		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(&rpc.GetAccountInfoResult{
				Value: &rpc.Account{
					Lamports:   0,
					Owner:      solana.PublicKey{},
					Data:       mintDataBytes,
					Executable: false,
					RentEpoch:  0,
				},
			}, nil).
			Times(1)

		assert.NoError(t, h.PostMint(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, "invalid mint, J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer is not MintAuthority", res.JSON500.Error)
	})

	t.Run("should return an error when failing to mint", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		mintAuth := solana.MustPublicKeyFromBase58("J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")
		mintData, err := encodeMintToBase64(token.Mint{
			MintAuthority:   &mintAuth,
			Supply:          1,
			Decimals:        6,
			IsInitialized:   true,
			FreezeAuthority: &mintAuth,
		})
		assert.NoError(t, err)
		mintDataBytes, err := rpc.DataBytesOrJSONFromBase64(mintData)
		assert.NoError(t, err)

		m.
			EXPECT().
			GetWalletPubKey().
			Return(solana.MustPublicKeyFromBase58("J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")).
			Times(1)
		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(&rpc.GetAccountInfoResult{
				Value: &rpc.Account{
					Lamports:   0,
					Owner:      solana.PublicKey{},
					Data:       mintDataBytes,
					Executable: false,
					RentEpoch:  0,
				},
			}, nil).
			Times(1)
		m.
			EXPECT().
			MintToWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("some error")).
			Times(1)

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON500)
		assert.Equal(t, "failed to mint", res.JSON500.Error)
	})

	t.Run("should return success", func(t *testing.T) {
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		c, rec := unittest.GetTestRequestRecorder(e, strings.NewReader(string(reqBody)))

		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))

		mintAuth := solana.MustPublicKeyFromBase58("J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")
		mintData, err := encodeMintToBase64(token.Mint{
			MintAuthority:   &mintAuth,
			Supply:          1,
			Decimals:        6,
			IsInitialized:   true,
			FreezeAuthority: &mintAuth,
		})
		assert.NoError(t, err)
		mintDataBytes, err := rpc.DataBytesOrJSONFromBase64(mintData)
		assert.NoError(t, err)

		m.
			EXPECT().
			GetWalletPubKey().
			Return(solana.MustPublicKeyFromBase58("J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer")).
			Times(1)
		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(&rpc.GetAccountInfoResult{
				Value: &rpc.Account{
					Lamports:   0,
					Owner:      solana.PublicKey{},
					Data:       mintDataBytes,
					Executable: false,
					RentEpoch:  0,
				},
			}, nil).
			Times(1)
		m.
			EXPECT().
			MintToWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return("some tx hash", nil).
			Times(1)

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		res, err := apispec.ParsePostMintResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, res.JSON200)
		assert.Equal(t, "some tx hash", res.JSON200.TxHash)
	})
}

func encodeMintToBase64(mint token.Mint) (string, error) {
	buf := new(bytes.Buffer)
	err := bin.NewBorshEncoder(buf).Encode(mint)
	b64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return b64Data, err
}
