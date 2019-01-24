// Copyright © 2019 Coleman Word <coleman.word@gofunct.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/gofunct/mamba/pkg/function"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile, userLicense, in, out, pkg string
	logger                             *logrus.Logger
	rootCmd                            = &cobra.Command{
		Use:   "mamba",
		Short: "A generator for Mamba based Applications 🐍",
		Long: `Mamba is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Mamba application.`,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("failed to execute:%s/n", err)
	}
}

func init() {
	cobra.OnInitialize(
		func() { logger = logrus.New() },
		function.InitConfig(cfgFile),
	)
	{
		rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
		rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
		rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
		rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
		rootCmd.PersistentFlags().StringVarP(&in, "input", "i", ".", "use Viper for configuration")
		rootCmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "use Viper for configuration")
		rootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "use Viper for configuration")
	}

	{
		viper.BindPFlags(rootCmd.Flags())
		viper.BindPFlags(rootCmd.PersistentFlags())
		viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
		viper.SetDefault("license", "apache")
		viper.SetDefault("input", ".")
		viper.SetDefault("output", ".")
	}

	{
		rootCmd.AddCommand(addCmd)
		rootCmd.AddCommand(initCmd)
		rootCmd.AddCommand(testCmd)
		rootCmd.AddCommand(walkCmd)
		rootCmd.AddCommand(htmlCmd)
	}
}
