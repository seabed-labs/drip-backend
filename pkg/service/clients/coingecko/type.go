package coinGecko

const (
	callsPerSecond = 10
	baseUrl        = "https://api.coingecko.com/api/v3"
)

type CoinGeckoMetadataResponse struct {
	Id     string
	Symbol string
	Name   string
}
