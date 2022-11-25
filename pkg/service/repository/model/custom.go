package model

type CurrentTVL struct {
	TotalUSDValue float64 `json:"totalUsdValue" db:"total_usd_value"`
}
