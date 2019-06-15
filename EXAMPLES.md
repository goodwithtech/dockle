# Detailed Usages and Examples

## TOC

- Common Examples
  - Scan an image
  - Scan an image file
  - Get or Save the results as JSON
  - Specify exit code
  -  Ignore the specified checkpoints
  - Clear image caches
- Continuous Integration (CI)
  - Travis CI
  - CircleCI
- Authorization for Private Docker Registry
  - Docker Hub
  - Amazon ECR (Elastic Container Registry)
  - GCR (Google Container Registry)
  - Self Hosted Registry (BasicAuth)

## Common Examples

### Scan an image

Simply specify an image name (and a tag).

```bash
$ dockle goodwithtech/test-image:v1
```

<details>
<summary>Result</summary>

```
FATAL   - CIS-DI-0001: Create a user for the container
        * Last user should not be root
WARN    - CIS-DI-0005: Enable Content trust for Docker
        * export DOCKER_CONTENT_TRUST=1 before docker pull/build
FATAL   - CIS-DI-0006: Add HEALTHCHECK instruction to the container image
        * not found HEALTHCHECK statement
FATAL   - CIS-DI-0007: Do not use update instructions alone in the Dockerfile
        * Use 'Always combine RUN 'apt-get update' with 'apt-get install' : /bin/sh -c apt-get update && apt-get install -y git
FATAL   - CIS-DI-0008: Remove setuid and setgid permissions in the images
        * Found setuid file: etc/passwd grw-r--r--
        * Found setuid file: usr/lib/openssh/ssh-keysign urwxr-xr-x
        * Found setuid file: app/hoge.txt ugrw-r--r--
        * Found setuid file: app/hoge.txt ugrw-r--r--
        * Found setuid file: etc/shadow urw-r-----
FATAL   - CIS-DI-0009: Use COPY instead of ADD in Dockerfile
        * Use COPY : /bin/sh -c #(nop) ADD file:81c0a803075715d1a6b4f75a29f8a01b21cc170cfc1bff6702317d1be2fe71a3 in /app/credentials.json
FATAL   - CIS-DI-0010: Do not store secrets in ENVIRONMENT variables
        * Suspicious ENV key found : MYSQL_PASSWD
FATAL   - CIS-DI-0010: Do not store secret files
        * Suspicious filename found : app/credentials.json
PASS    - DKL-DI-0001: Avoid sudo command
FATAL   - DKL-DI-0002: Avoid sensitive directory mounting
        * Avoid mounting sensitive dirs : /usr
PASS    - DKL-DI-0003: Avoid apt-get/apk/dist-upgrade
PASS    - DKL-DI-0004: Use apk add with --no-cache
FATAL   - DKL-DI-0005: Clear apt-get caches
        * Use 'apt-get clean && rm -rf /var/lib/apt/lists/*' : /bin/sh -c apt-get update && apt-get install -y git
PASS    - DKL-DI-0006: Avoid latest tag
FATAL   - DKL-LI-0001: Avoid empty password
        * No password user found! username : nopasswd
PASS    - DKL-LI-0002: Be unique UID
PASS    - DKL-LI-0002: Be unique GROUP
```
</details>

### Scan an image file

```bash
$ docker save alpine:latest -o alpine.tar
$ dockle --input alpine.tar
```

### Get or Save the results as JSON

```bash
$ dockle -f json goodwithtech/test-image:v1
$ dockle -f json -o results.json goodwithtech/test-image:v1
```

<details>
<summary>Result</summary>

```json
{
  "summary": {
    "fatal": 6,
    "warn": 2,
    "info": 2,
    "pass": 7
  },
  "details": [
    {
      "code": "CIS-DI-0001",
      "title": "Create a user for the container",
      "level": "WARN",
      "alerts": [
        "Last user should not be root"
      ]
    },
    {
      "code": "CIS-DI-0005",
      "title": "Enable Content trust for Docker",
      "level": "INFO",
      "alerts": [
        "export DOCKER_CONTENT_TRUST=1 before docker pull/build"
      ]
    },
    {
      "code": "CIS-DI-0006",
      "title": "Add HEALTHCHECK instruction to the container image",
      "level": "WARN",
      "alerts": [
        "not found HEALTHCHECK statement"
      ]
    },
    {
      "code": "CIS-DI-0008",
      "title": "Remove setuid and setgid permissions in the images",
      "level": "INFO",
      "alerts": [
        "Found setuid file: usr/lib/openssh/ssh-keysign urwxr-xr-x"
      ]
    },
    {
      "code": "CIS-DI-0009",
      "title": "Use COPY instead of ADD in Dockerfile",
      "level": "FATAL",
      "alerts": [
        "Use COPY : /bin/sh -c #(nop) ADD file:81c0a803075715d1a6b4f75a29f8a01b21cc170cfc1bff6702317d1be2fe71a3 in /app/credentials.json "
      ]
    },
    {
      "code": "CIS-DI-0010",
      "title": "Do not store secrets in ENVIRONMENT variables",
      "level": "FATAL",
      "alerts": [
        "Suspicious ENV key found : MYSQL_PASSWD"
      ]
    },
    {
      "code": "CIS-DI-0010",
      "title": "Do not store secret files",
      "level": "FATAL",
      "alerts": [
        "Suspicious filename found : app/credentials.json "
      ]
    },
    {
      "code": "DKL-DI-0002",
      "title": "Avoid sensitive directory mounting",
      "level": "FATAL",
      "alerts": [
        "Avoid mounting sensitive dirs : /usr"
      ]
    },
    {
      "code": "DKL-DI-0005",
      "title": "Clear apt-get caches",
      "level": "FATAL",
      "alerts": [
        "Use 'rm -rf /var/lib/apt/lists' after 'apt-get install' : /bin/sh -c apt-get update \u0026\u0026 apt-get install -y git"
      ]
    },
    {
      "code": "DKL-LI-0001",
      "title": "Avoid empty password",
      "level": "FATAL",
      "alerts": [
        "No password user found! username : nopasswd"
      ]
    }
  ]
}
```

</details>

### Specify exit code

By default, `Dockle` exits with code `0` even if there are some problems.

Use the `--exit-code` option to exit with a non-zero exit code if any alert were found.

```bash
$ dockle  -exist-code 1 [IMAGE_NAME]
```

### Ignore the specified checkpoints

Use `.dockleignore` file.

```bash
$ cat .dockleignore
# set root to default user because we want to run nginx
CIS-DI-0001
# Use latest tag because to check the image inside only
DKL-DI-0006
```

### Clear image caches

The `--clear-cache` option removes image caches. Use this option if the image with the same tag is updated (such as `latest` tag).

```bash
$ dockle --clear-cache python:3.7
```

## Continuous Integration (CI)

You can scan your built image with `Dockle` in Travis CI/CircleCI.

In these examples, the test will fail with if any warnings were found.

Though, you can ignore the specified target checkpoints by using `.dockleignore` file.

Or, if you just want the results to display and not let the test fail for this, specify `--exit-code` to `0` in `dockle` command.

### Travis CI

```yaml
services:
  - docker

env:
  global:
    - COMMIT=${TRAVIS_COMMIT::8}

before_install:
  - docker build -t dockle-ci-test:${COMMIT} .
  - export VERSION=$(curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
  - wget https://github.com/goodwithtech/dockle/releases/download/v${VERSION}/dockle_${VERSION}_Linux-64bit.tar.gz
  - tar zxvf dockle_${VERSION}_Linux-64bit.tar.gz
script:
  - ./dockle dockle-ci-test:${COMMIT}
  - ./dockle --exit-code 1 dockle-ci-test:${COMMIT}
```

- Example: https://travis-ci.org/goodwithtech/dockle-ci-test
- Repository: https://github.com/goodwithtech/dockle-ci-test

### CircleCI

```yaml
jobs:
  build:
    docker:
      - image: docker:18.09-git
    steps:
      - checkout
      - setup_remote_docker
      - run:
          name: Build image
          command: docker build -t dockle-ci-test:${CIRCLE_SHA1} .
      - run:
          name: Install dockle
          command: |
            apk add --update curl
            VERSION=$(
                curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
                grep '"tag_name":' | \
                sed -E 's/.*"v([^"]+)".*/\1/'
            )
            wget https://github.com/goodwithtech/dockle/releases/download/v${VERSION}/dockle_${VERSION}_Linux-64bit.tar.gz
            tar zxvf dockle_${VERSION}_Linux-64bit.tar.gz
            mv dockle /usr/local/bin
      - run:
          name: Scan the local image with dockle
          command: dockle --exit-code 1 dockle-ci-test:${CIRCLE_SHA1}
workflows:
  version: 2
  release:
    jobs:
      - build
```

- Example: https://circleci.com/gh/goodwithtech/dockle-ci-test
- Repository: https://github.com/goodwithtech/dockle-ci-test

## Authorization for Private Docker Registry

`Dockle` can download images from a private registry, without installing `Docker` or any other 3rd party tools. It's designed so for ease of use in a CI process.

All you have to do is: install `Dockle` and set ENVIRONMENT variables.

- NOTE: I don't recommend using ENV vars in your local machine.

### Docker Hub

To download the private repository from Docker Hub, you need to set `DOCKLE_AUTH_URL`, `DOCKLE_USERNAME` and `DOCKLE_PASSWORD` ENV vars.


```bash
export DOCKLE_AUTH_URL=https://registry.hub.docker.com
export DOCKLE_USERNAME={DOCKERHUB_USERNAME}
export DOCKLE_PASSWORD={DOCKERHUB_PASSWORD}
```

- NOTE: You don't need to set ENV vars when downloading from the public repository.

### Amazon ECR (Elastic Container Registry)

`Dockle` uses the AWS SDK. You don't need to install `aws` CLI tool.

Use [AWS CLI's ENVIRONMENT variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

```bash
export AWS_ACCESS_KEY_ID={AWS ACCESS KEY}
export AWS_SECRET_ACCESS_KEY={SECRET KEY}
export AWS_DEFAULT_REGION={AWS REGION}
```

### GCR (Google Container Registry)

`Dockle` uses the Google Cloud SDK. So, you don't need to install `gcloud` command.

If you'd like to use the target project's repository, you can settle via `GOOGLE_APPLICATION_CREDENTIAL`.

```bash
# must set DOCKLE_USERNAME empty char
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credential.json
```

### Self Hosted Registry (BasicAuth)

BasicAuth server needs `DOCKLE_USERNAME` and `DOCKLE_PASSWORD`.

```bash
export DOCKLE_USERNAME={USERNAME}
export DOCKLE_PASSWORD={PASSWORD}

# if you'd like to use 80 port, use NonSSL
export DOCKLE_NON_SSL=true
```