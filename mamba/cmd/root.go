package cmd

import (
	"fmt"
	"os"

	"github.com/gofunct/common/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	// Used for flags.
	cfgFile, userLicense string

	rootCmd = &cobra.Command{
		Use:   "mamba",
		Short: "A generator for Mamba based Applications",
		Long: utils.Red(`oooo     oooo                          oooo                   
 8888o   888   ooooooo   oo ooo oooo    888ooooo     ooooooo  
 88 888o8 88   ooooo888   888 888 888   888    888   ooooo888 
 88  888  88 888    888   888 888 888   888    888 888    888 
o88o  8  o88o 88ooo88 8o o888o888o888o o888ooo88    88ooo88 8o`),
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Print("failed to execute command", err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mamba.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	if err := viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author")); err != nil {
		log.Print("failed to bind to flag", err)
	}
	if err := viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper")); err != nil {
		log.Print("failed to bind to flag", err)
	}

	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mamba")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
