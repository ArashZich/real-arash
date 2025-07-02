package payping

import (
	httpClient "gitag.ir/armogroup/armo/services/reality/cashier/pkg/http"
)

// Driver configures the payping driver
type Driver struct {
	Token    string
	Callback string
}

// Const's for payping
const (
	APIPurchaseURL = "https://api.payping.ir/v3/pay"
	APIPaymentURL  = "https://api.payping.ir/v3/pay/start/"
	APIVerifyURL   = "https://api.payping.ir/v3/pay/verify"
)

var client httpClient.Client

func init() {
	client = httpClient.NewHTTP()
}

// GetDriverName returns the name of the driver
func (Driver) GetDriverName() string {
	return "PayPing"
}

// SetClient sets the http client
func (Driver) SetClient(c httpClient.Client) {
	client = c
}

func NewDriver(token, callback string) *Driver {
	return &Driver{
		Token:    token,
		Callback: callback,
	}
}
