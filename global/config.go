package global

import (
	"mysql-anywhere-syncer/utils/logger"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"pass" validate:"required"`
	Charset  string `mapstructure:"charset" validate:"required"`
	SlaveID  uint32 `mapstructure:"slave_id" validate:"required"`

	RuleConfigs []*Rule `mapstructure:"rule" validate:"required,dive"`

	// ------------------- MongoDB -----------------
	MongodbHost     string `mapstructure:"mongodb_host"`
	MongodbUsername string `mapstructure:"mongodb_username"`
	MongodbPassword string `mapstructure:"mongodb_password"`
	MongodbPort     int    `mapstructure:"mongodb_port"`

	// ------------------- Redis -----------------
	RedisHost     string `mapstructure:"redis_host"`
	RedisPort     int    `mapstructure:"redis_port"`
	RedisPass     string `mapstructure:"redis_pass"`
	RedisDatabase int    `mapstructure:"redis_db"`
}

var (
	GbConfig *Config
	RulesMap map[string][]*Rule
)

func Initialize() {
	logger.NewLogger()
	GbConfig = &Config{}
	if err := viper.Unmarshal(GbConfig); err != nil {
		logger.Logger.Fatal("Failed to unmarshal config")
	}

	validate := validator.New()
	err := validate.Struct(GbConfig)
	if err != nil {
		logger.Logger.Fatal("Failed to validate config: ", err)
	}

}
