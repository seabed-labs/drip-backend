// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

func Use(db *gorm.DB) *Query {
	return &Query{
		db:                       db,
		AccountUpdateQueueItem:   newAccountUpdateQueueItem(db),
		OrcaWhirlpool:            newOrcaWhirlpool(db),
		OrcaWhirlpoolDeltaBQuote: newOrcaWhirlpoolDeltaBQuote(db),
		Position:                 newPosition(db),
		ProtoConfig:              newProtoConfig(db),
		SchemaMigration:          newSchemaMigration(db),
		SourceReference:          newSourceReference(db),
		Token:                    newToken(db),
		TokenAccountBalance:      newTokenAccountBalance(db),
		TokenPair:                newTokenPair(db),
		TokenSwap:                newTokenSwap(db),
		Vault:                    newVault(db),
		VaultPeriod:              newVaultPeriod(db),
		VaultWhitelist:           newVaultWhitelist(db),
	}
}

type Query struct {
	db *gorm.DB

	AccountUpdateQueueItem   accountUpdateQueueItem
	OrcaWhirlpool            orcaWhirlpool
	OrcaWhirlpoolDeltaBQuote orcaWhirlpoolDeltaBQuote
	Position                 position
	ProtoConfig              protoConfig
	SchemaMigration          schemaMigration
	SourceReference          sourceReference
	Token                    token
	TokenAccountBalance      tokenAccountBalance
	TokenPair                tokenPair
	TokenSwap                tokenSwap
	Vault                    vault
	VaultPeriod              vaultPeriod
	VaultWhitelist           vaultWhitelist
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                       db,
		AccountUpdateQueueItem:   q.AccountUpdateQueueItem.clone(db),
		OrcaWhirlpool:            q.OrcaWhirlpool.clone(db),
		OrcaWhirlpoolDeltaBQuote: q.OrcaWhirlpoolDeltaBQuote.clone(db),
		Position:                 q.Position.clone(db),
		ProtoConfig:              q.ProtoConfig.clone(db),
		SchemaMigration:          q.SchemaMigration.clone(db),
		SourceReference:          q.SourceReference.clone(db),
		Token:                    q.Token.clone(db),
		TokenAccountBalance:      q.TokenAccountBalance.clone(db),
		TokenPair:                q.TokenPair.clone(db),
		TokenSwap:                q.TokenSwap.clone(db),
		Vault:                    q.Vault.clone(db),
		VaultPeriod:              q.VaultPeriod.clone(db),
		VaultWhitelist:           q.VaultWhitelist.clone(db),
	}
}

type queryCtx struct {
	AccountUpdateQueueItem   *accountUpdateQueueItemDo
	OrcaWhirlpool            *orcaWhirlpoolDo
	OrcaWhirlpoolDeltaBQuote *orcaWhirlpoolDeltaBQuoteDo
	Position                 *positionDo
	ProtoConfig              *protoConfigDo
	SchemaMigration          *schemaMigrationDo
	SourceReference          *sourceReferenceDo
	Token                    *tokenDo
	TokenAccountBalance      *tokenAccountBalanceDo
	TokenPair                *tokenPairDo
	TokenSwap                *tokenSwapDo
	Vault                    *vaultDo
	VaultPeriod              *vaultPeriodDo
	VaultWhitelist           *vaultWhitelistDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		AccountUpdateQueueItem:   q.AccountUpdateQueueItem.WithContext(ctx),
		OrcaWhirlpool:            q.OrcaWhirlpool.WithContext(ctx),
		OrcaWhirlpoolDeltaBQuote: q.OrcaWhirlpoolDeltaBQuote.WithContext(ctx),
		Position:                 q.Position.WithContext(ctx),
		ProtoConfig:              q.ProtoConfig.WithContext(ctx),
		SchemaMigration:          q.SchemaMigration.WithContext(ctx),
		SourceReference:          q.SourceReference.WithContext(ctx),
		Token:                    q.Token.WithContext(ctx),
		TokenAccountBalance:      q.TokenAccountBalance.WithContext(ctx),
		TokenPair:                q.TokenPair.WithContext(ctx),
		TokenSwap:                q.TokenSwap.WithContext(ctx),
		Vault:                    q.Vault.WithContext(ctx),
		VaultPeriod:              q.VaultPeriod.WithContext(ctx),
		VaultWhitelist:           q.VaultWhitelist.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	return &QueryTx{q.clone(q.db.Begin(opts...))}
}

type QueryTx struct{ *Query }

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
