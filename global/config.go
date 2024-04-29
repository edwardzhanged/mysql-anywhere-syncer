package global

import (
	"mysql-mongodb-syncer/utils/logger"

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

	// ------------------- MONGODB -----------------
	MongodbHost     string `mapstructure:"mongodb_host"`     //mongodb地址，多个用逗号分隔
	MongodbUsername string `mapstructure:"mongodb_username"` //mongodb用户名，默认为空
	MongodbPassword string `mapstructure:"mongodb_password"` //mongodb密码，默认为空
	MongodbPort     int    `mapstructure:"mongodb_port"`     //mongodb端口，默认27017
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
		logger.Logger.Fatal("Config file validation failed")
	}

}
