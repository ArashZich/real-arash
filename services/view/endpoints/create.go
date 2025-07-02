package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CreateViewRequest struct {
	Name           string    `json:"name,omitempty"`
	Ip             string    `json:"ip"`
	BrowserAgent   string    `json:"browser_agent"`
	OperatingSys   string    `json:"operating_sys"`
	Device         string    `json:"device"`
	IsAR           bool      `json:"is_ar"`
	Is3D           bool      `json:"is_3d"`
	IsVR           bool      `json:"is_vr"`
	Url            string    `json:"url"`
	ProductUID     uuid.UUID `json:"product_uid"`
	OrganizationID int       `json:"organization_id"` // Add this line
	VisitUID       uuid.UUID `json:"visit_uid"`
}

type IPAPIResponse struct {
	RegionName string `json:"regionName"` // Ensure this matches the API's response
}

func GetCityByIP(ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipInfo IPAPIResponse
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}

	return ipInfo.RegionName, nil
}

func (c *CreateViewRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":            govalidity.New("name").Optional().MinMaxLength(2, 100),
		"ip":              govalidity.New("ip").Optional(),
		"browser_agent":   govalidity.New("browser_agent").Optional().MinMaxLength(2, 100),
		"operating_sys":   govalidity.New("operating_sys").Optional().MinMaxLength(2, 100),
		"device":          govalidity.New("device").Optional().MinMaxLength(2, 100),
		"is_ar":           govalidity.New("is_ar").Optional(),
		"is_3d":           govalidity.New("is_3d").Optional(),
		"is_vr":           govalidity.New("is_vr").Optional(),
		"url":             govalidity.New("url").Optional(),
		"product_uid":     govalidity.New("product_uid").Required(),
		"organization_id": govalidity.New("organization_id").Required(), // Add validation rule
		"visit_uid":       govalidity.New("visit_uid").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":            "نام",
			"ip":              "آی پی",
			"browser_agent":   "مرورگر",
			"operating_sys":   "سیستم عامل",
			"device":          "دستگاه",
			"is_ar":           "AR",
			"is_3d":           "3D",
			"is_vr":           "VR",
			"url":             "آدرس",
			"product_uid":     "محصول",
			"organization_id": "سازمان", // Add this label
			"visit_uid":       "شناسه بازدید",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Create(ctx context.Context, input CreateViewRequest) (models.View, response.ErrorResponse) {

	// Id := policy.ExtractIdClaim(ctx)
	// id, _ := strconv.Atoi(Id)

	var organization models.Organization
	// Retrieve and validate the organization using input.OrganizationID
	err := s.db.WithContext(ctx).First(&organization, input.OrganizationID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	// Initialize regionName as empty and fetch only if IP is provided
	var regionName string
	if input.Ip != "" {
		var err error
		regionName, err = GetCityByIP(input.Ip)
		if err != nil {
			s.logger.With(ctx).Error("Failed to fetch geolocation data: ", err)
			// Optionally handle the error further or log it
		}
	}

	var product models.Product
	view := models.View{
		// UserID:       id,
		Name:           input.Name,
		Ip:             dtp.NullString{String: input.Ip, Valid: input.Ip != ""},
		BrowserAgent:   input.BrowserAgent,
		OperatingSys:   input.OperatingSys,
		Device:         input.Device,
		IsAR:           input.IsAR,
		Is3D:           input.Is3D,
		IsVR:           input.IsVR,
		Url:            dtp.NullString{String: input.Url, Valid: input.Url != ""},
		ProductUID:     input.ProductUID,
		VisitUID:       input.VisitUID,
		OrganizationID: input.OrganizationID, // Set the OrganizationID
		RegionName:     regionName,           // Add this line

	}

	err = s.db.WithContext(ctx).Create(&view).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "خطایی در ایجاد بازرسی رخ داد")
	}

	err = s.db.WithContext(ctx).First(&product, "product_uid", input.ProductUID).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "خطایی در یافتن محصول رخ داده است")
	}
	view.ProductID = product.ID
	// update product view
	product.ViewCount = product.ViewCount + 1
	err = s.db.WithContext(ctx).Save(&product).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.View{}, response.GormErrorResponse(err, "خطایی در بروزرسانی محصول رخ داده است")
	}
	return view, response.ErrorResponse{}
}
