package validate

var fields = map[string]string{
	"Email":    "E-posta",
	"Password": "Şifre",
	"Name":     "isim",
	"Surname":  "soyad",
}

func GetFields() *map[string]string {
	return &fields
}
