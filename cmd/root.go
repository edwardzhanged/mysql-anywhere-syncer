package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mysql-mongodb-syncer/global"
	"mysql-mongodb-syncer/services"
	"mysql-mongodb-syncer/syncer"

	"mysql-mongodb-syncer/utils/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	dumpFlag  bool
	firstTime = true

	rootCmd = &cobra.Command{
		Use:   "syncer-cli",
		Short: "Sync mysql to mongodb",
		Run:   rootCmdRun,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./app.yml)")
	rootCmd.PersistentFlags().BoolVar(&dumpFlag, "dump", false, "dump data from mysql to file(default is false)")

}

func initConfig() {
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath("./")
	}

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.WithError(err).Fatal("Failed to read config file")
	}
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	global.Initialize()

	global.RulesMap = make(map[string][]*global.Rule)
	targets := initializeRulesMap()
	for _, rule := range targets {
		switch rule.Target {
		case global.TargetMongoDB:
			syncer.NewMongo(&syncer.ConnectOptions{
				Host:     global.GbConfig.MongodbHost,
				Port:     global.GbConfig.MongodbPort,
				Username: global.GbConfig.MongodbUsername,
				Password: global.GbConfig.MongodbPassword,
			})
			if err := syncer.MongoInstance.Connect(); err != nil {
				logger.Logger.WithError(err).Fatal("Failed to connect to mongodb")
			}
		default:
		}
	}

	services.NewListener()
	services.ListenerService.Start(dumpFlag, firstTime)
	firstTime = false
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		global.Initialize()
		initializeRulesMap()
		services.ListenerService.Reload(dumpFlag, firstTime)
		color.Yellowln("Config file changed:", e.Name)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			color.Redln("Received an interrupt, stopping services...")
			return
		default:
		}
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func initializeRulesMap() []*global.Rule {
	global.RulesMap = make(map[string][]*global.Rule)
	targets := make([]*global.Rule, 0)
	for _, rule := range global.GbConfig.RuleConfigs {
		schemaTable := fmt.Sprintf("%s.%s", rule.Schema, rule.Table)
		global.RulesMap[schemaTable] = append(global.RulesMap[schemaTable], rule)
		targets = append(targets, rule)
	}
	return targets
}
