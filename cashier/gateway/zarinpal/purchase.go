package zarinpal

import (
	"encoding/json"

	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/templates"
)

// Purchase sends a request to Zarinpal to purchase an bill.
func (d *Driver) Purchase(bill *bill.Bill) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant_id":  d.MerchantID,
		"callback_url": d.Callback,
		"description":  bill.GetDescription(),
		"amount":       bill.GetAmount(),
		"metadata":     bill.GetDetails(),
	}
	resp, _ := client.Post(APIPurchaseURL, reqBody, nil)
	if resp.StatusCode() != 100 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status() + " purchase failed",
		}
	}
	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["data"].(map[string]interface{})["authority"].(string), nil
}

// PayURL returns the url to redirect the user to in order to pay the bill.
func (*Driver) PayURL(bill *bill.Bill) string {
	return APIPaymentURL + bill.GetTransactionID()
}

// PayMethod returns the Request Method to be used to pay the bill.
func (*Driver) PayMethod() string {
	return "GET"
}

// RenderRedirectForm renders the html form for redirect to payment page.
func (d *Driver) RenderRedirectForm(bill *bill.Bill) (string, error) {
	return templates.RenderRedirectTemplate(d.PayMethod(), d.PayURL(bill), nil)
}
