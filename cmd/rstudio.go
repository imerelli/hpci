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
	"os"
  "os/exec"
  "text/template"
	"github.com/spf13/cobra"
  "log"
  homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path"
)

var rstudioFile string
var Rversion string
var project string

// declaring struct 
type SlurmData struct { 
    SingularityBin string 
    SingularityImage string
    PortRange string
    UrlDomain string 
    PrjDir string
    Rversion string
    Password string
} 

// rstudioCmd represents the rstudio command
var rstudioCmd = &cobra.Command{
	Use:   "rstudio",
	Short: "hpci rstudio spawn application",
	Long: `hpci rstudio spawn application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mandatoryConfig()
		RunRstudio()
	},
}

func init() {
	cobra.OnInitialize(initRstudio)
	rootCmd.AddCommand(rstudioCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rstudioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rstudioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rstudioCmd.Flags().StringVarP(&rstudioFile, "file", "f", "", "config file (default is $HOME/.hpci-rstudio.yaml)")
	rstudioCmd.Flags().StringVarP(&Rversion, "Rversion", "r", "", "path to R (absolute path)")
	viper.BindPFlag("Rversion", rstudioCmd.Flags().Lookup("Rversion"))
	rstudioCmd.Flags().StringVarP(&project,"project", "p", "", "project directory (absolute path)")
	viper.BindPFlag("project", rstudioCmd.Flags().Lookup("project"))
}

// initRstudio reads in config file 
func initRstudio() {
	if rstudioFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rstudioFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hpci" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hpci-rstudio")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.MergeInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func mandatoryConfig() {
	var mandatory []string

	// Administrator  Values
	if !viper.IsSet("SingularityBin") {mandatory = append(mandatory, "Singularity Bin Path")}
	if !viper.IsSet("RstudioImage") {mandatory = append(mandatory, "Rstudio image")}
	if mandatory != nil {
		fmt.Println("ERROR: The following value are mandatory: ",mandatory)
		fmt.Println("       sudo permissions are needed to edit /etc/hpci/hpci.conf file")
		os.Exit(1)
	}
	
	// Program Values
	if !viper.IsSet("Rversion") {mandatory = append(mandatory, "Rversion")}
	if !viper.IsSet("project") {mandatory = append(mandatory, "project")}
	if mandatory != nil {
		fmt.Println("ERROR: The following value are mandatory: ",mandatory)
		os.Exit(1)
	}
}

func RunRstudio() {
	password,_:=RandomHex(20)
	slurmdata := SlurmData{
		viper.GetString("SingularityBin"),
		viper.GetString("RstudioImage"),
		viper.GetString("RstudioPorts"),
		viper.GetString("ClusterUrl"),
		viper.GetString("project"),
		viper.GetString("Rversion"),
		password}

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create .Rprofile file
	t, err := template.New("Rprofile").Parse(ProfileTemplate())
	check(err)
	f, err := os.Create(viper.GetString("project")+"/.Rprofile")
	check(err)
	err = t.Execute(f, slurmdata)
	check(err)
	f.Close()
		
	// Create Project.Rproj file
  f, err = os.Create(viper.GetString("project")+"/"+path.Base(viper.GetString("project"))+".Rproj")
	check(err)
	_, err = f.WriteString(ProjectFile())
	check(err)
	f.Close()

	// Create Sbatch file
	t, err = template.New("slurm-job").Parse(SlurmTemplate())
	check(err)
	filename,_:=RandomHex(20)
	f, err = os.Create("/tmp/"+filename)
	check(err)
	err = t.Execute(f, slurmdata)
	check(err)
	f.Close()
  cmd := exec.Command("sbatch","/tmp/"+filename)
	
  // run command
  if otuput, err := cmd.Output(); err != nil {
    fmt.Println( "Error:", err )
  } else {
  	fmt.Printf( "Output: %s\n", otuput )
  }

	err = os.Remove("/tmp/"+filename)
	check(err)

}