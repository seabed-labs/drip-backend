package controller

import (
	"fmt"
	"math"
	"math/big"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const image = "data:image/svg+xml;charset=UTF-8,%3Csvg width='626' height='671' viewBox='0 0 626 671' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cg clip-path='url(%23clip0_309_116)'%3E%3Cpath d='M514.84 352C501.9 326.84 487.07 302.73 472.94 278.22C460.88 257.3 448.69 236.45 436.99 215.33C414.86 175.33 393.88 134.82 376.06 92.72C365.595 68.3918 356.651 43.4371 349.28 18C347.59 12 345.99 6 344.36 0C342.611 2.59619 341.576 5.6071 341.36 8.73C340.98 11.73 339.12 14.2 338.21 17.04C328.666 46.4611 317.461 75.3174 304.65 103.47C283.537 150.228 259.655 195.685 233.13 239.6C222.41 257.46 211.937 275.46 201.71 293.6C196.86 302.26 192.09 310.98 187.47 319.77C172.59 348.06 158.94 376.89 149.55 407.56C139.35 440.87 136 474.63 143.55 508.96C153.31 553.59 175.07 591.1 210.17 620.69C220.968 629.67 232.544 637.671 244.76 644.6C282.523 665.934 326.155 674.536 369.19 669.13C399.68 665.13 427.52 653.73 453.36 637.22C462.96 631.053 472.085 624.178 480.66 616.65C501.52 598.4 517.93 576.83 529.15 551.37C542.32 521.47 546.66 490.07 544.65 457.73C542.27 420.42 531.93 385.2 514.84 352ZM278.49 586.46C269.017 582.928 260.068 578.124 251.89 572.18C250.59 571.25 248.89 570.66 248.41 568.5C254.7 569.73 260.83 571.09 266.99 572.05C287.644 575.397 308.74 574.892 329.21 570.56C368.9 562.04 400.67 540.87 425.45 509.07C444.37 484.82 454.45 457.02 457.19 426.41C458.751 408.05 457.948 389.565 454.8 371.41C449.93 342.03 440.56 314 428.33 286.9C425.15 279.81 421.67 272.85 418.39 265.8C417.93 264.8 416.8 263.97 417.9 262.28C438.37 295.18 454.57 329.83 465.54 365.95C465.67 366.117 465.825 366.262 466 366.38C466.263 367.702 466.732 368.974 467.39 370.15C469.58 380.15 471.86 389.86 473.65 399.6C477.348 418.647 478.692 438.076 477.65 457.45C476.86 471.06 474.9 484.45 470.65 497.45C461.47 525.34 444.54 547.67 420.92 564.9C397.07 582.26 370.27 591.9 340.92 594.62C319.76 596.675 298.411 593.885 278.49 586.46V586.46ZM486.49 564C481.02 573.11 475.12 582 468.07 589.84C448.75 611.27 425.43 626.78 397.6 635.36C380.31 640.66 362.68 642.55 344.76 642.94L337.34 643.08C337.499 641.97 337.522 640.845 337.41 639.73C345.09 637.81 353.18 636.73 360.99 634.33C399.92 622.21 434.64 602.99 463.26 573.57C481.692 554.676 496.32 532.415 506.35 508C507.15 506.11 507.91 504.16 508.67 502.17C509.427 502.634 510.227 503.023 511.06 503.33C505.33 524.51 497.87 545 486.51 564H486.49Z' fill='%2362AAFF'/%3E%3Cpath d='M127.33 512.54C119.71 477.62 121.8 441.71 133.72 402.75C144.32 368.1 160.19 336.1 172.82 312.08C177.19 303.73 181.93 295.08 187.26 285.51C191.04 278.777 194.893 271.983 198.82 265.13C192.76 252.68 186.96 240.13 181.56 227.3C173.52 208.647 166.653 189.51 161 170C159.71 165.43 158.48 160.79 157.22 156.22C155.89 158.211 155.113 160.52 154.97 162.91C154.67 165.23 153.21 167.12 152.55 169.27C145.224 191.85 136.62 213.996 126.78 235.6C110.566 271.454 92.2405 306.315 71.9 340C63.6467 353.687 55.61 367.503 47.79 381.45C44.0433 388.077 40.3967 394.77 36.85 401.53C25.42 423.23 15 445.36 7.77 468.88C2.38419e-06 494.42 -2.64 520.33 3.13 546.66C10.61 580.92 27.31 609.71 54.25 632.4C62.5324 639.293 71.4106 645.437 80.78 650.76C102.909 663.562 127.94 670.503 153.502 670.925C179.064 671.346 204.311 665.235 226.85 653.17C217.264 647.232 208.108 640.625 199.45 633.4C162.28 602.06 138 561.37 127.33 512.54ZM105.46 638.67C103.08 638.573 100.752 637.941 98.65 636.82C92.6789 633.77 87.3624 629.582 83 624.49C76.3497 616.962 70.8623 608.481 66.72 599.33C64.88 595.43 63.31 591.43 61.72 587.43C54.21 567.9 48.18 547.87 42.12 527.85C37.25 511.78 32.78 495.6 27.7 479.6C27.5 478.93 27.37 478.23 27.19 477.6C27.2597 477.591 27.3303 477.591 27.4 477.6C32.12 487.7 36.74 497.86 41.58 507.92C50.35 526.23 59.29 544.44 69.85 561.81C72.96 566.94 76.47 571.81 80.02 576.58C85.76 584.37 91.1 592.42 96.14 600.66C100.822 608.004 104.838 615.751 108.14 623.81C109.337 626.746 110.127 629.831 110.49 632.98C110.94 636.42 109 638.69 105.46 638.67Z' fill='%2362AAFF'/%3E%3Cpath d='M625.43 312.47C623.961 290.405 617.902 268.889 607.64 249.3C599.92 234.3 591.07 219.91 582.64 205.3C575.45 192.81 568.17 180.3 561.18 167.73C547.96 143.88 535.43 119.73 524.8 94.5499C518.547 80.0192 513.203 65.1138 508.8 49.9199C507.8 46.3799 506.8 42.7699 505.85 39.1899C504.818 40.7525 504.202 42.5527 504.06 44.4199C503.83 46.2499 502.74 47.6699 502.21 49.4199C496.506 66.9711 489.829 84.191 482.21 101C469.594 128.911 455.332 156.049 439.5 182.27C439.2 182.77 438.91 183.27 438.64 183.76C442.72 191.413 447.003 199.277 451.49 207.35C460.61 223.75 470.18 240.35 479.49 256.35L487.34 269.97C491.18 276.59 495.043 283.193 498.93 289.78C509.27 307.44 520 325.72 529.64 344.44C543.133 370.431 552.595 398.323 557.7 427.16C568.361 421.866 578.296 415.221 587.26 407.39C599.679 396.694 609.567 383.373 616.21 368.39C624 350.51 626.62 331.76 625.43 312.47ZM533 199.76C532.2 201.388 531.112 202.858 529.79 204.1C526.24 207.38 522.33 207.38 518.72 204.1C517.129 202.686 515.86 200.947 515 199C509.63 187.37 509.93 175.71 514.86 164.08C516.62 160.01 519.24 156.27 524.14 156.2C529.04 156.13 531.56 160.01 533.52 163.95C536.4 169.68 537.26 175.84 537.06 182.44C537.45 188.43 535.76 194.22 533 199.76Z' fill='%2362AAFF'/%3E%3C/g%3E%3Cdefs%3E%3CclipPath id='clip0_309_116'%3E%3Crect width='625.7' height='670.89' fill='white'/%3E%3C/clipPath%3E%3C/defs%3E%3C/svg%3E"

var defaultTokenMetadata = apispec.TokenMetadata{
	Collection: struct {
		Family string `json:"family"`
		Name   string `json:"name"`
	}{
		Name:   "Drip Position",
		Family: "Drip Position v1",
	},
	Name:        "Drip Position v1",
	Description: "Drip Position v1",
	ExternalUrl: "https://drip.dcaf.so",
	Symbol:      "DP",
	Image:       image,
}

// TODO(Mocha): optimize
func (h Handler) GetV1DripPositionPubkeyPathMetadata(
	c echo.Context, mint apispec.PubkeyPathParam,
) error {
	log := logrus.WithField("mint", string(mint))
	position, err := h.repo.GetPositionByNFTMint(c.Request().Context(), string(mint))
	if err != nil {
		log.WithError(err).Error("failed to find drip position")
		return c.JSON(http.StatusOK, defaultTokenMetadata)
	}
	log = log.WithField("position", position.Pubkey).WithField("vault", position.Vault)
	vault, err := h.repo.AdminGetVaultByAddress(c.Request().Context(), position.Vault)
	if err != nil {
		log.WithError(err).Error("failed to find vault for position")
		return c.JSON(http.StatusOK, defaultTokenMetadata)
	}
	tokenMints, err := h.repo.GetTokensByAddresses(c.Request().Context(), vault.TokenAMint, vault.TokenBMint)
	if err != nil {
		log.WithError(err).Error("failed to get tokenMints")
	}
	if len(tokenMints) != 2 {
		log.WithField("len(tokenMints)", len(tokenMints)).Error("invalid number of mints returned")
	}
	var tokenA model.Token
	var tokenB model.Token
	for _, tokenMint := range tokenMints {
		if tokenMint.Pubkey == vault.TokenAMint {
			tokenA = *tokenMint
		} else if tokenMint.Pubkey == vault.TokenBMint {
			tokenB = *tokenMint
		}
	}
	depositAmount := *big.NewFloat(float64(position.DepositedTokenAAmount))
	if tokenA.Decimals != 0 {
		scale := math.Pow(10, float64(tokenA.Decimals))
		depositAmount = *new(big.Float).Quo(&depositAmount, big.NewFloat(scale))
	}
	var tokenASymbol string
	if tokenA.Symbol != nil {
		tokenASymbol = *tokenA.Symbol
	} else {
		tokenASymbol = tokenA.Pubkey[0:4] + "..." + tokenA.Pubkey[len(tokenA.Pubkey)-4:len(tokenA.Pubkey)]
	}
	var tokenBSymbol string
	if tokenB.Symbol != nil {
		tokenBSymbol = *tokenB.Symbol
	} else {
		tokenBSymbol = tokenB.Pubkey[0:4] + "..." + tokenB.Pubkey[len(tokenB.Pubkey)-4:len(tokenB.Pubkey)]
	}
	var nofMSwaps string
	endPeriodID := position.DcaPeriodIDBeforeDeposit + position.NumberOfSwaps
	if vault.LastDcaPeriod > endPeriodID {
		nofMSwaps = fmt.Sprintf("%d/%d", position.NumberOfSwaps, position.NumberOfSwaps)
	} else {
		nofMSwaps = fmt.Sprintf("%d/%d", endPeriodID-vault.LastDcaPeriod, position.NumberOfSwaps)
	}
	res := defaultTokenMetadata
	res.Description = fmt.Sprintf("Drip %s %s to %s\nDeposited on %s, with %s drips remaining",
		depositAmount.Text('e', 2), tokenASymbol, tokenBSymbol, position.DepositTimestamp.String(), nofMSwaps)
	return c.JSON(http.StatusOK, res)
}
