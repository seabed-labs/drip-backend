package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/sirupsen/logrus"
)

func (h Handler) PostMint(
	c echo.Context,
) error {
	var mintRequest Swagger.MintRequest
	if err := c.Bind(&mintRequest); err != nil {
		return c.JSON(http.StatusBadRequest, Swagger.ErrorResponse{Error: "invalid request body"})
	}
	amount, err := strconv.ParseFloat(mintRequest.Amount, 64)
	if err != nil {
		errMsg := "invalid amount, must be float64"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: errMsg})
	}
	resp, err := h.solanaClient.GetAccountInfo(c.Request().Context(), solana.MustPublicKeyFromBase58(mintRequest.Mint))
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
	if *mint.MintAuthority != h.solanaClient.GetWalletPubKey() {
		err := fmt.Errorf("invalid mint, %s is not MintAuthority", h.solanaClient.GetWalletPubKey())
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: err.Error()})
	}
	txHash, err := h.solanaClient.MintToWallet(c.Request().Context(), mintRequest.Mint, mintRequest.Wallet, getBaseAmountWithDecimals(amount, mint.Decimals))
	if err != nil {
		errMsg := "failed to mint"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, Swagger.ErrorResponse{Error: errMsg})
	}
	return c.JSON(http.StatusOK, Swagger.MintResponse{TxHash: txHash})
}

func getBaseAmountWithDecimals(amount float64, decimals uint8) uint64 {
	if decimals <= 1 {
		return uint64(amount)
	}
	return uint64(amount * math.Pow(10, float64(decimals)))
}
