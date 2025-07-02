package asanpardakht

import (
	"time"

	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/templates"
)

// Purchase send purchase request to asanpardakht gateway
func (d *Driver) Purchase(bill *bill.Bill) (string, error) {
	var reqBody = map[string]interface{}{
		"callbackURL":             d.Callback,
		"additionalData":          bill.GetDescription(),
		"amountInRials":           bill.GetAmount() * 10,
		"localInvoiceId":          bill.GetUUID(),
		"serviceTypeId":           1,
		"localDate":               time.Now().Format("20060102 150405"),
		"merchantConfigurationId": d.MerchantConfigID,
		"paymentId":               0,
	}

	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"usr": d.Username,
		"pwd": d.Password,
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status(),
		}
	}
	return string(resp.Body()), nil
}

// PayURL return pay url
func (*Driver) PayURL(_ *bill.Bill) string {
	return APIPaymentURL
}

// PayMethod returns the Request Method to be used to pay the bill.
func (*Driver) PayMethod() string {
	return "POST"
}

// RenderRedirectForm renders the html form for redirect to payment page.
func (d *Driver) RenderRedirectForm(bill *bill.Bill) (string, error) {
	return templates.RenderRedirectTemplate(d.PayMethod(), d.PayURL(bill), map[string]string{
		"RefId": bill.GetUUID(),
	})
}
