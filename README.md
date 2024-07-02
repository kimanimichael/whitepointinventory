# White Point Inventory
White Point inventory is an inventory service for management of supplier transactions for White Point limited company. This aims to solve inventory management problems while also provide insights from data collected from previous transactions.

## Quick Start
For local development:

### Prerequisites
The following are required to run this service:

1. [Go](https://go.dev/doc/install) (version 1.21)
2. [PostgreSQL](https://www.postgresql.org/download/) (version 14.12)

Once prerequisites are installed:

### Database Creation

Create a database and name it appropriately. Check this [reference](https://github.com/Mike-Kimani/whitepointinventory/blob/master/.env#L2) .The database is named `whitepointinventory` in this case

### Apply all Database Migrations
``
cd sql/schema
``

``
goose postgres postgres://{userName}:{password}@localhost:5432/{databaseName} up
``

### Build and Start the Server
``
go build && ./whitepointinventory
``

## Documentation
The documentation is found [here](documentation)

