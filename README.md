# drip-backend

[![Maintainability](https://api.codeclimate.com/v1/badges/252814ca6aba27f4dc3d/maintainability)](https://codeclimate.com/repos/62c905f4ef60b5563400002b/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/252814ca6aba27f4dc3d/test_coverage)](https://codeclimate.com/repos/62c905f4ef60b5563400002b/test_coverage)
[![CI](https://github.com/dcaf-labs/drip-backend/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/dcaf-labs/drip-backend/actions/workflows/ci.yaml)

| 	            | Devnet                                                                                                                                      	 | Mainnet                                                                                                                                 	 |
|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| Staging    	 | [![Better Uptime Badge](https://betteruptime.com/status-badges/v1/monitor/isy8.svg)](https://betteruptime.com/?utm_source=status_badge)     	 | n/a                                                                                                                                     	 |
| Production 	 | [ ![Better Uptime Badge](https://betteruptime.com/status-badges/v1/monitor/g7cf.svg) ]( https://betteruptime.com/?utm_source=status_badge ) 	 | [![Better Uptime Badge](https://betteruptime.com/status-badges/v1/monitor/goyh.svg)](https://betteruptime.com/?utm_source=status_badge) 	 |

## Deploy Process 

> TODO: Update production env's via a github action 

| 	            | Devnet                                                                                                                                      	 | Mainnet                                                                                                                                 	 |
|--------------|-----------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| Staging    	 | Merge to `main`    	                                                                                                                          | n/a                                                                                                                                     	 |
| Production 	 | Rebase `devnet` with `main` 	                                                                                                                 | Rebase `mainnet` with `main` 	                                                                                                            |

## Setup

- [^1] install go 1.18 
- [^2] install all packages with
- [^3] install the `oapi-codegen` cli, version `v1.10.1`
- [^4] install `docker` and `docker-compose`
- [^5] setup a `.env` file

## Start the Server
- [^6] start the database
- [^8] start the api and event server
  - This also runs db migrations  

## Tests

- [^13] run all tests
- [^14] If an interface is changed, the associated mocks need to be re-generated

## API Docs

- [^3] use [oapi-codegen](https://github.com/deepmap/oapi-codegen) to code gen our api go types.
- [^15] if the api docs are updated, we need to update the api codegen'd files
- API docs are viewable at `http://localhost:8080/swagger.json`.

## Database

- [^6] spin up a local database via docker
- [^7] stop the local docker database

> **_NOTE:_** TODO(mocha) do codegen in a dockerized db that is setup and destroyed automatically in the script.

### Migrations

- located in `internal/database/psql/migrations`
- all migrations that are a number larger than what the db version is will be run automatically on startup or during codegen
- [^12] running the migrations will automatically increment the schema version in the db

> **_NOTE:_** The [^6]database must be running prior to running this script.

### Codegen

- Database [models](app/internal/data/psql/generated) are generated using the database schema via [go-gorm/gen](https://github.com/go-gorm/gen)
- Before generating the models, the [^6]database needs to be running
- The [^11]codegen script will also run migrations if new migrations are detected 

```bash
go run cmd/codegen/main.go
```

> **_NOTE:_** The [^6]database must be running prior to running this script.

### Process for Creating/Updating Database Models

- create a migration file under `internal/database/psql/migrations`, and name it appropriately (ex. `2_new_migration.up.sql`)
- [^11] Run the migration + codegen script 

> **_NOTE:_** The [^6]database must be running prior to running this script.

### Manually backfill the DB with known data/addresses

- [^10] Backfilling is done via the backfill script
- View the code to see how it is done

[^1]: Install go 
    [brew formula](https://formulae.brew.sh/formula/go)
    Verify with `go version`

[^2]: Install all go packages
    `go get -u ./...`

[^3]: Install the oapi-codegen cli
    `go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen`

[^4]: Install Docker/Docker-compose
    [Docs](https://docs.docker.com/compose/install/)

[^5]: Example `.env`
    ```env

    # Devnet Staging
    ENV="STAGING"
    NETWORK="DEVNET"
    DRIP_PROGRAM_ID="F1NyoZsUhJzcpGyoEqpDNbUMKVvCnSXcCki1nN3ycAeo"
    # random wallet, this is the mint auth
    DRIP_BACKEND_WALLET="[141,241,173,131,255,186,170,216,65,246,24,196,173,94,39,225,161,108,251,102,177,20,166,223,13,69,103,38,242,107,72,194,177,170,44,204,179,183,235,4,231,51,88,169,156,153,132,247,235,166,41,123,87,219,139,204,95,1,176,98,72,90,51,82]"
    
    # Only for local testing, this is false on all deployed env's
    SHOULD_BYPASS_ADMIN_AUTH=true 
    
    # Test Channel
    DISCORD_WEBHOOK_ID=1021592812954857492
    DISCORD_ACCESS_TOKEN=qPeOyI4e4k6kYah44k9_PXFQDsuLO7lbHcazLrsKcvzqvrQh_lr1PK21kB3GZCSTv2Xg
    
    PORT="8080"
    PSQL_USER="dcaf"
    PSQL_PASS="drip"
    PSQL_DBNAME="drip"
    PSQL_PORT="5432"
    PSQL_HOST="localhost"
    PSQL_SSLMODE=disable
    OUTPUT=./internal/pkg/repository/models
    GOOGLE_KEY="540992596258-sa2h4lmtelo44tonpu9htsauk5uabdon.apps.googleusercontent.com"
    GOOGLE_SECRET="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"
    ```
[^6]: Start the database via `docker-compose`
    `docker-compose --file ./build/docker-compose.yaml  --env-file ./.env up`

[^7]: Stop the database via `docker-compose`
    `docker-compose --file ./build/docker-compose.yaml  --env-file ./.env down`

[^8]: Start the api-server and event server
    `go run main.go`

[^9]: Start just the event server
    `go run cmd/event/main.go`

[^10]: Run the backfill data script
    `go run cmd/backfill/main.go`

[^11]: Run migrations and codegen db models
    `go run cmd/codegen/main.go`

[^12]: Just run db migrations
    `go run cmd/migrate/main.go`

[^13]: Run unit and integration tests
    `go test ./...`

[^14]: Update mock files
    `./scripts/create-mocks.sh`

[^15]: Update api-docs codegen'd files
    `oapi-codegen -package apispec ./docs/swagger.yaml > pkg/api/apispec/generated.go`

[//]: # ()
[//]: # (```bash)

[//]: # (# Add to ZSHRC)

[//]: # (export GOPRIVATE=github.com/dcaf-labs/solana-go-clients)

[//]: # (go get -u ./...)

[//]: # (```)

[//]: # ()
[//]: # (Add the following to ~/.gitconfig )

[//]: # ()
[//]: # (```txt)

[//]: # ([url "ssh://git@github.com/"])

[//]: # (insteadOf = https://github.com/)

[//]: # (```)