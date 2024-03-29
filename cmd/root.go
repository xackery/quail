package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quail",
	Short: "Ever[Q]uest [U]niversal [A]rchive, [I]mport, and [L]oader system",
	Long: `An Ever[Q]uest [U]niversal [A]rchive, [I]mport, and [L]oader system.
  - .ani animation files (inspect)
  - .eqg pfs archives (compress, extract, inspect)
  - .mod model files (inspect)
  - .ter terrain files (inspect)
  - .zon zone files (inspect)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.quail.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose debugging")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	isVerbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if isVerbose {
		log.SetLogLevel(0)
		log.Debugf("Verbose logging enabled")
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".quail" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".quail")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
