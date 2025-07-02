package endpoints

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/ARmo-BigBang/kit/exp"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type DocumentData struct {
	Title          string    `json:"title"`
	CategoryID     int       `json:"category_id"`
	PhoneNumber    string    `json:"phone_number"`
	CellPhone      string    `json:"cell_phone"`
	Website        string    `json:"website"`
	Telegram       string    `json:"telegram"`
	Instagram      string    `json:"instagram"`
	Linkedin       string    `json:"linkedin"`
	Location       string    `json:"location"`
	Size           string    `json:"size"`
	FileURI        string    `json:"file_uri"`
	AssetURI       string    `json:"asset_uri"`
	PreviewURI     string    `json:"preview_uri"`
	ProductID      int       `json:"product_id"`
	SizeMB         int       `json:"size_mb"`
	Order          int       `json:"order"`
	ProductUID     uuid.UUID `json:"product_uid"`
	OrganizationID int       `json:"organization_id"`
	OwnerID        int       `json:"owner_id"`
	OwnerType      string    `json:"owner_type"`
}

type CreateProductRequest struct {
	Name           string         `json:"name"`
	ThumbnailURI   string         `json:"thumbnail_uri"`
	CategoryID     int            `json:"category_id"`
	OrganizationID int            `json:"organization_id"`
	Documents      []DocumentData `json:"documents"` // تغییر به آرایه‌ای از داکیومنت‌ها
	Disabled       bool           `json:"disabled"`
}

func (d *DocumentData) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"title":           govalidity.New("title").Required().MinMaxLength(2, 200),
		"category_id":     govalidity.New("category_id").Required(),
		"phone_number":    govalidity.New("phone_number").Optional(),
		"cell_phone":      govalidity.New("cell_phone").Optional(),
		"website":         govalidity.New("website").Optional(),
		"telegram":        govalidity.New("telegram").Optional(),
		"instagram":       govalidity.New("instagram").Optional(),
		"linkedin":        govalidity.New("linkedin").Optional(),
		"location":        govalidity.New("location").Optional(),
		"size":            govalidity.New("size").Optional(),
		"file_uri":        govalidity.New("file_uri").Required(),
		"preview_uri":     govalidity.New("preview_uri").Required(),
		"product_id":      govalidity.New("product_id").Required(),
		"size_mb":         govalidity.New("size_mb").Required(),
		"order":           govalidity.New("order").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"title":           "عنوان",
			"category_id":     "شناسه دسته بندی",
			"phone_number":    "شماره تلفن",
			"cell_phone":      "شماره موبایل",
			"website":         "سایت",
			"telegram":        "تلگرام",
			"instagram":       "اینستاگرام",
			"linkedin":        "لینکدین",
			"location":        "موقعیت مکانی",
			"size":            "اندازه",
			"file_uri":        "آدرس فایل",
			"preview_uri":     "آدرس پیش نمایش",
			"product_id":      "شناسه محصول",
			"size_mb":         "اندازه فایل",
			"order":           "ترتیب",
			"organization_id": "شناسه سازمان",
		},
	)

	errr := govalidity.ValidateBody(r, schema, d)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (c *CreateProductRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":            govalidity.New("name").Required().MinMaxLength(2, 200),
		"thumbnail_uri":   govalidity.New("thumbnail_uri").Required(),
		"category_id":     govalidity.New("category_id").Required(),
		"organization_id": govalidity.New("organization_id").Required(),
		"disabled":        govalidity.New("disabled"),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":            "نام محصول",
			"thumbnail_uri":   "آیکون",
			"category_id":     "دسته محصول",
			"organization_id": "سازمان",
			"disabled":        "غیر فعال",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Create(ctx context.Context, input CreateProductRequest) (models.Product, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}
	if !policy.CanCreateOrganization(ctx, user) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد محصول را ندارید")
		return models.Product{}, response.ErrorForbidden("شما دسترسی ایجاد محصول را ندارید")
	}
	var products []models.Product

	var currentCategory models.Category
	err = s.db.WithContext(ctx).Where("id = ?", input.CategoryID).First(&currentCategory).Error
	if err != nil {
		s.logger.With(ctx).Error("خطایی در یافتن دسته محصول رخ داده است")
		return models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن دسته محصول رخ داده است")
	}

	parentCategoryID := exp.TerIf(currentCategory.ParentID.Valid, int(currentCategory.ParentID.Int64), int(currentCategory.ID))

	var currentOrganization models.Organization
	err = s.db.WithContext(ctx).Where("id = ?", input.OrganizationID).First(&currentOrganization).Error
	if err != nil {
		s.logger.With(ctx).Error("خطایی در یافتن سازمان رخ داده است")
		return models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	var pkg models.Package
	err = s.db.WithContext(ctx).Where("user_id =? AND organization_id = ?", id, currentOrganization.ID).
		Preload("Plan").
		Preload("Plan.Categories").
		Find(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن طرح رخ داده است")
	}

	found := false
	for _, category := range pkg.Plan.Categories {
		if category.ID == uint(parentCategoryID) {
			found = true
			break
		}
	}

	if !found {
		s.logger.With(ctx).Error("خطا در یافتن دسته بندی بسته")
		return models.Product{}, response.ErrorForbidden("خطا در یافتن دسته بندی بسته")
	}

	if pkg.ID == 0 {
		s.logger.With(ctx).Error("شما بسته برای افزودن این محصول خریداری نکردید")
		return models.Product{}, response.ErrorForbidden("شما بسته برای افزودن این محصول خریداری نکردید")
	}

	if pkg.ExpiredAt.Valid && pkg.ExpiredAt.Time.Before(time.Now()) {
		s.logger.With(ctx).Error("مدت اشتراک شما به اتمام رسیده است")
		return models.Product{}, response.ErrorForbidden("مدت اشتراک شما به اتمام رسیده است")
	}

	err = s.db.WithContext(ctx).Preload("Documents").Where("package_id = ?", pkg.ID).Find(&products).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطایی در یافتن محصولات بسته رخ داده است")
	}

	if len(products) >= pkg.Plan.ProductLimit {
		s.logger.With(ctx).Error("محدودیت ایجاد محصول برای این بسته به پایان رسیده است")
		return models.Product{}, response.ErrorForbidden("محدودیت ایجاد محصول برای این بسته به پایان رسیده است")
	}

	totalSizeMB := 0
	for _, doc := range input.Documents {
		totalSizeMB += doc.SizeMB
	}

	if totalSizeMB > pkg.StorageLimitMB {
		s.logger.With(ctx).Error("حجم فایل شما برای افزودن محصول بیش از حد مجاز است")
		return models.Product{}, response.ErrorBadRequest("حجم فایل شما برای افزودن محصول بیش از حد مجاز است")
	}

	productUID, err := uuid.NewUUID()
	if err != nil {
		s.logger.With(ctx).Error("خطا در تولید UUID برای محصول")
		return models.Product{}, response.ErrorResponse{Message: "خطا در تولید UUID برای محصول"}
	}

	var documents []models.Document
	for _, doc := range input.Documents {
		document := models.Document{
			Title:      doc.Title,
			CategoryID: input.CategoryID,
			PhoneNumber: dtp.NullString{
				String: doc.PhoneNumber,
				Valid:  doc.PhoneNumber != "",
			},
			CellPhone: dtp.NullString{
				String: doc.CellPhone,
				Valid:  doc.CellPhone != "",
			},
			Website: dtp.NullString{
				String: doc.Website,
				Valid:  doc.Website != "",
			},
			Telegram: dtp.NullString{
				String: doc.Telegram,
				Valid:  doc.Telegram != "",
			},
			Instagram: dtp.NullString{
				String: doc.Instagram,
				Valid:  doc.Instagram != "",
			},
			Linkedin: dtp.NullString{
				String: doc.Linkedin,
				Valid:  doc.Linkedin != "",
			},
			Location: dtp.NullString{
				String: doc.Location,
				Valid:  doc.Location != "",
			},
			Size: dtp.NullString{
				String: doc.Size,
				Valid:  doc.Size != "",
			},
			FileURI: doc.FileURI,
			AssetURI: dtp.NullString{
				String: doc.AssetURI,
				Valid:  doc.AssetURI != "",
			},
			PreviewURI:      doc.PreviewURI,
			SizeMB:          doc.SizeMB,
			Order:           doc.Order,
			UserID:          id,
			ProductUID:      productUID,
			Category:        &currentCategory,
			OrganizationID:  pkg.OrganizationID,
			OrganizationUID: currentOrganization.OrganizationUID, // اضافه کردن OrganizationUID
			OwnerID:         doc.OwnerID,
			OwnerType:       doc.OwnerType,
		}
		documents = append(documents, document)
	}

	product := models.Product{
		DisabledAt: dtp.NullTime{
			Time:  time.Now(),
			Valid: input.Disabled,
		},
		Name:            input.Name,
		ThumbnailURI:    input.ThumbnailURI,
		CategoryID:      input.CategoryID,
		PackageID:       int(pkg.ID),
		OrganizationID:  pkg.OrganizationID,
		OrganizationUID: currentOrganization.OrganizationUID, // اضافه کردن OrganizationUID
		Documents:       documents,
		ViewCount:       0,
		UserID:          id,
		ProductUID:      productUID,
		Category:        &currentCategory,
	}

	err = s.db.WithContext(ctx).Create(&product).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return product, response.GormErrorResponse(err, "خطایی در ایجاد محصول رخ داد")
	}

	pkg.StorageLimitMB -= totalSizeMB
	pkg.ProductLimit -= 1
	err = s.db.WithContext(ctx).Save(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Product{}, response.GormErrorResponse(err, "خطایی در بروزرسانی بسته رخ داده است")
	}

	return product, response.ErrorResponse{}
}
