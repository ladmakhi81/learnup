package env

type MinioEnvConfig struct {
	URL      string `koanf:"url"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Region   string `koanf:"region"`
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

type SmtpEnvConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type EnvConfig struct {
	Minio  MinioEnvConfig  `koanf:"minio"`
	Redis  RedisEnvConfig  `koanf:"redis"`
	MainDB MainDBEnvConfig `koanf:"main_db"`
	Smtp   SmtpEnvConfig   `koanf:"smtp"`
}
