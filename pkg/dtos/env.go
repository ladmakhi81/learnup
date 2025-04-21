package dtos

type MinioEnvConfig struct {
	URL       string `koanf:"url"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
	Region    string `koanf:"region"`
}

type RedisEnvConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type MainDBEnvConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
}

type TemporalEnvConfig struct {
	Host     string `koanf:"host"`
	Port     uint   `koanf:"port"`
	Endpoint string `koanf:"endpoint"`
}

type SmtpEnvConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type AppEnvConfig struct {
	Port           int    `koanf:"port"`
	TokenSecretKey string `koanf:"token_secret_key"`
	OPENAI_KEY     string `koanf:"openai_key"`
}

type ZarinpalEnvConfig struct {
	Merchant    string `koanf:"merchant"`
	CallbackURL string `koanf:"callback_url"`
}

type ZibalEnvConfig struct {
	Merchant    string `koanf:"merchant"`
	CallbackURL string `koanf:"callback_url"`
}

type StripeEnvConfig struct {
	Key         string `koanf:"key"`
	CallbackURL string `koanf:"callback_url"`
}

type EnvConfig struct {
	Minio    MinioEnvConfig    `koanf:"minio"`
	Redis    RedisEnvConfig    `koanf:"redis"`
	MainDB   MainDBEnvConfig   `koanf:"main_db"`
	Smtp     SmtpEnvConfig     `koanf:"smtp"`
	App      AppEnvConfig      `koanf:"app"`
	Temporal TemporalEnvConfig `koanf:"temporal"`
	Zarinpal ZarinpalEnvConfig `koanf:"zarinpal"`
	Zibal    ZibalEnvConfig    `koanf:"zibal"`
	Stripe   StripeEnvConfig   `koanf:"stripe"`
}
