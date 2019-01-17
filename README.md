![](https://github.com/gofunct/common/blob/master/logo/dark_logo_transparent_background.png?raw=true)

# Mamba

* Author: Coleman Word 
* Email: coleman.word@gofunct.com
* Download: `go get github.com/gofunct/mamba
* Description: gRPC service to execute shell commands transparently on a remote server

## Usage
```text
A general purpose scripting utility for developers and administrators

Usage:
  mamba [command]

Available Commands:
  exec        execute a script on a backend grpc server
  gcloud      
  help        Help about any command
  local       
  protoc      A brief description of your command
  replace     for all files with the .tmpl extension, replace the keys with the values present in the provided metadata
  serve       start a grpc server to handle remote script requests
  test        run an interactive web server to test code in current repository

Flags:
      --config string   config file (default is $HOME/.chronic.yaml)
  -h, --help            help for mamba

Use "mamba [command] --help" for more information about a command.

```

## Roadmap
High Priority:
 - [ ] Add .tf extension to supported config types

## Issues
- [ ] Step One
