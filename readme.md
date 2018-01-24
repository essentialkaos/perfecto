<p align="center"><a href="#readme"><img src="https://gh.kaos.io/perfecto.svg"/></a></p>

<p align="center"><a href="#installing">Installing</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/essentialkaos/perfecto"><img src="https://goreportcard.com/badge/github.com/essentialkaos/perfecto"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-perfecto-master"><img alt="codebeat badge" src="https://codebeat.co/badges/74af2307-8aa2-48eb-afd5-2ae3620a1149" /></a>
  <a href="https://travis-ci.org/essentialkaos/perfecto"><img src="https://travis-ci.org/essentialkaos/perfecto.svg"></a>
  <a href="#license"><img src="https://gh.kaos.io/ekol.svg"></a>
</p>

_perfecto_ is tool for checking perfectly written RPM specs.

![Screenshot](https://gh.kaos.io/perfecto.png)

### Installing

#### From sources

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.7+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/perfecto
```

For update to latest stable release, do:

```
go get -u github.com/essentialkaos/perfecto
```


#### From ESSENTIAL KAOS Public repo for RHEL6/CentOS6

```bash
[sudo] yum install -y https://yum.kaos.io/6/release/x86_64/kaos-repo-8.0-0.el6.noarch.rpm
[sudo] yum install perfecto
```

#### From ESSENTIAL KAOS Public repo for RHEL7/CentOS7

```bash
[sudo] yum install -y https://yum.kaos.io/7/release/x86_64/kaos-repo-8.0-0.el7.noarch.rpm
[sudo] yum install perfecto
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.io/perfecto/latest).

### Usage

```
Usage: spec-file {options}

Options

  --summary, -s      Print only summary info
  --no-lint, -nl     Disable rpmlint checks
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/perfecto.svg?branch=master)](https://travis-ci.org/essentialkaos/perfecto) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/perfecto.svg?branch=develop)](https://travis-ci.org/essentialkaos/perfecto) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
