package global

import (
	"github.com/spf13/viper"
)

type Config struct {
	Target string `mapstructure:"target"` // 目标类型，支持redis、mongodb

	Addr     string `mapstructure:"addr"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
	Charset  string `mapstructure:"charset"`
	SlaveID  uint32 `mapstructure:"slave_id"`

	RuleConfigs []*Rule `mapstructure:"rule"`

	// ------------------- MONGODB -----------------
	MongodbAddr     string `mapstructure:"mongodb_addrs"`    //mongodb地址，多个用逗号分隔
	MongodbUsername string `mapstructure:"mongodb_username"` //mongodb用户名，默认为空
	MongodbPassword string `mapstructure:"mongodb_password"` //mongodb密码，默认为空
}

var (
	GbConfig *Config
	RulesMap map[string]bool
)

func Initialize() {
	GbConfig = &Config{}
	if err := viper.Unmarshal(GbConfig); err != nil {
		panic(err)
	}
	for _, rule := range GbConfig.RuleConfigs {
		RulesMap[rule.Schema+"."+rule.Table] = true
	}
}
