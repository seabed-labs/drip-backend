package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GetAdminPositionsExpandParams string

const (
	expandPositionAll         = GetAdminPositionsExpandParams("all")
	expandPositionVault       = GetAdminPositionsExpandParams("vault")
	expandPositionProtoConfig = GetAdminPositionsExpandParams("protoConfig")
	expandPositionTokenA      = GetAdminPositionsExpandParams("tokenA")
	expandPositionTokenB      = GetAdminPositionsExpandParams("tokenB")
)

func (h Handler) GetV1AdminPositions(c echo.Context, params apispec.GetV1AdminPositionsParams) error {
	res := apispec.ListAdminPositions{}
	positions, err := h.repo.GetAdminPositions(
		c.Request().Context(),
		(*bool)(params.Enabled),
		repository.PositionFilterParams{
			IsClosed: (*bool)(params.IsClosed),
			Wallet:   nil,
		},
		getPaginationParamsFromAPI(params.Offset, params.Limit),
	)
	if err != nil {
		logrus.WithError(err).Error("failed to GetAdminPositions")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: err.Error()})
	}

	// Populate Base Result
	for _, position := range positions {
		res = append(res, apispec.ExpandedAdminPosition{
			Position: positionModelToAPI(position),
		})
	}
	if params.Expand == nil {
		return c.JSON(http.StatusOK, res)
	}
	if hasValue(*params.Expand, string(expandPositionAll)) {
		newParams := apispec.ExpandAdminPositionsQueryParam{string(expandPositionVault), string(expandPositionProtoConfig), string(expandPositionTokenA), string(expandPositionTokenB)}
		params.Expand = &newParams
	}
	shouldExpandVault := hasValue(*params.Expand, string(expandPositionVault))
	shouldExpandProtoConfig := hasValue(*params.Expand, string(expandPositionProtoConfig))
	shouldExpandTokenA := hasValue(*params.Expand, string(expandPositionTokenA))
	shouldExpandTokenB := hasValue(*params.Expand, string(expandPositionTokenB))

	var vaultPubkeys []string
	for i := range positions {
		position := positions[i]
		vaultPubkeys = append(vaultPubkeys, position.Vault)
	}

	var vaults []*model.Vault
	vaultsByPubkey := make(map[string]*model.Vault)

	var protoConfigPubkeys []string
	var protoConfigs []*model.ProtoConfig
	protoConfigsByPubkey := make(map[string]*model.ProtoConfig)

	var tokenPubkeys []string
	var tokens []*model.Token
	tokensByPubkey := make(map[string]*model.Token)

	vaults, err = h.repo.AdminGetVaultsByAddresses(c.Request().Context(), vaultPubkeys...)
	if err != nil {
		logrus.WithError(err).Error("failed to AdminGetVaultsByAddresses")
		return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
	}
	for i := range vaults {
		vault := vaults[i]
		vaultsByPubkey[vault.Pubkey] = vault
		protoConfigPubkeys = append(protoConfigPubkeys, vault.ProtoConfig)
		tokenPubkeys = append(tokenPubkeys, vault.TokenAMint, vault.TokenBMint)
	}

	if shouldExpandProtoConfig {
		protoConfigs, err = h.repo.GetProtoConfigsByAddresses(c.Request().Context(), protoConfigPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetProtoConfigsByAddresses")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
		}
		for i := range protoConfigs {
			protoConfigsByPubkey[protoConfigs[i].Pubkey] = protoConfigs[i]
		}
	}

	if shouldExpandTokenA || shouldExpandTokenB {
		tokens, err = h.repo.GetTokensByMints(c.Request().Context(), tokenPubkeys)
		if err != nil {
			logrus.WithError(err).Error("failed to GetTokensByMints")
			return c.JSON(http.StatusInternalServerError, apispec.ErrorResponse{Error: "internal server error"})
		}
		for i := range tokens {
			tokensByPubkey[tokens[i].Pubkey] = tokens[i]
		}
	}
	for i := range res {
		vault := vaultsByPubkey[res[i].Position.Vault]

		if shouldExpandVault && vault != nil {
			apiVault := vaultModelToAPI(vault)
			res[i].Vault = &apiVault
		}

		if shouldExpandTokenA && vault != nil {
			tokenA := tokensByPubkey[vault.TokenAMint]
			if tokenA != nil {
				apiToken := tokenModelToApi(tokenA)
				res[i].TokenA = &apiToken
			}
		}
		if shouldExpandTokenB && vault != nil {
			tokenB := tokensByPubkey[vault.TokenBMint]
			if tokenB != nil {
				apiToken := tokenModelToApi(tokenB)
				res[i].TokenB = &apiToken
			}
		}
		if shouldExpandProtoConfig && vault != nil {
			protoConfig := protoConfigsByPubkey[vault.ProtoConfig]
			apiProtoConfig := protoConfigModelToAPI(protoConfig)
			res[i].ProtoConfig = &apiProtoConfig
		}
	}

	return c.JSON(http.StatusOK, res)
}
