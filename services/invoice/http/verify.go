package transports

import (
	"encoding/json"
	"log"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

type FormData struct {
	ClientRefId  string      `json:"clientRefId"`
	PaymentCode  string      `json:"paymentCode"`
	PaymentRefId json.Number `json:"paymentRefId"`
	Amount       json.Number `json:"amount"`
	CardNumber   string      `json:"cardNumber"`
	CardHashPan  string      `json:"cardHashPan"`
}

func (r *resource) verify(ctx echo.Context) error {
	var input = &service.VerifyInvoiceRequest{}

	// مرحله ۲: چاپ تمام داده‌های فرم
	if err := ctx.Request().ParseForm(); err != nil {
		log.Printf("[ERROR] Error parsing form: %v", err)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid form data"))
	}

	// مرحله ۳: استخراج داده‌های فرم
	formData := ctx.Request().Form["data"]
	if len(formData) > 0 {
		var formValues FormData
		err := json.Unmarshal([]byte(formData[0]), &formValues)
		if err != nil {
			log.Printf("[ERROR] Error unmarshaling form data: %v", err)
			return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid form data"))
		}

		// مقداردهی `input`
		input.RefId = formValues.PaymentRefId.String()
		input.ClientRefId = formValues.ClientRefId

	} else {
		log.Printf("[ERROR] Missing form data 'data' key")
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Missing form data"))
	}

	// مرحله ۴: بررسی داده‌های ورودی
	errors := input.Validate(ctx.Request())
	if errors != nil {
		log.Printf("[ERROR] Validation Errors: %+v", errors)
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest(errors))
	}

	htmlResponse, err := r.Invoice.Verify(ctx.Request().Context(), *input)
	if err.StatusCode != 0 {
		log.Printf("[ERROR] Payment Verification Failed for RefId: %s, Error: %s", input.RefId, err.Message)
		return ctx.String(http.StatusBadRequest, err.Message)
	}

	// مرحله ۶: پرداخت موفق
	log.Printf("[SUCCESS] Payment Verification Successful for RefId: %s, ClientRefId: %s", input.RefId, input.ClientRefId)
	return ctx.HTML(http.StatusCreated, htmlResponse)
}
