package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

const tvlQuery = "select sum(token_usd_value) as total_usd_value from (select token.pubkey as token_mint, (sum(amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value from vault join token_account on vault.token_a_account = token_account.pubkey join token on token_account.mint = token.pubkey where amount != 0 group by token.pubkey\nunion \nselect token.pubkey as token_mint, (sum(amount)*token.ui_market_price_usd)/power(10, token.decimals) as token_usd_value from vault join token_account on vault.token_b_account = token_account.pubkey join token on token_account.mint = token.pubkey where amount != 0 group by token.pubkey) as q1;"

func (d repositoryImpl) GetCurrentTVL(ctx context.Context) (*model.CurrentTVL, error) {
	tvl := model.CurrentTVL{}
	if err := d.db.Get(&tvl, tvlQuery); err != nil {
		return nil, err
	}
	return &tvl, nil
}
