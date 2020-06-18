package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yitypes "github.com/yidaqiang/elasticsearch-manager/pkg/types"
	"os"
)

var (
	cfgFile = ""
	manager = yitypes.ElasticSearchManager{}
)

func main() {

	cmd := newRootCmd(os.Args[1:])
	cobra.OnInitialize(initConfig)

	if err := cmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("elasticsearch-manager.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&manager)
	if err != nil {
		panic(err)
	}
	manager.Init()
}
