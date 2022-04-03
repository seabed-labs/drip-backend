package api

import (
	"fmt"
	"net/http"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MintResponse struct {
	TxHash string `json:"txHash"`
	Error  string `json:"Error"`
}

type MintRequest struct {
	Mint               string `json:"mint" form:"mint" binding:"required"`
	Amount             uint64 `json:"amount" form:"amount" binding:"required"`
	DestinationWallet  string `json:"destinationWallet" form:"destinationWallet"`
	DestinationAccount string `json:"destinationAccount" form:"destinationAccount"`
}

// Mint godoc
// @Summary      Mint tokens (DEVNET ONLY)
// @Description  mint test tokens to a desired associated token account, or passed in token account
// @Tags         mint
// @Accept       json
// @Produce      json
// @Param        mint                query     string  true   "Mint base58 pubkey."
// @Param        amount              query     string  true   "Amount to be minted in base amount."
// @Param        destinationWallet   query     string  false  "Destination wallet. If specificed, the associated token account will be used. If it does not exist, one will be created. One of destinationAccount or destinationWallet MUST be specified."
// @Param        destinationAccount  query     string  false  "Destination token account. Must be initialized prior to calling. One of destinationAccount or destinationWallet MUST be specified."
// @Success      200                 {object}  api.MintResponse
// @Failure      400                 {object}  api.MintResponse
// @Failure      500                 {object}  api.MintResponse
// @Router       /mint [get]
func (h Handler) Mint(
	c *gin.Context,
) {
	var mintRequest MintRequest
	if err := c.ShouldBindQuery(&mintRequest); err != nil {
		panic("todo")
	}
	res := MintResponse{}
	logrus.WithField("request", fmt.Sprintf("%v", mintRequest)).Info("")
	if mintRequest.DestinationAccount == "" && mintRequest.DestinationWallet == "" {
		err := fmt.Errorf("one of [destinationAccount, destinationWallet] must be provided")
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	resp, err := h.solanaClient.Client.GetAccountInfo(c.Request.Context(), solana.MustPublicKeyFromBase58(mintRequest.Mint))
	if err != nil {
		err := fmt.Errorf("invalid mint")
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var mint token.Mint
	if err := bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(&mint); err != nil {
		err := fmt.Errorf("invalid mint")
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if *mint.MintAuthority != h.solanaClient.Wallet.PublicKey() {
		err := fmt.Errorf("invalid mint, %s is not MintAuthority", h.solanaClient.Wallet.PublicKey())
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var txHash string
	if mintRequest.DestinationAccount != "" {
		txHash, err = h.solanaClient.MintToAccount(c.Request.Context(), mintRequest.Mint, mintRequest.DestinationAccount, mintRequest.Amount)
	} else {
		txHash, err = h.solanaClient.MintToWallet(c.Request.Context(), mintRequest.Mint, mintRequest.DestinationWallet, mintRequest.Amount)
	}
	if err != nil {
		err := fmt.Errorf("failed to MintTo")
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.TxHash = txHash
	c.JSON(http.StatusOK, res)
}
