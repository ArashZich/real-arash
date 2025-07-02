// Package gocashier multi gateway cashier package for Golang
package cashier

import (
	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	"gitag.ir/armogroup/armo/services/reality/cashier/gateway"
	"gitag.ir/armogroup/armo/services/reality/cashier/gateway/payping"
	httpClient "gitag.ir/armogroup/armo/services/reality/cashier/pkg/http"
	"gitag.ir/armogroup/armo/services/reality/cashier/receipt"
	"gitag.ir/armogroup/armo/services/reality/config"
)

// Version is the version of the cashier module
const Version = "v1.9.0"

// Cashier is the payment main struct of the cashier module
type Cashier struct {
	driver gateway.Driver
	bill   *bill.Bill
}

// Amount set the amount of cashier bill
func (c *Cashier) SetAmount(amount int) *Cashier {
	c.bill.SetAmount(uint32(amount))
	return c
}

// Purchase send the purchase request to the cashier gateway
func (c *Cashier) Purchase() error {
	transID, err := c.driver.Purchase(c.bill)
	if err != nil {
		return err
	}
	c.bill.SetTransactionID(transID)
	return nil
}

// PayURL return the cashier URL
func (c *Cashier) PayURL() string {
	return c.driver.PayURL(c.bill)
}

// PayMethod returns the Request Method to be used to pay the bill.
func (c *Cashier) PayMethod() string {
	return c.driver.PayMethod()
}

// Client sets the driver http client.
func (c *Cashier) SetClient(client httpClient.Client) *Cashier {
	c.driver.SetClient(client)
	return c
}

// GetBill return the cashier bill
func (c *Cashier) GetBill() *bill.Bill {
	return c.bill
}

// GetTransactionID return the cashier transaction id
func (c *Cashier) GetTransactionID() string {
	return c.bill.GetTransactionID()
}

// Description set the cashier description
func (c *Cashier) SetDescription(description string) *Cashier {
	c.bill.SetDescription(description)
	return c
}

// Detail set the cashier detail
func (c *Cashier) SetDetail(key, value string) *Cashier {
	c.bill.Detail(key, value)
	return c
}

// RenderRedirectForm renders the html form for redirect to cashier page.
func (c *Cashier) RenderRedirectForm() (string, error) {
	return c.driver.RenderRedirectForm(c.bill)
}

func (c *Cashier) Verify(vReq interface{}) (*receipt.Receipt, error) {
	return c.driver.Verify(vReq)
}

// TODO: make config inside this package to later set it from outside

// NewCashier create a new cashier
func NewCashier(drv string) *Cashier {
	var (
		PayPingToken    = config.AppConfig.PayPingToken
		PayPingCallback = config.AppConfig.PayPingCallback
	)

	// var (
	// 	IDPayMerchantID = config.AppConfig.IDPayMerchantID
	// 	IDPayCallback   = config.AppConfig.IDPayCallback
	// 	IDPaySandbox    = config.AppConfig.IDPaySandbox
	// )

	bill := bill.NewBill()

	var driver gateway.Driver
	if drv == "payping" {
		driver = payping.NewDriver(PayPingToken, PayPingCallback)
	}
	// if drv == "idpay" {
	// 	driver = idpay.NewDriver(IDPayMerchantID, IDPayCallback, IDPaySandbox)
	// }

	return &Cashier{
		driver,
		bill,
	}
}
