// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// Delegation defines model for Delegation.
type Delegation struct {
	Amount    *int64     `json:"amount,omitempty"`
	Delegator *string    `json:"delegator,omitempty"`
	Level     *int64     `json:"level,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// GetXtzDelegationsParams defines parameters for GetXtzDelegations.
type GetXtzDelegationsParams struct {
	// Year Optional query parameter to filter results by year. Must be in the format YYYY and start with "20".
	Year *string `form:"year,omitempty" json:"year,omitempty"`

	// Page Page number for pagination. Defaults to 1 if not specified.
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of results per page. Defaults to 10 if not specified.
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Delegations
	// (GET /xtz/delegations)
	GetXtzDelegations(c *gin.Context, params GetXtzDelegationsParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetXtzDelegations operation middleware
func (siw *ServerInterfaceWrapper) GetXtzDelegations(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetXtzDelegationsParams

	// ------------- Optional query parameter "year" -------------

	err = runtime.BindQueryParameter("form", true, false, "year", c.Request.URL.Query(), &params.Year)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter year: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", c.Request.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter pageSize: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetXtzDelegations(c, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/xtz/delegations", wrapper.GetXtzDelegations)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RUTW/bMAz9KwTXo5s43gcw3wYEGHbYB9AdFrQdwNi0o8L6mMR0dQL/90FyFydtA2wn",
	"UyT19Pispz1WVjtr2EjAco+h2rCmFC6545ZEWRNXzlvHXhSnGmm7NRKjxnpNgiUqI+/eYIbSOx6X3LLH",
	"IcN6BLI+9j+Wg3hl2ljt+J67f0QSpTkIaXfSX5PwZSxNe/7CD4eMXd9xJTjElDKNTVyUdLH2nXc2wDQv",
	"XLG/V1XEu2cfkgK4mOWzPJKwjg05hSW+TqkMHckmyTJ/kN28PuCkXMtJp6heSn6qscSPLD9ktzzqjCie",
	"NAv7gOX1HmsOlVdu1B+/poA6+LVl38OhF8RCo7oYeQ7bTgKse+iZ/Aw+b4PAmkEZkA3DqBesVqsVkKkh",
	"CHmB30o2cINFfoMzjNJgiekMzNCQjvJENMwer0achR9Iu6RckRfFKICwj1t/Fvl1fvn+dl8MFy/9jqdz",
	"faOWwWz1mn0kCI5aZZIkM1hyQ2kisbAA1YCxAsFxpRrF9Tm2jlo+YVuPMFguson54vn1ek7uy8jLNgdt",
	"HSeK/IRc/l/srtTuHMO3xxTzFzjeZug5OGvCaMQiz+OnskZ4NCQ516kqSTi/C6N5p6NObVyTpKwS1ilx",
	"4bnBEl/Np1dh/vgkzI/eg8lW5D3153x2quYH6FSQqOaxQ2Jf2GpNvh99AcuT6jD8CQAA///O9q/KowQA",
	"AA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
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
	res := make(map[string]func() ([]byte, error))
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
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
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
