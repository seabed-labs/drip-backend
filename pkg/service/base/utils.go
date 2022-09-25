package base

import "github.com/dcaf-labs/drip/pkg/service/repository/model"

func filterVaultsWithZeroDripAmount(vaults []*model.Vault) []*model.Vault {
	res := []*model.Vault{}
	for i := range vaults {
		if vaults[i].DripAmount != 0 {
			res = append(res, vaults[i])
		}
	}
	return res
}
