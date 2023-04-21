package model

type CurrentTVL struct {
	TotalUSDValue float64 `json:"totalUsdValue" db:"total_usd_value"`
}

type LifeTimeVolume struct {
	TotalUSDVolume float64 `json:"totalUsdVolume" db:"total_usd_volume"`
}

type LifeTimeDeposit struct {
	TotalUSDDeposit float64 `json:"totalUsdDeposit" db:"total_usd_deposit"`
}

type LifeTimeWithdrawal struct {
	TotalUSDWithdrawal float64 `json:"totalUsdWithdrawal" db:"total_usd_withdrawal"`
}
