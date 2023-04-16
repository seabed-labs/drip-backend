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
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (p impl) ProcessTransactionUpateQueue(ctx context.Context) {
	var wg sync.WaitGroup
	ch := make(chan *model.TransactionUpdateQueueItem)
	defer func() {
		close(ch)
		logrus.Info("exiting ProcessTransactionUpateQueue...")
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
			} else if queueItem == nil {
				logrus.WithError(err).Error("failed to get next queue item")
				continue
			}
			//if depth, err := p.accountUpdateQueue.AccountUpdateQueueItemDepth(ctx); err != nil {
			//	logrus.WithError(err).Error("failed to get queue depth")
			//} else {
			//	logrus.WithField("depth", depth).Infof("queue depth")
			//}
			ch <- queueItem
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
			if err := p.processTransactionUpdateQueueItem(ctx, id, queueItem); err != nil {
				p.handleTransactionUpdateErr(ctx, queueItem, err)
			}
		}
	}
}

func (p impl) processTransactionUpdateQueueItem(ctx context.Context, id string, queueItem *model.TransactionUpdateQueueItem) error {
	log := logrus.WithField("id", id).WithField("signature", queueItem.Signature)
	var txRaw rpc.TransactionWithMeta
	if err := json.Unmarshal([]byte(queueItem.TxJSON), &txRaw); err != nil {
		log.WithError(err).Error("failed to unmarshall tx...")
		return err
	}
	tx, err := txRaw.GetTransaction()
	if err != nil {
		log.WithError(err).Error("failed to get transaction from unmarshalled tx")
		return err
	}
	version := drip.GetIdlVersion(txRaw.Slot)
	for i := range tx.Message.Instructions {
		ix := tx.Message.Instructions[i]
		accounts, err := ix.ResolveInstructionAccounts(&tx.Message)
		if err != nil {
			logrus.WithError(err).Error("failed ResolveInstructionAccounts")
			return err
		}
		blockTime := time.Unix(0, 0)
		if txRaw.BlockTime != nil {
			blockTime = time.Unix(int64(*txRaw.BlockTime), 0)
		}
		signature := tx.Signatures[0].String()
		switch version {
		case 1:
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV1DepositIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetDepositAccounts()
				return p.repo.UpsertDepositMetric(ctx, &model.DepositMetric{
					Signature:           signature,
					IxIndex:             int32(i),
					IxName:              ixName,
					IxVersion:           0,
					Slot:                int32(txRaw.Slot),
					Time:                blockTime,
					Vault:               parsedAccounts.Common.Vault.String(),
					Referrer:            pointer.ToString(parsedAccounts.Common.Referrer.String()),
					TokenADepositAmount: parsedIx.Params.TokenADepositAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV1DepositIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV1DepositWithMetadataIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetDepositWithMetadataAccounts()
				return p.repo.UpsertDepositMetric(ctx, &model.DepositMetric{
					Signature:           signature,
					IxIndex:             int32(i),
					IxName:              ixName,
					IxVersion:           0,
					Slot:                int32(txRaw.Slot),
					Time:                blockTime,
					Vault:               parsedAccounts.Common.Vault.String(),
					Referrer:            pointer.ToString(parsedAccounts.Common.Referrer.String()),
					TokenADepositAmount: parsedIx.Params.TokenADepositAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV1DepositWithMetadataIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV1DripOrcaWhirlpool(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetDripOrcaWhirlpoolAccounts()
				tokenASwapped := uint64(0)
				tokenBReceived := uint64(0)
				keeperTokenAReceived := uint64(0)
				if vaultTokenAToWhirlpool := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[0].PublicKey.String() == parsedAccounts.Common.VaultTokenAAccount.String()
				}); vaultTokenAToWhirlpool != nil {
					tokenASwapped = *vaultTokenAToWhirlpool.Amount
				}
				if vaultTokenAToKeeper := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.DripFeeTokenAAccount.String()
				}); vaultTokenAToKeeper != nil {
					keeperTokenAReceived = *vaultTokenAToKeeper.Amount
				}
				if whirlpoolToVaultTokenB := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.VaultTokenBAccount.String()
				}); whirlpoolToVaultTokenB != nil {
					tokenBReceived = *whirlpoolToVaultTokenB.Amount
				}

				return p.repo.UpsertDripMetric(ctx, &model.DripMetric{
					Signature:                  signature,
					IxIndex:                    int32(i),
					IxName:                     ixName,
					IxVersion:                  1,
					Slot:                       int32(txRaw.Slot),
					Time:                       blockTime,
					Vault:                      parsedAccounts.Common.Vault.String(),
					VaultTokenASwappedAmount:   tokenASwapped,
					VaultTokenBReceivedAmount:  tokenBReceived,
					KeeperTokenAReceivedAmount: keeperTokenAReceived,
				})

			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV1DripOrcaWhirlpool")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV1WithdrawIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetWithdrawBAccounts()
				userTokenBWithdrawAmount := uint64(0)
				treasuryTokenBReceivedAmount := uint64(0)
				referralTokenBReceivedAmount := uint64(0)
				if vaultTokenBToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.UserTokenBAccount.String()
				}); vaultTokenBToUser != nil {
					userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
				}
				if vaultTokenBToTreasury := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.VaultTreasuryTokenBAccount.String()
				}); vaultTokenBToTreasury != nil {
					treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
				}
				if vaultTokenBToReferral := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.Referrer.String()
				}); vaultTokenBToReferral != nil {
					referralTokenBReceivedAmount = *vaultTokenBToReferral.Amount
				}
				return p.repo.UpsertWithdrawMetric(ctx, &model.WithdrawalMetric{
					Signature:                    signature,
					IxIndex:                      int32(i),
					IxName:                       ixName,
					IxVersion:                    0,
					Slot:                         int32(txRaw.Slot),
					Time:                         blockTime,
					Vault:                        parsedAccounts.Common.Vault.String(),
					UserTokenAWithdrawAmount:     0,
					UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
					TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
					ReferralTokenBReceivedAmount: referralTokenBReceivedAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV1WithdrawIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV1ClosePositionIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetClosePositionAccounts()
				userTokenAWithdrawAmount := uint64(0)
				userTokenBWithdrawAmount := uint64(0)
				treasuryTokenBReceivedAmount := uint64(0)
				referralTokenBReceivedAmount := uint64(0)
				if vaultTokenAToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.UserTokenAAccount.String()
				}); vaultTokenAToUser != nil {
					userTokenAWithdrawAmount = *vaultTokenAToUser.Amount
				}
				if vaultTokenBToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.UserTokenBAccount.String()
				}); vaultTokenBToUser != nil {
					userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
				}
				if vaultTokenBToTreasury := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.VaultTreasuryTokenBAccount.String()
				}); vaultTokenBToTreasury != nil {
					treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
				}
				if vaultTokenBToReferral := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.Common.Referrer.String()
				}); vaultTokenBToReferral != nil {
					referralTokenBReceivedAmount = *vaultTokenBToReferral.Amount
				}
				return p.repo.UpsertWithdrawMetric(ctx, &model.WithdrawalMetric{
					Signature:                    signature,
					IxIndex:                      int32(i),
					IxName:                       ixName,
					IxVersion:                    0,
					Slot:                         int32(txRaw.Slot),
					Time:                         blockTime,
					Vault:                        parsedAccounts.Common.Vault.String(),
					UserTokenAWithdrawAmount:     userTokenAWithdrawAmount,
					UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
					TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
					ReferralTokenBReceivedAmount: referralTokenBReceivedAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV1ClosePositionIx")
				return err
			}
		case 0:
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV0DepositIx(accounts, ix.Data); parsedIx != nil {
				return p.repo.UpsertDepositMetric(ctx, &model.DepositMetric{
					Signature: signature,
					IxIndex:   int32(i),
					IxName:    ixName,
					IxVersion: 0,
					Slot:      int32(txRaw.Slot),
					Time:      blockTime,
					Vault:     parsedIx.GetVaultAccount().PublicKey.String(),
					// todo: remove
					//TokenAMint:          parsedIx.,
					//Referrer:            nil,
					TokenADepositAmount: parsedIx.Params.TokenADepositAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV0DepositIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV0DepositWithMetadataIx(accounts, ix.Data); parsedIx != nil {
				return p.repo.UpsertDepositMetric(ctx, &model.DepositMetric{
					Signature: signature,
					IxIndex:   int32(i),
					IxName:    ixName,
					IxVersion: 0,
					Slot:      int32(txRaw.Slot),
					Time:      blockTime,
					Vault:     parsedIx.GetVaultAccount().PublicKey.String(),
					// todo: remove
					//TokenAMint:          parsedIx.,
					//Referrer:            nil,
					TokenADepositAmount: parsedIx.Params.TokenADepositAmount,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV0DepositWithMetadataIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV0DripOrcaWhirlpool(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetDripOrcaWhirlpoolAccounts()
				tokenASwapped := uint64(0)
				tokenBReceived := uint64(0)
				keeperTokenAReceived := uint64(0)
				if vaultTokenAToWhirlpool := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[0].PublicKey.String() == parsedAccounts.VaultTokenAAccount.String()
				}); vaultTokenAToWhirlpool != nil {
					tokenASwapped = *vaultTokenAToWhirlpool.Amount
				}
				if vaultTokenAToKeeper := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.DripFeeTokenAAccount.String()
				}); vaultTokenAToKeeper != nil {
					keeperTokenAReceived = *vaultTokenAToKeeper.Amount
				}
				if whirlpoolToVaultTokenB := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.VaultTokenBAccount.String()
				}); whirlpoolToVaultTokenB != nil {
					tokenBReceived = *whirlpoolToVaultTokenB.Amount
				}

				return p.repo.UpsertDripMetric(ctx, &model.DripMetric{
					Signature: signature,
					IxIndex:   int32(i),
					IxName:    ixName,
					IxVersion: 1,
					Slot:      int32(txRaw.Slot),
					Time:      blockTime,
					Vault:     parsedAccounts.Vault.String(),
					// TODO: Remove
					//TokenAMint:                 "",
					//TokenBMint:                 "",
					VaultTokenASwappedAmount:   tokenASwapped,
					VaultTokenBReceivedAmount:  tokenBReceived,
					KeeperTokenAReceivedAmount: keeperTokenAReceived,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV0DripOrcaWhirlpool")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV0WithdrawIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetWithdrawBAccounts()
				userTokenBWithdrawAmount := uint64(0)
				treasuryTokenBReceivedAmount := uint64(0)
				if vaultTokenBToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.UserTokenBAccount.String()
				}); vaultTokenBToUser != nil {
					userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
				}
				if vaultTokenBToTreasury := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.VaultTreasuryTokenBAccount.String()
				}); vaultTokenBToTreasury != nil {
					treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
				}
				return p.repo.UpsertWithdrawMetric(ctx, &model.WithdrawalMetric{
					Signature: signature,
					IxIndex:   int32(i),
					IxName:    ixName,
					IxVersion: 0,
					Slot:      int32(txRaw.Slot),
					Time:      blockTime,
					Vault:     parsedAccounts.Vault.String(),
					// TODO: Remove
					//TokenAMint:                   nil,
					//TokenBMint:                   "",
					UserTokenAWithdrawAmount:     0,
					UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
					TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
					ReferralTokenBReceivedAmount: 0,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV0WithdrawIx")
				return err
			}
			if parsedIx, ixName, err := p.dripIxParser.MaybeParseV0ClosePositionIx(accounts, ix.Data); parsedIx != nil {
				parsedAccounts := parsedIx.GetClosePositionAccounts()
				userTokenAWithdrawAmount := uint64(0)
				userTokenBWithdrawAmount := uint64(0)
				treasuryTokenBReceivedAmount := uint64(0)
				if vaultTokenAToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.UserTokenAAccount.String()
				}); vaultTokenAToUser != nil {
					userTokenAWithdrawAmount = *vaultTokenAToUser.Amount
				}
				if vaultTokenBToUser := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.UserTokenBAccount.String()
				}); vaultTokenBToUser != nil {
					userTokenBWithdrawAmount = *vaultTokenBToUser.Amount
				}
				if vaultTokenBToTreasury := p.findInnerIxTokenTransfer(txRaw, tx.Message, i, log, func(transfer token.Transfer) bool {
					return transfer.Accounts[1].PublicKey.String() == parsedAccounts.VaultTreasuryTokenBAccount.String()
				}); vaultTokenBToTreasury != nil {
					treasuryTokenBReceivedAmount = *vaultTokenBToTreasury.Amount
				}
				return p.repo.UpsertWithdrawMetric(ctx, &model.WithdrawalMetric{
					Signature: signature,
					IxIndex:   int32(i),
					IxName:    ixName,
					IxVersion: 0,
					Slot:      int32(txRaw.Slot),
					Time:      blockTime,
					Vault:     parsedAccounts.Vault.String(),
					// TODO: Remove
					//TokenAMint:                   nil,
					//TokenBMint:                   "",
					UserTokenAWithdrawAmount:     userTokenAWithdrawAmount,
					UserTokenBWithdrawAmount:     userTokenBWithdrawAmount,
					TreasuryTokenBReceivedAmount: treasuryTokenBReceivedAmount,
					ReferralTokenBReceivedAmount: 0,
				})
			} else if err != nil {
				log.WithError(err).Error("failed MaybeParseV0ClosePositionIx")
				return err
			}
		default:
			log.WithField("version", version).Warn("unknown drip program version")
		}

	}
	return nil
}

func (p impl) findInnerIxTokenTransfer(txRaw rpc.TransactionWithMeta, msg solana.Message, ixIndex int, log *logrus.Entry, condition func(transfer token.Transfer) bool) *token.Transfer {
	innerIXS := lo.Filter(txRaw.Meta.InnerInstructions, func(innerIx rpc.InnerInstruction, idx int) bool {
		return innerIx.Index == uint16(ixIndex)
	})
	instructions := lo.Flatten(lo.Map(innerIXS, func(innerInnerIx rpc.InnerInstruction, innerInnerIdx int) []solana.CompiledInstruction {
		return innerInnerIx.Instructions
	}))

	tokenTransfers := lo.FilterMap(instructions, func(ix solana.CompiledInstruction, idx int) (token.Transfer, bool) {
		accounts, err := ix.ResolveInstructionAccounts(&msg)
		if err != nil {
			log.WithError(err).Error("failed to  innerInnerIx.ResolveInstructionAccounts(&tx.Message)")
			return token.Transfer{}, false
		}
		if transfer, _, err := p.dripIxParser.MaybeParseTokenTransfer(accounts, ix.Data); err != nil {
			log.WithError(err).Error("failed to p.dripIxParser.MaybeParseTokenTransfer(accounts, innerInnerIx.Data);")
			return token.Transfer{}, false
		} else {
			return *transfer, condition(*transfer)
		}
	})
	if len(tokenTransfers) == 1 {
		return &tokenTransfers[0]
	}
	log.WithField("len(tokenTransfers)", len(tokenTransfers)).Warn("expected 1 token transfer, but got 0 or more then 1")
	return nil
}

func (p impl) handleTransactionUpdateErr(ctx context.Context, queueItem *model.TransactionUpdateQueueItem, err error) {
	logrus.WithError(err).Error("failed to process update")
	if requeueErr := p.transactionUpdateQueue.ReQueueTransactionUpdateItem(ctx, queueItem); requeueErr != nil {
		logrus.WithError(requeueErr).Error("failed to add item to queue")
	}
}
