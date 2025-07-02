package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/filter"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
)

type DocumentQueryFilterType struct {
	Title       filter.FilterValue[string] `json:"title,omitempty"`
	PhoneNumber filter.FilterValue[string] `json:"phone_number"`
	CellPhone   filter.FilterValue[string] `json:"cell_phone"`
	Website     filter.FilterValue[string] `json:"website"`
	Telegram    filter.FilterValue[string] `json:"telegram"`
	Instagram   filter.FilterValue[string] `json:"instagram"`
	Linkedin    filter.FilterValue[string] `json:"linkedin"`
	Location    filter.FilterValue[string] `json:"location"`
	Size        filter.FilterValue[string] `json:"size"`
	FileURI     filter.FilterValue[string] `json:"file_uri"`
	PreviewURI  filter.FilterValue[string] `json:"Preview_uri"`
	CategoryID  filter.FilterValue[int]    `json:"category_id"`
	OwnerID     filter.FilterValue[int]    `json:"owner_id"`
	Order       filter.FilterValue[int]    `json:"order"`
	SizeMB      filter.FilterValue[int]    `json:"size_mb"`
	CreatedAt   filter.FilterValue[string] `json:"created_at"`
}

type DocumentQueryRequestParams struct {
	ID      string                  `json:"id,omitempty"`
	Order   string                  `json:"order,omitempty"`
	OrderBy string                  `json:"order_by,omitempty"`
	Query   string                  `json:"query,omitempty"`
	Filters DocumentQueryFilterType `json:"filters,omitempty"`
}

func (data *DocumentQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"id":       govalidity.New("id").Optional(),
		"query":    govalidity.New("query").Optional(),
		"order":    govalidity.New("order").Optional(),
		"order_by": govalidity.New("order_by").Optional(),
		"filters": govalidity.Schema{
			"title": govalidity.Schema{
				"op":    govalidity.New("filter.title.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.title.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"phone_number": govalidity.Schema{
				"op":    govalidity.New("filter.phone_number.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.phone_number.value").Optional(),
			},
			"website": govalidity.Schema{
				"op":    govalidity.New("filter.website.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.website.value").Optional(),
			},
			"telegram": govalidity.Schema{
				"op":    govalidity.New("filter.telegram.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.telegram.value").Optional(),
			},
			"instagram": govalidity.Schema{
				"op":    govalidity.New("filter.instagram.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.instagram.value").Optional(),
			},
			"linkedin": govalidity.Schema{
				"op":    govalidity.New("filter.linkedin.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.linkedin.value").Optional(),
			},
			"cell_phone": govalidity.Schema{
				"op":    govalidity.New("filter.cell_phone.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.cell_phone.value").Optional(),
			},
			"location": govalidity.Schema{
				"op":    govalidity.New("filter.location.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.location.value").Optional(),
			},
			"size": govalidity.Schema{
				"op":    govalidity.New("filter.size.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.size.value").Optional(),
			},
			"file_uri": govalidity.Schema{
				"op":    govalidity.New("filter.file_uri.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.file_uri.value").Optional(),
			},
			"Preview_uri": govalidity.Schema{
				"op":    govalidity.New("filter.Preview_uri.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.Preview_uri.value").Optional(),
			},
			"order": govalidity.Schema{
				"op":    govalidity.New("filter.order.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.order.value").Optional(),
			},
			"size_mb": govalidity.Schema{
				"op":    govalidity.New("filter.size_mb.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.size_mb.value").Optional(),
			},
			"product_id": govalidity.Schema{
				"op":    govalidity.New("filter.owner_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.owner_id.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
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

func makeFilters(params DocumentQueryRequestParams) []string {
	var where []string

	if params.ID != "" {
		where = append(where, fmt.Sprintf("id = %s", params.ID))
	}
	if params.Filters.Title.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, params.Filters.Title.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}
	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("category_id %s %s", opValue.Operator, val))
	}
	if params.Filters.PhoneNumber.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.PhoneNumber.Op, params.Filters.PhoneNumber.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("phone_number %s %s", opValue.Operator, val))
	}
	if params.Filters.CellPhone.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CellPhone.Op, params.Filters.CellPhone.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("cell_phone %s %s", opValue.Operator, val))
	}
	if params.Filters.Website.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Website.Op, params.Filters.Website.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("website %s %s", opValue.Operator, val))
	}
	if params.Filters.Telegram.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Telegram.Op, params.Filters.Telegram.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("telegram %s %s", opValue.Operator, val))
	}
	if params.Filters.Instagram.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Instagram.Op, params.Filters.Instagram.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("instagram %s %s", opValue.Operator, val))
	}
	if params.Filters.FileURI.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.FileURI.Op, params.Filters.FileURI.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("file_uri %s %s", opValue.Operator, val))
	}
	if params.Filters.PreviewURI.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.PreviewURI.Op, params.Filters.PreviewURI.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("preview_uri %s %s", opValue.Operator, val))
	}
	if params.Filters.Order.Op != "" {
		OrderStr := strconv.Itoa(params.Filters.Order.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.Order.Op, OrderStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("order %s %s", opValue.Operator, val))
	}
	if params.Filters.SizeMB.Op != "" {
		SizeMBStr := strconv.Itoa(params.Filters.SizeMB.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.SizeMB.Op, SizeMBStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("size_mb %s %s", opValue.Operator, val))
	}
	if params.Filters.OwnerID.Op != "" {
		OwnerIDStr := strconv.Itoa(params.Filters.OwnerID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OwnerID.Op, OwnerIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("owner_id %s %s", opValue.Operator, val))
	}
	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		where = append(where, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}

	return where

}

func (s *service) Query(
	ctx context.Context, offset, limit int, params DocumentQueryRequestParams,
) ([]models.Document, response.ErrorResponse) {
	var documents []models.Document
	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order)).
		Preload("User").Preload("Product").Preload("Document.Category").Preload("Category")

	where := makeFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Find(&documents).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Document{}, response.GormErrorResponse(err, "خطایی در یافتن سند رخ داده است")
	}

	return documents, response.ErrorResponse{}
}
