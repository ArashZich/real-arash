package idpay

import (
	"encoding/json"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/templates"
)

// Purchase send purchase request to idpay gateway
func (d *Driver) Purchase(bill *bill.Bill) (string, error) {
	var reqBody = map[string]interface{}{
		"callback": d.Callback,
		"desc":     bill.GetDescription(),
		"amount":   bill.GetAmount(),
		"order_id": bill.GetUUID(),
	}
	if d := bill.GetDetail("phone"); d != "" {
		reqBody["phone"] = d
	}
	if d := bill.GetDetail("email"); d != "" {
		reqBody["mail"] = d
	}
	if d := bill.GetDetail("name"); d != "" {
		reqBody["name"] = d
	}
	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"X-API-KEY": d.MerchantID,
		"X-SANDBOX": strconv.FormatBool(d.Sandbox),
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 201 {
		return "", e.ErrPurchaseFailed{
			Message: resp.Status(),
		}
	}
	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	return res["id"].(string), nil
}

// PayURL return pay url
func (d *Driver) PayURL(bill *bill.Bill) string {
	if d.Sandbox {
		return APISandBoxPaymentURL + bill.GetTransactionID()
	}
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
