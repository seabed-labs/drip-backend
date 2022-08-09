package controller

import (
	"net/http"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetV1DripPubkeyPathTokenmetadata(c echo.Context, pubkeyPath Swagger.PubkeyPathParam) error {
	switch pubkeyPath {
	// USDT
	case "8ULDKGmKJJaZa32eiL36ARr6cFaZaoAXAosWeg5r17ra":
		return c.JSON(http.StatusOK, Swagger.TokenMetadata{
			Collection: struct {
				Family string `json:"family"`
				Name   string `json:"name"`
			}{
				Name:   "",
				Family: "",
			},
			Description: "",
			ExternalUrl: "https://drip.dcaf.so",
			Name:        "Drip Devnet USDT",
			Symbol:      "DUSDT",
			Image:       "data:image/svg+xml;charset=UTF-8,%3Csvg width='32' height='32' viewBox='0 0 32 32' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M16 32C7.163 32 0 24.837 0 16C0 7.163 7.163 0 16 0C24.837 0 32 7.163 32 16C32 24.837 24.837 32 16 32ZM17.922 13.793V11.427H23.336V7.819H8.595V11.427H14.009V13.792C9.609 13.994 6.3 14.866 6.3 15.91C6.3 16.954 9.609 17.825 14.009 18.028V25.61H17.922V18.026C22.315 17.824 25.616 16.953 25.616 15.91C25.616 14.867 22.315 13.996 17.922 13.793ZM17.922 17.383V17.381C17.812 17.389 17.245 17.423 15.98 17.423C14.97 17.423 14.259 17.393 14.009 17.381V17.384C10.121 17.213 7.219 16.536 7.219 15.726C7.219 14.917 10.121 14.24 14.009 14.066V16.71C14.263 16.728 14.991 16.771 15.997 16.771C17.204 16.771 17.809 16.721 17.922 16.711V14.068C21.802 14.241 24.697 14.918 24.697 15.726C24.697 16.536 21.802 17.211 17.922 17.383V17.383Z' fill='%2362AAFF'/%3E%3C/svg%3E%0A",
		})
	// BTC
	case "5nY3xT4PJe7NU41zqBx5UACHDckrimmfwznv4uLenrQg":
		return c.JSON(http.StatusOK, Swagger.TokenMetadata{
			Collection: struct {
				Family string `json:"family"`
				Name   string `json:"name"`
			}{
				Name:   "",
				Family: "",
			},
			Description: "",
			ExternalUrl: "https://drip.dcaf.so",
			Name:        "Drip Devnet BTC",
			Symbol:      "DBTC",
			Image:       "data:image/svg+xml;charset=UTF-8,%3Csvg width='32' height='32' viewBox='0 0 32 32' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M16 32C7.163 32 0 24.837 0 16C0 7.163 7.163 0 16 0C24.837 0 32 7.163 32 16C32 24.837 24.837 32 16 32ZM23.189 14.02C23.503 11.924 21.906 10.797 19.724 10.045L20.432 7.205L18.704 6.775L18.014 9.54C17.56 9.426 17.094 9.32 16.629 9.214L17.324 6.431L15.596 6L14.888 8.839C14.512 8.753 14.142 8.669 13.784 8.579L13.786 8.57L11.402 7.975L10.942 9.821C10.942 9.821 12.225 10.115 12.198 10.133C12.898 10.308 13.024 10.771 13.003 11.139L12.197 14.374C12.245 14.386 12.307 14.404 12.377 14.431L12.194 14.386L11.064 18.918C10.978 19.13 10.761 19.449 10.271 19.328C10.289 19.353 9.015 19.015 9.015 19.015L8.157 20.993L10.407 21.554C10.825 21.659 11.235 21.769 11.638 21.872L10.923 24.744L12.65 25.174L13.358 22.334C13.83 22.461 14.288 22.579 14.736 22.691L14.03 25.519L15.758 25.949L16.473 23.083C19.421 23.641 21.637 23.416 22.57 20.75C23.322 18.604 22.533 17.365 20.982 16.558C22.112 16.298 22.962 15.555 23.189 14.02V14.02ZM19.239 19.558C18.706 21.705 15.091 20.544 13.919 20.253L14.869 16.448C16.041 16.741 19.798 17.32 19.239 19.558ZM19.774 13.989C19.287 15.942 16.279 14.949 15.304 14.706L16.164 11.256C17.139 11.499 20.282 11.952 19.774 13.989Z' fill='%2362AAFF'/%3E%3C/svg%3E%0A",
		})
	default:
		return c.JSON(http.StatusBadRequest, Swagger.ErrorResponse{Error: "invalid mint"})
	}

}
