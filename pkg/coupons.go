package rewerse

import (
	"fmt"
	"strings"
	"time"
)

// OneScanCoupons is a struct for the raw data returned by the OneScan-Coupons endpoint
type OneScanCoupons struct {
	AppID   string `json:"appId"`
	Code    string `json:"code"`
	Coupons []struct {
		ImageURL       string `json:"imageUrl"`
		Description    string `json:"description"`
		Banner         string `json:"banner"`
		Title          string `json:"title"`
		Subtitle       string `json:"subtitle"`
		Provider       string `json:"provider"`
		Validity       string `json:"validity"`
		InformationURL string `json:"informationUrl"`
		RawValues      struct {
			ID string `json:"id"`
		} `json:"rawValues"`
	} `json:"coupons"`
}

func (o OneScanCoupons) String() string {
	coupons := fmt.Sprintf("Got %d OneScan-Coupons\n\n", len(o.Coupons))
	for _, c := range o.Coupons {
		coupons += c.Title + "\n"
		coupons += strings.ReplaceAll(c.Description, "\n", " ")
		coupons += c.Validity + "\n"
		coupons += c.Provider + "\n"
		if c.InformationURL != "" {
			coupons += c.InformationURL + "\n"
		}
		coupons += "\n"
	}
	return coupons
}

// Coupons is a struct for the raw data returned by the Coupons endpoint
type Coupons struct {
	Data struct {
		GetCoupons struct {
			CouponStatus string `json:"couponStatus"`
			PaybackEwe   any    `json:"paybackEwe"`
			Coupons      []struct {
				CouponID   int    `json:"couponId"`
				CouponType string `json:"couponType"`
				Title      string `json:"title"`
				Subtitle   string `json:"subtitle"`
				Validity   struct {
					ValidFrom time.Time `json:"validFrom"`
					ValidTo   time.Time `json:"validTo"`
				} `json:"validity"`
				ProductLogo string `json:"productLogo"`
				Description struct {
					RedeemDescription any    `json:"redeemDescription"`
					Combinability     string `json:"combinability"`
					Validity          any    `json:"validity"`
					TermsOfUse        string `json:"termsOfUse"`
					Disclaimer        any    `json:"disclaimer"`
				} `json:"description"`
				OfferTitle            string `json:"offerTitle"`
				DisplayClassification any    `json:"displayClassification"`
				CouponDetails         any    `json:"couponDetails"`
				Activated             any    `json:"activated"`
				Preview               any    `json:"preview"`
				IsNew                 any    `json:"isNew"`
			} `json:"coupons"`
		} `json:"getCoupons"`
	} `json:"data"`
	Extensions struct {
		HTTP []struct {
			Path         []string `json:"path"`
			Message      string   `json:"message"`
			StatusCode   int      `json:"statusCode"`
			ResponseBody any      `json:"responseBody"`
		} `json:"http"`
	} `json:"extensions"`
}

func (cs Coupons) String() string {
	coupons := fmt.Sprintf("Got %d Coupons\n\n", len(cs.Data.GetCoupons.Coupons))
	for _, c := range cs.Data.GetCoupons.Coupons {
		coupons += c.Title + " // " + c.OfferTitle + "\n"
		coupons += c.Subtitle + "\n"
		coupons += "\n"
	}
	return coupons
}

// GetCoupons returns all available coupons
// There are two differerent coupons-structs: OneScanCoupons and Coupons
// I'm not sure what the difference is and I don't really care either. Please let me know!
func GetCoupons() (oc OneScanCoupons, c Coupons, err error) {
	req, err := BuildCustomRequest(apiHost, "v3/onescan/coupons")
	if err != nil {
		return
	}

	err = DoRequest(req, &oc)
	if err != nil {
		return
	}

	reqC, err := BuildCustomRequest(clientHost, "coupons")
	if err != nil {
		return
	}

	err = DoRequest(reqC, &c)
	if err != nil {
		return
	}

	return
}
