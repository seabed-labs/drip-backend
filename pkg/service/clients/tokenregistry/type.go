package tokenregistry

const (
	callsPerSecond = 10
	url            = "https://cdn.jsdelivr.net/gh/solana-labs/token-list@latest/src/tokens/solana.tokenlist.json"
)

type Token struct {
	ChainID    int      `json:"chainId"`
	Address    string   `json:"address"`
	Symbol     string   `json:"symbol"`
	Name       string   `json:"name"`
	Decimals   int      `json:"decimals"`
	LogoURI    string   `json:"logoURI"`
	Tags       []string `json:"tags,omitempty"`
	Extensions struct {
		Facebook string `json:"facebook"`
		Twitter  string `json:"twitter"`
		Website  string `json:"website"`
	} `json:"extensions,omitempty"`
}

type TokenRegistryResponse struct {
	Tokens []Token `json:"tokens"`
}
