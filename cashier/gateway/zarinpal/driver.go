package zarinpal

import (
	httpClient "gitag.ir/armogroup/armo/services/reality/cashier/pkg/http"
)

// Driver config struct for zarinpal driver
type Driver struct {
	MerchantID string
	Callback   string
}

// Const's for zarinpal driver
const (
	APIPurchaseURL = "https://api.zarinpal.com/pg/v4/payment/request.json"
	APIVerifyURL   = "https://api.zarinpal.com/pg/v4/payment/verify.json"
	APIPaymentURL  = "https://www.zarinpal.com/pg/StartPay/"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns driver name
func (Driver) GetDriverName() string {
	return "ZarinPal"
}

// SetClient sets the http client
func (Driver) SetClient(c httpClient.Client) {
	client = c
}
