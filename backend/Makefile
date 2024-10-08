.PHONY: run
include .env
export $(shell sed 's/=.*//' .env)

SWAGGER_FILE = ../swagger.yaml
OUTPUT_DIR = models

run:
	go run cmd/main.go -cfg configs/local.json

tests:
	go test -v ./integration_tests/...

postgres.start:
	if [ ! "$(shell docker ps -q -f name=$(POSTGRES_CONTAINER))" ]; then \
        if [ "$(shell docker ps -aq -f status=exited -f name=$(POSTGRES_CONTAINER))" ]; then \
            docker rm $(POSTGRES_CONTAINER); \
        fi; \
        docker run --restart unless-stopped -d -p $(POSTGRES_PORT):5432 \
            -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
            --name $(POSTGRES_CONTAINER) postgres:$(POSTGRES_VERSION) ;\
        sleep 10; \
    fi;
	-docker exec $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -c "create database $(POSTGRES_DB)"
	-docker exec $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -c "grant all privileges on database $(POSTGRES_DB) to $(POSTGRES_USER)"

postgres.test.start:
	if [ ! "$(shell docker ps -q -f name=$(POSTGRES_TEST_CONTAINER))" ]; then \
        if [ "$(shell docker ps -aq -f status=exited -f name=$(POSTGRES_TEST_CONTAINER))" ]; then \
            docker rm $(POSTGRES_TEST_CONTAINER); \
        fi; \
        docker run --restart unless-stopped -d -p $(POSTGRES_TEST_PORT):5432 \
            -e POSTGRES_USER=$(POSTGRES_TEST_USER) -e POSTGRES_PASSWORD=$(POSTGRES_TEST_PASSWORD) \
            --name $(POSTGRES_TEST_CONTAINER) postgres:$(POSTGRES_TEST_VERSION) ;\
        sleep 10; \
    fi;
	-docker exec $(POSTGRES_TEST_CONTAINER) psql -U $(POSTGRES_TEST_USER) -c "create database $(POSTGRES_TEST_DB)"
	-docker exec $(POSTGRES_TEST_CONTAINER) psql -U $(POSTGRES_TEST_USER) -c "grant all privileges on database $(POSTGRES_TEST_DB) to $(POSTGRES_TEST_USER)"


postgres.stop:
	docker stop $(POSTGRES_CONTAINER)
	docker rm $(POSTGRES_CONTAINER)

postgres.test.stop:
	docker stop $(POSTGRES_TEST_CONTAINER)
	docker rm $(POSTGRES_TEST_CONTAINER)	

postgres.restart:
	make postgres.stop
	make postgres.start
	make migrate.up

migrate.up:
	goose -dir migrations -allow-missing postgres "host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" up

migrate.down:
	goose -dir migrations postgres "host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" down

migrate.status:
	goose -dir migrations postgres "host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" status

build-models:
	oapi-codegen -generate types -o $(OUTPUT_DIR)/types.gen.go -package models $(SWAGGER_FILE)

clean-models:
	rm -f $(OUTPUT_DIR)/types.gen.go 	