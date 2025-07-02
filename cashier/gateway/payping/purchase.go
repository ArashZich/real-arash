package payping

import (
	"encoding/json"
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/cashier/bill"
	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/templates"
	"github.com/labstack/gommon/log"
)

// Purchase send purchase request to payping
func (d *Driver) Purchase(bill *bill.Bill) (string, error) {
	var reqBody = map[string]interface{}{
		"returnUrl":   d.Callback,
		"description": bill.GetDescription(),
		"amount":      bill.GetAmount(),
		"clientRefId": bill.GetUUID(),
	}
	if d := bill.GetDetail("phone"); d != "" {
		reqBody["payerIdentity"] = d
	} else if d := bill.GetDetail("email"); d != "" {
		reqBody["payerIdentity"] = d
	}
	if d := bill.GetDetail("name"); d != "" {
		reqBody["payerName"] = d
	}
	fmt.Printf("%v", reqBody)
	resp, err := client.Post(APIPurchaseURL, reqBody, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})

	if err != nil {
		log.Info(fmt.Sprintf("PayPing: Purchase failed with error %s", err.Error()))
		return "", err
	}
	if resp.StatusCode() != 200 {
		log.Info(fmt.Sprintf("PayPing: Purchase failed with status code %d", resp.StatusCode()))
		log.Info(fmt.Sprintf("PayPing: Purchase failed with body %s", resp.Body()))
		return "", e.ErrPurchaseFailed{
			Message: "Purchase failed",
		}
	}
	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return "", err
	}
	paymentCode := res["paymentCode"].(string)
	// paymentURL := res["url"].(string)
	return paymentCode, nil
}

// PayURL return pay url
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
