postgresinit:
	docker run --name postgres1 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine
postgres:
	docker exec -it postgres1 bin/bash

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres postgres16

dropdb:
	docker exec -it postgres16 postgres16

.PHONY: postgresinit postgres createdb dropdb
