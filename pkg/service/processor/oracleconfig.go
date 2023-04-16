package processor

//func (p impl) UpsertOracleConfigByAddress(ctx context.Context, address string) error {
//	var oracleConfig drip.OracleConfig
//	if err := p.solanaClient.GetAccount(ctx, address, &oracleConfig); err != nil {
//		return err
//	}
//	return p.repo.UpsertOracleConfigs(ctx, &model.OracleConfig{
//		Pubkey:          address,
//		Enabled:         oracleConfig.Enabled,
//		Source:          int16(oracleConfig.Source),
//		UpdateAuthority: oracleConfig.UpdateAuthority.String(),
//		TokenAMint:      oracleConfig.TokenAMint.String(),
//		TokenAPrice:     oracleConfig.TokenAPrice.String(),
//		TokenBMint:      oracleConfig.TokenBMint.String(),
//		TokenBPrice:     oracleConfig.TokenBPrice.String(),
//	})
//}
//
//func (p impl) ensureOracleConfig(ctx context.Context, address string) (*model.OracleConfig, error) {
//	oracleConfig, err := p.repo.GetOracleConfigByAddress(ctx, address)
//	if err != nil && err.Error() == repository.ErrRecordNotFound {
//		if err := p.UpsertOracleConfigByAddress(ctx, address); err != nil {
//			return nil, err
//		}
//		return p.repo.GetOracleConfigByAddress(ctx, address)
//	}
//	return oracleConfig, nil
//}
