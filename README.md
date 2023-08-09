<!-- Add warning that this is a VERY basic POC -->
# WARNING

**This is a very basic proof of concept repository for running Nextflow workflows remotely via Golang api. Use with caution!**


# Overview

This repository contains Golang code to execute Nextflow workflows through a basic api as subprocesses. The goal is to demonstrate a simple way to invoke Nextflow workflows from Go code. 

The main logic is in `main.go` which does the following:

- Construct a command to invoke a Nextflow workflow 
- Execute the command as a subprocess
- Stream stdout and stderr to log and to the http response. 

There is also a simple API defined:

- `/exec` - Execute the hello workflow
- `/exec/{workflow}` - Execute a specific workflow by name

Overall this shows a straightforward pattern for running Nextflow remotely from Golang, but lacks robustness for production use. Treat as a simple proof of concept only!


# Dev/ test

Either using the Dockerfile or just download nextflow binary and have JDK on your workstation. 

Update the nextflow binary location in code. 

If you have `nodemon` installed you can have a hot reloading dev environment set up like: 

## Dependencies
```
go mod download
```

## Run dev
```
nodemon --exec go run main.go --signal SIGTERM
```