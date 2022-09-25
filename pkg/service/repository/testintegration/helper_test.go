package testintegration

//
//import (
//	"testing"
//	"time"
//
//	"github.com/gagliardetto/solana-go"
//
//	"github.com/google/uuid"
//	"github.com/jmoiron/sqlx"
//	"github.com/test-go/testify/assert"
//)
//
//func truncateDB(t *testing.T, db *sqlx.DB) {
//	_, err := db.Exec("truncate proto_config cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec("truncate token_pair cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate token cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate vault cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate vault_period cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate position cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate token_account_balance cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate token_price cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate token_swap cascade")
//	assert.NoError(t, err)
//	_, err = db.Exec(" truncate user_position cascade")
//	assert.NoError(t, err)
//}
//
///////////////////////////////////////////////////////////////////////////////
//// Seed
///////////////////////////////////////////////////////////////////////////////
//type seedTokensParams struct {
//	tokenAPubkey *string
//	tokenBPubkey *string
//}
//
//type seedTokensResult struct {
//	tokenAPubkey string
//	tokenBPubkey string
//}
//
////nolint:funlen
//func seedTokens(t *testing.T, db *sqlx.DB, params seedTokensParams) seedTokensResult {
//	seedTokensResult := seedTokensResult{}
//
//	if params.tokenAPubkey == nil {
//		seedTokensResult.tokenAPubkey = uuid.New().String()
//		_, err := db.Exec(
//			`insert into
//    						token(pubkey, symbol, decimals, icon_url)
//							values
//							    ($1, $2, $3, $4)`,
//			seedTokensResult.tokenAPubkey, "btc", 8, nil,
//		)
//		assert.NoError(t, err)
//	} else {
//		seedTokensResult.tokenAPubkey = *params.tokenAPubkey
//	}
//
//	if params.tokenBPubkey == nil {
//		seedTokensResult.tokenBPubkey = uuid.New().String()
//		_, err := db.Exec(
//			`insert into
//    						token(pubkey, symbol, decimals, icon_url)
//							values
//							    ($1, $2, $3, $4)`,
//			seedTokensResult.tokenBPubkey, "eth", 18, nil,
//		)
//		assert.NoError(t, err)
//	} else {
//		seedTokensResult.tokenBPubkey = *params.tokenBPubkey
//	}
//
//	return seedTokensResult
//}
//
//type seedTokenPairParams struct {
//	seedTokensParams
//	tokenPairID *string
//}
//
//type seedTokenPairResult struct {
//	seedTokensResult
//	tokenPairID string
//}
//
////nolint:funlen
//func seedTokenPair(t *testing.T, db *sqlx.DB, params seedTokenPairParams) seedTokenPairResult {
//	seedTokenPairResult := seedTokenPairResult{}
//	seedTokenPairResult.seedTokensResult = seedTokens(t, db, params.seedTokensParams)
//
//	if params.tokenPairID == nil {
//		seedTokenPairResult.tokenPairID = uuid.New().String()
//		_, err := db.Exec(
//			`insert into
//    						token_pair(id, token_a, token_b)
//							values
//							    ($1, $2, $3)`,
//			seedTokenPairResult.tokenPairID, seedTokenPairResult.tokenAPubkey, seedTokenPairResult.tokenBPubkey,
//		)
//		assert.NoError(t, err)
//	} else {
//		seedTokenPairResult.tokenPairID = *params.tokenPairID
//	}
//
//	return seedTokenPairResult
//}
//
//type seedProtoConfigParams struct {
//	protoConfigPubkey *string
//}
//
//type seedProtoConfigResult struct {
//	protoConfigPubkey string
//}
//
////nolint:funlen
//func seedProtoConfig(t *testing.T, db *sqlx.DB, params seedProtoConfigParams) seedProtoConfigResult {
//	seedProtoConfigReult := seedProtoConfigResult{}
//	if params.protoConfigPubkey == nil {
//		seedProtoConfigReult.protoConfigPubkey = uuid.New().String()
//		_, err := db.Exec(
//			`insert into
//    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread)
//							values($1, $2, $3, $4)`,
//			seedProtoConfigReult.protoConfigPubkey, 1, 2, 3)
//		assert.NoError(t, err)
//	} else {
//		seedProtoConfigReult.protoConfigPubkey = *params.protoConfigPubkey
//	}
//
//	return seedProtoConfigReult
//}
//
//type seedVaultParams struct {
//	seedProtoConfigParams
//	seedTokenPairParams
//
//	tokenAAcount    *string
//	tokenBAccount   *string
//	treasuryAccount *string
//}
//
//type seedVaultResult struct {
//	seedProtoConfigResult
//	seedTokenPairResult
//	vaultPubkey string
//
//	tokenAAcount    string
//	tokenBAccount   string
//	treasuryAccount string
//}
//
////nolint:funlen
//func seedVault(t *testing.T, db *sqlx.DB, params seedVaultParams) seedVaultResult {
//	seedVaultResult := seedVaultResult{}
//	seedVaultResult.seedProtoConfigResult = seedProtoConfig(t, db, params.seedProtoConfigParams)
//	seedVaultResult.seedTokenPairResult = seedTokenPair(t, db, params.seedTokenPairParams)
//
//	if params.tokenAAcount == nil {
//		seedVaultResult.tokenAAcount = uuid.New().String()
//	} else {
//		seedVaultResult.tokenAAcount = *params.tokenAAcount
//	}
//	if params.tokenBAccount == nil {
//		seedVaultResult.tokenBAccount = uuid.New().String()
//	} else {
//		seedVaultResult.tokenBAccount = *params.tokenBAccount
//	}
//	if params.treasuryAccount == nil {
//		seedVaultResult.treasuryAccount = uuid.New().String()
//	} else {
//		seedVaultResult.treasuryAccount = *params.treasuryAccount
//	}
//
//	seedVaultResult.vaultPubkey = uuid.New().String()
//	_, err := db.Exec(
//		`insert into
//			vault(pubkey, proto_config, token_a_account, token_b_account, treasury_token_b_account, last_dca_period, drip_amount, dca_activation_timestamp, enabled, token_pair_id)
//			values
//				($1, $2, $3, $4,$5, $6,$7,$8,$9,$10)`,
//		seedVaultResult.vaultPubkey, seedVaultResult.protoConfigPubkey, seedVaultResult.tokenAAcount, seedVaultResult.tokenBAccount, seedVaultResult.treasuryAccount, 0, 0, time.Time{}, false, seedVaultResult.tokenPairID,
//	)
//	assert.NoError(t, err)
//	return seedVaultResult
//}
//
//type seedVaultPeriodParams struct {
//	seedVaultParams
//	vaultPubkey *string
//}
//
//type seedVaultPeriodResult struct {
//	seedVaultResult
//	vaultPeriodPubkey string
//}
//
//func seedVaultPeriod(t *testing.T, db *sqlx.DB, params seedVaultPeriodParams) seedVaultPeriodResult {
//	seedVaultPeriodResult := seedVaultPeriodResult{
//		seedVaultResult: seedVault(t, db, params.seedVaultParams),
//	}
//	seedVaultPeriodResult.vaultPeriodPubkey = uuid.New().String()
//	_, err := db.Exec(
//		`insert into
//    						vault_period(pubkey, vault, period_id, twap, dar)
//							values($1, $2, $3, $4, $5)`,
//		seedVaultPeriodResult.vaultPeriodPubkey, seedVaultPeriodResult.vaultPubkey, 0, 0, 0)
//	assert.NoError(t, err)
//	return seedVaultPeriodResult
//}
//
//type seedPositionParams struct {
//	seedVaultParams
//}
//
//type seedPositionResult struct {
//	seedVaultResult
//	positionPubkey string
//}
//
//func seedPosition(t *testing.T, db *sqlx.DB, params seedPositionParams) seedPositionResult {
//	seedVaultPeriodResult := seedPositionResult{
//		seedVaultResult: seedVault(t, db, params.seedVaultParams),
//	}
//	seedVaultPeriodResult.positionPubkey = uuid.New().String()
//	_, err := db.Exec(
//		`insert into
//			position(pubkey, vault, authority, deposited_token_a_amount, withdrawn_token_b_amount, deposit_timestamp, dca_period_id_before_deposit, number_of_swaps, periodic_drip_amount, is_closed)
//			values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
//		seedVaultPeriodResult.positionPubkey, seedVaultPeriodResult.vaultPubkey, uuid.New().String(), 0, 0, time.Time{}, 0, 0, 0, false)
//	assert.NoError(t, err)
//	return seedVaultPeriodResult
//}
//
//type seedTokenSwapParams struct {
//	seedTokenPairParams
//	tokenSwapPubkey *string
//}
//
//type seedTokenSwapResult struct {
//	seedTokenPairResult
//	tokenSwapID     string
//	tokenSwapPubkey string
//}
//
//func seedTokenSwap(t *testing.T, db *sqlx.DB, params seedTokenSwapParams) seedTokenSwapResult {
//	seedTokenSwapResult := seedTokenSwapResult{
//		seedTokenPairResult: seedTokenPair(t, db, params.seedTokenPairParams),
//	}
//
//	if params.tokenSwapPubkey == nil {
//		seedTokenSwapResult.tokenSwapPubkey = solana.NewWallet().PublicKey().String()
//		seedTokenSwapResult.tokenSwapID = uuid.New().String()
//		_, err := db.Exec(
//			`insert into
//			token_swap(id, pubkey, mint, authority, fee_account, token_a_account, token_b_account, token_pair_id, token_a_mint, token_b_mint)
//			values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
//			seedTokenSwapResult.tokenSwapID, seedTokenSwapResult.tokenSwapPubkey,
//			solana.NewWallet().PublicKey().String(),
//			solana.NewWallet().PublicKey().String(),
//			solana.NewWallet().PublicKey().String(),
//			solana.NewWallet().PublicKey().String(),
//			solana.NewWallet().PublicKey().String(),
//			seedTokenSwapResult.tokenPairID,
//			solana.NewWallet().PublicKey().String(),
//			solana.NewWallet().PublicKey().String())
//		assert.NoError(t, err)
//	} else {
//		seedTokenSwapResult.tokenSwapPubkey = *params.tokenSwapPubkey
//		seedTokenSwapResult.tokenSwapID = uuid.New().String()
//	}
//	return seedTokenSwapResult
//}
//
//type seedTokenAccountBalanceParams struct {
//	seedTokenPairParams
//	tokenAccountBalancePubkey *string
//}
//
//type seedTokenAccountBalanceResult struct {
//	seedTokenPairResult
//	tokenAccountBalancePubkey string
//}
//
//func seedTokenAccountBalance(t *testing.T, db *sqlx.DB, params seedTokenAccountBalanceParams) seedTokenAccountBalanceResult {
//	seedTokenAccountBalanceResult := seedTokenAccountBalanceResult{
//		seedTokenPairResult: seedTokenPair(t, db, params.seedTokenPairParams),
//	}
//
//	if params.tokenAccountBalancePubkey == nil {
//		seedTokenAccountBalanceResult.tokenAccountBalancePubkey = solana.NewWallet().PublicKey().String()
//		_, err := db.Exec(
//			`insert into
//				token_account_balance(pubkey, mint, owner, amount, state)
//				values($1, $2, $3, $4, $5)`,
//			seedTokenAccountBalanceResult.tokenAccountBalancePubkey,
//			seedTokenAccountBalanceResult.tokenAPubkey,
//			solana.NewWallet().PublicKey().String(),
//			9,
//			"initialized")
//		assert.NoError(t, err)
//	} else {
//		seedTokenAccountBalanceResult.tokenAccountBalancePubkey = *params.tokenAccountBalancePubkey
//	}
//	return seedTokenAccountBalanceResult
//}