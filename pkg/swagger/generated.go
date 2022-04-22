// Package Swagger provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.10.1 DO NOT EDIT.
package Swagger

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ErrorResponse defines model for errorResponse.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MintResponse defines model for mintResponse.
type MintResponse struct {
	TxHash string `json:"txHash"`
}

// PingResponse defines model for pingResponse.
type PingResponse struct {
	Message string `json:"message"`
}

// Amount defines model for amount.
type Amount string

// Mint defines model for mint.
type Mint string

// Wallet defines model for wallet.
type Wallet string

// GetMintParams defines parameters for GetMint.
type GetMintParams struct {
	// Mint base58 public key.
	Mint Mint `json:"mint"`

	// Wallet base58 public key.
	Wallet Wallet `json:"wallet"`

	// Amount to be minted in base amount.
	Amount Amount `json:"amount"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Get request
	Get(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetMint request
	GetMint(ctx context.Context, params *GetMintParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) Get(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMint(ctx context.Context, params *GetMintParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMintRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetRequest generates requests for Get
func NewGetRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetMintRequest generates requests for GetMint
func NewGetMintRequest(server string, params *GetMintParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/mint")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "mint", runtime.ParamLocationQuery, params.Mint); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "wallet", runtime.ParamLocationQuery, params.Wallet); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "amount", runtime.ParamLocationQuery, params.Amount); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// Get request
	GetWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetResponse, error)

	// GetMint request
	GetMintWithResponse(ctx context.Context, params *GetMintParams, reqEditors ...RequestEditorFn) (*GetMintResponse, error)
}

type GetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PingResponse
}

// Status returns HTTPResponse.Status
func (r GetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetMintResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *MintResponse
	JSON400      *ErrorResponse
	JSON500      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetMintResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetMintResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetWithResponse request returning *GetResponse
func (c *ClientWithResponses) GetWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetResponse, error) {
	rsp, err := c.Get(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetResponse(rsp)
}

// GetMintWithResponse request returning *GetMintResponse
func (c *ClientWithResponses) GetMintWithResponse(ctx context.Context, params *GetMintParams, reqEditors ...RequestEditorFn) (*GetMintResponse, error) {
	rsp, err := c.GetMint(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMintResponse(rsp)
}

// ParseGetResponse parses an HTTP response from a GetWithResponse call
func ParseGetResponse(rsp *http.Response) (*GetResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PingResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetMintResponse parses an HTTP response from a GetMintWithResponse call
func ParseGetMintResponse(rsp *http.Response) (*GetMintResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMintResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest MintResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health Check
	// (GET /)
	Get(ctx echo.Context) error
	// Mint tokens (DEVNET ONLY)
	// (GET /mint)
	GetMint(ctx echo.Context, params GetMintParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Get converts echo context to params.
func (w *ServerInterfaceWrapper) Get(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Get(ctx)
	return err
}

// GetMint converts echo context to params.
func (w *ServerInterfaceWrapper) GetMint(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetMintParams
	// ------------- Required query parameter "mint" -------------

	err = runtime.BindQueryParameter("form", true, true, "mint", ctx.QueryParams(), &params.Mint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter mint: %s", err))
	}

	// ------------- Required query parameter "wallet" -------------

	err = runtime.BindQueryParameter("form", true, true, "wallet", ctx.QueryParams(), &params.Wallet)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter wallet: %s", err))
	}

	// ------------- Required query parameter "amount" -------------

	err = runtime.BindQueryParameter("form", true, true, "amount", ctx.QueryParams(), &params.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter amount: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMint(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/", wrapper.Get)
	router.GET(baseURL+"/mint", wrapper.GetMint)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8RVTW/TQBD9K6uBA0hWEj4qIZ+gtKIVtEUFgVCVw2Q9iZd4P7q7LkSV/zvasa3ENFFa",
	"hOASJTOz783Om325BWm1s4ZMDJDfgkOPmiJ5/oXa1iambwUF6ZWLyhrI4Q3HRbRiRkIrE6kQyogZBhLt",
	"mRFkoFLpdU1+BRkY1AR5j5iBp+taeSogj76mDIIsSWOiiiuXKkP0yiygaTJIDHebOFMmMuXBK+HqWaWk",
	"WNJqFzFjPIz2B1YVbSH+yvH7U3c4DyFv+iTLQN5bf0nBWROIVfLWkY+K1untV1gzXnVl026eu+HizxMM",
	"5X68ri4BOmUWuwE1hYAL2o/YF055ANKaUOsEcQXoXKUkJgHG34M1ME2jntuEKa2JKFkn0qgqyKGQONdW",
	"lvjaeRutSeGRtHqtyZHEuThLJdBkv+l75JUTM5RLMoUI5G+UpKRtVLGiPn/Y5iGDG/KhPfhsNBlNEp4l",
	"Z9ApyOEFhzKwjroIZOAwljyZcfpYtDu277rDHj8qs+DeyI9aeM/lpwXk8I63zXlb1HI3oO8E406eTyb9",
	"LKl9bHeO5LcbC/vY0xxyeDRe28e429jxYBtYymHvF+9Z+VBrjX4FOZwQVrEUb0uSyzRnXHDTLPA0lY57",
	"B/ijWaXDIlJIhrUkE5JvoSgopL0TGIKVCpOFcVqglMmjMmG9cBhC622D3LaJn7UGs2mgV9vntC5p79Vk",
	"e+s6C7lHZWewzfQ/LsDAX3YsQAYv/yLj0CG3UB5iIS7puqaQpggH/5L71ETyBivxid+rOGYjHj4B/jPr",
	"1vPJ0fGX8+PP4uL8w7enG++Bl6X1xn3Scm9tuozRwbT5FQAA//8oxKG26QcAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
