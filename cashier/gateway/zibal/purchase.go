package zibal

import (
	"encoding/json"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/templates"
)

// Purchase sends a request to zibal to purchase an bill.
func (d *Driver) Purchase(bill *bill.Bill) (string, error) {
	var reqBody = map[string]interface{}{
		"merchant":    d.Merchant,
		"callbackUrl": d.Callback,
		"description": bill.GetDescription(),
		"orderId":     bill.GetUUID(),
		"amount":      bill.GetAmount(),
	}
	if d := bill.GetDetail("phone"); d != "" {
		reqBody["mobile"] = d
	}
	resp, _ := client.Post(APIPurchaseURL, reqBody, nil)
	var res map[string]interface{}
	err := json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 || res["result"].(float64) != 100 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status() + " " + res["message"].(string),
		}
	}

	return strconv.Itoa(int(res["trackId"].(float64))), nil
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
