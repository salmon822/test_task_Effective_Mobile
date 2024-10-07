.PHONY: run

SWAGGER_FILE = ../swagger.yaml
OUTPUT_DIR = models

run:
	go run cmd/main.go -cfg configs/local.json

postgres.start:
	if [ ! "$(shell docker ps -q -f name=song-library-postgres)" ]; then \
        if [ "$(shell docker ps -aq -f status=exited -f name=song-library-postgres)" ]; then \
            docker rm song-library-postgres; \
        fi; \
        docker run --restart unless-stopped -d -p 2264:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=1337 --name song-library-postgres postgres:16 ;\
        sleep 10; \
    fi;
	-docker exec song-library-postgres psql -U admin -c "create database song_library"
	-docker exec song-library-postgres psql -U admin -c "grant all privileges on database song_library to admin"

postgres.stop:
	docker stop song-library-postgres
	docker rm song-library-postgres

postgres.restart:
	make postgres.stop
	make postgres.start
	make migrate.up

migrate.up:
	goose -dir migrations -allow-missing postgres "host=localhost port=2264 user=admin password=1337 dbname=song_library sslmode=disable" up

migrate.down:
	goose -dir migrations postgres "host=localhost port=2264 user=admin password=1337 dbname=song_library sslmode=disable" down

migrate.status:
	goose -dir migrations postgres "host=localhost port=2264 user=admin password=1337 dbname=song_library sslmode=disable" status


build-models:
	oapi-codegen -generate types -o $(OUTPUT_DIR)/types.gen.go -package models $(SWAGGER_FILE)

clean-models:
	rm -f $(OUTPUT_DIR)/types.gen.go 	