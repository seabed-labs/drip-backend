package ixparser

import (
	dripV0 "github.com/dcaf-labs/solana-drip-go/pkg/v0"
	dripV1 "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

type IxParser interface {
	MaybeParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, string, error)
	MaybeParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, string, error)
	MaybeParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, string, error)
	MaybeParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, string, error)
	MaybeParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, string, error)
	MaybeParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, string, error)
	MaybeParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, string, error)
	MaybeParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, string, error)
	MaybeParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, string, error)
	MaybeParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, string, error)
	MaybeParseTokenTransfer(accounts []*solana.AccountMeta, data []byte) (*token.Transfer, string, error)
}

type ixParser struct {
}

func NewIxParser() IxParser {
	return ixParser{}
}

func (p ixParser) MaybeParseTokenTransfer(accounts []*solana.AccountMeta, data []byte) (*token.Transfer, string, error) {
	decodedIx, err := token.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := token.InstructionIDToName(decodedIx.TypeID.Uint8()); ixName != "Transfer" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*token.Transfer), ixName, nil
	}
}

func (p ixParser) MaybeParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, string, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV0.Deposit), ixName, nil
	}
}

func (p ixParser) MaybeParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, string, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV0.DepositWithMetadata), ixName, nil

	}
}

func (p ixParser) MaybeParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, string, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV0.DripOrcaWhirlpool), ixName, nil
	}
}

func (p ixParser) MaybeParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, string, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV0.ClosePosition), ixName, nil
	}
}

func (p ixParser) MaybeParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, string, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV0.WithdrawB), ixName, nil

	}
}

func (p ixParser) MaybeParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, string, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV1.Deposit), ixName, nil
	}
}

func (p ixParser) MaybeParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, string, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV1.DepositWithMetadata), ixName, nil

	}
}

func (p ixParser) MaybeParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, string, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV1.DripOrcaWhirlpool), ixName, nil
	}
}

func (p ixParser) MaybeParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, string, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV1.ClosePosition), ixName, nil
	}
}

func (p ixParser) MaybeParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, string, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, "", err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, "", nil
	} else {
		return decodedIx.Impl.(*dripV1.WithdrawB), ixName, nil

	}
}
