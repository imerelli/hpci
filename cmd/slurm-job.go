package cmd

import (
	"math/rand"
	"time"
	"strings"
)


func SlurmTemplate() string {

	const tpl = `#!/bin/sh
#SBATCH --time=08:00:00
#SBATCH --signal=USR2
#SBATCH --ntasks=2
#SBATCH --job-name=Rstudio
#SBATCH --cpus-per-task=2
#SBATCH --mem=8192
#SBATCH --error={{.WorkDir}}/rstudio-server.%j.err
#SBATCH --output={{.WorkDir}}/rstudio-server.%j.out
#SBATCH --partition=p-cpus

export PASSWORD=$(openssl rand -base64 15)

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


{{.SingularityBin}} exec --bind {{.WorkDir}}/.Rprofile:/usr/local/lib64/R/etc/Rprofile.site {{.SingularityImage}} rserver --www-port=${PORT} --auth-none=0  --auth-pam-helper-path=pam-helper 


printf 'rserver exited' 1>&2
`
	return tpl
}

func RandomString() string{
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
	    "abcdefghijklmnopqrstuvwxyzåäö" +
	    "0123456789")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
	    b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
