package zarinpal

import (
	"encoding/json"

	e "gitag.ir/armogroup/armo/services/reality/cashier/errors"
	"gitag.ir/armogroup/armo/services/reality/cashier/receipt"
)

// VerifyRequest is the request struct for verify
type VerifyRequest struct {
	Amount    string `json:"Amount"`
	Authority string `json:"Authority"`
}

// Verify is the function to verify a payment
func (d *Driver) Verify(vReq interface{}) (*receipt.Receipt, error) {
	verifyReq, ok := vReq.(*VerifyRequest)
	if !ok {
		return nil, e.ErrInternal{
			Message: "vReq is not of type VerifyRequest",
		}
	}
	resp, err := client.Post(APIVerifyURL, verifyReq, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 100 {
		return nil, e.ErrInvalidPayment{
			Message: resp.Status() + " Invalid payment",
		}
	}

	var res map[string]interface{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}
	rec := receipt.NewReceipt(res["data"].(map[string]interface{})["ref_id"].(string), d.GetDriverName())
	return rec, nil
}
