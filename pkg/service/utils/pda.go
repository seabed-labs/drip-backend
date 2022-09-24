package utils

import (
	"fmt"
	"strconv"

	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
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

func GetTokenMetadataPDA(
	mint string,
) (string, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			solana.TokenMetadataProgramID.Bytes(),
			solana.MustPublicKeyFromBase58(mint).Bytes(),
		},
		solana.TokenMetadataProgramID,
	)
	return addr.String(), err
}

func GetWhirlpoolPDA(whirlpoolAddr string) (string, error) {
	addr, _, err := solana.FindProgramAddress([][]byte{
		[]byte("oracle"),
		solana.MustPublicKeyFromBase58(whirlpoolAddr).Bytes(),
	}, whirlpool.ProgramID)
	return addr.String(), err
}
