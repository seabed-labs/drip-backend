package repository

//
//import (
//	"context"
//
//	"github.com/volatiletech/sqlboiler/v4/boil"
//
//	"github.com/dcaf-protocol/drip/internal/pkg/repository/models"
//	"github.com/jmoiron/sqlx"
//)
//
//type VaultRepository interface {
//	GetVault(context.Context, string) (*models.Vault, error)
//	UpsertVault(context.Context, models.Vault) error
//}
//
//func NewVaultRepository(
//	db *sqlx.DB,
//) (VaultRepository, error) {
//	return newVaultRepositoryImpl(db)
//}
//
//type vaultRepositoryImpl struct {
//	db *sqlx.DB
//}
//
//func newVaultRepositoryImpl(
//	db *sqlx.DB,
//) (vaultRepositoryImpl, error) {
//	return vaultRepositoryImpl{db: db}, nil
//}
//
//func (repo vaultRepositoryImpl) GetVault(ctx context.Context, pubkey string) (*models.Vault, error) {
//	return models.FindVault(ctx, repo.db, pubkey)
//}
//
//func (repo vaultRepositoryImpl) UpsertVault(ctx context.Context, vault models.Vault) error {
//	//boil.Whitelist(
//	//	models.VaultColumns.ProtoConfig,
//	//	models.VaultColumns.TokenAMint,
//	//	models.VaultColumns.TokenBMint,
//	//	models.VaultColumns.TokenAAccount,
//	//	models.VaultColumns.TokenBAccount,
//	//	models.VaultColumns.TreasuryTokenBAccount,
//	//	models.VaultColumns.LastDcaPeriod,
//	//	models.VaultColumns.DripAmount,
//	//	models.VaultColumns.DcaActivationTimestamp,
//	//),
//	return vault.Upsert(
//		ctx, repo.db, true,
//		[]string{models.VaultColumns.Pubkey},
//		boil.Blacklist(
//			models.VaultColumns.Pubkey,
//		),
//		boil.Infer())
//}
