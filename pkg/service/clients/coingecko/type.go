package coingecko

const (
	callsPerSecond        = 10
	baseUrl               = "https://api.coingecko.com/api/v3"
	coinsMarketsPath      = "/coins/markets"
	CoinsMarketsPathLimit = 100
)

type CoinGeckoMetadataResponse struct {
	ID     string
	Symbol string
	Name   string
}
type CoinGeckoTokensMarketPriceResponse []CoinGeckoTokenMarketPriceResponse

type CoinGeckoTokenMarketPriceResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	CurrentPrice  float64 `json:"current_price"`
	MarketCapRank *int32  `json:"market_cap_rank"`
}
