package env

type Environment struct {
	AppUrl               string `env:"APP_URL" envDefault:"http://localhost:3000"`
	PostgresUser         string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword     string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresDb           string `env:"POSTGRES_DB" envDefault:"postgres"`
	PostgresConn         string `env:"POSTGRES_CONNECTION_STRING,expand" envDefault:"postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_NAME}?sslmode=disable"`
	HttpHost             string `env:"HTTP_HOST" envDefault:"0.0.0.0:3001"`
	HttpCors             string `env:"HTTP_CORS" envDefault:""`
	AuthSecret           string `env:"AUTH_SECRET" envDefault:""`
	AuthAdminEmail       string `env:"AUTH_ADMIN_EMAIL" envDefault:""`
	EuropeBaseUrl        string `env:"EUROPE_BASE_URL" envDefault:""`
	CurrencyRateEndpoint string `env:"CURRENCY_RATE_ENDPOINT" envDefault:""`
}
