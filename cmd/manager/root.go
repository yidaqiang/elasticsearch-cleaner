package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newRootCmd(args []string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "es-manager [COMMAND] [ARGS]",
		Short: "ElasticSearch clean Tool",
	}

	cmd.AddCommand(
		newDeleteIndicesCmd(&manager, args),
	)
	flags := cmd.PersistentFlags()
	addRootFlag(flags)

	return cmd
}

func addRootFlag(fs *pflag.FlagSet) {
	fs.StringVarP(&cfgFile, "config", "c", "",
		"配置文件路径，默认当前目录下 elasticsearch-manager.yml 文件")
}
