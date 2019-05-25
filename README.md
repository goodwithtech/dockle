# Lyon
A Simple Security and Filesystem auditing tool for Containers, Suitable for CI


# checkpoints

- General
  - [ ] detect os
  - [ ] use official container on the base
- Do not write secrets
  - [ ] check ENV vars
    - credential information
    - service environment
      - not includes production, stage, dev
  - [ ] check credential files
- Users, Groups and Authentication
  - [ ] Default user not a root user
  - [ ] Administrator accounts
  - [ ] Unique UIDs
  - [ ] Unique group IDs
  - [ ] Unique group names
  - [ ] Unnecessary priviledge escalation(setuid, setgid)
- File systems
  - [ ] Check /tmp
  - [ ] Check /var/tmp
  - [ ] check mount points
  - [ ] check package cache files
- Check /etc/hosts
  - [ ] duplicates
  - [ ] hostname
  - [ ] localhost
- Packages
  - [ ] Package managers
- Networking
  - [ ] Check listening ports
- File Permissions
  - [ ] Insecure permission
- Processes
  - [ ] Single CMD
  - [ ] Single ENTRYPOINT
- Image Size
  - [ ] check large size container
- Container Tag
  - [ ] check `latest` tag
  - [ ] Avoid `latest` in base container
- Check PHP ini file
