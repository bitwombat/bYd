// Package types provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package types

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Comment defines model for Comment.
type Comment struct {
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Id          string    `json:"id"`
	PostId      int       `json:"postId"`
}

// Post defines model for Post.
type Post struct {
	Comments    *Comment  `json:"comments,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
	Id          int       `json:"id"`
	IsAnon      bool      `json:"isAnon"`
	Title       string    `json:"title"`
	Type        int       `json:"type"`
	UserId      string    `json:"userId"`
}

// State defines model for State.
type State struct {
	Code string `json:"code"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

// User defines model for User.
type User struct {
	// CreatedAt The date that the user was created.
	CreatedAt *openapi_types.Date  `json:"createdAt,omitempty"`
	Email     *openapi_types.Email `json:"email,omitempty"`
	FirstName string               `json:"firstName"`

	// Id Unique identifier for the given user.
	Id         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
	LastName   string `json:"lastName"`
	State      State  `json:"state"`
}

// CreatePost defines model for CreatePost.
type CreatePost struct {
	Description string `json:"description"`
	IsAnon      bool   `json:"isAnon"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

// CreateNewPostJSONBody defines parameters for CreateNewPost.
type CreateNewPostJSONBody struct {
	Description string `json:"description"`
	IsAnon      bool   `json:"isAnon"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

// CreateNewPostJSONRequestBody defines body for CreateNewPost for application/json ContentType.
type CreateNewPostJSONRequestBody CreateNewPostJSONBody

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
	// GetAllPosts request
	GetAllPosts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateNewPost request with any body
	CreateNewPostWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateNewPost(ctx context.Context, body CreateNewPostJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAllPosts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAllPostsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateNewPostWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateNewPostRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateNewPost(ctx context.Context, body CreateNewPostJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateNewPostRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAllPostsRequest generates requests for GetAllPosts
func NewGetAllPostsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/posts")
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

// NewCreateNewPostRequest calls the generic CreateNewPost builder with application/json body
func NewCreateNewPostRequest(server string, body CreateNewPostJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateNewPostRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateNewPostRequestWithBody generates requests for CreateNewPost with any type of body
func NewCreateNewPostRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/posts")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

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
	// GetAllPosts request
	GetAllPostsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllPostsResponse, error)

	// CreateNewPost request with any body
	CreateNewPostWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateNewPostResponse, error)

	CreateNewPostWithResponse(ctx context.Context, body CreateNewPostJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateNewPostResponse, error)
}

type GetAllPostsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON2XX      *struct {
		Posts *[]Post `json:"posts,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r GetAllPostsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAllPostsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateNewPostResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Post
}

// Status returns HTTPResponse.Status
func (r CreateNewPostResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateNewPostResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAllPostsWithResponse request returning *GetAllPostsResponse
func (c *ClientWithResponses) GetAllPostsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetAllPostsResponse, error) {
	rsp, err := c.GetAllPosts(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAllPostsResponse(rsp)
}

// CreateNewPostWithBodyWithResponse request with arbitrary body returning *CreateNewPostResponse
func (c *ClientWithResponses) CreateNewPostWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateNewPostResponse, error) {
	rsp, err := c.CreateNewPostWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateNewPostResponse(rsp)
}

func (c *ClientWithResponses) CreateNewPostWithResponse(ctx context.Context, body CreateNewPostJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateNewPostResponse, error) {
	rsp, err := c.CreateNewPost(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateNewPostResponse(rsp)
}

// ParseGetAllPostsResponse parses an HTTP response from a GetAllPostsWithResponse call
func ParseGetAllPostsResponse(rsp *http.Response) (*GetAllPostsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAllPostsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode/100 == 2:
		var dest struct {
			Posts *[]Post `json:"posts,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON2XX = &dest

	}

	return response, nil
}

// ParseCreateNewPostResponse parses an HTTP response from a CreateNewPostWithResponse call
func ParseCreateNewPostResponse(rsp *http.Response) (*CreateNewPostResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateNewPostResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Post
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Your GET endpoint
	// (GET /v1/posts)
	GetAllPosts(w http.ResponseWriter, r *http.Request)
	// Your POST endpoint
	// (POST /v1/posts)
	CreateNewPost(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetAllPosts operation middleware
func (siw *ServerInterfaceWrapper) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllPosts(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateNewPost operation middleware
func (siw *ServerInterfaceWrapper) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateNewPost(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/posts", wrapper.GetAllPosts)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/posts", wrapper.CreateNewPost)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xXb2+cuBP+Ksi/30uSBTYb2H116V1VVTq10aU9tar6wpgBTMA2tmGBaL/7yYb9l+ym",
	"2969A5uZ8TzPMzPmCRFeCc6AaYVWT0hC3YDSb3hCwS78LgFruOdKmzfCmQZmH7EQJSVYU85mheLMrCmS",
	"Q4XNk5BcgNSTkwQUkVSYb81rhbs/gWU6R6vA8zwX6V4AWiGlJWUZclF3pTQXJc1yG4omaIXWcSPovA7y",
	"x6It0GbjIqru2Ohwso85LwGzMw56FaZEw0LFN0NgHWiqSziwfzV+FeWikni58AePjubW6jJrHcVlHfi9",
	"n5ZNijbG3EBNJSRo9W06yS4l9wixKdL3HU48LoBo48X4GUEfyeJVNdFzTACxLCZ3divlssIarVCCNVxp",
	"WgG6jIJq6Hlcg1rc6mZhIXhG7CVO8ow8rsO09LOwhpHH5FLb5JEMfZlni6ReM2sruNLvD+0p05CBPOMg",
	"TfKyIRQ3QUe7lzTQBO1cPudgD+H3nXJ2gD+n5nT0uV/PmzwOIWdxbY+/ratnbI1e7fP/JaRohf4329fp",
	"bGJ8to2+cf8jgtsl9oohSsq2DuNfJdjjSgyx388LHnUvCX6doDi8Jd3gN31QJPxXqryUaRIOrKtxzKKf",
	"r/Ii9Hnk9Qsd0KI9VeWvH19EC0mrhKUak9G8USDfXyxwmtOgm4tHWpO0PaPPyeNBs9h2j3OCPWggk26t",
	"7i4TbY1Jq5YZl6UIb21KDxprOKXa5GKYa2gTUmSkwRXvfrIJsOU8kcvhNpnnvW9tGa4ujtzdQEzr4KZv",
	"w2x5BmHrzx0TOgBtzPsy1GRZ1PWcaKHDnNhDflYgzTZ0uBKlgezbky3Qj+kbKs0sRP5yGV753tXcRy6C",
	"CtMSrRAuKYFrVVGd/5aZtWvCq+3+3yBpSs3RtWzARSmVSn+wcKA7Y4hGZP2bwEUl3u09GHfIRYpm7LP4",
	"w/KJAs9fXnnRVXCDNt/d1ybIUVdAn3JwTCKOzrF2dA6OkaizxsqZrK6Re9yULuxHxQKKym8X4jHA48ie",
	"UDnocePKc3+bIyyeXu6OejtO5DOjdQMOTYBpA6t0Ui5tQhltgdm0rk+FompPxIV9qgn9QZeLdbAgYuxT",
	"e3pOHFdta+61eTAK9KSk92AcBDo6+DbGgeCtZE/cOlxEWcq390FMTFbmjE1VYdmjFfrr7cMn5+7+vbIA",
	"xl8TtHcaf02cu/3d0XkA2YI0XyMXtSDVyIV/7Zm8uQCGBTXD89osuUhgnVtFzlp/Zma1fcnAQmsUa92a",
	"lovegb4ry3v7jcFECc7UKOfgy5d/cZ/dhaUaqh+Oadtsd5MEYSlxf0YVsMgiuRayyLpygvoE+seqfWgI",
	"AaUs73sKvvJGOu/efnKAJYJTdq5RFUPq10MTMho1dHelelkb43+Aw2DtiHF4HGM97n+A9TRa9j8S/Tl8",
	"jv41Zgc/GpvnZHneT5H1YzJeovgB1o7Z27esE3jef3z4IaA+vmm1CKNCSVqh6YpuNT62/Eaarp5rLVaz",
	"WckJLnOu9CryPM+23ZO3osfqtlGNjLLOMzP4nwAAAP//reKHfrYNAAA=",
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