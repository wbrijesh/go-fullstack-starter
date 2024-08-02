package schema

type ConfigType struct {
	JwtKey                     string `yaml:"jwt_secret"`
	Port                       int    `yaml:"port"`
	RateLimitRequestsPerSecond int    `yaml:"rate_limit_rps"`
}
