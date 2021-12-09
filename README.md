<p align="center"><a href="#readme"><img src="https://gh.kaos.st/perfecto.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/perfecto"><img src="https://kaos.sh/r/perfecto.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/perfecto"><img src="https://kaos.sh/b/74af2307-8aa2-48eb-afd5-2ae3620a1149.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/perfecto/ci"><img src="https://kaos.sh/w/perfecto/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/perfecto/codeql"><img src="https://kaos.sh/w/perfecto/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/c/perfecto"><img src="https://kaos.sh/c/perfecto.svg" alt="Coverage Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#checks">Checks</a> • <a href="#installing">Installing</a> • <a href="#using-with-github-actions">Using with Github Actions</a> • <a href="#using-with-docker">Using with Docker</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

_perfecto_ is tool for checking perfectly written RPM specs. Currently, _perfecto_ used by default for checking specs for [EK Public Repository](https://yum.kaos.st).

![Screenshot](https://gh.kaos.st/perfecto.png)

![Screenshot](https://gh.kaos.st/perfecto2.png)

### Checks

You can find additional information about every _perfecto_ check in [project wiki](https://github.com/essentialkaos/perfecto/wiki).

### Installing

#### From sources

Make sure you have a working Go 1.16+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go install github.com/essentialkaos/perfecto
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install perfecto
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/perfecto/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) perfecto
```

### Using with Github Actions

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
        uses: actions/checkout@v2

      # https://docs.docker.com/docker-hub/download-rate-limit/
      - name: Login to DockerHub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Run Perfecto docker image
        uses: docker://essentialkaos/perfecto:centos7
        with:
          args: --version

      - name: Install perfecto-docker
        run: |
          wget https://kaos.sh/perfecto/perfecto-docker
          chmod +x perfecto-docker

      - name: Run Perfecto check
        env:
          IMAGE: essentialkaos/perfecto:centos7
        run: ./perfecto-docker your-app.spec

```

### Using with Docker

Install latest version of Docker, then:

```bash
curl -fL# -o perfecto-docker https://kaos.sh/perfecto/perfecto-docker
chmod +x perfecto-docker
sudo mv perfecto-docker /usr/bin/
perfecto-docker PATH_TO_YOUR_SPEC_HERE
```

### Usage

```
Usage: perfecto {options} file…

Options

  --absolve, -A id…          Disable some checks by their ID
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
| `master` | [![CI](https://kaos.sh/w/perfecto/ci.svg?branch=master)](https://kaos.sh/w/perfecto/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/perfecto/ci.svg?branch=develop)](https://kaos.sh/w/perfecto/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
