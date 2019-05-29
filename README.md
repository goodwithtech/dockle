# Lyon
A Simple Security and Filesystem auditing tool for Containers, Suitable for CI


# checkpoints

- manifest parse
  - [ ] Use absolute WORKDIR.
  - [x] Last user should not be root
  - [ ] Use the --no-cache switch to avoid the need to use --update and remove /var/cache/apk/* when done installing packages
  - [ ] log to STDERR
  
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
  - [ ] check mount points : Never
  - [ ] check package cache files : 
- Check /etc/hosts
  - [ ] duplicates
  - [ ] hostname
  - [ ] localhost
- Packages
  - [ ] Package managers
- Networking
  - [ ] Check listening ports
    - /etc/services : all port
- File Permissions
  - [ ] Insecure permission
- Image Size
  - [ ] check large size container
- Container Tag
  - [ ] check `latest` tag
  - [ ] Avoid `latest` in base container
- Check PHP ini file
