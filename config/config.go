package config

type (
	Config struct {
		DB      DB     `yaml:"db"`
		Auth    Auth   `yaml:"auth"`
		Aws     Aws    `yaml:"aws"`
		Sys     Sys    `yaml:"sys"`
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

	Aws struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
		Region string `yaml:"region"`
		Bucket string `yaml:"bucket"`
	}

	Sys struct {
		AttachmentBase string `yaml:"attachment_base"`
	}
)
