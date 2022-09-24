package orcawhirlpool

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

type QuoteEstimate struct {
	EstimatedAmountIn      string `json:"estimatedAmountIn"`
	EstimatedAmountOut     string `json:"estimatedAmountOut"`
	EstimatedEndTickIndex  int    `json:"estimatedEndTickIndex"`
	EstimatedEndSqrtPrice  string `json:"estimatedEndSqrtPrice"`
	EstimatedFeeAmount     string `json:"estimatedFeeAmount"`
	Amount                 string `json:"amount"`
	AmountSpecifiedIsInput bool   `json:"amountSpecifiedIsInput"`
	AToB                   bool   `json:"aToB"`
	OtherAmountThreshold   string `json:"otherAmountThreshold"`
	SqrtPriceLimit         string `json:"sqrtPriceLimit"`
	TickArray0             string `json:"tickArray0"`
	TickArray1             string `json:"tickArray1"`
	TickArray2             string `json:"tickArray2"`
	Error                  string `json:"error"`
}

const scriptPath = "./pkg/service/orcawhirlpool/orcaWhirlpoolQuoteEstimate.ts"

func EvaluateOrcaWhirlpool(whirlpool string, vault *model.Vault, network configs.Network) (uint64, error) {
	inputAmount := uint64(100)
	// todo: we should exit early if the drip amount is 0
	if vault.DripAmount != 0 {
		inputAmount = vault.DripAmount
	}
	swapEstimate, err := getOrcaWhirlpoolQuoteEstimate(whirlpool, vault.TokenAMint, inputAmount, solana.GetURL(network))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(swapEstimate.Amount, 10, 64)
}

func getOrcaWhirlpoolQuoteEstimate(
	whirlpool string,
	inputTokenMint string,
	inputTokenAmount uint64,
	connectionUrl string,
) (QuoteEstimate, error) {
	root := configs.GetProjectRoot()
	script := fmt.Sprintf("%s/%s", root, scriptPath)
	command := fmt.Sprintf("npx ts-node %s", script) +
		fmt.Sprintf(" %s", whirlpool) +
		fmt.Sprintf(" %s", inputTokenMint) +
		fmt.Sprintf(" %d", inputTokenAmount) +
		fmt.Sprintf(" %s", connectionUrl)
	parts := strings.Fields(command)
	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return QuoteEstimate{}, err
	}
	var quote QuoteEstimate
	if err := json.Unmarshal(data, &quote); err != nil {
		return QuoteEstimate{}, fmt.Errorf("failed to unmarshal quote estimate %w", err)
	}
	if quote.Error != "" {
		return QuoteEstimate{}, fmt.Errorf("%s", quote.Error)
	}
	return quote, nil
}
