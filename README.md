# White Point Inventory
White Point inventory is an inventory service for management of supplier transactions for White Point limited company. This aims to solve inventory management problems while also provide insights from data collected from previous transactions.

## Prerequisites
The following are required to run this service:

1. [Docker](https://docs.docker.com/get-docker/) (version 26.0.0)

To develop whitepointinventory, you will need:

1. [Go](https://go.dev/doc/install) (version 1.21)
2. [PostgreSQL](https://www.postgresql.org/download/) (version 14.12)

## Installation
Ensure all  pre-requisites are satisfied before carrying out installation and clone the repo

### Using Docker
Execute this command from the project root:

`docker compose -f docker/docker-compose.yml up`

### Without Docker
#### Database Creation

Create a PostgreSQL database and name it appropriately. Check this [reference](https://github.com/Mike-Kimani/whitepointinventory/blob/master/.env#L2) .The database is named `whitepointinventory` in this case

#### Apply all Database Migrations
``
cd sql/schema
``

``
goose postgres postgres://{userName}:{password}@localhost:5432/{databaseName} up
``

#### Build and Start the Server
``
go build && ./whitepointinventory
``

## Documentation
The documentation is found [here](documentation)

