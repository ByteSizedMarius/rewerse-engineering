package rewerse

// basketResponse wraps the API response for basket operations
type basketResponse struct {
	Data struct {
		Basket Basket `json:"basket"`
	} `json:"data"`
}

// Basket represents a shopping basket/cart
type Basket struct {
	// ID is the server-generated basket UUID: "9d134eb6-d29a-40b2-be46-94f3ea3424d0"
	ID string `json:"id"`
	// Version is used for optimistic locking, increments on each change
	Version int `json:"version"`
	// OrderID is set when basket is converted to an order
	OrderID *string `json:"orderId,omitempty"`
	// DeviceID is the server-assigned device identifier
	DeviceID string `json:"deviceId"`

	ServiceSelection      ServiceSelection      `json:"serviceSelection"`
	ServiceConfiguration  ServiceConfiguration  `json:"serviceConfiguration"`
	Staggerings           Staggerings           `json:"staggerings"`
	LineItems             []LineItem            `json:"lineItems"`
	Violations            []Violation           `json:"violations"`
	Summary               BasketSummary         `json:"summary"`
	TimeSlotInformation   TimeSlotInformation   `json:"timeSlotInformation"`
}

// ServiceSelection contains the selected market and service type
type ServiceSelection struct {
	// WWIdent is the market ID: "831002"
	WWIdent string `json:"wwIdent"`
	// ServiceType is "PICKUP" or "DELIVERY"
	ServiceType string `json:"serviceType"`
	// ZipCode is the customer's postal code: "67065"
	ZipCode string `json:"zipCode"`
}

// ServiceConfiguration contains service-specific settings
type ServiceConfiguration struct {
	// MinimumOrderAmount in cents, e.g. 5000 = 50.00 EUR for delivery
	MinimumOrderAmount int `json:"minimumOrderAmount"`
}

// Staggerings contains pricing tier information
type Staggerings struct {
	ReachedStaggering *Staggering `json:"reachedStaggering,omitempty"`
	NextStaggering    *Staggering `json:"nextStaggering,omitempty"`
}

// Staggering represents a minimum order amount threshold
type Staggering struct {
	// ArticlePriceThreshold is the minimum order amount in cents
	ArticlePriceThreshold int `json:"articlePriceThreshold"`
	// RemainingArticlePrice in cents (how much more needed to reach minimum)
	RemainingArticlePrice int `json:"remainingArticlePrice,omitempty"`
	// DisplayText is the German description: "Mindestbestellwert: 40,00 €"
	DisplayText string `json:"displayText"`
}

// LineItem represents a product in the basket
type LineItem struct {
	// Quantity is the number of items
	Quantity int `json:"quantity"`
	// Price is the unit price in cents: 459 = 4.59 EUR
	Price int `json:"price"`
	// TotalPrice is quantity * price in cents
	TotalPrice int `json:"totalPrice"`
	// Grammage is the weight/volume info: "150g (1 kg = 30,60 EUR)"
	Grammage string `json:"grammage"`
	// Product contains product details
	Product LineItemProduct `json:"product"`
	// Violations are issues with this line item
	Violations []Violation `json:"violations"`
}

// LineItemProduct contains product info within a line item
type LineItemProduct struct {
	ProductID   string `json:"productId"`
	Title       string `json:"title"`
	ImageURL    string `json:"imageURL"`
	ArticleID   string `json:"articleId"`
	NAN         string `json:"nan"`
	OrderLimit  int    `json:"orderLimit"`
	Listing     Listing `json:"listing"`
	Attributes  ProductAttributes `json:"attributes"`
}

// Listing contains listing-specific info (market-specific pricing)
type Listing struct {
	// ListingID is used for basket operations: "8-FP05LLPR-rewe-online-services|48465001-320516"
	ListingID          string `json:"listingId"`
	ListingVersion     int    `json:"listingVersion"`
	CurrentRetailPrice int    `json:"currentRetailPrice"`
	TotalRefundPrice   int    `json:"totalRefundPrice"`
	Grammage           string `json:"grammage"`
}

// ProductAttributes contains product flags
type ProductAttributes struct {
	IsBulkyGood     bool `json:"isBulkyGood"`
	IsOrganic       bool `json:"isOrganic"`
	IsVegan         bool `json:"isVegan"`
	IsVegetarian    bool `json:"isVegetarian"`
	IsDairyFree     bool `json:"isDairyFree"`
	IsGlutenFree    bool `json:"isGlutenFree"`
	IsAgeRestricted bool `json:"isAgeRestricted"`
	IsRegional      bool `json:"isRegional"`
	IsNew           bool `json:"isNew"`
}

// Violation represents an issue with the basket or line item
type Violation struct {
	// ID is the violation type: "minimum.delivery.not.reached"
	ID string `json:"id"`
	// Message is the user-facing text: "Mindestbestellwert 50 EUR nicht erreicht!"
	Message string `json:"message"`
	// DetailMessage provides additional context
	DetailMessage *string `json:"detailMessage,omitempty"`
}

// BasketSummary contains basket totals
type BasketSummary struct {
	// ArticleCount is the number of distinct products
	ArticleCount int `json:"articleCount"`
	// ArticlePrice is the subtotal in cents
	ArticlePrice int `json:"articlePrice"`
	// TotalArticleQuantity is the sum of all quantities
	TotalArticleQuantity int `json:"totalArticleQuantity"`
	// TotalPrice is the final price including fees in cents
	TotalPrice int `json:"totalPrice"`
	// TotalPriceExclServiceFee is total minus service fee
	TotalPriceExclServiceFee int `json:"totalPriceExclServiceFee"`
	// Fees contains fee breakdowns
	Fees BasketFees `json:"fees"`
}

// BasketFees contains optional fee components
type BasketFees struct {
	BeverageSurcharge     *int `json:"beverageSurcharge,omitempty"`
	ReusableBagSurcharge  *int `json:"reusableBagSurcharge,omitempty"`
	TransportBoxSurcharge *int `json:"transportBoxSurcharge,omitempty"`
	ServiceFee            *int `json:"serviceFee,omitempty"`
	Refund                *int `json:"refund,omitempty"`
}

// TimeSlotInformation contains delivery/pickup time slot info
type TimeSlotInformation struct {
	StartTime     *string `json:"startTime,omitempty"`
	EndTime       *string `json:"endTime,omitempty"`
	TimeSlotPrice *int    `json:"timeSlotPrice,omitempty"`
	TimeSlotText  string  `json:"timeSlotText"`
}

// createBasketRequest is the POST body for creating a basket
type createBasketRequest struct {
	IncludeTimeslot bool `json:"includeTimeslot"`
}

// setQuantityRequest is the POST body for setting item quantity
type setQuantityRequest struct {
	BasketVersion   int  `json:"basketVersion"`
	Quantity        int  `json:"quantity"`
	IncludeTimeslot bool `json:"includeTimeslot"`
}
