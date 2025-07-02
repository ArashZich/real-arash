package notification

type Template struct {
	Name    string
	Driver  string
	Message string
}
type Templates map[string]Template

var templates = Templates{
	"verify_phone": Template{
		Name:    "",
		Driver:  "",
		Message: "",
	},
}


func ValidateTemplate(t string) Template {
	return templates[t]
}
