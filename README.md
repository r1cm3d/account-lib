# account-lib

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ricardomedeirosdacostajunior/aws-poc)
[![license](https://img.shields.io/badge/license-BSD--2-red)](./LICENSE)

**TL;DR:**
```console
make
```

## Prerequisites
[![Docker](https://img.shields.io/badge/Docker-19.03.9-blue)](https://www.docker.com/)
[![Docker-compose](https://img.shields.io/badge/Docker--compose-1.29.2-blue)](https://github.com/docker/compose/releases)
[![GNU Make](https://img.shields.io/badge/GNU%20Make-4.2.1-lightgrey)](https://www.gnu.org/software/make/)
[![GNU Bash](https://img.shields.io/badge/GNU%20Bash-4.2.1-lightgrey)](https://www.gnu.org/software/bash/)

## Table of Contents
* [TL;DR](#account-lib)
* [Prerequisites](#prerequisites)
* [About the Project](#about-the-project)
* [Testing](#testing)

## About The Project

This Go client library abstract an Account HTTP API that is simulated using this [docker-compose](./docker-compose.yaml).
Aiming to be simple as possible, this project follows the same structure as 
[decimal](https://github.com/shopspring/decimal) library. It also applies the
patterns found in [Uber Style Guide](https://github.com/uber-go/guide/blob/master/style.md). This guide has 
useful insights about of writing Go code in general.

### Testing
#### Unit tests
[comment]: <> (TODO: add description here)
```sh
make unit-test
```

#### Integration tests
[comment]: <> (TODO: add description here)
```sh
make integration-test
```

#### All tests
[comment]: <> (TODO: add description here)
```sh
make test
```