/*
Copyright © 2026 Zoom theoldzoom@proton.me

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"os"

	"github.com/theOldZoom/gofm/internal/config"
	"github.com/theOldZoom/gofm/internal/tui/setup"
	"github.com/theOldZoom/gofm/internal/verbose"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gofm",
	Short: "A CLI for Last.fm",

	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/gofm/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "verbose output")
	cobra.CheckErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))
}

func initConfig() {
	if cfgFile != "" {
		verbose.Printf("using explicit config file: %s", cfgFile)
		config.SetPath(cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		path, err := config.Path()
		cobra.CheckErr(err)
		verbose.Printf("resolved config path: %s", path)
		config.SetPath(path)
		viper.SetConfigFile(path)
	}

	viper.AutomaticEnv()
	verbose.Printf("automatic environment loading enabled")

	err := viper.ReadInConfig()
	if err == nil {
		verbose.Printf("config loaded from: %s", viper.ConfigFileUsed())
		return
	}

	var notFound viper.ConfigFileNotFoundError
	if !os.IsNotExist(err) && !errors.As(err, &notFound) {
		cobra.CheckErr(err)
	}

	verbose.Printf("config not found, starting interactive setup")
	cfg, err := setup.Run()
	cobra.CheckErr(err)

	viper.Set("username", cfg.Username)
	viper.Set("api_key", cfg.ApiKey)
	verbose.Printf("setup completed, configuration stored in memory")
}
