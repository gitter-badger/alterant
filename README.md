# Alterant
Alterant is a lightweight provisioning tool built with ease of use, security, and flexibility in mind.

_Alter your machine with ease_.

### Features
* Encrypts sensitive data with OpenPGP keys
* Allows flexible organization of configurations
* Installs dotfiles with symlinks
* Executes scripts defined within the YAML
* Automatically resolves dependecies
* Intelligently performs updates
* Easy installation with zero dependencies

## Installation
Install the latest [release](https://github.com/autonomy/alterant/releases) in your `$PATH`.

## Documentation
For usage and examples see [Alterant](http://autonomy.github.io/alterant).

## Hacking
Compiling from source:
````bash
$ go get -d github.com/autonomy/alterant
$ go get -u github.com/FiloSottile/gvt
$ cd $GOPATH/src/github.com/autonomy/alterant
$ make deps
$ make [linux|darwin]
````
Before any pull request is accepted be sure to follow the guidelines outlined in [CONTRIBUTING.md](CONTRIBUTING.md).

## Built With

* Go
* Git2Go - libgit2 bindings
* OpenPGP - Golang implementation
* OpenSSL

## Authors

* **[The Autonomy Team](https://github.com/orgs/autonomy/people)**

As well as the [contributors](https://github.com/autonomy/alterant/contributors).

## License

This project is licensed under the Apache License 2.0 - see [LICENSE.md](LICENSE.md) for details.
