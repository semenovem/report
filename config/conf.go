package config

const (
	IndexHTMLFileName          = "/index.html"
	ForbiddenHTMLFileName      = "/forbidden.html"
	TooManyRequestHTMLFileName = "/too-many-request.html"
)

type Main struct {
	Base Base
	Rest Rest

	HTMLDir    string `env:"HTML_DIR,required"`
	AccessCode string `env:"ACCESS_CODE,required"`

	Ozon struct {
		Path string `env:"OZON_PATH,required"`

		ClientID1 string `env:"OZON_CLIENT_ID_1,required"`
		ClientID2 string `env:"OZON_CLIENT_ID_2,required"`

		APIKey1 string `env:"OZON_API_KEY_1,required"`
		APIKey2 string `env:"OZON_API_KEY_2,required"`
	}
}
