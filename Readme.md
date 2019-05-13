![MGW Logo](https://github.com/affiliator/mgw/raw/master/.github/banner.png "MGW Logo")

**Note**: This project is under heavy development, more information coming soon. 

[![Build Status](https://travis-ci.org/affiliator/mgw.svg?branch=master)](https://travis-ci.org/affiliator/mgw)
[![Coverage Status](https://coveralls.io/repos/github/affiliator/mgw/badge.svg?branch=master)](https://coveralls.io/github/affiliator/mgw?branch=master)

## Requirements
You need [golang](https://github.com/golang/go) in version 1.11+ and [dep](https://github.com/golang/dep) to handle the dependencies.
Also for productive deployments you may want to use Nginx as reverse proxy to force authentication. Otherwise you should make sure your service is not bound to a public interface.

## Usage
After making sure your target system met all requirements, you can proceed with building & deploying the software.

Install dependencies
```
$ make vendor
```

You can build from source using make:
```
$ make build
```

Then you need to configure everything by copying the `config.example.json` to `config.json` or by running:
```
$ make prepare
```

To start the daemon you can either use the make command:
```
$ make run 
```

Or execute the `mgw` binary directly:
```
$ ./mgw serve --config="config.json"
```

# Todo
 - [ ] --pid (-p) param will be ignored
 - [ ] When the `--config (-c)` is missing, everything crashes.  
 - [x] More documentation & add readme
 - [x] Add Makefile functionality
 - [ ] Add Mailgun as processor
 - [ ] Add dokku deployment
 - [ ] Add a proper documentation
 - [ ] Store received e-mails in a database
 - [ ] Store E-Mail state in database
 - [ ] Add API to read E-Mails State
 - [ ] Make it possible to add add e-mail properties through header
 - [x] general refactoring, this is my first go project. 
 - [x] Add Jan as Contributor
 - [ ] Add proper logging adapter
 - [ ] Allow configuration overriding via env variables
 - [ ] Add unit tests
 - [ ] Build and test in travis