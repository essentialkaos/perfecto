<p align="center"><a href="#readme"><img src="https://gh.kaos.st/perfecto.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/perfecto"><img src="https://kaos.sh/r/perfecto.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/l/perfecto"><img src="https://kaos.sh/l/19f019d1310c2cb69b29.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/perfecto"><img src="https://kaos.sh/b/74af2307-8aa2-48eb-afd5-2ae3620a1149.svg" alt="Codebeat badge" /></a>
  <br/>
  <a href="https://kaos.sh/c/perfecto"><img src="https://kaos.sh/c/perfecto.svg" alt="Coverage Status" /></a>
  <a href="https://kaos.sh/w/perfecto/ci"><img src="https://kaos.sh/w/perfecto/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/perfecto/codeql"><img src="https://kaos.sh/w/perfecto/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#checks">Checks</a> • <a href="#installing">Installing</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

_perfecto_ is tool for checking perfectly written RPM specs. Currently, _perfecto_ used by default for checking specs for [EK Public Repository](https://yum.kaos.st).

![Screenshot](https://gh.kaos.st/perfecto.png)

![Screenshot](https://gh.kaos.st/perfecto2.png)

### Checks

You can find additional information about every _perfecto_ check in [project wiki](https://github.com/essentialkaos/perfecto/wiki).

### Installing

#### From sources

Make sure you have a working Go 1.19+ workspace ([instructions](https://go.dev/doc/install)), then:

```bash
go install github.com/essentialkaos/perfecto@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://pkgs.kaos.st)

```bash
sudo yum install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm

# EL7 (OracleLinux/CentOS 7)
sudo yum install perfecto

# EL8 (OracleLinux/Alma/Rocky 8)
sudo dnf install perfecto

# Alma/Rocky 9
sudo dnf --enablerepo=crb install perfecto
# OracleLinux 9
sudo dnf --enablerepo=ol9_codeready_builder install perfecto
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/perfecto/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) perfecto
```

#### Container image

Official _perfecto_ images available on [GitHub Container Registry](https://kaos.sh/p/perfecto) and [Docker Hub](https://kaos.sh/d/perfecto). Install the latest version of [Podman](https://podman.io/getting-started/installation.html) or [Docker](https://docs.docker.com/engine/install/), then:

```bash
curl -#L -o perfecto-container https://kaos.sh/perfecto/perfecto-container
chmod +x perfecto-container
sudo mv perfecto-container /usr/bin/perfecto
perfecto your.spec
```

Official container images with _perfecto_:

- [`ghcr.io/essentialkaos/perfecto:micro`](https://kaos.sh/p/perfecto)
- [`ghcr.io/essentialkaos/perfecto:ol7`](https://kaos.sh/p/perfecto)
- [`ghcr.io/essentialkaos/perfecto:ol8`](https://kaos.sh/p/perfecto)
- [`ghcr.io/essentialkaos/perfecto:ol9`](https://kaos.sh/p/perfecto)
- [`essentialkaos/perfecto:micro`](https://kaos.sh/d/perfecto)
- [`essentialkaos/perfecto:ol7`](https://kaos.sh/d/perfecto)
- [`essentialkaos/perfecto:ol8`](https://kaos.sh/d/perfecto)
- [`essentialkaos/perfecto:ol9`](https://kaos.sh/d/perfecto)

#### Using with Github Actions

For using latest stable version _perfecto_ with Github Actions use this `perfecto.yml` file or add it to your workflow:

```yaml
name: Perfecto

on:
  push:
    branches: [master, develop]
  pull_request:
    branches: [master]

jobs:
  Perfecto:
    name: Perfecto
    runs-on: ubuntu-latest

    steps:
      - name: Code checkout
        uses: actions/checkout@v3

      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Check specs with Perfecto
        uses: essentialkaos/perfecto-action@v1
        with:
          files: myapp.spec
```

Additional information about action configuration can be found on [the official GitHub action page](https://github.com/marketplace/actions/ek-perfecto).

### Usage

```
Usage: perfecto {options} file…

Options

  --ignore, -I id…           Disable one or more checks by their ID
  --format, -f format        Output format (summary|tiny|short|github|json|xml)
  --lint-config, -c file     Path to RPMLint configuration file
  --error-level, -e level    Return non-zero exit code if alert level greater than given (notice|warning|error|critical)
  --quiet, -q                Suppress all normal output
  --no-lint, -nl             Disable RPMLint checks
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
| `master` | [![CI](https://kaos.sh/w/perfecto/ci.svg?branch=master)](https://kaos.sh/w/perfecto/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/perfecto/ci.svg?branch=develop)](https://kaos.sh/w/perfecto/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
