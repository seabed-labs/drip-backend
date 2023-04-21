package ixparser

import (
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"github.com/AlekSi/pointer"
	dripV0 "github.com/dcaf-labs/solana-drip-go/pkg/v0"
	dripV1 "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

type IxParser interface {
	GetDripIxName(version int, accounts []*solana.AccountMeta, data []byte) *string

	ParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, error)
	ParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, error)
	ParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, error)
	ParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, error)
	ParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, error)

	ParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, error)
	ParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, error)
	ParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, error)
	ParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, error)
	ParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, error)
	ParseTokenTransfer(accounts []*solana.AccountMeta, data []byte) (*token.Transfer, error)

	FindInnerIxTokenTransfer(tx rpc.GetTransactionResult, msg solana.Message, ixIndex int, condition func(transfer token.Transfer) bool) (*token.Transfer, error)
}

type ixParser struct {
}

func (p ixParser) GetDripIxName(version int, accounts []*solana.AccountMeta, data []byte) *string {
	if version == 1 {
		decodedIx, err := dripV1.DecodeInstruction(accounts, data)
		if err != nil {
			return nil
		}
		return pointer.ToString(dripV0.InstructionIDToName(decodedIx.TypeID))
	} else {
		decodedIx, err := dripV0.DecodeInstruction(accounts, data)
		if err != nil {
			return nil
		}
		return pointer.ToString(dripV0.InstructionIDToName(decodedIx.TypeID))
	}
}

func NewIxParser() IxParser {
	return ixParser{}
}

func (p ixParser) ParseTokenTransfer(accounts []*solana.AccountMeta, data []byte) (*token.Transfer, error) {
	decodedIx, err := token.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := token.InstructionIDToName(decodedIx.TypeID.Uint8()); ixName != "Transfer" {
		return nil, fmt.Errorf("not a Transfer ix")
	} else {
		return decodedIx.Impl.(*token.Transfer), nil
	}
}

func (p ixParser) ParseV0DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.Deposit, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, fmt.Errorf("not a v0 Deposit ix")
	} else {
		return decodedIx.Impl.(*dripV0.Deposit), nil
	}
}

func (p ixParser) ParseV0DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.DepositWithMetadata, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, fmt.Errorf("not a v0 DepositWithMetadata ix")
	} else {
		return decodedIx.Impl.(*dripV0.DepositWithMetadata), nil

	}
}

func (p ixParser) ParseV0DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV0.DripOrcaWhirlpool, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, fmt.Errorf("not a v0 DripOrcaWhirlpool ix")
	} else {
		return decodedIx.Impl.(*dripV0.DripOrcaWhirlpool), nil
	}
}

func (p ixParser) ParseV0ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.ClosePosition, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, fmt.Errorf("not a v0 ClosePosition ix")
	} else {
		return decodedIx.Impl.(*dripV0.ClosePosition), nil
	}
}

func (p ixParser) ParseV0WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV0.WithdrawB, error) {
	decodedIx, err := dripV0.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV0.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, fmt.Errorf("not a v0 WithdrawB ix")
	} else {
		return decodedIx.Impl.(*dripV0.WithdrawB), nil

	}
}

func (p ixParser) ParseV1DepositIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.Deposit, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "Deposit" {
		return nil, fmt.Errorf("not a v1 Depost ix")
	} else {
		return decodedIx.Impl.(*dripV1.Deposit), nil
	}
}

func (p ixParser) ParseV1DepositWithMetadataIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.DepositWithMetadata, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DepositWithMetadata" {
		return nil, fmt.Errorf("not a v1 DepositWithMetadata ix")
	} else {
		return decodedIx.Impl.(*dripV1.DepositWithMetadata), nil

	}
}

func (p ixParser) ParseV1DripOrcaWhirlpool(accounts []*solana.AccountMeta, data []byte) (*dripV1.DripOrcaWhirlpool, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "DripOrcaWhirlpool" {
		return nil, fmt.Errorf("not a v1 DripOrcaWhirlpool ix")
	} else {
		return decodedIx.Impl.(*dripV1.DripOrcaWhirlpool), nil
	}
}

func (p ixParser) ParseV1ClosePositionIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.ClosePosition, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "ClosePosition" {
		return nil, fmt.Errorf("not a v1 ClosePosition ix")
	} else {
		return decodedIx.Impl.(*dripV1.ClosePosition), nil
	}
}

func (p ixParser) ParseV1WithdrawIx(accounts []*solana.AccountMeta, data []byte) (*dripV1.WithdrawB, error) {
	decodedIx, err := dripV1.DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	if ixName := dripV1.InstructionIDToName(decodedIx.TypeID); ixName != "WithdrawB" {
		return nil, fmt.Errorf("not a v1 WithdrawB ix")
	} else {
		return decodedIx.Impl.(*dripV1.WithdrawB), nil

	}
}

func (p ixParser) FindInnerIxTokenTransfer(txRaw rpc.GetTransactionResult, msg solana.Message, ixIndex int, condition func(transfer token.Transfer) bool) (*token.Transfer, error) {
	innerIXS := lo.Filter(txRaw.Meta.InnerInstructions, func(innerIx rpc.InnerInstruction, idx int) bool {
		return innerIx.Index == uint16(ixIndex)
	})
	instructions := lo.Flatten(lo.Map(innerIXS, func(innerInnerIx rpc.InnerInstruction, innerInnerIdx int) []solana.CompiledInstruction {
		return innerInnerIx.Instructions
	}))
	tokenTransfers := lo.FilterMap(instructions, func(ix solana.CompiledInstruction, idx int) (token.Transfer, bool) {
		accounts, err := ix.ResolveInstructionAccounts(&msg)
		if err != nil {
			logrus.WithError(err).Error("failed to ix.ResolveInstructionAccounts(&msg)")
			return token.Transfer{}, false
		}
		if transfer, err := p.ParseTokenTransfer(accounts, ix.Data); err != nil {
			return token.Transfer{}, false
		} else {
			return *transfer, condition(*transfer)
		}
	})
	if len(tokenTransfers) > 1 {
		return nil, fmt.Errorf("expected 0 or 1 transfers, but got %d", len(tokenTransfers))
	}
	if len(tokenTransfers) == 1 {
		return &tokenTransfers[0], nil
	}
	return nil, nil
}
