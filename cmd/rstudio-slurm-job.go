package cmd

import (
  "crypto/rand"
  "encoding/hex"
)

func RandomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }
  return hex.EncodeToString(bytes), nil
}

func SlurmTemplate() string {

	const tpl = `#!/bin/sh
#SBATCH --time=08:00:00
#SBATCH --signal=USR2
#SBATCH --ntasks=2
#SBATCH --job-name=Rstudio
#SBATCH --cpus-per-task=2
#SBATCH --mem=8192
#SBATCH --output={{.PrjDir}}/rstudio-server.%j.out
#SBATCH --partition=master

export PASSWORD={{.Password}}
export HOMEDIR=$( getent passwd "$USER" | cut -d: -f6 )

PORT=$(shuf -i {{.PortRange}} -n 1)
cat <<EOF
Log in to RStudio Server:

   url: $SLURM_NODELIST.{{.UrlDomain}}:${PORT}
   user: ${USER}
   password: ${PASSWORD}

When done using RStudio Server, terminate the job by:

1. Exit the RStudio Session ("power" button in the top right corner of the RStudio window)
2. Issue the following command on the login node:

      scancel ${SLURM_JOB_ID}

EOF

{{.SingularityBin}} exec -c --bind {{.PrjDir}} --bind {{.PrjDir}}/.Rprofile:${HOMEDIR}/.Rprofile --bind {{.Rversion}} {{.SingularityImage}} rserver --www-port=${PORT} --auth-none=0  --auth-pam-helper-path=pam-helper --rsession-ld-library-path={{.Rversion}}/lib --rsession-which-r={{.Rversion}}/bin/R 

printf 'rserver exited' 1>&2
`
	return tpl
}

func ProfileTemplate() string {
  const tpl =`setHook("rstudio.sessionInit", function(newSession) {
  if (newSession && is.null(rstudioapi::getActiveProject()))
    rstudioapi::openProject("{{.PrjDir}}/")
}, action = "append")
`
  
  return tpl
}


func ProjectFile() string {
  
  const tpl =`Version: 1.0

RestoreWorkspace: Default
SaveWorkspace: Default
AlwaysSaveHistory: Default

EnableCodeIndexing: Yes
UseSpacesForTab: Yes
NumSpacesForTab: 2
Encoding: UTF-8

RnwWeave: Sweave
LaTeX: pdfLaTeX`
  
  return tpl
}

