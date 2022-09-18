package utils

import (
	"fmt"
	"strconv"

	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/gagliardetto/solana-go"
)

func GetVaultPeriodPDA(vaultAddress string, vaultPeriodID int64) (string, error) {
	vaultPubkey, err := solana.PublicKeyFromBase58(vaultAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get convert vault string addr to pubkey, err: %w", err)
	}
	vaultPeriod, _, err := solana.FindProgramAddress([][]byte{
		[]byte("vault_period"),
		vaultPubkey[:],
		[]byte(strconv.FormatInt(vaultPeriodID, 10)),
	}, drip.ProgramID)
	if err != nil {
		return "", err
	}
	return vaultPeriod.String(), nil
}
