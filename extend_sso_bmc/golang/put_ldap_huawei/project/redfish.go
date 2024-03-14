package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RedfishClient struct {
	Address  string
	Username string
	Password string
}

func newRedfishClient(username, password, address string, port int32) *RedfishClient {
	port, protocol := preparePortAndProtocol(port)

	return &RedfishClient{
		Address:  fmt.Sprintf("%s%s:%d", protocol, address, port),
		Username: username,
		Password: password,
	}
}

func preparePortAndProtocol(in int32) (int32, string) {
	port := checkPort(in)

	protocol := HTTP_PROTOCOL
	if port == HTTPS_PORT {
		protocol = HTTPS_PROTOCOL
	}

	return port, protocol
}

func checkPort(port int32) int32 {
	if port <= 0 {
		return HTTPS_PORT
	}

	return port
}

func (r *RedfishClient) PatchData(endpoint string, body string) error {
	resp, err := r.runRequest(endpoint, http.MethodPatch, nil, strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (r *RedfishClient) GetPage(endpoint string) ([]byte, error) {
	resp, err := r.runRequest(endpoint, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %s", err)
	}

	return b, nil
}

func (r *RedfishClient) PostData(endpoint string, body string) error {
	resp, err := r.runRequest(endpoint, http.MethodPost, nil, strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (r *RedfishClient) GetAllHeadersFromPage(endpoint string, method string, body io.Reader) (http.Header, error) {
	resp, err := r.runRequest(endpoint, method, nil, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Header, nil
}

func (r *RedfishClient) DeleteByPage(endpoint string) error {
	resp, err := r.runRequest(endpoint, http.MethodDelete, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (r *RedfishClient) GetHeaderFromPage(endpoint string, headerKey string, method string, body io.Reader) (string, error) {
	resp, err := r.runRequest(endpoint, method, nil, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	header := resp.Header.Get(headerKey)
	if header == "" {
		return "", fmt.Errorf("empty header")
	}

	return header, nil
}

func (r *RedfishClient) PatchDataWithHeaders(endpoint string, body string, customHeaders map[string]string) error {
	resp, err := r.runRequest(endpoint, http.MethodPatch, customHeaders, strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (r *RedfishClient) PostDataWithHeaders(endpoint string, body string, customHeaders map[string]string) error {
	resp, err := r.runRequest(endpoint, http.MethodPost, customHeaders, strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (r *RedfishClient) runRequest(endpoint, method string, customHeaders map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, r.Address+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request error: %s", err)
	}

	// общие хэдеры
	req.Header.Set("User-Agent", "gofish/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// ворзможно добавить обработку '@' символа
	for k, v := range customHeaders {
		if k == "" || v == "" {
			continue
		}

		req.Header.Set(k, v)
	}

	req.SetBasicAuth(r.Username, r.Password)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error: %s", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got redfish error, but can not read response body: %s", err)
		}
		defer resp.Body.Close()
		return nil, fmt.Errorf("got redfish error: %s", string(payload))
	}

	return resp, nil
}
