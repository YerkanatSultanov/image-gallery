postgresinit:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine
postgres:
	docker exec -it postgres16 bin/bash

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres authDatabase

.PHONY: postgresinit postgres createdb
