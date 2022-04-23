package api

import (
	"fmt"
	"net/http"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) GetMint(
	c echo.Context, params Swagger.GetMintParams,
) error {
	resp, err := h.solanaClient.Client.GetAccountInfo(c.Request().Context(), solana.MustPublicKeyFromBase58(string(params.Mint)))
	if err != nil {
		errMsg := "failed to get account info"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: errMsg})
	}
	var mint token.Mint
	if err := bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(&mint); err != nil {
		errMsg := "failed to decode mint"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: errMsg})
	}
	if *mint.MintAuthority != h.solanaClient.Wallet.PublicKey() {
		err := fmt.Errorf("invalid mint, %s is not MintAuthority", h.solanaClient.Wallet.PublicKey())
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: err.Error()})
	}
	txHash, err := h.solanaClient.MintToWallet(c.Request().Context(), string(params.Mint), string(params.Wallet), 1)
	if err != nil {
		errMsg := "failed to mint"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: errMsg})
	}
	return c.JSON(http.StatusOK, Swagger.MintResponse{TxHash: txHash})
}
