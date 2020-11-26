/*
Copyright © 2020 Marco Moscatelli  <marco.moscatelli@itb.cnr.it>

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
  "path/filepath"
  "text/template"
	"github.com/spf13/cobra"
  "log"
  homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path"
	"strings"
  "os/exec"
	"io"
	"bytes"
)

var cfgRstudio = viper.New()
var cfgRstudioFile string

type Rstudio struct {
	Rversion string `mapstructure:"rversion"`
	Project string `mapstructure:"project"`
	Directory []string `mapstructure:"directory"`
	ExtraBind string `mapstructure:"extrabind"`
	HomeDirectory string `mapstructure:"homedirectory"`
}

var RstudioCfg Rstudio

// declaring struct 
type SlurmData struct { 
    Password string
    GeneralConf GeneralConfig
    RstudioConf Rstudio
} 

// rstudioCmd represents the rstudio command
var rstudioCmd = &cobra.Command{
	Use:   "rstudio",
	Short: "hpci rstudio spawn application",
	Long: `hpci rstudio spawn application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// mandatoryConfig()
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

	rstudioCmd.Flags().StringVarP(&cfgRstudioFile, "file", "f", "", "config file (default is $HOME/.hpci-rstudio.yaml)")
	rstudioCmd.Flags().StringVarP(&RstudioCfg.Rversion, "Rversion", "r", "", "path to R (absolute path)")
	cfgRstudio.BindPFlag("Rversion", rstudioCmd.Flags().Lookup("Rversion"))
	rstudioCmd.Flags().StringVarP(&RstudioCfg.Project,"project", "p", "", "project directory (absolute path)")
	cfgRstudio.BindPFlag("project", rstudioCmd.Flags().Lookup("project"))
	rstudioCmd.Flags().StringSliceVarP(&RstudioCfg.Directory,"directory", "d", []string{}, "add extra directory as  (absolute path separate with comma or repeat flag)")
	cfgRstudio.BindPFlag("directory", rstudioCmd.Flags().Lookup("directory"))
}

// initRstudio reads in config file 
func initRstudio() {
	cfgGeneral.SetDefault("rstudio.job.cpus", 10)
	if cfgRstudioFile != "" {
		// Use config file from the flag.
		cfgRstudio.SetConfigFile(cfgRstudioFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
		
		// Search config in home directory with name ".hpci" (without extension).
		cfgRstudio.AddConfigPath(home)
		cfgRstudio.SetConfigName(".hpci-rstudio")
	}

	// If a config file is found, read it in.
	if err := cfgRstudio.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cfgRstudio.ConfigFileUsed())
	} 
	if err := cfgRstudio.Unmarshal(&RstudioCfg); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
	RstudioCfg.HomeDirectory = home
}

func mandatoryConfig() {
	var mandatory []string

	// Administrator  Values
	if !cfgGeneral.IsSet("singularity.binary") {mandatory = append(mandatory, "Singularity Bin Path")}
	if !cfgGeneral.IsSet("rstudio.sif") {mandatory = append(mandatory, "Rstudio image")}
	if mandatory != nil {
		fmt.Println("ERROR: The following value are mandatory: ",mandatory)
		fmt.Println("       sudo permissions are needed to edit /etc/hpci/hpci.conf file")
		os.Exit(1)
	}
	
	// Program Values
	if !cfgRstudio.IsSet("r.version") {mandatory = append(mandatory, "Rversion")}
	if !cfgRstudio.IsSet("project") {mandatory = append(mandatory, "project")}
	if mandatory != nil {
		fmt.Println("ERROR: The following value are mandatory: ",mandatory)
		os.Exit(1)
	}
	_, err := os.Stat(cfgRstudio.GetString("project"))
		fmt.Println("ERROR: PATH not exist?",err)			
	if _, err := os.Stat(cfgRstudio.GetString("project")); !os.IsNotExist(err) {
		fmt.Println("ERROR: PATH not exist",cfgRstudio.GetString("project"))		
	}
}

func CheckError(err error) {
	if err != nil {
			log.Fatal(err)
	}
}

func RunRstudio() {
	password,_:=RandomHex(20)
	if len(RstudioCfg.Directory) > 0 { RstudioCfg.ExtraBind = "--bind "+strings.Join(RstudioCfg.Directory[:]," --bind ") }
	slurmdata := SlurmData{password,GeneralCfg,RstudioCfg}
	
	// Create .Rprofile file
	t, err := template.New("Rprofile").Parse(ProfileTemplate())
	CheckError(err)
	f, err := os.Create(cfgRstudio.GetString("project")+"/.Rprofile")
	CheckError(err)
	err = t.Execute(f, slurmdata)
	CheckError(err)
	f.Close()
		
	// Create Project.Rproj file
  f, err = os.Create(cfgRstudio.GetString("project")+"/"+path.Base(cfgRstudio.GetString("project"))+".Rproj")
	CheckError(err)
	_, err = f.WriteString(ProjectFile())
	CheckError(err)
	f.Close()

	// Create Slurm Output Directory
	outPath:= filepath.Join(RstudioCfg.HomeDirectory,GeneralCfg.Rstudio.Job.OutDir)
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
    os.Mkdir(outPath, 0755)
	}
	
	// Create .Rprofile file in homedirectory
	Rprofilepath:=filepath.Join(RstudioCfg.HomeDirectory,".Rprofile")
	if _, err := os.Stat(Rprofilepath); os.IsNotExist(err) {
		RprofileFile, err := os.Create(Rprofilepath)
		if err != nil {
			log.Fatal(err)
		}
		RprofileFile.Close()
	}
	
	// Create Sbatch file
	t, err = template.New("slurm-job").Parse(SlurmTemplate())
	CheckError(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, slurmdata)
	CheckError(err)
	f.Close()

	// f, err = os.Create("/tmp/slurm.test")
	// CheckError(err)
	// err = t.Execute(f, slurmdata)
	// CheckError(err)
	// f.Close()

	cmd := exec.Command("sbatch")
	slurmjob, err := cmd.StdinPipe()
	CheckError(err)

	go func() {
		defer slurmjob.Close()
		io.WriteString(slurmjob, tpl.String())
	}()

 output, err := cmd.CombinedOutput()
 CheckError(err)
 fmt.Printf( "%s", output )
 fmt.Printf( "Output Job Information: %s\n", outPath )

}