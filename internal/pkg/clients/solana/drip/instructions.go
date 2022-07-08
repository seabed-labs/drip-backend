// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package drip

import (
	"bytes"
	"fmt"

	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "Drip"

func init() {
	if !ProgramID.IsZero() {
		ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}

var (
	Instruction_InitVaultProtoConfig = ag_binary.TypeID([8]byte{195, 96, 99, 29, 46, 21, 146, 219})

	Instruction_InitVault = ag_binary.TypeID([8]byte{77, 79, 85, 150, 33, 217, 52, 106})

	Instruction_InitVaultPeriod = ag_binary.TypeID([8]byte{46, 103, 251, 142, 95, 43, 55, 27})

	Instruction_ClosePosition = ag_binary.TypeID([8]byte{123, 134, 81, 0, 49, 68, 98, 98})

	Instruction_Deposit = ag_binary.TypeID([8]byte{242, 35, 198, 137, 82, 225, 242, 182})

	Instruction_WithdrawB = ag_binary.TypeID([8]byte{28, 146, 254, 247, 183, 161, 195, 149})

	Instruction_TriggerDca = ag_binary.TypeID([8]byte{151, 43, 127, 87, 124, 189, 232, 35})
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id ag_binary.TypeID) string {
	switch id {
	case Instruction_InitVaultProtoConfig:
		return "InitVaultProtoConfig"
	case Instruction_InitVault:
		return "InitVault"
	case Instruction_InitVaultPeriod:
		return "InitVaultPeriod"
	case Instruction_ClosePosition:
		return "ClosePosition"
	case Instruction_Deposit:
		return "Deposit"
	case Instruction_WithdrawB:
		return "WithdrawB"
	case Instruction_TriggerDca:
		return "TriggerDca"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.AnchorTypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"init_vault_proto_config", (*InitVaultProtoConfig)(nil),
		},
		{
			"init_vault", (*InitVault)(nil),
		},
		{
			"init_vault_period", (*InitVaultPeriod)(nil),
		},
		{
			"close_position", (*ClosePosition)(nil),
		},
		{
			"deposit", (*Deposit)(nil),
		},
		{
			"withdraw_b", (*WithdrawB)(nil),
		},
		{
			"trigger_dca", (*TriggerDca)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst *Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteBytes(inst.TypeID.Bytes(), false)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBorshDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}
