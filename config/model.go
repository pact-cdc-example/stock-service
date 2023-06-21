package config

type config struct {
	Postgres    Postgres    `mapstructure:"postgres"`
	Server      Server      `mapstructure:"server"`
	ExternalURL ExternalURL `mapstructure:"externalURL"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type ExternalURL struct {
	ProductAPI string `mapstructure:"productApi"`
}
