<p align="center"><a href="#readme"><img src="https://gh.kaos.st/perfecto.svg"/></a></p>

<p align="center"><a href="#installing">Installing</a> • <a href="#using-on-travisci">Using on TravisCI</a> • <a href="#using-with-docker">Using with Docker</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/essentialkaos/perfecto"><img src="https://goreportcard.com/badge/github.com/essentialkaos/perfecto"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-perfecto-master"><img alt="codebeat badge" src="https://codebeat.co/badges/74af2307-8aa2-48eb-afd5-2ae3620a1149" /></a>
  <a href="https://travis-ci.org/essentialkaos/perfecto"><img src="https://travis-ci.org/essentialkaos/perfecto.svg"></a>
  <a href='https://coveralls.io/github/essentialkaos/perfecto'><img src='https://coveralls.io/repos/github/essentialkaos/perfecto/badge.svg' alt='Coverage Status' /></a>
  <a href="#license"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

_perfecto_ is tool for checking perfectly written RPM specs. Currently, _perfecto_ used by default for checking specs for [EK Public Repository](https://yum.kaos.st).

![Screenshot](https://gh.kaos.st/perfecto.png)

![Screenshot](https://gh.kaos.st/perfecto2.png)

### Installing

#### From sources

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.8+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/perfecto
```

For update to latest stable release, do:

```
go get -u github.com/essentialkaos/perfecto
```

#### From ESSENTIAL KAOS Public repo for RHEL6/CentOS6

```bash
[sudo] yum install -y https://yum.kaos.st/6/release/x86_64/kaos-repo-9.1-0.el6.noarch.rpm
[sudo] yum install perfecto
```

#### From ESSENTIAL KAOS Public repo for RHEL7/CentOS7

```bash
[sudo] yum install -y https://yum.kaos.st/7/release/x86_64/kaos-repo-9.1-0.el7.noarch.rpm
[sudo] yum install perfecto
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/perfecto/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) perfecto
```

### Using on TravisCI

For using latest stable version _perfecto_ on TravisCI use this `.travis.yml` file:

```yaml
language: bash

cache: apt

before_install:
  - echo "deb http://us.archive.ubuntu.com/ubuntu xenial main universe" | sudo tee -a /etc/apt/sources.list
  - sudo apt-get update -qq
  - sudo apt-get install -y rpmlint
  - sudo ln -sf /usr/bin/python2.7 /usr/bin/python2.6
  - wget https://apps.kaos.st/perfecto/latest/linux/x86_64/perfecto
  - chmod +x perfecto
  - ./perfecto -v

script:
  - ./perfecto PATH_TO_YOUR_SPEC_HERE

```

or this:

```yaml
services:
  - docker

env:
  global:
    - IMAGE=essentialkaos/perfecto:centos7

before_install:
  - docker pull "$IMAGE"
  - wget https://kaos.sh/perfecto/perfecto-docker
  - chmod +x perfecto-docker

script:
  - ./perfecto-docker PATH_TO_YOUR_SPEC_HERE

```

or this:

```yaml
services:
  - docker

env:
  global:
    - IMAGE=essentialkaos/perfecto:centos7

before_install:
  - docker pull "$IMAGE"

script:
  - bash <(curl -fsSL https://kaos.sh/perfecto/perfecto-docker) PATH_TO_YOUR_SPEC_HERE

```

### Using with Docker

Install latest version of Docker, then:

```bash
curl -o perfecto-docker https://kaos.sh/perfecto/perfecto-docker
chmod +x perfecto-docker
[sudo] mv perfecto-docker /usr/bin/
perfecto-docker PATH_TO_YOUR_SPEC_HERE
```

### Usage

```
Usage: perfecto {options} file…

Options

  --format, -f format        Output format (summary|tiny|short|json|xml)
  --lint-config, -c file     Path to rpmlint configuration file
  --error-level, -e level    Return non-zero exit code if alert level greater than given (notice|warning|error|critical)
  --quiet, -q                Suppress all normal output
  --no-lint, -nl             Disable rpmlint checks
  --no-color, -nc            Disable colors in output
  --help, -h                 Show this help message
  --version, -v              Show version

Examples

  perfecto app.spec
  Check spec and print extended report

  perfecto --no-lint app.spec
  Check spec without rpmlint and print extended report

  perfecto --format tiny app.spec
  Check spec and print tiny report

  perfecto --format summary app.spec
  Check spec and print summary

  perfecto --format json app.spec 1> report.json
  Check spec, generate report in JSON format and save as report.json

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/perfecto.svg?branch=master)](https://travis-ci.org/essentialkaos/perfecto) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/perfecto.svg?branch=develop)](https://travis-ci.org/essentialkaos/perfecto) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
