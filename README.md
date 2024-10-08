# Online Song Library - Test Task

This repository is a test task implementing an **Online Song Library**. It allows you to manage songs, including creating, updating, deleting, and retrieving song data. It also includes features like song text pagination and filtering.

## Table of Contents

- [Task Requirements](#task-requirements)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Testing](#testing)
- [API Endpoints](#api-endpoints)

## Task Requirements

The test task includes the following functionalities:
- Create a new song
- Update an existing song
- Delete a song
- Retrieve a list of songs with filtering by group name, song title, and release date
- Pagination of song lyrics by verses

## Prerequisites

Before running the application, ensure you have the following installed:

- **Go**: Version 1.18 or higher
- **PostgreSQL**: Version 14 or higher
- **Docker** (optional, for running PostgreSQL in a container)
- **Make** (optional, for simplified commands)

### PostgreSQL Setup

If PostgreSQL is not already installed locally, you can run it in Docker, but make sure you fill the .env file:

```bash
make postgres.start
```

.env.example file contains examples of variables that must be filled in

## Installation

Clone the repository
```bash 
git clone https://github.com/salmon822/online-song-library-test.git
cd online-song-library-test
```

Install Go dependencies:

```bash
go mod tidy
```

## Running the Application

Before you can run the application, you need to either run the container with PostgreSQL or run it locally without the container.

Running PostgreSQL in a Docker container:
```bash
make postgres.start
```

Also you have to generate models from Swagger:
```bash
make build-models
```

To run the application locally:
```bash
go run cmd/main.go -cfg configs/local.json
```
OR
```bash
make run
```
This will start the server on http://localhost:8080.

## Testing

The repository includes integration tests.

Before you can run tests, you need to either run the container with PostgreSQL for Tests or run it locally without the container.

Running PostgreSQL for Tests in a Docker container:
```bash
make postgres.test.start
```

Also make sure you fill the .env file in /integrations_tests directory, this example also presents in .env.example

To run the tests:
```bash
make tests
```

## API Endpoints

Base URL: http://localhost:8080
- POST /songs/create: Create a new song.
- GET /songs/filter: Retrieve a list of songs with filtering and pagination.
- GET /songs/{id}: Get song text by its ID.
- PATCH /songs/{id}/update: Update an existing song by its ID.
- DELETE /songs/{id}/delete: Delete a song by its ID.


## Notes

Make sure PostgreSQL is running before you start the application.
Same works for tests.
Ensure that models from OpenAPI are generated.
