package rfish

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/stmcginnis/gofish"
)

var (
	ErrPowerOffHost    = errors.New("host must be turned off before starting update firmware")
	ErrEmptySystemList = errors.New("systems list can't be empty")
	ErrCreateSession   = errors.New("can't create new RedFish Session")
	ErrGetAuthToken    = errors.New("failed get auth token from Session")
	ErrPostingFirmware = errors.New("can't posting firmware, status code not valid")
)

type Option interface {
	apply(*options)
}

const AuthHeader = "X-Auth-Token"

type (
	Credentials  credentials
	EndpointOpt  string
	InsecureOpt  bool
	BasicAuthOpt bool
)

type credentials struct {
	username string
	password string
}

type options struct {
	endpoint  string
	username  string
	password  string
	insecure  bool
	basicAuth bool
}

type RedfishAdapter struct {
	client *gofish.APIClient
	config *gofish.ClientConfig
}

func New() RedfishAdapter {
	return RedfishAdapter{}
}

func (ra *RedfishAdapter) Connect(opts ...Option) error {
	options := &options{
		endpoint:  "http://localhost:80",
		username:  "admin",
		password:  "admin",
		insecure:  false,
		basicAuth: false,
	}

	for _, o := range opts {
		o.apply(options)
	}

	cfg := gofish.ClientConfig{
		Endpoint:            options.endpoint,
		Username:            options.username,
		Password:            options.password,
		Insecure:            options.insecure,
		BasicAuth:           options.basicAuth,
		TLSHandshakeTimeout: 70,
	}

	client, err := gofish.Connect(cfg)
	if err != nil {
		return fmt.Errorf("gofish failed connect to target device: config: %+v: error: %w", cfg, err)
	}

	ra.client = client
	ra.config = &cfg

	return nil
}

func (ra *RedfishAdapter) GetTaskMonitor(taskMonitor string) (*http.Response, error) {
	resp, err := ra.client.Get(taskMonitor)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ra *RedfishAdapter) PostWithHeaders(headers map[string]string, updateURL, body string) (*http.Response, error) {
	payload := make(map[string]any)
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	resp, err := ra.client.PostWithHeaders(updateURL, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("error when try request to update: %w", err)
	}

	return resp, nil
}

func (ra *RedfishAdapter) CreateSession(headers map[string]string, username, password string) (map[string]string, string, error) {
	p := fmt.Sprintf("{\"UserName\": \"%s\", \"Password\": \"%s\"}", username, password)
	payload := make(map[string]any)
	if err := json.Unmarshal([]byte(p), &payload); err != nil {
		return nil, "", err
	}
	sessionURL := "/redfish/v1/SessionService/Sessions"
	resp, err := ra.client.PostWithHeaders(sessionURL, payload, headers)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	var result = make(map[string]any)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	var id string
	x, ok := result["Id"]
	if ok {
		id = x.(string)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, "", ErrCreateSession
	}

	authToken := resp.Header.Get(AuthHeader)
	if authToken == "" {
		return nil, "", ErrGetAuthToken
	}

	headers[AuthHeader] = authToken

	return headers, id, nil
}

func (ra *RedfishAdapter) DeleteSession(headers map[string]string, id string) error {

	resp, err := ra.client.DeleteWithHeaders("/redfish/v1/SessionService/Sessions/"+id, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func WithInsecure(o bool) Option {
	return InsecureOpt(o)
}

func WithBasicAuth(o bool) Option {
	return BasicAuthOpt(o)
}

func WithEndpoint(endpoint string) Option {
	return EndpointOpt(endpoint)
}

func WithCredentials(username, password string) Option {
	return Credentials{
		username: username,
		password: password,
	}
}

func (in InsecureOpt) apply(opts *options) {
	opts.insecure = bool(in)
}

func (b BasicAuthOpt) apply(opts *options) {
	opts.basicAuth = bool(b)
}

func (c Credentials) apply(opts *options) {
	opts.username = c.username
	opts.password = c.password
}

func (e EndpointOpt) apply(opts *options) {
	opts.endpoint = string(e)
}
