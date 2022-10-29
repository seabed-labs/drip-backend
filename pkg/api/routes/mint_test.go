package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dcaf-labs/drip/pkg/unittest"

	"github.com/dcaf-labs/drip/pkg/service/base"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/config"

	solanaClient "github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/repository"

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
	mockConfig := config.NewMockAppConfig(ctrl)
	mockConfig.EXPECT().GetWalletPrivateKey().Return(unittest.GetTestPrivateKey()).AnyTimes()
	mockConfig.EXPECT().GetNetwork().Return(config.DevnetNetwork).AnyTimes()
	mockConfig.EXPECT().GetEnvironment().Return(config.StagingEnv).AnyTimes()
	mockConfig.EXPECT().GetServerPort().Return(8080).AnyTimes()
	e := echo.New()

	t.Run("should return an error when providing invalid amount", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "xyz",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"invalid amount, must be float64\"}\n", rec.Body.String())
	})

	t.Run("should return an error when failing to get account info", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		m.
			EXPECT().
			GetAccountInfo(gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("some error")).
			AnyTimes()

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"failed to get account info\"}\n", rec.Body.String())
	})

	t.Run("should return an error when failing to decode borsh", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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
			AnyTimes()

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"failed to decode mint\"}\n", rec.Body.String())
	})

	t.Run("should return an error when api wallet is not mint authority", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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
			AnyTimes()
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
			AnyTimes()

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"invalid mint, J15X2DWTRPVTJsofDrf5se4zkNv1sD1eJPgEHwvuNJer is not MintAuthority\"}\n", rec.Body.String())
	})

	t.Run("should return an error when failing to mint", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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
			AnyTimes()
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
			AnyTimes()
		m.
			EXPECT().
			MintToWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("some error")).
			AnyTimes()

		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Equal(t, "{\"error\":\"failed to mint\"}\n", rec.Body.String())
	})

	t.Run("should return success", func(t *testing.T) {
		m := solanaClient.NewMockSolana(ctrl)
		h := NewHandler(mockConfig, m, base.NewMockBase(ctrl), repository.NewMockRepository(ctrl))
		reqBody, err := json.Marshal(apispec.MintRequest{
			Amount: "100",
			Mint:   mint,
			Wallet: solana.NewWallet().PublicKey().String(),
		})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

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
			AnyTimes()
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
			AnyTimes()
		m.
			EXPECT().
			MintToWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return("some tx hash", nil).
			AnyTimes()
		err = h.PostMint(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, "{\"txHash\":\"some tx hash\"}\n", rec.Body.String())
	})
}

func encodeMintToBase64(mint token.Mint) (string, error) {
	buf := new(bytes.Buffer)
	err := bin.NewBorshEncoder(buf).Encode(mint)
	b64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return b64Data, err
}
