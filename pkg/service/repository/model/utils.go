package model

func GetTokenSwapsByTokenPairID(tokenSwaps []*TokenSwap) map[string][]*TokenSwap {
	res := make(map[string][]*TokenSwap)
	for i := range tokenSwaps {
		if _, ok := res[tokenSwaps[i].TokenPairID]; !ok {
			res[tokenSwaps[i].TokenPairID] = []*TokenSwap{}
		}
		res[tokenSwaps[i].TokenPairID] = append(res[tokenSwaps[i].TokenPairID], tokenSwaps[i])
	}
	return res
}

func GetTokenAccountBalancesByPubkey(balances []*TokenAccountBalance) map[string]*TokenAccountBalance {
	res := make(map[string]*TokenAccountBalance)
	for i := range balances {
		res[balances[i].Pubkey] = balances[i]
	}
	return res
}

func GetVaultWhitelistsByVault(whitelists []*VaultWhitelist) map[string][]*VaultWhitelist {
	res := make(map[string][]*VaultWhitelist)
	for i := range whitelists {
		if _, ok := res[whitelists[i].VaultPubkey]; !ok {
			res[whitelists[i].VaultPubkey] = []*VaultWhitelist{}
		}
		res[whitelists[i].VaultPubkey] = append(res[whitelists[i].VaultPubkey], whitelists[i])
	}
	return res
}

func GetVaultWhitelistsBySwap(whitelists []*VaultWhitelist) map[string][]*VaultWhitelist {
	res := make(map[string][]*VaultWhitelist)
	for i := range whitelists {
		if _, ok := res[whitelists[i].TokenSwapPubkey]; !ok {
			res[whitelists[i].TokenSwapPubkey] = []*VaultWhitelist{}
		}
		res[whitelists[i].TokenSwapPubkey] = append(res[whitelists[i].TokenSwapPubkey], whitelists[i])
	}
	return res
}

func GetTokenAccountPubkeysForTokenSwaps(tokenSwaps []*TokenSwap) []string {
	res := []string{}
	for i := range tokenSwaps {
		res = append(res, tokenSwaps[i].TokenAAccount, tokenSwaps[i].TokenBAccount)
	}
	return res
}

func GetVaultPubkeys(vaults []*Vault) []string {
	var vaultPubkeys []string
	for i := range vaults {
		vaultPubkeys = append(vaultPubkeys, vaults[i].Pubkey)
	}
	return vaultPubkeys
}

func GetVaultsByPubkey(vaults []*Vault) map[string]*Vault {
	res := make(map[string]*Vault)
	for i := range vaults {
		res[vaults[i].Pubkey] = vaults[i]
	}
	return res
}

func GetTokenPairIDsForVaults(vaults []*Vault) []string {
	var tokenPairIDs []string
	for i := range vaults {
		tokenPairIDs = append(tokenPairIDs, vaults[i].TokenPairID)
	}
	return tokenPairIDs
}
