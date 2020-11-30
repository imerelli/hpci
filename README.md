High Performance Computing Interface
==============

Descrizione
-----------

High Performance Computing Interface to spawn Singularity Application on Slurm Cluster

**Note**: Require Singularity and Nginx

Features
--------
* Runs Rstudio server on an HPC node

Requirements
------------

* Singularity 3.6.4
* Slurm 19.05.5
* Nginx 1.19.4


Usage
=====

Admin Configuration file (/etc/hpci/hpci.conf):

```yaml
---
singularity:
  binary: /opt/singularity/3.6.4/bin/singularity

cluster:
  url: hpc.bioinformatics.itb.cnr.it

ldaps:
  host: freeipa.cluster.bioinformatics.itb.cnr.it
  user-dn: uid=%s,cn=users,cn=accounts,dc=cluster,dc=bioinformatics,dc=itb,dc=cnr,dc=it
  ca-file: /etc/ipa/ca.crt
    
rstudio:
  sif: /opt/singularity/softwares/Rstudio/rstudio-1.2.5033.sif
  ports:
    from: 9000
    to: 9099
  auth: ldaps # password  
  job:
    name: Rstudio
    time: 08:00:00
    partition: master
    cpus: 2
    memory: 8192
    outdir: .rstudiojob
    output: rstudio-server.%j.out  
```

Nginx for Node Configuration:

```conf
server {
        listen 9000-9099;
        resolver RESOLVER_IP  valid=60s;
        server_name ~^(?<node>.+)-job(.+)\.hpc\.bioinformatics\.itb\.cnr\.it;

        location / {

                proxy_pass http://$node.cluster.bioinformatics.itb.cnr.it:$server_port;
                proxy_set_header   X-Forwarded-Host $host;
                proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto http;
                proxy_set_header Host $http_host;
                proxy_redirect off;
                proxy_set_header X-Real-IP $remote_addr;
        }
}
```

User Configuration file:
```yaml
---
Rversion: /opt/R/3.6.1
project: /datasets/999/P1
directory:
  - dir1
  - dir2 
```

Command Line Interface:
```console
hpci rstudio -r /path/to/R -p /path/to/project -d /path/to/extra/directory

hpci rstudio -f /path/to/user/configuration/file
```

Flags
=========

```console
hpci rstudio -h
hpci rstudio spawn application.

Usage:
  hpci rstudio [flags]

Flags:
  -r, --Rversion string     path to R (absolute path)
  -d, --directory strings   add extra directory (absolute path separate with comma or repeat flag)
  -f, --file string         config file 
  -h, --help                help for rstudio
  -p, --project string      project directory (absolute path)
```

Flag | Description |
-------- | ----------- 
`-r` | path to R (absolute path)
`-d` | add extra directory (absolute path separate with comma or repeat flag)
`-f` | config file
`-p` | project directory (absolute path) 

Create Executable
=========
The Makefile use docker to create the specific executable
The executable is created in the bin directory

```console
make PLATFORM=linux/amd64 
```

***Note*** Is important that  BuildKit is  enabled. (export  DOCKER_BUILDKIT=1)   

Create Singularity Image
=========
Singularity image for Rstudio Server

```console
singularity build rstudio-1.2.5033.sif Rstudio-1.2.5033.def
```


Authors
=======

Written by Marco Moscatelli <marco.moscatelli@itb.cnr.it>

Revised and Tested by Matteo Gnocchi and Ivan Merelli
