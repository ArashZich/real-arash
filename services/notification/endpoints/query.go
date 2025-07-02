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

type NotificationQueryFilterType struct {
	Title          filter.FilterValue[string] `json:"title,omitempty"`
	Message        filter.FilterValue[string] `json:"message,omitempty"`
	Type           filter.FilterValue[string] `json:"type,omitempty"`
	UserID         filter.FilterValue[int]    `json:"user_id,omitempty"`
	CategoryID     filter.FilterValue[int]    `json:"category_id,omitempty"`
	OrganizationID filter.FilterValue[int]    `json:"organization_id,omitempty"`
	CreatedAt      filter.FilterValue[string] `json:"created_at,omitempty"`
	IsRead         filter.FilterValue[int]    `json:"is_read,omitempty"` // تغییر داده شده به int
}

type NotificationQueryRequestParams struct {
	ID      string                      `json:"id,omitempty"`
	Order   string                      `json:"order,omitempty"`
	OrderBy string                      `json:"order_by,omitempty"`
	Query   string                      `json:"query,omitempty"`
	Filters NotificationQueryFilterType `json:"filters,omitempty"`
}

func (data *NotificationQueryRequestParams) ValidateQueries(r *http.Request) govalidity.ValidityResponseErrors {
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
			"message": govalidity.Schema{
				"op":    govalidity.New("filter.message.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.message.value").Optional(),
			},
			"type": govalidity.Schema{
				"op":    govalidity.New("filter.type.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.type.value").Optional(),
			},
			"user_id": govalidity.Schema{
				"op":    govalidity.New("filter.user_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.user_id.value").Optional(),
			},
			"category_id": govalidity.Schema{
				"op":    govalidity.New("filter.category_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.category_id.value").Optional(),
			},
			"organization_id": govalidity.Schema{
				"op":    govalidity.New("filter.organization_id.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.organization_id.value").Optional(),
			},
			"created_at": govalidity.Schema{
				"op":    govalidity.New("filter.created_at.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.created_at.value").Optional(),
			},
			"is_read": govalidity.Schema{
				"op":    govalidity.New("filter.is_read.op").FilterOperators().Optional(),
				"value": govalidity.New("filter.is_read.value").Optional(),
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

func makeFilters(params NotificationQueryRequestParams) string {
	var whereOr []string
	var whereAnd []string

	if params.ID != "" {
		whereAnd = append(whereAnd, fmt.Sprintf("id = %s", params.ID))
	}

	if params.Filters.Title.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Title.Op, params.Filters.Title.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("title %s %s", opValue.Operator, val))
	}
	if params.Filters.Message.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Message.Op, params.Filters.Message.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("message %s %s", opValue.Operator, val))
	}
	if params.Filters.Type.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.Type.Op, params.Filters.Type.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("type %s %s", opValue.Operator, val))
	}
	if params.Filters.UserID.Op != "" {
		UserIDStr := strconv.Itoa(params.Filters.UserID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.UserID.Op, UserIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("user_id %s %s", opValue.Operator, val))
	}
	if params.Filters.CategoryID.Op != "" {
		CategoryIDStr := strconv.Itoa(params.Filters.CategoryID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.CategoryID.Op, CategoryIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("category_id %s %s", opValue.Operator, val))
	}
	if params.Filters.OrganizationID.Op != "" {
		OrganizationIDStr := strconv.Itoa(params.Filters.OrganizationID.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.OrganizationID.Op, OrganizationIDStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("organization_id %s %s", opValue.Operator, val))
	}
	if params.Filters.CreatedAt.Op != "" {
		opValue := filter.GetDBOperatorAndValue(params.Filters.CreatedAt.Op, params.Filters.CreatedAt.Value)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereOr = append(whereOr, fmt.Sprintf("created_at %s %s", opValue.Operator, val))
	}
	if params.Filters.IsRead.Op != "" {
		IsReadStr := strconv.Itoa(params.Filters.IsRead.Value)
		opValue := filter.GetDBOperatorAndValue(params.Filters.IsRead.Op, IsReadStr)
		val := exp.TerIf(opValue.Value == "", "", fmt.Sprintf("'%s'", opValue.Value))
		whereAnd = append(whereAnd, fmt.Sprintf("is_read %s %s", opValue.Operator, val))
	}

	// Combine OR filters with OR condition
	orCondition := strings.Join(whereOr, " OR ")

	// Combine AND filters with AND condition
	andCondition := strings.Join(whereAnd, " AND ")

	// Combine OR and AND conditions
	if orCondition != "" && andCondition != "" {
		return fmt.Sprintf("(%s) AND (%s)", orCondition, andCondition)
	} else if orCondition != "" {
		return orCondition
	} else {
		return andCondition
	}
}

func (s *service) Query(
	ctx context.Context, offset, limit int, params NotificationQueryRequestParams,
) ([]models.Notification, response.ErrorResponse) {
	var notifications []models.Notification
	where := makeFilters(params)

	tx := s.db.WithContext(ctx).Offset(offset).Limit(limit).
		Order(fmt.Sprintf("%s %s", params.OrderBy, params.Order))

	if where != "" {
		tx = tx.Where(where)
	}

	err := tx.Find(&notifications).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []models.Notification{}, response.GormErrorResponse(err, "خطایی در یافتن نوتیفیکیشن‌ها رخ داده است")
	}

	return notifications, response.ErrorResponse{}
}
