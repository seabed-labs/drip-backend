package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gorilla/mux"
)

func (h Handler) GetMint(
	w http.ResponseWriter,
	r *http.Request,
) {
	mintPubKey := mux.Vars(r)["publicKey"]
	resp, err := h.solanaClient.Client.GetAccountInfo(r.Context(), solana.MustPublicKeyFromBase58(mintPubKey))
	if err != nil {
		panic(err)
	}
	var mint token.Mint
	err = bin.NewBorshDecoder(resp.Value.Data.GetBinary()).Decode(&mint)
	if err != nil {
		panic(err)
	}
	// if mint.MintAuthority != h.solanaClient.Wallet.PublicKey().ToPointer() {
	// 	panic(fmt.Errorf("I am not the auth :("))
	// }
	// h.solanaClient.MintTo(context.TODO(), mintPubKey, )
	json.NewEncoder(w).Encode(fmt.Sprintf("GetMint %s %s", mintPubKey, mint.MintAuthority))
}
