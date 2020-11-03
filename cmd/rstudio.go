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
)

// rstudioCmd represents the rstudio command
var rstudioCmd = &cobra.Command{
	Use:   "rstudio",
	Short: "hpci rstudio spawn application",
	Long: `hpci rstudio spawn application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rstudio called")
		RunRstudio()
	},
}

func init() {
	rootCmd.AddCommand(rstudioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rstudioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rstudioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// declaring struct 
type SlurmData struct { 
    SingularityBin string 
    SingularityImage string
    PortRange string
    UrlDomain string 
    WorkDir string
} 


// /opt/singularity/bin/singularity
func RunRstudio() {
	slurmdata := SlurmData{"/opt/singularity/bin/singularity", "/opt/Rstudio/R.3.6.3.simg", "9000-9099","hpc.bioinformatics.itb.cnr.it","R2"}
	t, err := template.ParseFiles("../templates/slurm-job.sh") 

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	tplstruct := SlurmTemplate()

	t, err = template.New("slurm-job").Parse(tplstruct)
	check(err)
	// Create the file
	filename:=RandomString()
	f, err := os.Create("/tmp/"+filename)
	check(err)

	err = t.Execute(f, slurmdata)
	check(err)

	// Close the file when done.
	f.Close()

    cmd := exec.Command("sbatch","/tmp/"+filename)

    // run command
    if otuput, err := cmd.Output(); err != nil {
        fmt.Println( "Error:", err )
    } else {
        fmt.Printf( "Output: %s\n", otuput )
    }

}