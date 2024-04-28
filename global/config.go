package global

import (
	"fmt"
	"log"
	"mysql-mongodb-syncer/syncer"

	"github.com/spf13/viper"
)

type Config struct {
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
	MongodbPort     int    `mapstructure:"mongodb_port"`     //mongodb端口，默认27017
}

var (
	GbConfig *Config
	RulesMap map[string][]*Rule
)

func Initialize() {
	GbConfig = &Config{}
	if err := viper.Unmarshal(GbConfig); err != nil {
		panic(err)
	}

	RulesMap = make(map[string][]*Rule)
	targets := make([]*Rule, 0)
	for _, rule := range GbConfig.RuleConfigs {
		schemaTable := fmt.Sprintf("%s.%s", rule.Schema, rule.Table)
		RulesMap[schemaTable] = append(RulesMap[schemaTable], rule)
		targets = append(targets, rule)
	}
	for _, rule := range targets {
		switch rule.Target {
		case TargetMongoDB:
			syncer.NewMongo(&syncer.ConnectOptions{
				Host:     GbConfig.MongodbAddr,
				Port:     GbConfig.MongodbPort,
				Username: GbConfig.MongodbUsername,
				Password: GbConfig.MongodbPassword,
			})
			if err := syncer.MongoInstance.Connect(); err != nil {
				log.Fatal(err)
			}
		default:
		}
	}

}

// TODO: 校验配置逻辑
