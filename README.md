# Goline

## About

Goline is an open source project for CI/CD Pipeline based on Jenkins 2.0.

It supports following project types:
- Maven
- Gradle
- Shell
- Batch

## Documentation

### Quick Start
This project requires the Go 1.5+. If you use the Go 1.5, please remember to set the environment `GO15VENDOREXPERIMENT=1` before build the project.

```sh
$ go get github.com/supereagle/goline
$ cd $GOPATH/src/github.com/supereagle/goline/cmd
$ go build -v -a -o goline
$ $EDITOR config.json # Add the config file refer to the example config file `config-example.json`
$ ./goline
```

## Licensing

goline is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/supereagle/goline/blob/master/LICENSE) for the full
license text.
