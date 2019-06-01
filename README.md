# docker-guard
A Simple Security and Filesystem auditing tool for Containers, Suitable for CI


# checkpoints

- manifest parse
  - [x] Last user should not be root
  - [ ] apk : Use the --no-cache
  - [ ] dpkg : Use the --no-install-recommends
  - [ ] log to STDERR
  - [x] avoid mount sensitive dir 
  
- General
  - [ ] detect os
  - [ ] use official container on the base (need db) : Future
  - [x] use Docker Content Trust
- Do not write secrets
  - [x] check ENV vars
    - credential information
    - service environment
      - not includes production, stage, dev
  - [ ] check credential files
- Users, Groups and Authentication
  - [x] Default user not a root user
  - [x] Set password
  - [x] Unique UIDs
  - [x] Unique group names
  - [ ] Unnecessary priviledge escalation(setuid, setgid) : Future support
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
- File systems
  - [ ] Check /tmp : Future
  - [ ] Check /var/tmp : Future
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
- Container Tag
  - [ ] check `latest` tag
  - [ ] Avoid `latest` in base container
- Check PHP ini file
  
  
## if docker running
- Networking
  - [ ] `docker port container` if docker running
  - [ ] by file
    - /proc/1/net/tcp : openning port (if running)
- Volume mount
  - mount dangerous 
    - /boot, /dev, /etc, /lib
