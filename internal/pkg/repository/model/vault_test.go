package model

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/test-go/testify/assert"
)

func TestVaultModel(t *testing.T) {

	vault := Vault{
		Pubkey:                 "123",
		ProtoConfig:            "13",
		TokenAMint:             "12",
		TokenBMint:             "1123",
		TokenAAccount:          "123",
		TokenBAccount:          "123",
		TreasuryTokenBAccount:  "123",
		LastDcaPeriod:          123,
		DripAmount:             123,
		DcaActivationTimestamp: time.Now(),
	}
	blob, err := json.Marshal(vault)
	assert.NoError(t, err)
	fmt.Println(string(blob))
}
