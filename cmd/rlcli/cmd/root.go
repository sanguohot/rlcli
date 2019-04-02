package cmd

import (
	"github.com/sanguohot/rlcli/pkg/common/log"
	"github.com/sanguohot/rlcli/pkg/rlcli"
	"github.com/spf13/cobra"
)

var (
	addr  string
	limit string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rlcli",
	Short: "use to serve rate limit case.",
	Long:  `a command tool to serve simple rate limit case with gin.
use case: rlcli -l 10-M, means 10 reqs/minute`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ss := rlcli.New(limit, addr)
		ss.Serve()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&limit, "limit", "l", "2-M", "the rate limit config, default '2-M' which means 2 reqs/minute")
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "H", "localhost:8080", "the local address to serve, default 'localhost:8080'")
}
