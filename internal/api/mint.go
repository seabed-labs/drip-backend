package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	log "github.com/sirupsen/logrus"
)

type MintRequest struct {
	Mint               string `schema:"mint,required"`
	Amount             uint64 `schema:"amount,required"`
	DestinationWallet  string `schema:"destination"`
	DestinationAccount string `schema:"destination"`
}
type MintResponse struct {
	TxHash string `schema:"txHash,required"`
	Error  string `schema:"error"`
}

// type CandleRequest struct {
// 	Base     string `schema:"base,required"`
// 	Start    uint64 `schema:"start,required"`
// 	End      uint64 `schema:"end,required"`
// 	Quote    string `schema:"quote"`
// 	Exchange string `schema:"exchange"`
// }

// type CandleResponse struct {
// 	OHLCV []*models.OHLCVMarketData `schema:"ohlcv,required"`
// 	Error null.String               `schema:"message"`
// }

func (h Handler) GetMint(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	var mintRequest MintRequest
	res := MintResponse{}
	if err := h.decoder.Decode(&mintRequest, r.URL.Query()); err != nil {
		log.WithField("err", err.Error()).Errorf("decode")
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if mintRequest.DestinationAccount == "" && mintRequest.DestinationWallet == "" {
		err := fmt.Errorf("one of <DestinationAccount, DestinationWallet> must be set")
		log.WithField("err", err.Error()).Errorf("decode")
		res.Error = fmt.Sprintf("invalid query, err=%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	resp, err := h.solanaClient.Client.GetAccountInfo(r.Context(), solana.MustPublicKeyFromBase58(mintRequest.Mint))
	if err != nil {
		log.
			WithField("mint", mintRequest.Mint).
			WithField("err", err.Error()).
			Errorf("MustPublicKeyFromBase58")
		res.Error = fmt.Sprintf("invalid mint, err=%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	var mint token.Mint
	if err := bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(&mint); err != nil {
		log.
			WithField("err", err.Error()).
			Errorf("decode")
		res.Error = fmt.Errorf("invalid mint data, err=%s", err.Error()).Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		log.WithField("err", err.Error()).Errorf("decode")
		res.Error = fmt.Sprintf("invalid mint, err=%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	if *mint.MintAuthority != h.solanaClient.Wallet.PublicKey() {
		err := fmt.Errorf("server is not auth")
		log.
			WithField("authority", mint.MintAuthority).
			WithField("wallet", h.solanaClient.Wallet.PublicKey()).
			WithField("err", err.Error()).
			Errorf("invalid mint")
		res.Error = fmt.Sprintf("invalid mint, err=%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	var txHash string
	if mintRequest.DestinationAccount != "" {
		txHash, err = h.solanaClient.MintToAccount(r.Context(), mintRequest.Mint, mintRequest.DestinationAccount, mintRequest.Amount)
	} else {
		txHash, err = h.solanaClient.MintToWallet(r.Context(), mintRequest.Mint, mintRequest.DestinationWallet, mintRequest.Amount)
	}
	if err != nil {
		log.WithField("err", err.Error()).Errorf("mintTo")
		res.Error = fmt.Sprintf("failed to mint, err=%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	res.TxHash = txHash
	json.NewEncoder(w).Encode(res)
}
