package processor

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	drip "github.com/dcaf-labs/solana-drip-go/pkg"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (p impl) ProcessTransactionUpdateQueue(ctx context.Context) {
	var wg sync.WaitGroup
	ch := make(chan *model.TransactionUpdateQueueItem)
	defer func() {
		close(ch)
		logrus.Info("exiting ProcessTransactionUpdateQueue...")
		wg.Wait()
	}()

	for i := 0; i < processConcurrency; i++ {
		wg.Add(1)
		go p.processTransactionUpdateQueueItemWorker(ctx, strconv.FormatInt(int64(i), 10), &wg, ch)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			queueItem, err := p.transactionUpdateQueue.PopTransactionUpdateQueueItem(ctx)
			if err != nil && err == gorm.ErrRecordNotFound {
				continue
			} else if err != nil {
				logrus.WithError(err).Error("failed to fetch transaction from queue")
				continue
			} else if queueItem == nil {
				logrus.WithError(err).Error("failed to get next queue item")
				continue
			} else {
				ch <- queueItem
			}
		}
	}
}

func (p impl) processTransactionUpdateQueueItemWorker(ctx context.Context, id string, wg *sync.WaitGroup, queueCh chan *model.TransactionUpdateQueueItem) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logrus.Info("exiting processTransactionUpdateQueueItemWorker")
			return
		case queueItem := <-queueCh:
			if queueItem != nil {
				var txRaw rpc.GetTransactionResult
				if err := json.Unmarshal([]byte(queueItem.TxJSON), &txRaw); err != nil {
					logrus.WithError(err).Error("failed to unmarshall tx...")
					p.handleTransactionUpdateErr(ctx, queueItem, err)
				}
				if err := p.ProcessTransaction(ctx, txRaw); err != nil {
					p.handleTransactionUpdateErr(ctx, queueItem, err)
				}
			}
		}
	}
}

func (p impl) ProcessTransaction(ctx context.Context, txRaw rpc.GetTransactionResult) error {
	tx, err := txRaw.Transaction.GetTransaction()
	if err != nil {
		logrus.WithError(err).Error("failed to get transaction from unmarshalled tx")
		return err
	}
	version := drip.GetIdlVersion(txRaw.Slot)
	blockTime := time.Unix(0, 0)
	if txRaw.BlockTime != nil {
		blockTime = time.Unix(int64(*txRaw.BlockTime), 0)
	}
	signature := tx.Signatures[0].String()
	log := logrus.WithField("signature", signature).WithField("version", version)

	for i := range tx.Message.Instructions {
		ix := tx.Message.Instructions[i]
		accounts, err := ix.ResolveInstructionAccounts(&tx.Message)
		if err != nil {
			logrus.WithError(err).Error("failed ResolveInstructionAccounts")
			return err
		}
		ixName := p.ixParser.GetDripIxName(version, accounts, ix.Data)
		if ixName == nil {
			continue
		}
		log = log.WithField("ixName", *ixName)
		log.Info("starting to parse ix")

		switch *ixName {
		case "Deposit":
			metric, err := p.getDepositMetric(txRaw, ix, accounts, version, i, signature, *ixName, blockTime)
			if err != nil {
				return err
			}
			return p.repo.UpsertDepositMetric(ctx, metric)
		case "DepositWithMetadata":
			metric, err2 := p.getDepositWithMetadataMetric(txRaw, ix, accounts, version, i, signature, *ixName, blockTime)
			if err2 != nil {
				return err2
			}
			return p.repo.UpsertDepositMetric(ctx, metric)

		case "DripOrcaWhirlpool":
			metric, err2 := p.getDripMetric(txRaw, tx.Message, ix, accounts, version, i, signature, *ixName, blockTime)
			if err2 != nil {
				return err2
			}
			return p.repo.UpsertDripMetric(ctx, metric)
		case "ClosePosition":
			metric, err2 := p.getClosePositionMetric(txRaw, tx.Message, ix, accounts, version, i, signature, *ixName, blockTime)
			if err2 != nil {
				return err2
			}
			return p.repo.UpsertWithdrawMetric(ctx, metric)
		case "WithdrawB":
			metric, err2 := p.getWithdrawBMetric(txRaw, tx.Message, ix, accounts, version, i, signature, *ixName, blockTime)
			if err2 != nil {
				return err2
			}
			return p.repo.UpsertWithdrawMetric(ctx, metric)
		}
	}
	return nil
}

func (p impl) getWithdrawBMetric(
	txRaw rpc.GetTransactionResult, txMsg solana.Message, ix solana.CompiledInstruction, accounts []*solana.AccountMeta,
	version int, ixIndex int, signature string, ixName string, blockTime time.Time,
) (*model.WithdrawalMetric, error) {
	log := logrus.WithField("signature", signature).WithField("ixName", ixName).WithField("version", version).WithField("ixIndex", ixIndex)
	var (
		vault                      string
		vaultTreasuryTokenBAccount string
		userTokenBAccount          string

		userTokenAWithdrawAmount     uint64
		userTokenBWithdrawAmount     uint64
		treasuryTokenBReceivedAmount uint64
		referralTokenBReceivedAmount uint64
	)
	referrerIsAlsoTreasury := false
	if version == 1 {
		if parsedIx, err := p.ixParser.ParseV1WithdrawIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetWithdrawBAccounts()
			vault = parsedAccounts.Common.Vault.String()
			vaultTreasuryTokenBAccount = parsedAccounts.Common.VaultTreasuryTokenBAccount.String()
			userTokenBAccount = parsedAccounts.Common.UserTokenBAccount.String()
			if parsedAccounts.Common.Referrer.String() == parsedAccounts.Common.VaultTreasuryTokenBAccount.String() {
				referrerIsAlsoTreasury = true
			}
			foundTransfer := false
			if vaultTokenBToReferral, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
				destEqualsReferrer := transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.Referrer.String()
				// when the referrer and treasury are the same, we expect up to two token transfers to this token account
				if referrerIsAlsoTreasury && destEqualsReferrer && !foundTransfer {
					foundTransfer = true
					return true
				}
				return false
			}); err != nil {
				log.WithError(err).Error("failed to find vaultTokenBToReferral")
				return nil, err
			} else if vaultTokenBToReferral != nil {
				referralTokenBReceivedAmount = *vaultTokenBToReferral.Amount
			}
		} else if err != nil {
			return nil, err
		}
	} else {
		if parsedIx, err := p.ixParser.ParseV0WithdrawIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetWithdrawBAccounts()
			vault = parsedAccounts.Vault.String()
			vaultTreasuryTokenBAccount = parsedAccounts.VaultTreasuryTokenBAccount.String()
			userTokenBAccount = parsedAccounts.UserTokenBAccount.String()
		} else if err != nil {
			return nil, err
		}
	}
	if vaultTokenBToUser, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[1].PublicKey.String() == userTokenBAccount
	}); err != nil {
		log.WithError(err).Error("failed to find vaultTokenBToUser")
		return nil, err
	} else if vaultTokenBToUser != nil {
		userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
	}
	tokenTransferMatchCount := 0
	if vaultTokenBToTreasury, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		destEqualsVaultTreasury := transfer.Accounts[1].PublicKey.String() == vaultTreasuryTokenBAccount
		if destEqualsVaultTreasury {
			tokenTransferMatchCount = tokenTransferMatchCount + 1
		}
		// If the referrer is also the treasury we expect 0-2 transfers to the same token b account
		// if there is 1, then it is always the referrer account
		// if there is 2, then the first is the referrer transfer, the second is the treasury transfer
		if referrerIsAlsoTreasury {
			return destEqualsVaultTreasury && tokenTransferMatchCount == 2
		}
		return destEqualsVaultTreasury
	}); err != nil {
		log.WithError(err).Error("failed to find vaultTokenBToTreasury")
		return nil, err
	} else if vaultTokenBToTreasury != nil {
		treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
	}
	metric := &model.WithdrawalMetric{
		Signature:                    signature,
		IxIndex:                      int32(ixIndex),
		IxName:                       ixName,
		IxVersion:                    int32(version),
		Slot:                         int32(txRaw.Slot),
		Time:                         blockTime,
		Vault:                        vault,
		UserTokenAWithdrawAmount:     userTokenAWithdrawAmount,
		UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
		TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
		ReferralTokenBReceivedAmount: referralTokenBReceivedAmount,
	}
	return metric, nil
}

func (p impl) getClosePositionMetric(
	txRaw rpc.GetTransactionResult, txMsg solana.Message, ix solana.CompiledInstruction, accounts []*solana.AccountMeta,
	version int, ixIndex int, signature string, ixName string, blockTime time.Time,
) (*model.WithdrawalMetric, error) {
	log := logrus.WithField("signature", signature).WithField("ixName", ixName).WithField("version", version).WithField("ixIndex", ixIndex)
	var (
		vault                      string
		userTokenAAcount           string
		userTokenBAccount          string
		vaultTreasuryTokenBAccount string

		userTokenAWithdrawAmount     uint64
		userTokenBWithdrawAmount     uint64
		treasuryTokenBReceivedAmount uint64
		referralTokenBReceivedAmount uint64
	)
	referrerIsAlsoTreasury := false
	if version == 1 {
		if parsedIx, err := p.ixParser.ParseV1ClosePositionIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetClosePositionAccounts()
			vault = parsedAccounts.Common.Vault.String()
			userTokenAAcount = parsedAccounts.UserTokenAAccount.String()
			userTokenBAccount = parsedAccounts.Common.UserTokenBAccount.String()
			vaultTreasuryTokenBAccount = parsedAccounts.Common.VaultTreasuryTokenBAccount.String()
			if parsedAccounts.Common.Referrer.String() == parsedAccounts.Common.VaultTreasuryTokenBAccount.String() {
				referrerIsAlsoTreasury = true
			}
			foundTransfer := false
			if vaultTokenBToReferral, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
				destEqualsReferrer := transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.Referrer.String()
				// when the referrer and treasury are the same, we expect up to two token transfers to this token account
				if referrerIsAlsoTreasury && destEqualsReferrer && !foundTransfer {
					foundTransfer = true
					return true
				}
				return false
			}); err != nil {
				log.WithError(err).Error("failed to find vaultTokenBToReferral")
				return nil, err
			} else if vaultTokenBToReferral != nil {
				referralTokenBReceivedAmount = *vaultTokenBToReferral.Amount
			}
		} else if err != nil {
			return nil, err
		}
	} else {
		if parsedIx, err := p.ixParser.ParseV0ClosePositionIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetClosePositionAccounts()
			vault = parsedAccounts.Vault.String()
			userTokenAAcount = parsedAccounts.UserTokenAAccount.String()
			userTokenBAccount = parsedAccounts.UserTokenBAccount.String()
			vaultTreasuryTokenBAccount = parsedAccounts.VaultTreasuryTokenBAccount.String()
		} else if err != nil {
			return nil, err
		}
	}
	if vaultTokenAToUser, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[1].PublicKey.String() == userTokenAAcount
	}); err != nil {
		log.WithError(err).Error("failed to find vaultTokenAToUser")
		return nil, err
	} else if vaultTokenAToUser != nil {
		userTokenAWithdrawAmount = *vaultTokenAToUser.Amount
	}
	if vaultTokenBToUser, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[1].PublicKey.String() == userTokenBAccount
	}); err != nil {
		log.WithError(err).Error("failed to find vaultTokenBToUser")
		return nil, err
	} else if vaultTokenBToUser != nil {
		userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
	}
	tokenTransferMatchCount := 0
	if vaultTokenBToTreasury, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		destEqualsVaultTreasury := transfer.Accounts[1].PublicKey.String() == vaultTreasuryTokenBAccount
		if destEqualsVaultTreasury {
			tokenTransferMatchCount = tokenTransferMatchCount + 1
		}
		// If the referrer is also the treasury we expect 0-2 transfers to the same token b account
		// if there is 1, then it is always the referrer account
		// if there is 2, then the first is the referrer transfer, the second is the treasury transfer
		if referrerIsAlsoTreasury {
			return destEqualsVaultTreasury && tokenTransferMatchCount == 2
		}
		return destEqualsVaultTreasury
	}); err != nil {
		log.WithError(err).Error("failed to find vaultTokenBToTreasury")
		return nil, err
	} else if vaultTokenBToTreasury != nil {
		treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
	}
	metric := &model.WithdrawalMetric{
		Signature:                    signature,
		IxIndex:                      int32(ixIndex),
		IxName:                       ixName,
		IxVersion:                    int32(version),
		Slot:                         int32(txRaw.Slot),
		Time:                         blockTime,
		Vault:                        vault,
		UserTokenAWithdrawAmount:     userTokenAWithdrawAmount,
		UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
		TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
		ReferralTokenBReceivedAmount: referralTokenBReceivedAmount,
	}
	return metric, nil
}

func (p impl) getDripMetric(
	txRaw rpc.GetTransactionResult, txMsg solana.Message, ix solana.CompiledInstruction, accounts []*solana.AccountMeta,
	version int, ixIndex int, signature string, ixName string, blockTime time.Time,
) (*model.DripMetric, error) {
	log := logrus.WithField("signature", signature).WithField("ixName", ixName).WithField("version", version).WithField("ixIndex", ixIndex)
	var (
		vault                string
		vaultTokenAAccount   string
		vaultTokenBAccount   string
		dripFeeTokenAAccount string
	)
	if version == 1 {
		if parsedIx, err := p.ixParser.ParseV1DripOrcaWhirlpool(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDripOrcaWhirlpoolAccounts()
			vault = parsedAccounts.Common.Vault.String()
			vaultTokenAAccount = parsedAccounts.Common.VaultTokenAAccount.String()
			vaultTokenBAccount = parsedAccounts.Common.VaultTokenBAccount.String()
			dripFeeTokenAAccount = parsedAccounts.Common.DripFeeTokenAAccount.String()
		} else if err != nil {
			return nil, err
		}
	} else {
		if parsedIx, err := p.ixParser.ParseV0DripOrcaWhirlpool(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDripOrcaWhirlpoolAccounts()
			vault = parsedAccounts.Vault.String()
			vaultTokenAAccount = parsedAccounts.VaultTokenAAccount.String()
			vaultTokenBAccount = parsedAccounts.VaultTokenBAccount.String()
			dripFeeTokenAAccount = parsedAccounts.DripFeeTokenAAccount.String()
		} else if err != nil {
			return nil, err
		}
	}
	tokenASwapped := uint64(0)
	tokenBReceived := uint64(0)
	keeperTokenAReceived := uint64(0)
	if vaultTokenAToWhirlpool, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[0].PublicKey.String() == vaultTokenAAccount &&
			transfer.Accounts[1].PublicKey.String() != dripFeeTokenAAccount
	}); err != nil {
		log.WithError(err).Errorf("failed to find vaultTokenAToWhirlpool")
		return nil, err
	} else if vaultTokenAToWhirlpool != nil {
		tokenASwapped = *vaultTokenAToWhirlpool.Amount
	}
	if vaultTokenAToKeeper, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[1].PublicKey.String() == dripFeeTokenAAccount
	}); err != nil {
		log.WithError(err).Errorf("failed to find vaultTokenAToKeeper")
		return nil, err
	} else if vaultTokenAToKeeper != nil {
		keeperTokenAReceived = *vaultTokenAToKeeper.Amount
	}
	if whirlpoolToVaultTokenB, err := p.ixParser.FindInnerIxTokenTransfer(txRaw, txMsg, ixIndex, func(transfer token.Transfer) bool {
		return transfer.Accounts[1].PublicKey.String() == vaultTokenBAccount
	}); err != nil {
		log.WithError(err).Errorf("failed to find whirlpoolToVaultTokenB")
		return nil, err
	} else if whirlpoolToVaultTokenB != nil {
		tokenBReceived = *whirlpoolToVaultTokenB.Amount
	}
	metric := &model.DripMetric{
		Signature:                  signature,
		IxIndex:                    int32(ixIndex),
		IxName:                     ixName,
		IxVersion:                  int32(version),
		Slot:                       int32(txRaw.Slot),
		Time:                       blockTime,
		Vault:                      vault,
		VaultTokenASwappedAmount:   tokenASwapped,
		VaultTokenBReceivedAmount:  tokenBReceived,
		KeeperTokenAReceivedAmount: keeperTokenAReceived,
		TokenAUsdPriceDay:          nil,
		TokenBUsdPriceDay:          nil,
	}
	return metric, nil
}

func (p impl) getDepositWithMetadataMetric(
	txRaw rpc.GetTransactionResult, ix solana.CompiledInstruction, accounts []*solana.AccountMeta,
	version int, ixIndex int, signature string, ixName string, blockTime time.Time,
) (*model.DepositMetric, error) {
	var (
		vault               string
		depositor           string
		referrer            *string
		tokenADepositAmount uint64
	)
	if version == 1 {
		if parsedIx, err := p.ixParser.ParseV1DepositWithMetadataIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDepositWithMetadataAccounts()
			vault = parsedAccounts.Common.Vault.String()
			referrer = pointer.ToString(parsedAccounts.Common.Referrer.String())
			tokenADepositAmount = parsedIx.Params.TokenADepositAmount
			depositor = parsedAccounts.Common.Depositor.String()
		} else if err != nil {
			return nil, err
		}
	} else {
		if parsedIx, err := p.ixParser.ParseV0DepositWithMetadataIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDepositWithMetadataAccounts()
			vault = parsedAccounts.Vault.String()
			referrer = nil
			tokenADepositAmount = parsedIx.Params.TokenADepositAmount
			depositor = parsedAccounts.Depositor.String()
		} else if err != nil {
			return nil, err
		}
	}
	metric := &model.DepositMetric{
		Signature:           signature,
		IxIndex:             int32(ixIndex),
		IxName:              ixName,
		IxVersion:           int32(version),
		Slot:                int32(txRaw.Slot),
		Time:                blockTime,
		Vault:               vault,
		Referrer:            referrer,
		TokenADepositAmount: tokenADepositAmount,
		Depositor:           pointer.ToString(depositor),
		TokenAUsdPriceDay:   nil,
	}
	return metric, nil
}

func (p impl) getDepositMetric(
	txRaw rpc.GetTransactionResult, ix solana.CompiledInstruction, accounts []*solana.AccountMeta,
	version int, ixIndex int, signature string, ixName string, blockTime time.Time,
) (*model.DepositMetric, error) {
	var (
		vault               string
		depositor           string
		referrer            *string
		tokenADepositAmount uint64
	)

	if version == 1 {
		if parsedIx, err := p.ixParser.ParseV1DepositIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDepositAccounts()
			vault = parsedAccounts.Common.Vault.String()
			referrer = pointer.ToString(parsedAccounts.Common.Referrer.String())
			tokenADepositAmount = parsedIx.Params.TokenADepositAmount
			depositor = parsedAccounts.Common.Depositor.String()
		} else if err != nil {
			return nil, err
		}
	} else {
		if parsedIx, err := p.ixParser.ParseV0DepositIx(accounts, ix.Data); parsedIx != nil {
			parsedAccounts := parsedIx.GetDepositAccounts()
			vault = parsedAccounts.Vault.String()
			referrer = nil
			tokenADepositAmount = parsedIx.Params.TokenADepositAmount
			depositor = parsedAccounts.Depositor.String()
		} else if err != nil {
			return nil, err
		}
	}
	metric := &model.DepositMetric{
		Signature:           signature,
		IxIndex:             int32(ixIndex),
		IxName:              ixName,
		IxVersion:           int32(version),
		Slot:                int32(txRaw.Slot),
		Time:                blockTime,
		Vault:               vault,
		Referrer:            referrer,
		TokenADepositAmount: tokenADepositAmount,
		Depositor:           pointer.ToString(depositor),
		TokenAUsdPriceDay:   nil,
	}
	return metric, nil
}

func (p impl) handleTransactionUpdateErr(ctx context.Context, queueItem *model.TransactionUpdateQueueItem, err error) {
	logrus.WithError(err).Error("failed to process update")
	if requeueErr := p.transactionUpdateQueue.ReQueueTransactionUpdateItem(ctx, queueItem, err.Error()); requeueErr != nil {
		logrus.WithError(requeueErr).Error("failed to add item to queue")
	}
}
