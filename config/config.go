package config

type (
	Config struct {
		DB      DB     `yaml:"db"`
		Auth    Auth   `yaml:"auth"`
		Mixpay  Mixpay `yaml:"mixpay"`
		AppName string `yaml:"appname"`
	}

	Auth struct {
		JwtSecret         string `json:"jwt_secret"`
		MixinClientSecret string `json:"mixin_client_secret"`
	}

	DB struct {
		Driver     string `json:"driver"`
		Datasource string `json:"datasource"`
	}

	Mixpay struct {
		PayeeID     string `json:"payee_id"`
		CallbackURL string `json:"callback_url"`
	}
)
