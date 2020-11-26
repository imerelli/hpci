/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
  "log"
	"github.com/spf13/viper"
)

var cfgGeneral = viper.New()
//var cfgFile string

type GeneralConfig struct {
	Singularity struct {
		Binary string `mapstructure:"binary"`
	} `mapstructure:"singularity"`
	Cluster struct { 
		Url string `mapstructure:"url"`
	} `mapstructure:"cluster"`
	Ldaps struct {
		Host string `mapstructure:"host"`
		UserDN string `mapstructure:"user-dn"`
		CAfile string `mapstructure:"ca-file"`
	} `mapstructure:"ldaps"`
	Rstudio struct {
		Sif string `mapstructure:"sif"`
		Ports struct {
			From string `mapstructure:"from"`
			To string `mapstructure:"to"`
		} `mapstructure:"ports"`
		Auth string `mapstructure:"auth"`
		Job struct {
			Name string `mapstructure:"name"`
			Time string `mapstructure:"time"`
			Partition string `mapstructure:"partition"`
			Cpus string `mapstructure:"cpus"`
			Memory string `mapstructure:"memory"`
			OutDir string `mapstructure:"outdir"`
			Output string `mapstructure:"output"`
		} `mapstructure:"job"`
	} `mapstructure:"rstudio"`
}

var GeneralCfg GeneralConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hpci",
	Short: "High Performance Computing Interface",
	Long: `High Performance Computing Interface to spawn Singularity Application on HPC`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
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

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hpci.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	
}

// initConfig reads in config file /etc/hpci/hpci.conf.
func initConfig() {
  cfgGeneral.SetConfigType("yaml")
	cfgGeneral.SetConfigFile("/etc/hpci/hpci.conf")
	cfgGeneral.ReadInConfig()
	if err := cfgGeneral.Unmarshal(&GeneralCfg); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
