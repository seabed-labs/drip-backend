package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h Handler) PostMint(
	c echo.Context,
) error {
	var mintRequest apispec.MintRequest
	if err := c.Bind(&mintRequest); err != nil {
		return c.JSON(http.StatusBadRequest, apispec.ErrorResponse{Error: "invalid request body"})
	}
	amount, err := strconv.ParseFloat(mintRequest.Amount, 64)
	if err != nil {
		errMsg := "invalid amount, must be float64"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: errMsg})
	}
	resp, err := h.solanaClient.GetAccountInfo(c.Request().Context(), solana.MustPublicKeyFromBase58(mintRequest.Mint))
	if err != nil {
		errMsg := "failed to get account info"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: errMsg})
	}
	var mint token.Mint
	if err := bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(&mint); err != nil {
		errMsg := "failed to decode mint"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: errMsg})
	}
	if *mint.MintAuthority != h.solanaClient.GetWalletPubKey() {
		err := fmt.Errorf("invalid mint, %s is not MintAuthority", h.solanaClient.GetWalletPubKey())
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
	}
	txHash, err := h.solanaClient.MintToWallet(c.Request().Context(), mintRequest.Mint, mintRequest.Wallet, getBaseAmountWithDecimals(amount, mint.Decimals))
	if err != nil {
		errMsg := "failed to mint"
		logrus.WithError(err).Errorf(errMsg)
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: errMsg})
	}
	return c.JSON(http.StatusOK, apispec.MintResponse{TxHash: txHash})
}

func getBaseAmountWithDecimals(amount float64, decimals uint8) uint64 {
	if decimals <= 1 {
		return uint64(amount)
	}
	return uint64(amount * math.Pow(10, float64(decimals)))
}
