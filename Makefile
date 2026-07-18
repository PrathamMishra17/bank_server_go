

mysql-start:
	docker run --name mysql-docker -p 3307:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:latest
mysql-createdb:
	docker exec -it mysql-docker mysql -u  root    -ppassword  -e "CREATE DATABASE simple_bank"
mysql-dropdb:
	docker exec -it mysql-docker mysql -u root -ppassword -e "DROP DATABASE IF EXISTS simple_bank"
migrate-up:
	migrate -path db/migration -database "mysql://root:password@tcp(127.0.0.1:3307)/simple_bank" -verbose up
migrate-down:
	migrate -path db/migration -database "mysql://root:password@tcp(127.0.0.1:3307)/simple_bank" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: mysql-start mysql-createdb mysql-dropdb migrate-up migrate-down sqlc test