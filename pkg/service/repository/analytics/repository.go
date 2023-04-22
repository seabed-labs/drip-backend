package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/jmoiron/sqlx"
)

type AnalyticsRepository interface {
	GetDepositMetricBySignature(ctx context.Context, signature string) (*model.DepositMetric, error)
	GetDripMetricBySignature(ctx context.Context, signature string) (*model.DripMetric, error)
	GetWithdrawalMetricBySignature(ctx context.Context, signature string) (*model.WithdrawalMetric, error)

	GetCurrentTVL(ctx context.Context) (*model.CurrentTVL, error)
	GetUniqueDepositorCount(ctx context.Context) (int64, error)
	GetLifeTimeDepositNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeDeposit, error)
	GetLifeTimeVolumeNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeVolume, error)
	GetLifeTimeWithdrawalNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeWithdrawal, error)
}

type analyticsRepositoryImpl struct {
	repo *query.Query
	db   *sqlx.DB
}

func NewAnalyticsRepository(
	repo *query.Query,
	db *sqlx.DB,
) AnalyticsRepository {
	return analyticsRepositoryImpl{
		repo: repo,
		db:   db,
	}
}

func (d analyticsRepositoryImpl) GetUniqueDepositorCount(ctx context.Context) (int64, error) {
	return d.repo.
		DepositMetric.
		WithContext(ctx).
		Distinct(d.repo.DepositMetric.Depositor).
		Where(d.repo.DepositMetric.Depositor.IsNotNull()).
		Count()
}

const tvlQuery = "select sum(token_usd_value) as total_usd_value from (select token.pubkey as token_mint, (sum(amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value from vault join token_account on vault.token_a_account = token_account.pubkey join token on token_account.mint = token.pubkey where amount != 0 group by token.pubkey\nunion \nselect token.pubkey as token_mint, (sum(amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value from vault join token_account on vault.token_b_account = token_account.pubkey join token on token_account.mint = token.pubkey where amount != 0 group by token.pubkey) as q1;"

func (d analyticsRepositoryImpl) GetCurrentTVL(ctx context.Context) (*model.CurrentTVL, error) {
	tvl := model.CurrentTVL{}
	if err := d.db.GetContext(ctx, &tvl, tvlQuery); err != nil {
		return nil, err
	}
	return &tvl, nil
}

const lifeTimeVolumeNormalizedToCurrentPriceQuery = "select \n\tsum(token_usd_value) as \"total_usd_volume\"\nfrom \n\t(\n\tselect \n\t\ttoken.pubkey as token_mint, (sum(drip_metric.vault_token_a_swapped_amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value \n\tfrom \n\t\tvault \n\tjoin \n\t\tdrip_metric \n\ton \n\t\tvault.pubkey = drip_metric.vault \n\tjoin \n\t\ttoken \n\ton \n\t\tvault.token_a_mint = token.pubkey \n\tgroup by \n\t\ttoken.pubkey\n) as q1\n;"

func (d analyticsRepositoryImpl) GetLifeTimeVolumeNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeVolume, error) {
	volume := model.LifeTimeVolume{}
	if err := d.db.GetContext(ctx, &volume, lifeTimeVolumeNormalizedToCurrentPriceQuery); err != nil {
		return nil, err
	}
	return &volume, nil
}

const lifeTimeDepositNormalizedToCurrentPriceQuery = "select \n\tsum(token_usd_value) as \"total_usd_deposit\"\nfrom \n\t(\n\tselect \n\t\ttoken.pubkey as token_mint, (sum(deposit_metric.token_a_deposit_amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value \n\tfrom \n\t\tdeposit_metric \n\tjoin \n\t\tvault\n\ton\n\t\tdeposit_metric.vault = vault.pubkey\n\tjoin\n\t\ttoken \n\ton \n\t\tvault.token_a_mint = token.pubkey \n\tgroup by \n\t\ttoken.pubkey\n) as q1\n;\n"

func (d analyticsRepositoryImpl) GetLifeTimeDepositNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeDeposit, error) {
	deposit := model.LifeTimeDeposit{}
	if err := d.db.GetContext(ctx, &deposit, lifeTimeDepositNormalizedToCurrentPriceQuery); err != nil {
		return nil, err
	}
	return &deposit, nil
}

const lifeTimeWithdrawalNormalizedToCurrentPriceyQuery = "select \n\tsum(token_usd_value) as \"total_usd_withdrawal\"\nfrom \n\t(\n\tselect \n\t\ttoken.pubkey as token_mint, (sum(withdrawal_metric.user_token_a_withdraw_amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value \n\tfrom \n\t\twithdrawal_metric \n\tjoin\n\t\tvault\n\ton\n\t\twithdrawal_metric.vault = vault.pubkey\n\tjoin \n\t\ttoken \n\ton \n\t\tvault.token_a_mint = token.pubkey \n\tgroup by \n\t\ttoken.pubkey\n\tunion\n\tselect \n\t\ttoken.pubkey as token_mint, (sum(withdrawal_metric.user_token_b_withdraw_amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value \n\tfrom \n\t\twithdrawal_metric \n\tjoin\n\t\tvault\n\ton\n\t\twithdrawal_metric.vault = vault.pubkey\n\tjoin \n\t\ttoken \n\ton \n\t\tvault.token_b_mint = token.pubkey \n\tgroup by \n\t\ttoken.pubkey\n) as q1\n;"

func (d analyticsRepositoryImpl) GetLifeTimeWithdrawalNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeWithdrawal, error) {
	withdrawal := model.LifeTimeWithdrawal{}
	if err := d.db.GetContext(ctx, &withdrawal, lifeTimeWithdrawalNormalizedToCurrentPriceyQuery); err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

func (d analyticsRepositoryImpl) GetWithdrawalMetricBySignature(ctx context.Context, signature string) (*model.WithdrawalMetric, error) {
	return d.repo.
		WithdrawalMetric.
		WithContext(ctx).
		Where(d.repo.WithdrawalMetric.Signature.Eq(signature)).
		First()
}

func (d analyticsRepositoryImpl) GetDepositMetricBySignature(ctx context.Context, signature string) (*model.DepositMetric, error) {
	return d.repo.
		DepositMetric.
		WithContext(ctx).
		Where(d.repo.DepositMetric.Signature.Eq(signature)).
		First()
}

func (d analyticsRepositoryImpl) GetDripMetricBySignature(ctx context.Context, signature string) (*model.DripMetric, error) {
	return d.repo.
		DripMetric.
		WithContext(ctx).
		Where(d.repo.DripMetric.Signature.Eq(signature)).
		First()
}
