# LDAP Test

ldap crud

start a ldap sever

```yaml
version: "3.8"
services:
  openldap:
    image: osixia/openldap:1.5.0
    container_name: openldap-server
    ports:
      - "389:389" # LDAP standard port
      - "636:636" # LDAPS (secure) port
    environment:
      LDAP_ORGANISATION: "My Awesome Company"
      LDAP_DOMAIN: "example.org"
      LDAP_ADMIN_PASSWORD: "adminpassword"
    volumes:
      - /mnt/user/appdata/work/ldap/data:/var/lib/ldap
      - /mnt/user/appdata/work/ldap/config:/etc/ldap/slapd.d

  phpldapadmin:
    image: osixia/phpldapadmin:0.9.0
    container_name: phpldapadmin-ui
    ports:
      - "8801:80" # Web UI port
    environment:
      PHPLDAPADMIN_LDAP_HOSTS: "openldap-server"
      PHPLDAPADMIN_HTTPS: "false"
    depends_on:
      - openldap
```
