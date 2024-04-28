package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mysql-mongodb-syncer/global"

	"mysql-mongodb-syncer/services"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "syncer-cli",
		Short: "Sync mysql to mongodb",
		Run:   rootCmdRun,
	}
)

func rootCmdRun(cmd *cobra.Command, args []string) {
	global.Initialize()
	services.NewListener()
	services.ListenerService.Start()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("aaaaaaaaaaaa")
		global.Initialize()
		services.ListenerService.Reload()
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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./app.yml)")
}

func initConfig() {
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default config file.
		viper.SetConfigName("app")
		viper.AddConfigPath("./")
	}

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}
