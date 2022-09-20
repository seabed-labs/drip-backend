package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dave/jennifer/jen"
)

type allWhirlpoolsResponse []whirlpoolConfig

type whirlpoolConfig struct {
	Address         string  `json:"address"`
	Whitelisted     bool    `json:"whitelisted"`
	TokenMintA      string  `json:"tokenMintA"`
	TokenMintB      string  `json:"tokenMintB"`
	TickSpacing     int     `json:"tickSpacing"`
	Price           float64 `json:"price"`
	LpsFeeRate      float64 `json:"lpsFeeRate"`
	ProtocolFeeRate float64 `json:"protocolFeeRate"`
	PriceHistory    struct {
		Day struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"day"`
		Week struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"week"`
		Month struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"month"`
	} `json:"priceHistory"`
	TokenAPriceUSD struct {
		Price     float64 `json:"price"`
		Dex       float64 `json:"dex"`
		Coingecko float64 `json:"coingecko"`
	} `json:"tokenAPriceUSD"`
	TokenBPriceUSD struct {
		Price     float64 `json:"price"`
		Dex       float64 `json:"dex"`
		Coingecko float64 `json:"coingecko"`
	} `json:"tokenBPriceUSD"`
	Tvl float64 `json:"tvl"`
}

func main() {
	resp, err := http.Get("https://mainnet-zp2-v2.orca.so/pools")
	if err != nil {
		log.Panic(err.Error())
	}
	defer resp.Body.Close()
	var allPools allWhirlpoolsResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err.Error())
	}
	if err = json.Unmarshal(body, &allPools); err != nil {
		log.Panic(err.Error())
	}

	f := jen.NewFile("controller")
	f.Comment("Code generated by cmd/whirlpoolcodegen DO NOT EDIT.").Line().Line()
	f.Var().Id("mainnetOrcaWhirlpools").Op("=").Index().String().ValuesFunc(func(g *jen.Group) {
		for _, config := range allPools {
			g.Line().Lit(config.Address)
			g.Line().Comment(fmt.Sprintf("MintA: %s, MintB: %s", config.TokenMintA, config.TokenMintB))
		}
		g.Line()
	})
	f.Var().Id("mainnetOrcaWhirlpoolsMap").Op("=").Map(jen.String()).Bool().Values(jen.DictFunc(func(d jen.Dict) {
		for _, config := range allPools {
			d[jen.Lit(config.Address)] = jen.Lit(true).Op(",").Op(fmt.Sprintf("// MintA: %s, MintB: %s", config.TokenMintA, config.TokenMintB))
		}
	}))
	if err := f.Save("./pkg/api/routes/whirlpools.go"); err != nil {
		log.Panic(err.Error())
	}
	fmt.Printf("%#v", f)
}