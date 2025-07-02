package endpoints

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
)

type ViewQueryFilterType struct {
	Name           filter.FilterValue[string] `json:"name,omitempty"`
	Ip             filter.FilterValue[string] `json:"ip"`
	BrowserAgent   filter.FilterValue[string] `json:"browser_agent"`
	OperatingSys   filter.FilterValue[string] `json:"operating_sys"`
	Device         filter.FilterValue[string] `json:"device"`
	IsAR           filter.FilterValue[string] `json:"is_ar"`
	Is3D           filter.FilterValue[string] `json:"is_3d"`
	IsVR           filter.FilterValue[string] `json:"is_vr"`
	Url            filter.FilterValue[string] `json:"url"`
	ProductID      filter.FilterValue[int]    `json:"product_id"`
	CreatedAt      filter.FilterValue[string] `json:"created_at"`
	ProductUID     filter.FilterValue[string] `json:"product_uid"`
	OrganizationID filter.FilterValue[int]    `json:"organization_id"` // Add this line

}

type ViewQueryRequestParams struct {
	ID       string              `json:"id,omitempty"`
	Order    string              `json:"order,omitempty"`
	OrderBy  string              `json:"order_by,omitempty"`
	Query    string              `json:"query,omitempty"`
	Duration string              `json:"duration,omitempty"`
	Filters  ViewQueryFilterType `json:"filters,omitempty"`
}

type ExtendedViewResponse struct {
	Views         []models.View `json:"views"`
	Is3DLen       int           `json:"is_3d_len"`
	IsARLen       int           `json:"is_ar_len"`
	Total         int           `json:"total"`
	Browsers      []string      `json:"browsers"`
	OperatingSys  []string      `json:"operating_sys"`
	IPs           []string      `json:"ips"`
	VisitDuration int64         `json:"visit_duration"`
	RegionName    []string      `json:"region_name"` // Add this line
}

func (data *ViewQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"duration": govalidity.New("duration").In([]string{"one_week", "one_month", "three_months", "six_months", "one_year"}).Optional(),
		"filters": govalidity.Schema{
			"name": govalidity.Schema{
				"op":    govalidity.New("filter.name.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.name.value").Optional(),
			},
			"ip": govalidity.Schema{
				"op":    govalidity.New("filter.ip.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.ip.value").Optional(),
			},
			"browser_agent": govalidity.Schema{
				"op":    govalidity.New("filter.browser_agent.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.browser_agent.value").Optional(),
			},
			"operating_sys": govalidity.Schema{
				"op":    govalidity.New("filter.operating_sys.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.operating_sys.value").Optional(),
			},
			"device": govalidity.Schema{
				"op":    govalidity.New("filter.device.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.device.value").Optional(),
			},
			"is_ar": govalidity.Schema{
				"op":    govalidity.New("filter.is_ar.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.is_ar.value").Optional(),
			},
			"is_3d": govalidity.Schema{
				"op":    govalidity.New("filter.is_3d.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.is_3d.value").Optional(),
			},
			"is_vr": govalidity.Schema{
				"op":    govalidity.New("filter.is_vr.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.is_vr.value").Optional(),
			},
			"url": govalidity.Schema{
				"op":    govalidity.New("filter.url.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.url.value").Optional(),
			},
			"product_id": govalidity.Schema{
				"op":    govalidity.New("filter.product_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.product_id.value").Optional(),
			},
			"time_range": govalidity.Schema{
				"op":    govalidity.New("filter.time_range.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.time_range.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"product_uid": govalidity.Schema{
				"op":    govalidity.New("filter.product_uid.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.product_uid.value").Optional(),
			},
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
		},
	}
	errs := govalidity.ValidateQueries(r, schema, data)
	if len(errs) > 0 {
		dumpedErrors := govalidity.DumpErrors(errs)
		return dumpedErrors
	}

	return nil
}

func makeFilters(params ViewQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Name.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Name.Op, params.Filters.Name.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("name %s %s", opValue.Operator, val))
	}
	if params.Filters.Ip.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Ip.Op, params.Filters.Ip.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("ip %s %s", opValue.Operator, val))
	}
	if params.Filters.BrowserAgent.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.BrowserAgent.Op, params.Filters.BrowserAgent.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("browser_agent %s %s", opValue.Operator, val))
	}
	if params.Filters.OperatingSys.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.OperatingSys.Op, params.Filters.OperatingSys.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("operating_sys %s %s", opValue.Operator, val))
	}
	if params.Filters.Device.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Device.Op, params.Filters.Device.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("device %s %s", opValue.Operator, val))
	}
	if params.Filters.IsAR.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.IsAR.Op, params.Filters.IsAR.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("is_ar %s %s", opValue.Operator, val))
	}
	if params.Filters.Is3D.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Is3D.Op, params.Filters.Is3D.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("is_3d %s %s", opValue.Operator, val))
	}
	if params.Filters.IsVR.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.IsVR.Op, params.Filters.IsVR.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("is_vr %s %s", opValue.Operator, val))
	}
	if params.Filters.Url.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Url.Op, params.Filters.Url.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("url %s %s", opValue.Operator, val))
	}
	if params.Filters.ProductID.Op != "" {
		ProductIDStr := strconv.Itoa(params.Filters.ProductID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.ProductID.Op, ProductIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("product_id %s %s", opValue.Operator, val))
	}
	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}
	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}
	if params.Filters.ProductUID.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.ProductUID.Op, params.Filters.ProductUID.Value)

		// Assuming ProductUID is a string representation of a UUID
		// Validate or convert it to a proper UUID format as needed
		uuidVal, err := uuid.Parse(opValue.Value)
		if err != nil {
			log.Printf("Error parsing ProductUID: %v", err)
			return nil // or handle the error as required
		}

		where = append(where, fmt.Sprintf("product_uid %s '%s'", opValue.Operator, uuidVal.String()))
	}

	return where

}

func makeDuration(params ViewQueryRequestParams) []string {
	var where []string
	if params.Duration != "" {
		now := time.Now()
		var start, end string

		switch params.Duration {
		case "one_week":
			start = now.AddDate(0, 0, -7).Format("2006-01-02") + " 00:00:00"
			end = now.Format("2006-01-02") + " 23:59:59"
		case "one_month":
			start = now.AddDate(0, -1, 0).Format("2006-01-02") + " 00:00:00"
			end = now.Format("2006-01-02") + " 23:59:59"
		case "three_months":
			start = now.AddDate(0, -3, 0).Format("2006-01-02") + " 00:00:00"
			end = now.Format("2006-01-02") + " 23:59:59"
		case "six_months":
			start = now.AddDate(0, -6, 0).Format("2006-01-02") + " 00:00:00"
			end = now.Format("2006-01-02") + " 23:59:59"
		case "one_year":
			start = now.AddDate(-1, 0, 0).Format("2006-01-02") + " 00:00:00"
			end = now.Format("2006-01-02") + " 23:59:59"
		default:
			log.Printf("Invalid duration value: %s", params.Duration)
			return nil
		}

		where = append(where, fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", start, end))
	}
	return where
}

func (s *service) Query(
	ctx context.Context, offset, limit int, params ViewQueryRequestParams,
) (ExtendedViewResponse, response.ErrorResponse) {
	var views []models.View

	// Check permissions
	if !policy.CanGetViews(ctx) {
		s.logger.With(ctx).Error("شما دسترسی مشاهده بازدید را ندارید")
		return ExtendedViewResponse{}, response.ErrorForbidden("شما دسترسی مشاهده بازدید را ندارید")
	}

	// Log the received OrganizationID for debugging purposes
	// log.Printf("Received OrganizationID: %d", params.Filters.OrganizationID.Value)

	tx := s.db.WithContext(ctx).Model(&models.View{}).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order))

	var where []string
	if params.Duration != "" {
		durationWhere := makeDuration(params)
		if durationWhere != nil {
			where = append(where, durationWhere...)
		}
	}

	filtersWhere := makeFilters(params)
	if filtersWhere != nil {
		where = append(where, filtersWhere...)
	}

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	// Log the constructed SQL query
	// log.Printf("Constructed SQL Query: %s", tx.Statement.SQL.String())

	err := tx.Find(&views).Error
	// log.Printf("Executed SQL Query: %s", tx.Statement.SQL.String()) // Log the SQL query

	if err != nil {
		s.logger.With(ctx).Error(err)
		return ExtendedViewResponse{}, response.GormErrorResponse(err, "خطایی در یافتن بازرسی رخ داده است")
	}

	// Compute additional fields
	var is3DLen, isARLen int
	var visitDuration int64
	var browsers, operatingSystems, regionName, ips []string

	for _, view := range views {
		if view.Is3D {
			is3DLen++
		}
		if view.IsAR {
			isARLen++
		}
		// Check if IP is valid and non-empty before appending
		if view.Ip.Valid && view.Ip.String != "" {
			ips = append(ips, view.Ip.String)
		}
		visitDuration += view.VisitDuration
		// Append non-empty region names only
		if view.RegionName != "" {
			regionName = append(regionName, view.RegionName)
		}
		// Append non-empty browser agents only
		if view.BrowserAgent != "" {
			browsers = append(browsers, view.BrowserAgent)
		}
		// Append non-empty operating systems only
		if view.OperatingSys != "" {
			operatingSystems = append(operatingSystems, view.OperatingSys)
		}
	}
	var validViewsCount int
	for _, view := range views {
		// Hypothetical condition: Only count views with non-empty RegionName
		if view.RegionName != "" {
			validViewsCount++
		}
		// Other aggregation logic...
	}

	// Here is the change: Set the Total field to the length of the views slice
	totalViews := validViewsCount // Get the count of total views

	// Wrap the original views and computed fields in the new response struct
	extendedResponse := ExtendedViewResponse{
		Views:         views,
		Is3DLen:       is3DLen,
		IsARLen:       isARLen,
		Browsers:      browsers,
		OperatingSys:  operatingSystems,
		IPs:           ips,
		Total:         totalViews, // Set the total count of views here
		VisitDuration: visitDuration,
		RegionName:    regionName, // Add this line
	}

	return extendedResponse, response.ErrorResponse{}
}
