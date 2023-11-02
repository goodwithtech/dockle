
# Checkpoint Details

## Docker Image Checkpoints

These checkpoints referred to [CIS Docker 1.13.0 Benchmark v1.0.0](https://www.cisecurity.org/benchmark/docker/).

### CIS-DI-0001
**Create a user for the container**

> Create a non-root user for the container in the Dockerfile for the container image.
>
> It is a good practice to run the container as a non-root user, if possible.

```
# Dockerfile
RUN useradd -d /home/dockle -m -s /bin/bash dockle
USER dockle

or

RUN addgroup -S dockle && adduser -S -G dockle dockle
USER dockle

```

### CIS-DI-0002
**Use trusted base images for containers**

Dockle checks [Content Trust](https://docs.docker.com/engine/security/trust/content_trust/).

### CIS-DI-0003
**Do not install unnecessary packages in the container**

Not supported.

### CIS-DI-0004
**Scan and rebuild the images to include security patches**

Not supported.
Please check with [Trivy](https://github.com/knqyf263/trivy).

### CIS-DI-0005
**Enable Content trust for Docker**

> Content trust is disabled by default. You should enable it.

```bash
$ export DOCKER_CONTENT_TRUST=1
```

- https://docs.docker.com/engine/security/trust/content_trust/#about-docker-content-trust-dct

    > Docker Content Trust (DCT) provides the ability to use digital signatures for data sent to and received from remote Docker registries.<br/>
    > Engine Signature Verification prevents the following:
    >
    >   - `$ docker container run` of an unsigned image.
    >   - `$ docker pull` of an unsigned image.
    >   - `$ docker build` where the FROM image is not signed or is not scratch.

### CIS-DI-0006
**Add `HEALTHCHECK` instruction to the container image**

> Add `HEALTHCHECK` instruction in your docker container images to perform the health check on running containers.<br/>
> Based on the reported health status, the docker engine could then exit non-working containers and instantiate new ones.

```
# Dockerfile
HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost/ || exit 1
```

### CIS-DI-0007
**Do not use `update` instructions alone in the Dockerfile**

> Do not use `update` instructions such as `apt-get update` alone or in a single line in the Dockerfile.<br/>
> Adding the `update` instructions in a single line on the Dockerfile will cache the update layer.

```bash
RUN apt-get update && apt-get install -y package-a
```

### CIS-DI-0008
**Confirm safety of `setuid` and `setgid` files**

> Removing `setuid` and `setgid` permissions in the images would prevent privilege escalation attacks in the containers.<br/>
> `setuid` and `setgid` permissions could be used for elevating privileges.

```bash
chmod u-s setuid-file
chmod g-s setgid-file
```

### CIS-DI-0009
**Use `COPY` instead of `ADD` in Dockerfile**

> Use `COPY` instruction instead of `ADD` instruction in the Dockerfile.<br/>
> `ADD` instruction introduces risks such as adding malicious files from URLs without scanning and unpacking procedure vulnerabilities.

```
# Dockerfile
ADD test.json /app/test.json
↓
COPY test.json /app/test.json
```

### CIS-DI-0010
**Do not store secrets in Dockerfiles**

> Do not store any secrets in Dockerfiles.<br/>
> the secrets within these Dockerfiles could be easily exposed and potentially be exploited.

`Dockle` checks ENVIRONMENT variables and credential files.

### CIS-DI-0011
**Install verified packages only**

Not supported.
It's better to use [Trivy](https://github.com/knqyf263/trivy).

## Dockle Checkpoints for Docker

These checkpoints referred to [Docker Best Practice](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/) and so on.

### DKL-DI-0001
**Avoid `sudo` command**

- https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user

    > Avoid installing or using sudo as it has unpredictable TTY and signal-forwarding behavior that can cause problems.

### DKL-DI-0002
**Avoid sensitive directory mounting**

A volume mount makes weak points. This depends on mounting volumes.

Currently, `Dockle` checks following directories:

 - `/dev`, `/proc`, `/sys`

`dockle` only checks `VOLUME` statements, since we can't check `docker run -v /lib:/lib ...`.


### DKL-DI-0003
**Avoid `apt-get dist-upgrade`**

https://github.com/docker/docker.github.io/pull/12571

~~https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get~~
~~Avoid `RUN apt-get upgrade` and `dist-upgrade`, as many of the “essential” packages from the parent images cannot upgrade inside an unprivileged container.~~

### DKL-DI-0004
**Use `apk add` with `--no-cache`**

- https://github.com/gliderlabs/docker-alpine/blob/master/docs/usage.md#disabling-cache

    > As of Alpine Linux 3.3 there exists a new `--no-cache` option for `apk`. It allows users to install packages with an index that is updated and used on-the-fly and not cached locally:<br/>
    > ...<br/>
    > This avoids the need to use `--update` and remove `/var/cache/apk/*` when done installing packages.

### DKL-DI-0005
**Clear `apt-get` caches**

Use `apt-get clean && rm -rf /var/lib/apt/lists/*` after `apt-get install`.

- https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get

    > In addition, when you clean up the `apt cache` by removing `/var/lib/apt/lists` it reduces the image size, since the apt cache is not stored in a layer. Since the `RUN` statement starts with `apt-get update`, the package cache is always refreshed prior to `apt-get install`.

### DKL-DI-0006
**Avoid `latest` tag**

- https://vsupalov.com/docker-latest-tag/

    > Docker images tagged with `:latest` have caused many people a lot of trouble.

## Dockle Checkpoints for Linux

These checkpoints referred to [Linux Best Practices](https://www.cyberciti.biz/tips/linux-security.html) and so on.

### DKL-LI-0001
**Avoid empty password**

- https://blog.aquasec.com/cve-2019-5021-alpine-docker-image-vulnerability

    > CVE-2019-5021: Alpine Docker Image "null root password" Vulnerability

### DKL-LI-0002
**Be unique UID/GROUPs**

- http://www.linfo.org/uid.html

    > Contrary to popular belief, it is not necessary that each entry in the UID field be unique. However, non-unique UIDs can cause security problems, and thus UIDs should be kept unique across the entire organization.

### DKL-LI-0003
**Only put necessary files**

Check `.cache`, `.git` and so on directories.
