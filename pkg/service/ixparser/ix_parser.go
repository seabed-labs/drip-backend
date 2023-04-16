package ixparser

import (
	dripV0 "github.com/dcaf-labs/solana-drip-go/pkg/v0"
	dripV1 "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/gagliardetto/solana-go"
)

type IxParser interface {
	MaybeParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, error)
	MaybeParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, error)
	MaybeParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, error)
	MaybeParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, error)
	MaybeParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, error)
	MaybeParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, error)
	MaybeParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, error)
	MaybeParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, error)
	MaybeParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, error)
	MaybeParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, error)
}

type ixParser struct {
}

func NewIxParser() IxParser {
	return ixParser{}
}

func (p ixParser) MaybeParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV0.Deposit), nil
}

func (p ixParser) MaybeParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV0.DepositWithMetadata), nil
}

func (p ixParser) MaybeParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV0.DripOrcaWhirlpool), nil
}

func (p ixParser) MaybeParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV0.ClosePosition), nil
}

func (p ixParser) MaybeParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV0.WithdrawB), nil
}

func (p ixParser) MaybeParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV1.Deposit), nil
}

func (p ixParser) MaybeParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV1.DepositWithMetadata), nil
}

func (p ixParser) MaybeParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV1.DripOrcaWhirlpool), nil
}

func (p ixParser) MaybeParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV1.ClosePosition), nil
}

func (p ixParser) MaybeParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, nil
	}
	return decodedIx.Impl.(*dripV1.WithdrawB), nil
}
