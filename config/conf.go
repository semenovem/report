package config

type Main struct {
	Base Base
	Rest Rest

	IndexPath string `env:"INDEX_HTML_PATH,required"`

	Ozon struct {
		Path string `env:"OZON_PATH,required"`

		ClientID1 string `env:"OZON_CLIENT_ID_1,required"`
		ClientID2 string `env:"OZON_CLIENT_ID_2,required"`

		APIKey1 string `env:"OZON_API_KEY_1,required"`
		APIKey2 string `env:"OZON_API_KEY_2,required"`
	}
}
