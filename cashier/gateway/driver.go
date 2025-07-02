package gateway

import (
	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	httpClient "gitag.ir/armogroup/armo/services/reality/cashier/pkg/http"
	"gitag.ir/armogroup/armo/services/reality/cashier/receipt"
)

// Driver is the interface that must be implemented by all drivers.
type Driver interface {
	// Purchase sends a purchase request to the driver's gateway.
	Purchase(bill *bill.Bill) (string, error)
	// PayURL returns the url to redirect the user to for payment.
	PayURL(bill *bill.Bill) string
	// GetDriverName returns the name of the driver.
	GetDriverName() string
	// Verify checks the payment status of the github.com/mohammadv184/gopayment/bill
	Verify(vReq interface{}) (*receipt.Receipt, error)
	// PayMethod returns the payment request method.
	PayMethod() string
	// SetClient sets the http client.
	SetClient(client httpClient.Client)
	// RenderRedirectForm renders the html form for redirect to payment page.
	RenderRedirectForm(bill *bill.Bill) (string, error)
}
