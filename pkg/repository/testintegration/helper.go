package testintegration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/test-go/testify/assert"
)

func truncateDB(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("truncate proto_config cascade")
	assert.NoError(t, err)
	_, err = db.Exec("truncate token_pair cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate token cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate vault cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate vault_period cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate position cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate token_account_balance cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate token_price cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate token_swap cascade")
	assert.NoError(t, err)
	_, err = db.Exec(" truncate user_position cascade")
	assert.NoError(t, err)
}

/////////////////////////////////////////////////////////////////////////////
// Seed
/////////////////////////////////////////////////////////////////////////////
type seedTokensParams struct {
	tokenAPubkey *string
	tokenBPubkey *string
}

type seedTokensResult struct {
	tokenAPubkey string
	tokenBPubkey string
}

//nolint:funlen
func seedTokens(t *testing.T, db *sqlx.DB, params seedTokensParams) seedTokensResult {
	seedTokensResult := seedTokensResult{}

	if params.tokenAPubkey == nil {
		seedTokensResult.tokenAPubkey = uuid.New().String()
		_, err := db.Exec(
			`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4)`,
			seedTokensResult.tokenAPubkey, "btc", 8, nil,
		)
		assert.NoError(t, err)
	} else {
		seedTokensResult.tokenAPubkey = *params.tokenAPubkey
	}

	if params.tokenBPubkey == nil {
		seedTokensResult.tokenBPubkey = uuid.New().String()
		_, err := db.Exec(
			`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4)`,
			seedTokensResult.tokenBPubkey, "eth", 18, nil,
		)
		assert.NoError(t, err)
	} else {
		seedTokensResult.tokenBPubkey = *params.tokenBPubkey
	}

	return seedTokensResult
}

type seedTokenPairParams struct {
	seedTokensParams
	tokenPairID *string
}

type seedTokenPairResult struct {
	seedTokensResult
	tokenPairID string
}

//nolint:funlen
func seedTokenPair(t *testing.T, db *sqlx.DB, params seedTokenPairParams) seedTokenPairResult {
	seedTokenPairResult := seedTokenPairResult{}
	seedTokenPairResult.seedTokensResult = seedTokens(t, db, params.seedTokensParams)

	if params.tokenPairID == nil {
		seedTokenPairResult.tokenPairID = uuid.New().String()
		_, err := db.Exec(
			`insert into 
    						token_pair(id, token_a, token_b) 
							values
							    ($1, $2, $3)`,
			seedTokenPairResult.tokenPairID, seedTokenPairResult.tokenAPubkey, seedTokenPairResult.tokenBPubkey,
		)
		assert.NoError(t, err)
	} else {
		seedTokenPairResult.tokenPairID = *params.tokenPairID
	}

	return seedTokenPairResult
}

type seedProtoConfigParams struct {
	protoConfigPubkey *string
}

type seedProtoConfigResult struct {
	protoConfigPubkey string
}

//nolint:funlen
func seedProtoConfig(t *testing.T, db *sqlx.DB, params seedProtoConfigParams) seedProtoConfigResult {
	seedProtoConfigReult := seedProtoConfigResult{}
	if params.protoConfigPubkey == nil {
		seedProtoConfigReult.protoConfigPubkey = uuid.New().String()
		_, err := db.Exec(
			`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
			seedProtoConfigReult.protoConfigPubkey, 1, 2, 3)
		assert.NoError(t, err)
	} else {
		seedProtoConfigReult.protoConfigPubkey = *params.protoConfigPubkey
	}

	return seedProtoConfigReult
}

type seedVaultParams struct {
	seedProtoConfigParams
	seedTokenPairParams
	protoConfigPubkey *string

	tokenAAcount    *string
	tokenBAccount   *string
	treasuryAccount *string
}

type seedVaultResult struct {
	seedProtoConfigResult
	seedTokenPairResult
	vaultPubkey string

	tokenAAcount    string
	tokenBAccount   string
	treasuryAccount string
}

//nolint:funlen
func seedVault(t *testing.T, db *sqlx.DB, params seedVaultParams) seedVaultResult {
	seedVaultResult := seedVaultResult{}
	seedVaultResult.seedProtoConfigResult = seedProtoConfig(t, db, params.seedProtoConfigParams)
	seedVaultResult.seedTokenPairResult = seedTokenPair(t, db, params.seedTokenPairParams)

	if params.tokenAAcount == nil {
		seedVaultResult.tokenAAcount = uuid.New().String()
	} else {
		seedVaultResult.tokenAAcount = *params.tokenAAcount
	}
	if params.tokenBAccount == nil {
		seedVaultResult.tokenBAccount = uuid.New().String()
	} else {
		seedVaultResult.tokenBAccount = *params.tokenBAccount
	}
	if params.treasuryAccount == nil {
		seedVaultResult.treasuryAccount = uuid.New().String()
	} else {
		seedVaultResult.treasuryAccount = *params.treasuryAccount
	}

	seedVaultResult.vaultPubkey = uuid.New().String()
	_, err := db.Exec(
		`insert into 
    						vault(pubkey, proto_config, token_a_account, token_b_account, treasury_token_b_account, last_dca_period, drip_amount, dca_activation_timestamp, enabled, token_pair_id) 
							values
							    ($1, $2, $3, $4,$5, $6,$7,$8,$9,$10)`,
		seedVaultResult.vaultPubkey, seedVaultResult.protoConfigPubkey, seedVaultResult.tokenAAcount, seedVaultResult.tokenBAccount, seedVaultResult.treasuryAccount, 0, 0, time.Time{}, false, seedVaultResult.tokenPairID,
	)
	assert.NoError(t, err)
	return seedVaultResult
}

type seedVaultPeriodParams struct {
	seedVaultParams
	vaultPubkey *string
}

type seedVaultPeriodResult struct {
	seedVaultResult
	vaultPeriodPubkey string
}

func seedVaultPeriod(t *testing.T, db *sqlx.DB, params seedVaultPeriodParams) seedVaultPeriodResult {
	seedVaultPeriodResult := seedVaultPeriodResult{
		seedVaultResult: seedVault(t, db, params.seedVaultParams),
	}
	seedVaultPeriodResult.vaultPeriodPubkey = uuid.New().String()
	_, err := db.Exec(
		`insert into 
    						vault_period(pubkey, vault, period_id, twap, dar) 
							values($1, $2, $3, $4, $5)`,
		seedVaultPeriodResult.vaultPeriodPubkey, seedVaultPeriodResult.vaultPubkey, 0, 0, 0)
	assert.NoError(t, err)
	return seedVaultPeriodResult
}
