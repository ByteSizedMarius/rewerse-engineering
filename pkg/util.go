// Package rewerse provides a simple interface to interact with the REWE API.
// Only Getters are implemented.
package rewerse

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const (
	apiHost    = "mobile-api.rewe.de"
	clientHost = "mobile-clients-api.rewe.de"
)

var (
	Client    *http.Client
	cert      tls.Certificate
	userAgent string
	rdfa      string
	set       = false

	userAgents = []string{
		"Phone/Samsung_SM-G975U", "Phone/Samsung_SM-N975U", "Phone/Samsung_SM-G973U", "Phone/OnePlus_HD1925",
		"Phone/Xiaomi_M2007J3SY", "Phone/LG_LM-G820", "Phone/Google_Pixel_8_Pro", "Phone/Google_Pixel_7_Pro",
		"Phone/Samsung_SM-S911B", "Phone/Samsung_SM-S918B", "Phone/OnePlus_AC2003", "Phone/Xiaomi_2201123G",
		"Phone/Google_Pixel_8", "Phone/Google_Pixel_7a", "Phone/Samsung_SM-F946B", "Phone/Samsung_SM-S901B",
	}
)

func BuildCustomRequest(host, path string) (req *http.Request, err error) {
	if !set {
		panic("certificates not set")
	}

	req, err = http.NewRequest(http.MethodGet, "https://"+host+"/api/"+path, nil)
	if err != nil {
		err = fmt.Errorf("error creating request: %v", err)
		return
	}

	// Optional Headers
	// just adding these to fit in :)
	id, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("error generating uuid: %v", err)
		return
	}
	req.Header.Add("rdfa", rdfa)
	req.Header.Add("correlation-id", id.String())
	req.Header.Add("rd-service-types", "UNKNOWN")
	req.Header.Add("x-rd-service-types", "UNKNOWN")
	req.Header.Add("rd-is-lsfk", "false")
	req.Header.Add("rd-customer-zip", "")
	req.Header.Add("rd-postcode", "")
	req.Header.Add("x-rd-customer-zip", "")
	req.Header.Add("rd-market-id", "")
	req.Header.Add("x-rd-market-id", "")
	req.Header.Add("a-b-test-groups", "productlist-citrusad")
	// todo: some requests have a ruleSet header, but for others it makes them go 404

	// Strictly required headers
	req.Header.Set("user-agent", fmt.Sprintf("REWE-Mobile-Client/3.18.5.33032 Android/14 %s", userAgent))
	req.Header.Set("Host", host)
	req.Header.Set("Connection", "Keep-Alive")

	return
}

func DoRequest(req *http.Request, dest any) (err error) {
	// Execute Request
	resp, err := Client.Do(req)
	if err != nil {
		err = fmt.Errorf("error making request: %v", err)
		return
	}
	defer CloseWithWrap(resp.Body, &err)

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading response: %v", err)
		return
	}
	//fmt.Println(string(body))

	// Unmarshal the body into the destination
	if err = json.Unmarshal(body, &dest); err != nil {
		if strings.HasPrefix(string(body), "<!DOCTYPE html>") {
			err = fmt.Errorf("error: response is html (cloudflared)")
			return
		}

		err = fmt.Errorf("error unmarshalling response: %v", err)
		return
	}

	return
}

func SetCertificate(clientCert, clientKey string) error {
	var err error
	cert, err = tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return fmt.Errorf("error loading certificates: %v", err)
	}
	set = true

	// randomize a user-agent for this session
	userAgent = userAgents[rand.Intn(len(userAgents))]

	// rdfa is a static header
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("error generating uuid: %v", err)
	}
	rdfa = id.String()

	Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

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
