package configs

import "github.com/spf13/viper"

type Conf struct {
	Port                 string `mapstructure:"PORT"`
	RedisAddress         string `mapstructure:"REDIS_ADDRESS"`
	RedisPort            string `mapstructure:"REDIS_PORT"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	RateLimiter          string `mapstructure:"RATE_LIMITER"`
	RequestsPerSecond    int    `mapstructure:"REQUESTS_PER_SECOND"`
	BlockedTimeInSeconds int    `mapstructure:"BLOCKED_TIME_IN_SECONDS"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
