* fun fact, pentru pgadmin cand te conectezi pe pgadmin trebuie sa cauti ip ul local din podman -> postgres -> inspect

* podman:
* dpage/pgadmin4 in podman, pull request
* podman pull postgres
* cand creezi containerul pentru pgadmin trebuie sa ai:

"PGADMIN_DEFAULT_EMAIL=iosefgabriel268@gmail.com",
"PGADMIN_DEFAULT_PASSWORD=admin",

* ports 443:9000 / 80:9001

* container postgres
* POSTGRES_PASSWORD=admin
* POSTGRES_USER=admin
