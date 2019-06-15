<img src="https://raw.githubusercontent.com/goodwithtech/dockle/master/imgs/logo.png" width="450">

[![GitHub release](https://img.shields.io/github/release/goodwithtech/dockle.svg)](https://github.com/goodwithtech/dockle/releases/latest)
[![CircleCI](https://circleci.com/gh/goodwithtech/dockle.svg?style=svg)](https://circleci.com/gh/goodwithtech/dockle)
[![Go Report Card](https://goreportcard.com/badge/github.com/goodwithtech/dockle)](https://goreportcard.com/report/github.com/goodwithtech/dockle)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

> Dockle - Simple Security Auditing and helping build the Best Docker Images

`Dockle` helps you:

1. Build secure Docker images
    - Checkpoints includes [CIS Benchmarks](https://www.cisecurity.org/cis-benchmarks/)
2. Build [Best Practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) Docker images

To check your Docker image, only run:

```bash
$ brew install goodwithtech/dockle/dockle
$ dockle [YOUR_IMAGE_NAME]
```

<img src="https://raw.githubusercontent.com/goodwithtech/dockle/master/imgs/usage_pass_light.png" width="800">
<img src="https://raw.githubusercontent.com/goodwithtech/dockle/master/imgs/usage_fail_light.png" width="800">


# TOC

- [Features](#features)
- [Comparison](#comparison)
- [Installation](#installation)
  - [Linuxbrew](#linuxbrew)
  - [RHEL/CentOS](#rhelcentos)
  - [Debian/Ubuntu](#debianubuntu)
  - [Mac OS X / Homebrew](#mac-os-x--homebrew)
  - [Windows](#windows)
  - [Binary](#binary)
  - [From source](#from-source)
  - [Use Docker](#use-docker)
- [Checkpoint Summary](#checkpoint-summary)
- [Quick Start](#quick-start)
  - [Basic](#basic)
  - [Docker](#docker)
- [Examples](EXAMPLES.md)
  - Scan an image
  - Scan an image file
  - Get or Save the results as JSON
  - Specify exit code
  - Ignore the specified checkpoints
  - Clear image caches
- [Continuous Integration](EXAMPLES.md#continuous-integration-ci)
  - Travis CI
  - CircleCI
  - Authorization for Private Docker Registry
- [Checkpoint Details](CHECKPOINT.md)
  - CIS's Docker Image Checkpoints
  - Dockle Checkpoints for Docker
  - Dockle Checkpoints for Linux
- [Credits](#credits)
- [Roadmap](#roadmap)

# Features

- Detect container's vulnerabilities
- Helping build best-practice Dockerfile
- Simple usage
  - Specify only the image name
  - See [Quick Start](#quick-start) and [Examples](EXAMPLES.md)
- CIS Benchmarks Support
  - High accuracy
- DevSecOps
  - Suitable for CI such as Travis CI, CircleCI, Jenkins, etc.
  - See [CI Example](EXAMPLES.md#continuous-integration-ci)

# Comparison

|  | [Dockle](https://github.com/goodwithtech/dockle) | [Hadolint](https://github.com/hadolint/hadolint) | [Docker Bench for Security](https://github.com/docker/docker-bench-security) |
|--- |---:|---:|---:|
| Target |  BuildImage | Dockerfile | Host<br/>DockerDaemon<br/>BuildImage<br/>ContainerRuntime |
| How to run | Binary | Binary | ShellScript |
| Dependency | No | No | Some dependencies |
| CI Suitable | ✓ | ✓ | x |
| Purpose |SecurityAudit<br/>DockerfileLint| DockerfileLint | SecurityAudit<br/>DockerfileLint |
| Covered CIS Benchmarks (Docker Image and Build File) | 7 | 3 | 5 |

<details>
<summary>Detail of CIS Benchmark</summary>

|  | [Dockle](https://github.com/goodwithtech/dockle) | [Docker Bench for Security](https://github.com/docker/docker-bench-security) | [Hadolint](https://github.com/hadolint/hadolint) |
|---|:---:|:---:|:---:|
| 1.  Create a user for the container | ✓ | ✓ | ✓ |
| 2.  Use trusted base images for containers | - | – | - |
| 3.  Do not install unnecessary packages in the container | - | - | - |
| 4.  Scan and rebuild the images to include security patches | - | - | - |
| 5.  Enable Content trust for Docker | ✓ | ✓ | - |
| 6.  Add `HEALTHCHECK` instruction to the container image | ✓ | ✓ | - |
| 7.  Do not use `update` instructions alone in the Dockerfile | ✓ | ✓ | ✓|
| 8.  Remove `setuid` and `setgid` permissions in the images | ✓ | - | - |
| 9.  Use `COPY` instead of `ADD` in Dockerfile | ✓ | ✓ | ✓|
| 10. Do not store secrets in Dockerfiles | ✓ | - | - |
| 11. Install verified packages only | -  |  - | - |
| |7|5|3|

All checkpoints [here](#checkpoint-summary)!

</details>

# Installation

## Linuxbrew

You can use [Homebrew](https://docs.brew.sh/Homebrew-on-Linux) on Linux and WSL (Windows Subsystem for Linux).

```bash
$ brew install goodwithtech/dockle/dockle
```

## RHEL/CentOS

```bash
$ VERSION=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
) && rpm -ivh https://github.com/goodwithtech/dockle/releases/download/v${VERSION}/dockle_${VERSION}_Linux-64bit.rpm
```

## Debian/Ubuntu

```bash
$ VERSION=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
) && curl -L -o dockle.deb https://github.com/goodwithtech/dockle/releases/download/v${VERSION}/dockle_${VERSION}_Linux-64bit.deb
$ sudo dpkg -i dockle.deb && rm dockle.deb
```

## Mac OS X / Homebrew

You can use [Homebrew](https://brew.sh/) on macOS.

```bash
$ brew install goodwithtech/dockle/dockle
```

## Windows

```bash
$ VERSION=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
) && curl -L -o dockle.zip https://github.com/goodwithtech/dockle/releases/download/v${VERSION}/dockle_${VERSION}_Windows-64bit.zip
$ unzip dockle.zip && rm dockle.zip
$ ./dockle.exe [IMAGE_NAME]
```

## Binary

You can get the latest version binary from [releases page](https://github.com/goodwithtech/dockle/releases/latest).

Download the archive file for your operating system/architecture. Unpack the archive, and put the binary somewhere in your `$PATH` (on UNIX-y systems, `/usr/local/bin` or the like).

- NOTE: Make sure that it's execution bits turned on. (`chmod +x dockle`)

## From source

```bash
$ GO111MODULE=off go get github.com/goodwithtech/dockle/cmd/dockle
$ cd $GOPATH/src/github.com/goodwithtech/dockle && GO111MODULE=on go build -o $GOPATH/bin/dockle cmd/dockle/main.go
```

## Use Docker

There's a [`Dockle` image on Docker Hub](https://hub.docker.com/r/goodwithtech/dockle) also. You can try `dockle` before installing the command.

```
$ VERSION=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
) && docker run --rm goodwithtech/dockle:${VERSION} [YOUR_IMAGE_NAME]
```

# Quick Start

Here's a quick start. For more detailed usage and samples, such as using `dockle` on CIs, see [EXAMPLES.md](./EXAMPLES.md).


## Basic

Simply specify an image name (and a tag).

```bash
$ dockle [YOUR_IMAGE_NAME]
```

<details>
<summary>Result</summary>

```
FATAL   - Create a user for the container
        * Last user should not be root
WARN    - Enable Content trust for Docker
        * export DOCKER_CONTENT_TRUST=1 before docker pull/build
FATAL   - Add HEALTHCHECK instruction to the container image
        * not found HEALTHCHECK statement
FATAL   - Do not use update instructions alone in the Dockerfile
        * Use 'Always combine RUN apt-get update with apt-get install' : /bin/sh -c apt-get update && apt-get install -y git
PASS    - Remove setuid and setgid permissions in the images
FATAL   - Use COPY instead of ADD in Dockerfile
        * Use COPY : /bin/sh -c #(nop) ADD file:81c0a803075715d1a6b4f75a29f8a01b21cc170cfc1bff6702317d1be2fe71a3 in /app/credentials.json
FATAL   - Do not store secrets in ENVIRONMENT variables
        * Suspicious ENV key found : MYSQL_PASSWD
FATAL   - Do not store secret files
        * Suspicious filename found : app/credentials.json
PASS    - Avoid sudo command
FATAL   - Avoid sensitive directory mounting
        * Avoid mounting sensitive dirs : /usr
PASS    - Avoid apt-get/apk/dist-upgrade
PASS    - Use apk add with --no-cache
FATAL   - Clear apt-get caches
        * Use 'apt-get clean && rm -rf /var/lib/apt/lists/*' : /bin/sh -c apt-get update && apt-get install -y git
PASS    - Avoid latest tag
FATAL   - Avoid empty password
        * No password user found! username : nopasswd
PASS    - Be unique UID
PASS    - Be unique GROUP
```

</details>

## Docker

Also, you can use Docker to use `dockle` command as follow.

```bash
$ export DOCKLE_LATEST=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
)
$ docker run --rm goodwithtech/dockle:${DOCKLE_LATEST} [YOUR_IMAGE_NAME]
```

For more suitable use, I suggest mounting a cache directory. Replace `[YOUR_CACHE_DIR]` below with the cache directory on your machine.

```bash
$ export DOCKLE_LATEST=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
)
$ docker run --rm -v [YOUR_CACHE_DIR]:/root/.cache/ goodwithtech/dockle:${DOCKLE_LATEST} [YOUR_IMAGE_NAME]
```

- Example for macOS:

    ```bash
    $ docker run --rm -v $HOME/Library/Caches:/root/.cache/ goodwithtech/dockle:${DOCKLE_LATEST} [YOUR_IMAGE_NAME]
    ```

- If you'd like to scan the image on your host machine, you need to mount `docker.sock`.

    ```bash
    $ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ...
    ```

# Checkpoint Summary

- Details of each checkpoint see [CHECKPOINT.md](CHECKPOINT.md)

| CODE | DESCRIPTION | LEVEL[※](#level) |
|---|---|:---:|
| | CIS's Docker Image Checkpoints | |
| CIS-DI-0001 | Create a user for the container | WARN |
| CIS-DI-0002 | Use trusted base images for containers | FATAL
| CIS-DI-0003 | Do not install unnecessary packages in the container | FATAL
| CIS-DI-0004 | Scan and rebuild the images to include security patches | FATAL
| CIS-DI-0006 | Add `HEALTHCHECK` instruction to the container image | WARN
| CIS-DI-0007 | Do not use `update` instructions alone in the Dockerfile | FATAL
| CIS-DI-0008 | Remove `setuid` and `setgid` permissions in the images | INFO
| CIS-DI-0009 | Use `COPY` instead of `ADD` in Dockerfile | FATAL
| CIS-DI-0010 | Do not store secrets in Dockerfiles | FATAL
| CIS-DI-0011 | Install verified packages only | INFO
|| Dockle Checkpoints for Docker |
| DKL-DI-0001 | Avoid `sudo` command | FATAL
| DKL-DI-0002 | Avoid sensitive directory mounting | FATAL
| DKL-DI-0003 | Avoid `apt-get upgrade`, `apk upgrade`, `dist-upgrade` | FATAL
| DKL-DI-0004 | Use `apk add` with `--no-cache` | FATAL
| DKL-DI-0005 | Clear `apt-get` caches | FATAL
| DKL-DI-0006 | Avoid `latest` tag | WARN
|| Dockle Checkpoints for Linux |
| DKL-LI-0001 | Avoid empty password | FATAL
| DKL-LI-0002 | Be unique UID/GROUPs | FATAL

## Level

`Dockle` has 5 check levels.

| LEVEL | DESCRIPTION |
|:---:|---|
| FATAL | Be practical and prudent |
| WARN | Be practical and prudent, but limited uses (official docker image ) |
| INFO | May negatively inhibit the utility or performance |
| SKIP | Not found target files |
| PASS | Not found any problems |

# Credits

Special Thanks to [@knqyf263](https://github.com/knqyf263) (Teppei Fukuda) and [Trivy](https://github.com/knqyf263/trivy)

# License

- AGPLv3

# Author

[@tomoyamachi](https://github.com/tomoyamachi) (Tomoya Amachi)

# Roadmap

- [x] JSON output
- [ ] Check php.ini file
- [ ] Check nginx.conf file
- [ ] create CI badges
- Check /etc/hosts
  - [ ] duplicates
  - [ ] hostname
  - [ ] localhost
- Packages
  - [ ] Package managers
- File Permissions
  - [ ] Insecure permission
- Image Size
  - [ ] check large size container

## if running docker daemon...

- Networking
  - [ ] `docker port container` if docker running
  - [ ] by file
    - `/proc/1/net/tcp` : openning port (if running)
- Volume mount
  - dangerous mount
    - `/boot`, `/dev`, `/etc`, `/lib`
