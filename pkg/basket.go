package rewerse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ServiceType represents the delivery/pickup service type
type ServiceType string

const (
	ServicePickup   ServiceType = "PICKUP"
	ServiceDelivery ServiceType = "DELIVERY"
)

func validateServiceType(st ServiceType) error {
	if st != ServicePickup && st != ServiceDelivery {
		return fmt.Errorf("invalid service type: %q (must be PICKUP or DELIVERY)", st)
	}
	return nil
}

// BasketSession holds the context for basket operations.
// Create one with CreateBasket, then use it for subsequent operations.
type BasketSession struct {
	// ID is the server-generated basket UUID
	ID string
	// DeviceID is the server-assigned device identifier
	DeviceID string
	// Version tracks basket changes for optimistic locking
	Version int
	// MarketID is the selected market
	MarketID string
	// ZipCode is the customer's postal code
	ZipCode string
	// ServiceType is "PICKUP" or "DELIVERY"
	ServiceType ServiceType
}

// CreateBasket creates a new basket session for the given market and service type.
// The server generates the basket ID and device ID.
func CreateBasket(marketID, zipCode string, serviceType ServiceType) (*BasketSession, error) {
	if err := validateServiceType(serviceType); err != nil {
		return nil, err
	}

	body, err := json.Marshal(createBasketRequest{IncludeTimeslot: true})
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	req, err := BuildPostRequest(clientHost, "baskets", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	setBasketHeaders(req, "", marketID, zipCode, serviceType)

	var res basketResponse
	if err := DoRequest(req, &res); err != nil {
		return nil, err
	}

	basket := res.Data.Basket
	return &BasketSession{
		ID:          basket.ID,
		DeviceID:    basket.DeviceID,
		Version:     basket.Version,
		MarketID:    marketID,
		ZipCode:     zipCode,
		ServiceType: serviceType,
	}, nil
}

// GetBasket retrieves the current basket state.
// Note: Updates s.Version as a side effect for subsequent operations.
func (s *BasketSession) GetBasket() (Basket, error) {
	body, err := json.Marshal(createBasketRequest{IncludeTimeslot: true})
	if err != nil {
		return Basket{}, fmt.Errorf("error marshalling request: %w", err)
	}

	req, err := BuildPostRequest(clientHost, "baskets", bytes.NewReader(body))
	if err != nil {
		return Basket{}, err
	}

	setBasketHeaders(req, s.ID, s.MarketID, s.ZipCode, s.ServiceType)

	var res basketResponse
	if err := DoRequest(req, &res); err != nil {
		return Basket{}, err
	}

	s.Version = res.Data.Basket.Version
	return res.Data.Basket, nil
}

// SetItemQuantity sets the quantity of an item in the basket.
// Use listingId from product search results (e.g., "8-FP05LLPR-rewe-online-services|48465001-320516").
// Setting quantity to 0 removes the item.
func (s *BasketSession) SetItemQuantity(listingID string, quantity int) (Basket, error) {
	if quantity < 0 {
		return Basket{}, fmt.Errorf("quantity must be non-negative")
	}

	body, err := json.Marshal(setQuantityRequest{
		BasketVersion:   s.Version,
		Quantity:        quantity,
		IncludeTimeslot: true,
	})
	if err != nil {
		return Basket{}, fmt.Errorf("error marshalling request: %w", err)
	}

	path := fmt.Sprintf("baskets/%s/listings/%s", s.ID, url.PathEscape(listingID))
	req, err := BuildPostRequest(clientHost, path, bytes.NewReader(body))
	if err != nil {
		return Basket{}, err
	}

	setBasketHeaders(req, s.ID, s.MarketID, s.ZipCode, s.ServiceType)

	var res basketResponse
	if err := DoRequest(req, &res); err != nil {
		return Basket{}, err
	}

	s.Version = res.Data.Basket.Version
	return res.Data.Basket, nil
}

// RemoveItem removes an item from the basket entirely
func (s *BasketSession) RemoveItem(listingID string) (Basket, error) {
	path := fmt.Sprintf("baskets/%s/listings/%s", s.ID, url.PathEscape(listingID))
	req, err := BuildDeleteRequest(clientHost, path)
	if err != nil {
		return Basket{}, err
	}

	setBasketHeaders(req, s.ID, s.MarketID, s.ZipCode, s.ServiceType)

	var res basketResponse
	if err := DoRequest(req, &res); err != nil {
		return Basket{}, err
	}

	s.Version = res.Data.Basket.Version
	return res.Data.Basket, nil
}

// setBasketHeaders adds the required headers for basket operations
func setBasketHeaders(req *http.Request, basketID, marketID, zipCode string, serviceType ServiceType) {
	if basketID != "" {
		req.Header.Set("x-rd-basket-id", basketID)
	}
	setDualHeader(req, "market-id", marketID)
	setDualHeader(req, "customer-zip", zipCode)
	setDualHeader(req, "service-types", string(serviceType))
	req.Header.Set("rd-postcode", zipCode)
	req.Header.Set("rd-is-lsfk", "false")
	setCommonHeaders(req)
}
