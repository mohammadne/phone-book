package repository

type Config struct {
	CursorSecret string `koanf:"cursor_secret"`
	Limit        struct {
		Min int `koanf:"min"`
		Max int `koanf:"max"`
	} `koanf:"limit"`
}
