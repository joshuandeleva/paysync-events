postgres:
	docker run --name postgres16 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Mwag9836 -d -p 5432:5432 postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres paysync-events
dropdb:
	docker exec -it postgres16 dropdb --username=postgres --owner=postgres paysync-events
test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: server createdb dropdb test postgres