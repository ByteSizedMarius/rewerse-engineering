package rewerse

import (
	cr "crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	apiHost    = "mobile-api.rewe.de"
	clientHost = "mobile-clients-api.rewe.de"
)

// NewUUID generates a random UUID v4 string
func NewUUID() (string, error) {
	uuid := make([]byte, 16)

	n, err := io.ReadFull(cr.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// Set version (4) and variant bits according to RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}

// clientConfig holds all immutable client configuration.
// Swapped atomically to avoid lock contention on every request.
type clientConfig struct {
	client    *http.Client
	userAgent string
	instanaId string
	rdfaId    string
}

var (
	ErrNotInitialized = errors.New("certificates not set; call SetCertificate first")

	config atomic.Pointer[clientConfig]

	userAgents = []string{
		"Phone/Samsung_SM-G975U", "Phone/Samsung_SM-N975U", "Phone/Samsung_SM-G973U", "Phone/OnePlus_HD1925",
		"Phone/Xiaomi_M2007J3SY", "Phone/LG_LM-G820", "Phone/Google_Pixel_8_Pro", "Phone/Google_Pixel_7_Pro",
		"Phone/Samsung_SM-S911B", "Phone/Samsung_SM-S918B", "Phone/OnePlus_AC2003", "Phone/Xiaomi_2201123G",
		"Phone/Google_Pixel_8", "Phone/Google_Pixel_7a", "Phone/Samsung_SM-F946B", "Phone/Samsung_SM-S901B",
	}
)

// BuildCustomRequest creates a request to /api/{path} on the given host
func BuildCustomRequest(host, path string) (req *http.Request, err error) {
	return BuildCustomRequestRaw(host, "/api/"+path)
}

// BuildCustomRequestRaw creates a GET request to a raw path (without /api prefix)
func BuildCustomRequestRaw(host, path string) (req *http.Request, err error) {
	return buildRequest(http.MethodGet, host, path, nil)
}

// BuildPostRequest creates a POST request to /api/{path} with JSON body
func BuildPostRequest(host, path string, body io.Reader) (req *http.Request, err error) {
	req, err = buildRequest(http.MethodPost, host, "/api/"+path, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return
}

// BuildDeleteRequest creates a DELETE request to /api/{path}
func BuildDeleteRequest(host, path string) (req *http.Request, err error) {
	return buildRequest(http.MethodDelete, host, "/api/"+path, nil)
}

func buildRequest(method, host, path string, body io.Reader) (req *http.Request, err error) {
	cfg := config.Load()
	if cfg == nil {
		return nil, ErrNotInitialized
	}

	req, err = http.NewRequest(method, "https://"+host+path, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-instana-android", cfg.instanaId)
	req.Header.Set("user-agent", fmt.Sprintf("REWE-Mobile-Client/5.7.3.47565 Android/14 %s", cfg.userAgent))
	req.Header.Set("Host", host)
	req.Header.Set("Connection", "Keep-Alive")

	return
}

func DoRequest(req *http.Request, dest any) (err error) {
	cfg := config.Load()
	if cfg == nil {
		return ErrNotInitialized
	}

	resp, err := cfg.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer CloseWithWrap(resp.Body, &err)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateBody(body, 200))
	}

	if strings.HasPrefix(string(body), "<!DOCTYPE html>") {
		return fmt.Errorf("error: response is html (cloudflared)")
	}

	if err = json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	return
}

func truncateBody(body []byte, maxLen int) string {
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "..."
}

func SetCertificate(clientCert, clientKey string) error {
	loadedCert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return fmt.Errorf("error loading certificates: %w", err)
	}

	instanaId, err := NewUUID()
	if err != nil {
		return fmt.Errorf("error generating uuid: %w", err)
	}

	rdfaId, err := NewUUID()
	if err != nil {
		return fmt.Errorf("error generating rdfa uuid: %w", err)
	}

	cfg := &clientConfig{
		userAgent: userAgents[rand.Intn(len(userAgents))],
		instanaId: instanaId,
		rdfaId:    rdfaId,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{loadedCert},
				},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}

	config.Store(cfg)
	return nil
}

type CloseError struct {
	OriginalError, CloseError error
}

func (c CloseError) Error() string {
	return fmt.Sprintf("OriginalError=%v CloseError=%v", c.OriginalError, c.CloseError)
}

func CloseWithWrap(f io.Closer, e *error) {
	err := f.Close()
	if err != nil {
		if *e != nil {
			*e = CloseError{*e, err}
		} else {
			*e = err
		}
	}
}

// setCommonHeaders adds tracking headers used by most endpoints.
// Callers must ensure SetCertificate was called first (buildRequest validates this).
func setCommonHeaders(req *http.Request) {
	cfg := config.Load()
	req.Header.Set("rdfa", cfg.rdfaId)
	req.Header.Set("rdtga", "payment-enable-google-pay,productlist-citrusad")
	correlationId, _ := NewUUID()
	req.Header.Set("correlation-id", correlationId)
}

// setDualHeader sets both rd-{key} and x-rd-{key} headers.
// REWE API requires both prefixes for certain headers.
func setDualHeader(req *http.Request, key, value string) {
	req.Header.Set("rd-"+key, value)
	req.Header.Set("x-rd-"+key, value)
}
