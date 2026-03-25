package app

type conf struct {
	Port          string `yaml:"port"`
	DSN           string `yaml:"dsn"`
	TimezoneLoc   string `yaml:"timezone_loc"`
	SessionSecret string `yaml:"session_secret"`
	Storage_path  string `yaml:"storage_path"`
	OCR           struct {
		Address string `yaml:"address"`
		Token   string `yaml:"token"`
	} `yaml:"ocr"`
	Todu struct {
		ClientID     string   `yaml:"client_id"`
		ClientSecret string   `yaml:"client_secret"`
		RedirectURL  string   `yaml:"redirect_url"`
		Scopes       []string `yaml:"scopes"`
		Endpoint     struct {
			AuthURL  string `yaml:"auth_url"`
			TokenURL string `yaml:"token_url"`
			UserInfo string `yaml:"user_info"`
		} `yaml:"endpoint"`
	} `yaml:"todu"`
}
