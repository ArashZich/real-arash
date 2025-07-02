package payping

import (
	"encoding/json"

	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	PaymentRefID string `json:"paymentRefId"`
}

// Verify is the function to verify payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq, ok := vReq.(*VerifyRequest)
	if !ok {
		return nil, e.ErrInternal{
			Message: "vReq is not of type VerifyRequest",
		}
	}
	resp, _ := client.Post(APIVerifyURL, verifyReq, map[string]string{
		"Authorization": "Bearer " + d.Token,
	})
	var res map[string]interface{}
	_ = json.Unmarshal(resp.Body(), &res)
	if resp.StatusCode() != 200 {
		if res == nil {
			return nil, e.ErrInvalidPayment{
				Message: "error in verify payment",
			}
		}

		for _, v := range res {
			return nil, e.ErrInvalidPayment{
				Message: v.(string),
			}
		}
	}
	rec := receipt.NewReceipt(verifyReq.PaymentRefID, d.GetDriverName())
	rec.Detail("cardNumber", res["cardNumber"].(string))

	return rec, nil
}
