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
	const tpl = `#!/bin/bash
#SBATCH --time={{.GeneralConf.Rstudio.Job.Time}}
#SBATCH --signal=USR2
#SBATCH --job-name={{.GeneralConf.Rstudio.Job.Name}}
#SBATCH --ntasks=1
#SBATCH --cpus-per-task={{.GeneralConf.Rstudio.Job.Cpus}}
#SBATCH --mem={{.GeneralConf.Rstudio.Job.Memory}}
#SBATCH --output={{.RstudioConf.HomeDirectory}}/{{.GeneralConf.Rstudio.Job.OutDir}}/{{.GeneralConf.Rstudio.Job.Output}}
#SBATCH --partition={{.GeneralConf.Rstudio.Job.Partition}}

if [[ "{{.GeneralConf.Rstudio.Auth}}" == "password" ]]; then
  export RSTUDIO_PASSWORD={{.Password}}
  RSTUDIO_HELPER=rstudio_auth
elif [[ "{{.GeneralConf.Rstudio.Auth}}" == "ldaps" ]]; then 
  export LDAP_HOST={{.GeneralConf.Ldaps.Host}}
  export LDAP_USER_DN={{.GeneralConf.Ldaps.UserDN}}
  if [[ ! -z "{{.GeneralConf.Ldaps.CAfile}}" ]]; then 
    export LDAP_CERT_FILE=/ca-certs.crt
    BIND_CERT="--bind {{.GeneralConf.Ldaps.CAfile}}:/ca-certs.crt"
  fi
  RSTUDIO_PASSWORD="ldap credential"
  RSTUDIO_HELPER=ldap_auth

fi

PORT=$(comm -23 <(seq {{.GeneralConf.Rstudio.Ports.From}} {{.GeneralConf.Rstudio.Ports.To}} | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1)

cat <<EOF
Log in to RStudio Server:

   url: $SLURM_NODELIST.{{.GeneralConf.Cluster.Url}}:${PORT}
   Rversion: {{.RstudioConf.Rversion}}
   Project Directory: {{.RstudioConf.Project}} 
   Extra Directory: {{.RstudioConf.Directory}} 
   user: ${USER}
   password: ${RSTUDIO_PASSWORD}

When done using RStudio Server, terminate the job by:

1. Exit the RStudio Session ("power" button in the top right corner of the RStudio window)
2. Issue the following command on the login node:

      scancel ${SLURM_JOB_ID}

EOF

{{.GeneralConf.Singularity.Binary}} exec -c {{.RstudioConf.ExtraBind}} ${BIND_CERT} --bind {{.RstudioConf.Project}} --bind {{.RstudioConf.Project}}/.Rprofile:{{.RstudioConf.HomeDirectory}}/.Rprofile --bind {{.RstudioConf.HomeDirectory}}/{{.GeneralConf.Rstudio.Job.OutDir}} --bind {{.RstudioConf.Rversion}} {{.GeneralConf.Rstudio.Sif}} rserver --www-port ${PORT} --auth-none 0  --auth-pam-helper ${RSTUDIO_HELPER} --rsession-ld-library-path {{.RstudioConf.Rversion}}/lib --rsession-which-r {{.RstudioConf.Rversion}}/bin/R
printf 'rserver exited' 1>&2
`
	return tpl
}

func ProfileTemplate() string {
  const tpl =`setHook("rstudio.sessionInit", function(newSession) {
  if (newSession && is.null(rstudioapi::getActiveProject()))
    rstudioapi::openProject("{{.RstudioConf.Project}}/")
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

