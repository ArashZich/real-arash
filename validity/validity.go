package validity

import (
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

func ApplyTranslations() {
	govalidity.SetDefaultErrorMessages(
		&govaliditym.Validations{
			IsEmail:           "{field} باید ایمیل معتبر باشد",
			IsUrl:             "{field} باید آدرس معتبر باشد",
			IsIp:              "{field} باید آی پی معتبر باشد",
			IsRequired:        "{field} باید مقدار داشته باشد",
			IsMinLength:       "{field} باید حداقل {min} کاراکتر باشد",
			IsMin:             "{field} باید حداقل {min} باشد",
			IsMaxLength:       "{field} باید حداکثر {max} کاراکتر باشد",
			IsMax:             "{field} باید حداکثر {max} باشد",
			IsMinMaxLength:    "{field} باید حداقل {min} و حداکثر {max} کاراکتر باشد",
			IsAlpha:           "{field} باید حروف الفبا باشد",
			IsAlphaNum:        "{field} باید حروف الفبا و اعداد باشد",
			IsNumber:          "{field} باید عدد باشد",
			IsMaxDate:         "{field} باید حداکثر {max} باشد",
			IsDNSName:         "{field} باید نام دامنه معتبر باشد",
			IsFloat:           "{field} باید عدد اعشاری باشد",
			IsFilepath:        "{field} باید مسیر فایل معتبر باشد",
			IsHost:            "{field} باید میزبان معتبر باشد",
			IsIn:              "{field} باید یکی از مقادیر {in} باشد",
			IsInt:             "{field} باید عدد صحیح باشد",
			IsInRange:         "{field} باید بین {from} و {to} باشد",
			IsIpV4:            "{field} باید آی پی ورژن 4 باشد",
			IsIpV6:            "{field} باید آی پی ورژن 6 باشد",
			IsJson:            "{field} باید یک رشته جیسون معتبر باشد",
			IsLatitude:        "{field} باید عرض جغرافیایی معتبر باشد",
			IsLogitude:        "{field} باید طول جغرافیایی معتبر باشد",
			IsLowerCase:       "{field} باید حروف کوچک باشد",
			IsMaxDateTime:     "{field} باید حداکثر {max} باشد",
			IsFilterOperators: "{field} باید یکی از عملگرهای فیلتر {in} باشد",
			IsHexColor:        "{field} باید کد رنگ معتبر باشد",
			IsMaxTime:         "{field} باید حداکثر {max} باشد",
			IsMinDate:         "{field} باید حداقل {min} باشد",
			IsPort:            "{field} باید پورت معتبر باشد",
			IsMinDateTime:     "{field} باید حداقل {min} باشد",
			IsMinTime:         "{field} باید حداقل {min} باشد",
			IsSlice:           "{field} باید یک آرایه باشد",
			IsUpperCase:       "{field} باید حروف بزرگ باشد",
		},
	)
}
