# docker-guard
A Simple scanner security and tool for Containers, Suitable for CI

# Abstract

`docker-guard` is 
1) simple security auditing tool that helps you build secure Docker images
2) check a docker configuration tool that helps you build `best practice` Docker images 

# How to use

```bash
guard [YOUR_IMAGE_NAME]
```

# Installation

## Mac OS X / Homebrew

You can use homebrew on Mac OS.

```
$ brew install goodwithtech/docker-guard
```

## Binary (Including Windows)

Get the latest version from [this page](https://github.com/goodwithtech/docker-guard/releases/latest), and download the archive file for your operating system/architecture. Unpack the archive, and put the binary somewhere in your `$PATH` (on UNIX-y systems, /usr/local/bin or the like). Make sure it has execution bits turned on.


## From source

```sh
$ go get -u github.com/goodwithtech/docker-guard
```


## Docker

Replace [YOUR_CACHE_DIR] with the cache directory on your machine.

```
$ docker run --rm -v [YOUR_CACHE_DIR]:/root/.cache/ goodwithtech/docker-guard [YOUR_IMAGE_NAME]
```

Example for macOS:

```
$ docker run --rm -v $HOME/Library/Caches:/root/.cache/ goodwithtech/docker-guard [YOUR_IMAGE_NAME]
```

If you would like to scan the image on your host machine, you need to mount `docker.sock`.

```
$ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ...
```

Please re-pull latest `goodwithtech/docker-guard` if an error occured.

# Checkpoints

## Security Checkpoints

### SC0001 : All user/group should set password

https://blog.aquasec.com/cve-2019-5021-alpine-docker-image-vulnerability

> CVE-2019-5021: Alpine Docker Image ‘null root password’ Vulnerability
> 

### SC0002 : Last user should not be root

[hadolint:DL3002](https://github.com/hadolint/hadolint/wiki/DL3002)

<details>
<summary>Correct code:</summary>

```
FROM busybox
USER root
RUN ...
USER guest
```

</details>

### SC0003 : Avoid sensitive directory mounting

A volume mount makes weakpoints. 
This depends on mounting volumes.
Currently, docker-guard check following directories.

`/boot`,`/dev`,`/etc`,`/lib','/proc`,`/sys`, `/usr`

docker-guard only checks `VOLUME`. We can't check `docker run -v /lib:/lib ...`.

### SC0004 : use DOCKER CONTENT TRUST

> Docker Content Trust (DCT) provides the ability to use digital signatures for data sent to and received from remote Docker registries.
> Engine Signature Verification prevents the following:
> - `$ docker container run` of an unsigned image.
> - `$ docker pull` of an unsigned image.
> - `$ docker build` where the FROM image is not signed or is not scratch.
>> https://docs.docker.com/engine/security/trust/content_trust/#about-docker-content-trust-dct

<details>
<summary>How to set :</summary>

```
EXPORT DOCKER_CONTENT_TRUST=1 
```
</details>

### SC0005 : Don’t store credentials in the image

Images should be cupsule. Image shouldn't have state.
All variables run via `docker run -env XXXXX`.


### SC0006 : Unique UID/GROUPs

http://www.linfo.org/uid.html

> Contrary to popular belief, it is not necessary that each entry in the UID field be unique. However, non-unique UIDs can cause security problems, and thus UIDs should be kept unique across the entire organization.

## Dockerfile Checkpoints

### DC0001 : Avoid `apt-get upgrade`, `apk upgrade`, `dist-upgrade`

https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get
 
> Avoid RUN apt-get upgrade and dist-upgrade, as many of the “essential” packages from the parent images cannot upgrade inside an unprivileged container.


### DC0002 : Avoid `sudo` command

https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
> Avoid installing or using sudo as it has unpredictable TTY and signal-forwarding behavior that can cause problems.


### DC0003 : Use apk add with `--no-cache`

https://github.com/gliderlabs/docker-alpine/blob/master/docs/usage.md#disabling-cache

> As of Alpine Linux 3.3 there exists a new --no-cache option for apk. It allows users to install packages with an index that is updated and used on-the-fly and not cached locally:
> This avoids the need to use --update and remove /var/cache/apk/* when done installing packages.

### DC0004 : Use apt-get minimize

Use “apt-get clearn && rm -rf /var/lib/apt/lists/*` if use apt-get install

https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get
> In addition, when you clean up the apt cache by removing /var/lib/apt/lists it reduces the image size, since the apt cache is not stored in a layer. Since the RUN statement starts with apt-get update, the package cache is always refreshed prior to apt-get install.


### DC0005 : Avoid `latest` tag

https://vsupalov.com/docker-latest-tag/

> Docker images tagged with :latest have caused many people a lot of trouble.
  
  
# Examples

## Scan an image

## Scan an image file

## Specify exit code
By default, `docker-guard` exits with code 0 even if there is some problems.
Use the --exit-code option if you want to exit with a non-zero exit code.

```bash
$ guard  -exist-code 1 [IMAGE_NAME]
```

## Ignore the specified rules

Use `.guardignore`.

```bash
$ cat .guardignore
# set root to default user because we want to run nginx
SC0001
# Use latest tag because only check for image inside
DC0005
```

### Clear image caches

The `--clear-cache` option removes image caches. This option is useful if the image which has the same tag is updated (such as when using `latest` tag).

```
$ guard --clear-cache python:3.7
```

# Continuous Integration (CI)

Scan your image built in Travis CI/CircleCI. 
The test will fail if a vulnerability is found. 
When you don't want to fail the test, specify `--exit-code 0`.





# Roadmap

- Users, Groups and Authentication
  - [ ] Unnecessary priviledge escalation(setuid, setgid)
    ```
		fi := hdr.FileInfo()
		fm := fi.Mode()
		if fm&os.ModeSetuid != 0 {
		    // suid
		}
		if fm&os.ModeSetgid != 0 {
			// gid
		}
    ```
- General
  - [ ] detect os
  - [ ] use official container on the base

- [ ] Check php.ini file
- [ ] Check nginx.conf file
- [ ] log to STDERR
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
    - /proc/1/net/tcp : openning port (if running)
- Volume mount
  - mount dangerous 
    - /boot, /dev, /etc, /lib
