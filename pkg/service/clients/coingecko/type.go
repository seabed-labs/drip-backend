package coingecko

const (
	callsPerSecond        = 1
	baseUrl               = "https://api.coingecko.com/api/v3"
	coinsMarketsPath      = "/coins/markets"
	coinsListPath         = "/coins/list"
	solanaContractPath    = "/coins/solana/contract"
	CoinsMarketsPathLimit = 100

	//CacheCoinsMarketsPath = coinsMarketsPath
	cacheSolanaContractPath = solanaContractPath
	cacheCoinsListPath      = coinsListPath
)

type CoinsListResponse []CoinResponse

type CoinResponse struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Platforms struct {
		Solana *string `json:"solana,omitempty"`
	} `json:"platforms"`
}

type CoinGeckoMetadataResponse struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Image  struct {
		Small *string `json:"small"`
	} `json:"image"`
	ContractAddress     string   `json:"contract_address"`
	MarketCapRank       *int32   `json:"market_cap_rank"`
	CoingeckoRank       *int32   `json:"coingecko_rank"`
	CoingeckoScore      *float64 `json:"coingecko_score"`
	DeveloperScore      *float64 `json:"developer_score"`
	CommunityScore      *float64 `json:"community_score"`
	LiquidityScore      *float64 `json:"liquidity_score"`
	PublicInterestScore *float64 `json:"public_interest_score"`
	MarketData          struct {
		CurrentPrice map[string]float64 `json:"current_price"`
	} `json:"market_data"`
}
type CoinGeckoTokensMarketPriceResponse []CoinGeckoTokenMarketPriceResponse

type CoinGeckoTokenMarketPriceResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	CurrentPrice  float64 `json:"current_price"`
	MarketCapRank *int32  `json:"market_cap_rank"`
}
